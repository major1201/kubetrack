package handler

import (
	"encoding/json"
	"fmt"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/google/go-cmp/cmp"
	"github.com/major1201/kubetrack/config"
	"github.com/major1201/kubetrack/kube"
	kubecache "github.com/major1201/kubetrack/kube/cache"
	"github.com/major1201/kubetrack/log"
	"github.com/major1201/kubetrack/output"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type EventHandler struct {
	GeneralHandler
}

func NewEventHandler(conf config.KubeTrackConfiguration, outputers []output.Output) *EventHandler {
	return &EventHandler{
		GeneralHandler: GeneralHandler{
			config:    conf,
			outputers: outputers,
		},
	}
}

func (h *EventHandler) OnAdd(cluster kubecache.Cluster, obj any) {
	unstrObj := obj.(*unstructured.Unstructured)
	if !h.config.Events.Match(unstrObj) {
		return
	}
	eventTime := time.Now()
	eventIf, _ := kube.GetScheme().ConvertToVersion(obj.(runtime.Object), runtime.GroupVersioner(schema.GroupVersions(kube.GetScheme().PrioritizedVersionsAllGroups())))
	event := eventIf.(*corev1.Event)

	// filter initial list
	if h.isHistoryAdd(event) {
		return
	}

	h.pruneObject(unstrObj)

	// main tree
	content := output.OutputStruct{
		EventTime: eventTime,
		ObjectRef: event.InvolvedObject,
		EventType: output.EventTypeAdd,
		Source:    output.SourceTypeEvent,
		Object:    unstrObj.Object,
		Message:   h.displayOfMessage(event),
	}

	// write output
	for _, outputer := range h.outputers {
		if err := outputer.Write(content); err != nil {
			log.L.Error(err, "writing output failed", "name", outputer.Name())
		}
	}
}

func (h *EventHandler) OnUpdate(cluster kubecache.Cluster, oldObj, newObj any) {
	oldUnstrObj := oldObj.(*unstructured.Unstructured)
	newUnstrObj := newObj.(*unstructured.Unstructured)
	if !h.config.Events.Match(newUnstrObj) {
		return
	}
	eventTime := time.Now()
	oldEventIf, _ := kube.GetScheme().ConvertToVersion(oldObj.(runtime.Object), runtime.GroupVersioner(schema.GroupVersions(kube.GetScheme().PrioritizedVersionsAllGroups())))
	oldEvent := oldEventIf.(*corev1.Event)
	newEventIf, _ := kube.GetScheme().ConvertToVersion(newObj.(runtime.Object), runtime.GroupVersioner(schema.GroupVersions(kube.GetScheme().PrioritizedVersionsAllGroups())))
	newEvent := newEventIf.(*corev1.Event)

	h.pruneObject(oldUnstrObj)
	h.pruneObject(newUnstrObj)

	// main tree
	content := output.OutputStruct{
		EventTime: eventTime,
		ObjectRef: oldEvent.InvolvedObject,
		EventType: output.EventTypeUpdate,
		Source:    output.SourceTypeEvent,
		Object:    newUnstrObj.Object,
		Message:   h.displayOfMessage(newEvent),
		Diff:      cmp.Diff(oldUnstrObj.Object, newUnstrObj.Object),
	}

	// save json patch
	oldJSON, err := json.Marshal(oldObj)
	if err != nil {
		log.L.Error(err, "marshal oldObj json failed")
		return
	}
	newJSON, err := json.Marshal(newObj)
	if err != nil {
		log.L.Error(err, "marshal newObj json failed")
		return
	}
	jp, err := jsonpatch.CreateMergePatch(oldJSON, newJSON)
	if err != nil {
		log.L.Error(err, "create merge patch onUpdate failed")
		return
	}
	content.JsonPatch = string(jp)

	// write output
	for _, outputer := range h.outputers {
		if err := outputer.Write(content); err != nil {
			log.L.Error(err, "writing output failed", "name", outputer.Name())
		}
	}
}

func (h *EventHandler) OnDelete(cluster kubecache.Cluster, obj any) {
	// we don't care about the deletion of events
}

func (h *EventHandler) displayOfMessage(ev *corev1.Event) string {
	if ev == nil {
		return "_empty"
	}

	var cnt int32
	if ev.Series != nil {
		cnt = ev.Series.Count
	} else if ev.Count > 1 {
		cnt = ev.Count
	}
	cntStr := Ternary(cnt > 0, fmt.Sprintf(" x%d", cnt), "")

	return fmt.Sprintf("%s %s%s %s", ev.Type, ev.Reason, cntStr, ev.Message)
}
