package controller

import (
	"context"
	"fmt"
	"time"

	"mongokube/pkg/apis/mongokube/beta1"
	mkclientset "mongokube/pkg/client/clientset/versioned"
	mkinformers "mongokube/pkg/client/informers/externalversions/mongokube/beta1"
	mklister "mongokube/pkg/client/listers/mongokube/beta1"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	port = "50051"
)

// Controller Struct which has attributes k8s standard clientset, Mk generated clientset
// generated lister, cache and workqueue
type Controller struct {
	k8sclient   kubernetes.Clientset
	mkClient    mkclientset.Interface
	mkLister    mklister.MkLister
	mkSynched   cache.InformerSynced //if cache has been synched with api server
	mkWorkQueue workqueue.RateLimitingInterface
}

// This struct will represent the data for mongodb and mongo express service
type MongoService struct {
	name        string
	label       map[string]string
	serviceType v1.ServiceType
	port        int32
	nodePort    int32
}

// Initialize the Controller struct and add event handler for registering
// handler functions for adding and deleting Mk resources.
func NewController(
	k8sclient kubernetes.Clientset,
	mkClient mkclientset.Interface,
	mkInformer mkinformers.MkInformer,

) *Controller {
	c := &Controller{
		k8sclient:   k8sclient,
		mkClient:    mkClient,
		mkLister:    mkInformer.Lister(),
		mkSynched:   mkInformer.Informer().HasSynced,
		mkWorkQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "mongokube"),
	}

	mkInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleAdd,
		DeleteFunc: c.handleDel,
	})

	return c
}

// Add objects to queue
func (c *Controller) handleAdd(obj interface{}) {
	c.mkWorkQueue.Add(obj)
	fmt.Printf("Handling a Mk resource\n")
}

// Delete objects to queue
func (c *Controller) handleDel(obj interface{}) {
	c.mkWorkQueue.Done(obj)
	fmt.Printf("Deleting a Mk resource\n")
}

// Specifying the receiver of the method to be of type pointer to controller
// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(channel <-chan struct{}) {
	// Takes receive-only channel as argument
	// wait for the cache inside the informer to be synched before starting workers
	if !cache.WaitForCacheSync(channel, c.mkSynched) {
		fmt.Print("Waiting for cache to be synched\n")
	}

	//Create goroutine to call the worker function after every 1 second till the channel is stopped
	go wait.Until(c.worker, time.Second, channel)

	//Wait until some object is added into channel
	<-channel
}

func (c *Controller) worker() {
	// loop till processItem returns true, on false it will wait for a second and then again this function will be called by run()
	for c.processNextItem() {

	}
}

// Process the items from queue
func (c *Controller) processNextItem() bool {
	fmt.Printf("Processing the items from queue %v\n", c.mkWorkQueue.Len())
	item, shutdown := c.mkWorkQueue.Get()

	// Delete the item from queue, so that we wont process it again
	defer c.mkWorkQueue.Forget(item)

	if shutdown {
		return false
	}

	// Generating key for each item in queue
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("Getting key from cache %s\n", err.Error())
		return false
	}

	// Getting namespace, name from genrated key
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("Getting namespace and name from MetaNamespaceKeyFunc %s\n", err.Error())
		return false
	}

	mkResource, err := c.mkLister.Mks(ns).Get(name)

	if err != nil {
		fmt.Printf("Error getting Mk resource %s\n", err.Error())
		return false
	}

	// %+v for printing struct
	fmt.Printf("Mk resource specs are :%+v\n", mkResource.Spec)

	// Handle Mk resource
	deployment := c.handleMkResource(mkResource)

	if deployment {
		c.mkWorkQueue.Forget(item)
	}

	return true
}

// Handle mk resource whenever it is created and added to queue
func (c *Controller) handleMkResource(mkResource *beta1.Mk) bool {
	fmt.Printf("Creating a secret for mk resource: %s\n", mkResource.Name)
	secret, err := c.createSecret(mkResource)
	if err != nil {
		fmt.Printf("Failed to create secret: %s\n", err.Error())
	}

	fmt.Printf("Creating MongoDB deployment for mk resource: %s\n", mkResource.Name)
	deployment, err := c.createMongoDeployment(mkResource, secret)
	if err != nil {
		fmt.Printf("Failed to create deployment: %s\n", err.Error())
	}

	mongodbService := &MongoService{
		name:        "mongodb-service",
		label:       deployment.Labels,
		serviceType: v1.ServiceTypeClusterIP,
		port:        27017,
	}

	fmt.Printf("Creating MongoDB internal service for mk resource: %s\n", mkResource.Name)
	mongoDbService, err := c.createMongoService(mkResource, *mongodbService)

	if err != nil {
		fmt.Printf("Failed to create mongo db service: %s\n", err.Error())
	}

	fmt.Printf("Creating MongoExpress deployment for mk resource: %s\n", mkResource.Name)
	mongoExpressDeployment, err := c.createMongoExpressDeployment(mkResource, secret, mongoDbService)

	if err != nil {
		fmt.Printf("Failed to create mongo db service: %s\n", err.Error())
	}

	mongoExpressService := &MongoService{
		name:        "mongoexpress-service",
		label:       mongoExpressDeployment.Labels,
		serviceType: v1.ServiceTypeLoadBalancer,
		port:        8081,
		nodePort:    31000,
	}

	fmt.Printf("Creating MongoExpress external service for mk resource: %s\n", mkResource.Name)
	_, err = c.createMongoService(mkResource, *mongoExpressService)

	if err != nil {
		fmt.Printf("Failed to create mongo express service: %s\n", err.Error())
	}

	return true
}

