package beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{
	Group:   "mongokube.wrd",
	Version: "beta1",
}

var (
	SchemeBuilder runtime.SchemeBuilder
	AddToScheme   = SchemeBuilder.AddToScheme
)

func init() {
	// This func is called only once as soon as the package (beta1) is called,
	SchemeBuilder.Register(addKnownTypes)
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	// Add the types P4 and P4List to scheme
	scheme.AddKnownTypes(SchemeGroupVersion, &Mk{}, &MkList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
