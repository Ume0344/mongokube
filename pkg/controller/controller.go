package controller

import (
	"fmt"
	"time"

	mkclientset "mongokube/pkg/client/clientset/versioned"
	mkinformers "mongokube/pkg/client/informers/externalversions/mongokube/beta1"
	mklister "mongokube/pkg/client/listers/mongokube/beta1"

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

	return true
}
