package handler

import (
	"encoding/json"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/google/go-cmp/cmp"
	"github.com/major1201/kubetrack/config"
	kubecache "github.com/major1201/kubetrack/kube/cache"
	"github.com/major1201/kubetrack/log"
	"github.com/major1201/kubetrack/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

const timeDiffDuringSyncingOnAdd = 10 * time.Second

type GeneralHandler struct {
	config    config.KubeTrackConfiguration
	outputers []output.Output
	synced    bool
}

func NewGeneralHandler(conf config.KubeTrackConfiguration, outputers []output.Output) *GeneralHandler {
	return &GeneralHandler{
		config:    conf,
		outputers: outputers,
	}
}

func (h *GeneralHandler) OnAdd(cluster kubecache.Cluster, obj any) {
	eventTime := time.Now()
	unstrObj := obj.(*unstructured.Unstructured)

	// filter initial list
	if h.isHistoryAdd(unstrObj) {
		return
	}

	rule := h.getRule(unstrObj)
	if rule == nil {
		return
	}
	eventAction := rule.OnCreate

	// main tree
	objRef := BuildObjectReference(unstrObj)
	h.pruneObject(unstrObj)
	content := output.OutputStruct{
		EventTime: eventTime,
		ObjectRef: objRef,
		EventType: output.EventTypeAdd,
		Source:    output.SourceTypeGeneral,
		Fields:    BuildFieldsMap(unstrObj, rule.CareFields),
	}

	if eventAction.SaveFullObject {
		content.Object = unstrObj.Object
	}

	// write output
	for _, outputer := range h.outputers {
		if err := outputer.Write(content); err != nil {
			log.L.Error(err, "writing output failed", "name", outputer.Name())
		}
	}
}

func (h *GeneralHandler) OnUpdate(cluster kubecache.Cluster, oldObj, newObj any) {
	eventTime := time.Now()
	oldUnstrObj := oldObj.(*unstructured.Unstructured)
	newUnstrObj := newObj.(*unstructured.Unstructured)

	rule := h.getRule(oldUnstrObj)
	if rule == nil {
		return
	}
	eventAction := rule.OnUpdate

	// main tree
	objRef := BuildObjectReference(oldUnstrObj)
	h.pruneObject(oldUnstrObj)
	h.pruneObject(newUnstrObj)
	content := output.OutputStruct{
		EventTime: eventTime,
		ObjectRef: objRef,
		EventType: output.EventTypeUpdate,
		Source:    output.SourceTypeGeneral,
		Fields:    BuildFieldsMap(newUnstrObj, rule.CareFields),
	}
	if eventAction.SaveFullObject {
		content.Object = newUnstrObj.Object
	}

	if eventAction.SaveCmp {
		content.Diff = cmp.Diff(oldUnstrObj.Object, newUnstrObj.Object)
	}

	if eventAction.SaveJsonPatch {
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
	}

	// write output
	for _, outputer := range h.outputers {
		if err := outputer.Write(content); err != nil {
			log.L.Error(err, "writing output failed", "name", outputer.Name())
		}
	}
}

func (h *GeneralHandler) OnDelete(cluster kubecache.Cluster, obj any) {
	switch objWithType := obj.(type) {
	case *unstructured.Unstructured:
		h.onDeleteUnstr(cluster, objWithType, false)
	case cache.DeletedFinalStateUnknown:
		h.onDeleteFinalStateUnknown(cluster, objWithType)
	default:
		log.L.Error(nil, "unknown type onDelete")
	}
}

func (h *GeneralHandler) onDeleteUnstr(_ kubecache.Cluster, unstrObj *unstructured.Unstructured, isTombstone bool) {
	eventTime := time.Now()

	rule := h.getRule(unstrObj)
	if rule == nil {
		return
	}
	eventAction := rule.OnDelete

	// main tree
	objRef := BuildObjectReference(unstrObj)
	h.pruneObject(unstrObj)
	content := output.OutputStruct{
		EventTime: eventTime,
		ObjectRef: objRef,
		EventType: output.EventTypeDelete,
		Source:    output.SourceTypeGeneral,
		Fields:    BuildFieldsMap(unstrObj, rule.CareFields),
		Message:   Ternary(isTombstone, " [tombstone]", ""),
	}

	if eventAction.SaveFullObject {
		content.Object = unstrObj.Object
	}

	// write output
	for _, outputer := range h.outputers {
		if err := outputer.Write(content); err != nil {
			log.L.Error(err, "writing output failed", "name", outputer.Name())
		}
	}
}

func (h *GeneralHandler) onDeleteFinalStateUnknown(cluster kubecache.Cluster, tombstone cache.DeletedFinalStateUnknown) {
	unstrObj, ok := tombstone.Obj.(*unstructured.Unstructured)
	if !ok {
		log.L.Error(nil, "tombstone contained object that is not a Unstructured", "obj", tombstone)
		return
	}
	h.onDeleteUnstr(cluster, unstrObj, true)
}

func (h *GeneralHandler) SetSyned(synced bool) {
	h.synced = synced
}

func (h *GeneralHandler) getRule(obj runtime.Object) *config.Rule {
	// get the first rule matches
	for _, rule := range h.config.Rules {
		if rule.Match(obj) {
			rule := rule
			return &rule
		}
	}
	return nil // not found
}

// func (h *GeneralHandler) pruneObject(obj metav1.Object) {
func (h *GeneralHandler) pruneObject(obj *unstructured.Unstructured) {
	if obj == nil {
		return
	}

	obj.SetResourceVersion("")
	obj.SetManagedFields(nil)
}

func (h *GeneralHandler) isHistoryAdd(obj metav1.Object) bool {
	if h.synced {
		return false
	}
	if time.Since(obj.GetCreationTimestamp().Time) > timeDiffDuringSyncingOnAdd {
		return true
	}
	return false
}
