package kube

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client kubernetes client interface
type Client interface {
	GetContext() context.Context
	SetContext(ctx context.Context)

	GetRESTConfig() *rest.Config
	GetKubeClient() kubernetes.Interface
	GetDynamicClient() dynamic.Interface
	GetDiscoveryClient() discovery.DiscoveryInterface
	GetRESTMapper() meta.RESTMapper

	// health
	Ping() error

	// types
	ResourceToMapping(gvr schema.GroupVersionResource) (mapping *meta.RESTMapping, err error)
	KindToMapping(gvk schema.GroupVersionKind) (mapping *meta.RESTMapping, err error)
	ModelToMapping(model runtime.Object) (mapping *meta.RESTMapping, err error)

	// resources
	IsReady(obj runtime.Object) (ready bool, err error)
	IsReadyWithUpdate(obj runtime.Object) (ready bool, retObj runtime.Object, err error)
	ListAllResources(namespace string) (objs []unstructured.Unstructured, err error)
}
