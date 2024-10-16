package cache

import (
	"github.com/major1201/kubetrack/kube"
	"github.com/pkg/errors"
)

type KubeClient struct {
	*kube.ClientImpl

	globalInformer *globalInformer
	clusterID      ClusterID
}

func NewKubeClientFromGlobalInformer(gi GlobalInformer, clusterID ClusterID) (client kube.Client, err error) {
	giInst, ok := gi.(*globalInformer)
	if !ok {
		err = errors.New("gi should be type of *globalInformer")
		return
	}

	c := gi.GetCluster(clusterID)
	if c == nil {
		err = errors.Errorf("cluster %s not found", clusterID)
		return
	}

	clientImpl, ok := c.KubeClient().(*kube.ClientImpl)
	if !ok {
		err = errors.Errorf("the kube client is not kube.ClientImpl")
		return
	}
	client = &KubeClient{
		ClientImpl:     clientImpl,
		globalInformer: giInst,
		clusterID:      clusterID,
	}
	return
}
