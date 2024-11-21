package main

import (
	"os"
	"time"

	"github.com/major1201/kubetrack/config"
	"github.com/major1201/kubetrack/handler"
	"github.com/major1201/kubetrack/kube"
	kubecache "github.com/major1201/kubetrack/kube/cache"
	"github.com/major1201/kubetrack/log"
	"github.com/major1201/kubetrack/output"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// Name inspects project name
	Name = "kubetrack"

	// Version inspects the project version, which would be injected by build the tool
	Version = "custom"
)

func init() {
	// start program
	log.L.Info("starting up", "name", Name, "version", Version)
}

var gi kubecache.GlobalInformer

var eventGVR = schema.GroupVersionResource{Version: "v1", Resource: "events"}

func runMain(c *cli.Context) error {
	configPath := c.String("config")
	ktconfig, err := config.LoadFromFile(configPath)
	if err != nil {
		return err
	}

	gi = kubecache.NewGlobalInformer(kube.GetScheme())

	client, err := kube.NewClientInCluster()
	if err != nil {
		log.L.Error(err, "create kube client failed")
		return err
	}

	// make outputs
	var out []output.Output
	for _, outConfig := range ktconfig.Output {
		switch {
		case outConfig.Log != nil:
			out = append(out, output.NewLogOutput(outConfig.Log))
		case outConfig.Mysql != nil:
			out = append(out, output.NewMysqlOutput(&ktconfig, outConfig.Mysql))
		case outConfig.Postgres != nil:
			out = append(out, output.NewPostgresOutput(&ktconfig, outConfig.Postgres))
		}
	}

	generalHandler := handler.NewGeneralHandler(ktconfig, out)
	eventHandler := handler.NewEventHandler(ktconfig, out)

	// make units
	var units []kubecache.ResourceUnitWithHandlers
	for _, rule := range ktconfig.Rules {
		gv, err := schema.ParseGroupVersion(rule.APIVersion)
		if err != nil {
			return errors.Wrapf(err, "parse groupversion failed: %s", rule.APIVersion)
		}
		gvk := gv.WithKind(rule.Kind)
		mp, err := client.KindToMapping(gvk)
		if err != nil {
			return err
		}

		if len(rule.Namespaces) == 0 {
			log.L.Info("loading resources", "resource", mp.Resource.String())
			units = append(units, kubecache.ResourceUnitWithHandlers{
				ResourceUnit: kubecache.ResourceUnit{
					Resource: mp.Resource,
				},
				ResourceEventHandlers: []kubecache.ClusterResourceEventHandler{generalHandler},
			})
		} else {
			for _, namespace := range rule.Namespaces {
				log.L.Info("loading resources", "namespace", namespace, "resource", mp.Resource.String())
				units = append(units, kubecache.ResourceUnitWithHandlers{
					ResourceUnit: kubecache.ResourceUnit{
						Namespace: namespace,
						Resource:  mp.Resource,
					},
					ResourceEventHandlers: []kubecache.ClusterResourceEventHandler{generalHandler},
				})
			}
		}
	}

	// watch events
	if len(ktconfig.Events.Namespaces) == 0 {
		// all namespaces
		log.L.Info("loading events")
		units = append(units, kubecache.ResourceUnitWithHandlers{
			ResourceUnit: kubecache.ResourceUnit{
				Resource: eventGVR,
			},
			ResourceEventHandlers: []kubecache.ClusterResourceEventHandler{eventHandler},
		})
	} else {
		// for each namespace
		for _, namespace := range ktconfig.Events.Namespaces {
			log.L.Info("loading events", "namespace", namespace)
			units = append(units, kubecache.ResourceUnitWithHandlers{
				ResourceUnit: kubecache.ResourceUnit{
					Namespace: namespace,
					Resource:  eventGVR,
				},
				ResourceEventHandlers: []kubecache.ClusterResourceEventHandler{eventHandler},
			})
		}
	}

	gi.AddCluster(kubecache.ClusterID("default"), client, 0, nil, units)

	go waitAllSynced(generalHandler)

	<-make(chan struct{})

	return nil
}

func waitAllSynced(handler *handler.GeneralHandler) {
	t := time.NewTicker(100 * time.Millisecond)
	startTime := time.Now()
	defer func() {
		t.Stop()
		log.L.Info("list watch all resources on all clusters has synced", "duration", time.Since(startTime).String())
	}()
Out:
	for {
		<-t.C
		if gi.AllSynced() {
			if handler != nil {
				handler.SetSyned(true)
			}
			break Out
		}
	}
}

func main() {
	if err := getCLIApp().Run(os.Args); err != nil {
		log.L.Error(err, "flag unexpected error")
		os.Exit(1)
	}
}
