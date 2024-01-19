/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package beta1

import (
	"context"
	mongokubebeta1 "mongokube/pkg/apis/mongokube/beta1"
	versioned "mongokube/pkg/client/clientset/versioned"
	internalinterfaces "mongokube/pkg/client/informers/externalversions/internalinterfaces"
	beta1 "mongokube/pkg/client/listers/mongokube/beta1"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// MkInformer provides access to a shared informer and lister for
// Mks.
type MkInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() beta1.MkLister
}

type mkInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewMkInformer constructs a new informer for Mk type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewMkInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredMkInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredMkInformer constructs a new informer for Mk type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredMkInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MongokubeBeta1().Mks(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MongokubeBeta1().Mks(namespace).Watch(context.TODO(), options)
			},
		},
		&mongokubebeta1.Mk{},
		resyncPeriod,
		indexers,
	)
}

func (f *mkInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredMkInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *mkInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&mongokubebeta1.Mk{}, f.defaultInformer)
}

func (f *mkInformer) Lister() beta1.MkLister {
	return beta1.NewMkLister(f.Informer().GetIndexer())
}
