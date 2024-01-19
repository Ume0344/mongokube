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

// Code generated by client-gen. DO NOT EDIT.

package beta1

import (
	"context"
	beta1 "mongokube/pkg/apis/mongokube/beta1"
	scheme "mongokube/pkg/client/clientset/versioned/scheme"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MksGetter has a method to return a MkInterface.
// A group's client should implement this interface.
type MksGetter interface {
	Mks(namespace string) MkInterface
}

// MkInterface has methods to work with Mk resources.
type MkInterface interface {
	Create(ctx context.Context, mk *beta1.Mk, opts v1.CreateOptions) (*beta1.Mk, error)
	Update(ctx context.Context, mk *beta1.Mk, opts v1.UpdateOptions) (*beta1.Mk, error)
	UpdateStatus(ctx context.Context, mk *beta1.Mk, opts v1.UpdateOptions) (*beta1.Mk, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*beta1.Mk, error)
	List(ctx context.Context, opts v1.ListOptions) (*beta1.MkList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *beta1.Mk, err error)
	MkExpansion
}

// mks implements MkInterface
type mks struct {
	client rest.Interface
	ns     string
}

// newMks returns a Mks
func newMks(c *MongokubeBeta1Client, namespace string) *mks {
	return &mks{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mk, and returns the corresponding mk object, and an error if there is any.
func (c *mks) Get(ctx context.Context, name string, options v1.GetOptions) (result *beta1.Mk, err error) {
	result = &beta1.Mk{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Mks that match those selectors.
func (c *mks) List(ctx context.Context, opts v1.ListOptions) (result *beta1.MkList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &beta1.MkList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mks.
func (c *mks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a mk and creates it.  Returns the server's representation of the mk, and an error, if there is any.
func (c *mks) Create(ctx context.Context, mk *beta1.Mk, opts v1.CreateOptions) (result *beta1.Mk, err error) {
	result = &beta1.Mk{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(mk).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a mk and updates it. Returns the server's representation of the mk, and an error, if there is any.
func (c *mks) Update(ctx context.Context, mk *beta1.Mk, opts v1.UpdateOptions) (result *beta1.Mk, err error) {
	result = &beta1.Mk{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mks").
		Name(mk.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(mk).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *mks) UpdateStatus(ctx context.Context, mk *beta1.Mk, opts v1.UpdateOptions) (result *beta1.Mk, err error) {
	result = &beta1.Mk{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mks").
		Name(mk.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(mk).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the mk and deletes it. Returns an error if one occurs.
func (c *mks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mks").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mks").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched mk.
func (c *mks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *beta1.Mk, err error) {
	result = &beta1.Mk{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mks").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
