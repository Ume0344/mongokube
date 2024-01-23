package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	mkclientset "mongokube/pkg/client/clientset/versioned"
	mkinformers "mongokube/pkg/client/informers/externalversions"

	"mongokube/pkg/controller"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config := getConfig()

	k8sclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting k8sclient, %s", err.Error())
	}

	listDeployments, err := k8sclient.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	for _, d := range listDeployments.Items {
		fmt.Printf("name:%v\n", d.Name)
	}

	mkclient, err := mkclientset.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error getting mkclient, %s", err.Error())
	}

	mkinformers := mkinformers.NewSharedInformerFactory(mkclient, 10*time.Minute)

	c := controller.NewController(*k8sclient, mkclient, mkinformers.Mongokube().Beta1().Mks())

	channel := make(chan struct{})

	mkinformers.Start(channel)

	c.Run(channel)
}

func getConfig() *rest.Config {
	// This function set the configuration for kubernetes
	var kubeconfigpath *string

	// create filepath of kube config file which is at /home/apmec/.kube/config
	if home := homedir.HomeDir(); home != "" {
		kubeconfigpath = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfigpath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// creates configuration based on config path
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigpath)
	if err != nil {
		fmt.Printf("Could not get the config file due to %s", err.Error())

		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error %s, getting incluster config", err.Error())
		}
	}
	return config
}
