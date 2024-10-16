package kube

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientImpl describes the a kubernetes client
type ClientImpl struct {
	ctx               context.Context
	restConfig        *rest.Config
	kubeClient        kubernetes.Interface
	dynamicClient     dynamic.Interface
	discoveryClient   discovery.DiscoveryInterface
	restMapper        meta.RESTMapper
	scaleKindResolver scale.ScaleKindResolver
}

// GetContext get current context
func (c *ClientImpl) GetContext() context.Context {
	return c.ctx
}

// SetContext set query context
func (c *ClientImpl) SetContext(ctx context.Context) {
	c.ctx = ctx
}

// GetRESTConfig returns the REST config
func (c *ClientImpl) GetRESTConfig() *rest.Config {
	return c.restConfig
}

// GetKubeClient returns the kubernetes client
func (c *ClientImpl) GetKubeClient() kubernetes.Interface {
	return c.kubeClient
}

// GetDynamicClient returns the dynamic client
func (c *ClientImpl) GetDynamicClient() dynamic.Interface {
	return c.dynamicClient
}

// GetDiscoveryClient returns the discovery client
func (c *ClientImpl) GetDiscoveryClient() discovery.DiscoveryInterface {
	return c.discoveryClient
}

// GetRESTMapper returns the rest mapper
func (c *ClientImpl) GetRESTMapper() meta.RESTMapper {
	return c.restMapper
}

// NewClientForConfig returns the kubernetes client with rest config
func NewClientForConfig(restConfig *rest.Config) (client Client, err error) {
	clientRet := &ClientImpl{ctx: context.TODO()}

	clientRet.restConfig = restConfig

	// set kubernetes client
	if clientRet.kubeClient, err = kubernetes.NewForConfig(restConfig); err != nil {
		err = errors.WithStack(err)
		return
	}

	// set dynamic client
	if clientRet.dynamicClient, err = dynamic.NewForConfig(restConfig); err != nil {
		err = errors.WithStack(err)
		return
	}

	// set discovery client
	if clientRet.discoveryClient, err = discovery.NewDiscoveryClientForConfig(restConfig); err != nil {
		err = errors.WithStack(err)
		return
	}

	// set rest mapper
	restMapperRes, err := restmapper.GetAPIGroupResources(clientRet.discoveryClient)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	clientRet.restMapper = restmapper.NewDiscoveryRESTMapper(restMapperRes)

	client = clientRet
	return
}

// NewClient would try in-cluster then out-cluster, returns the kubernetes client
func NewClient() (client Client, err error) {
	var restConfig *rest.Config
	if restConfig, err = rest.InClusterConfig(); err != nil {
		if err != rest.ErrNotInCluster {
			err = errors.WithStack(err)
			return
		}
		// try out cluster
		if restConfig, err = getOutClusterConfig("", ""); err != nil {
			return
		}
	}
	restConfig.QPS = 1000
	restConfig.Burst = 1000

	return NewClientForConfig(restConfig)
}

// NewClientInCluster returns the kubernetes client inside the cluster
func NewClientInCluster() (client Client, err error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return NewClientForConfig(restConfig)
}

// NewClientOutCluster returns the kubernetes client from the kubeconfig and context out of the cluster
func NewClientOutCluster(context, kubeconfig string) (client Client, err error) {
	restConfig, err := getOutClusterConfig(context, kubeconfig)
	if err != nil {
		return
	}
	return NewClientForConfig(restConfig)
}

func getOutClusterConfig(context, kubeconfig string) (*rest.Config, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := &clientcmd.ConfigOverrides{
		ClusterDefaults: clientcmd.ClusterDefaults,
	}
	if context != "" {
		overrides.CurrentContext = context
	}
	if kubeconfig != "" {
		rules.ExplicitPath = kubeconfig
	}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		err = errors.WithStack(err)
	}
	return config, err
}