// Create a secret for mongodb
func (c *Controller) createSecret(mkResource *beta1.Mk) (*v1.Secret, error) {
	secretData := map[string][]byte{
		"username": []byte(mkResource.Spec.DbUsername),
		"password": []byte(mkResource.Spec.DbPassword),
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mongodb-secret",
			Namespace: mkResource.Namespace,
		},
		Data: secretData,
	}
	createdSecret, err := c.k8sclient.CoreV1().Secrets(mkResource.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})

	return createdSecret, err
}

func (c *Controller) createMongoDeployment(mkResource *beta1.Mk, secret *v1.Secret) (*appsv1.Deployment, error) {
	// container data
	// label to connect with service
	replica := int32(2)
	var containerPort int32 = 27017

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mkResource.Name + "-deployment",
			Namespace: mkResource.Namespace,
			Labels:    map[string]string{"app": mkResource.Name + "db"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": mkResource.Name + "db"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": mkResource.Name + "db"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  mkResource.Name + "-container",
							Image: mkResource.Spec.MongoDbImage,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: containerPort,
								},
							},
							Env: []v1.EnvVar{
								{
									Name: "MONGO_INITDB_ROOT_USERNAME",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: secret.Name,
											},
											Key: c.getKey("username", secret),
										},
									},
								},
								{
									Name: "MONGO_INITDB_ROOT_PASSWORD",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: secret.Name,
											},
											Key: c.getKey("password", secret),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	deploymentResponse, err := c.k8sclient.AppsV1().Deployments(mkResource.Namespace).Create(context.Background(), deployment, metav1.CreateOptions{})

	return deploymentResponse, err
}

// Create mongo express deployment
func (c *Controller) createMongoExpressDeployment(mkResource *beta1.Mk, secret *v1.Secret, mongodbService *v1.Service) (*appsv1.Deployment, error) {
	// container data
	// label to connect with service
	replica := int32(2)
	var containerPort int32 = 8081

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mkResource.Name + "-express-deployment",
			Namespace: mkResource.Namespace,
			Labels:    map[string]string{"app": mkResource.Name + "express"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": mkResource.Name + "express"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": mkResource.Name + "express"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  mkResource.Name + "-express-container",
							Image: mkResource.Spec.MongoExpressImage,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: containerPort,
								},
							},
							Env: []v1.EnvVar{
								{
									Name: "ME_CONFIG_MONGODB_ADMINUSERNAME",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: secret.Name,
											},
											Key: c.getKey("username", secret),
										},
									},
								},
								{
									Name: "ME_CONFIG_MONGODB_ADMINPASSWORD",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: secret.Name,
											},
											Key: c.getKey("password", secret),
										},
									},
								},
								{
									Name:  "ME_CONFIG_MONGODB_SERVER",
									Value: mongodbService.Name,
								},
							},
						},
					},
				},
			},
		},
	}

	deploymentResponse, err := c.k8sclient.AppsV1().Deployments(mkResource.Namespace).Create(context.Background(), deployment, metav1.CreateOptions{})

	return deploymentResponse, err
}

// Get the desired key from secret
func (c *Controller) getKey(key string, secret *v1.Secret) string {
	var desiredKey string

	for k := range secret.Data {
		if k == key {
			desiredKey = k
		}
	}

	return desiredKey
}

// Create service for pods of mongodb or mongoexpress
func (c *Controller) createMongoService(mkResource *beta1.Mk, mongoStruct MongoService) (*v1.Service, error) {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: mongoStruct.name,
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceType(mongoStruct.serviceType),
			Selector: mongoStruct.label,
			Ports: []v1.ServicePort{
				{
					Port:     mongoStruct.port,
					NodePort: mongoStruct.nodePort,
				},
			},
		},
	}

	serviceCreated, err := c.k8sclient.CoreV1().Services(mkResource.Namespace).Create(context.Background(), service, metav1.CreateOptions{})

	return serviceCreated, err
}
