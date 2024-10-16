package output

import (
	"fmt"
	"time"

	"github.com/major1201/kubetrack/config"
	"github.com/major1201/kubetrack/log"
	corev1 "k8s.io/api/core/v1"
)

type LogOutput struct {
	conf *config.OutputLog
}

func NewLogOutput(conf *config.OutputLog) *LogOutput {
	return &LogOutput{
		conf: conf,
	}
}

func (lo *LogOutput) Name() string {
	return "log"
}

func (lo *LogOutput) Write(out OutputStruct) error {
	var keysAndValues []any
	keysAndValues = append(keysAndValues, "eventTime", out.EventTime.Format(time.RFC3339))
	keysAndValues = append(keysAndValues, "eventType", string(out.EventType))
	keysAndValues = append(keysAndValues, "objectRef", lo.objectRefString(out.ObjectRef))
	keysAndValues = append(keysAndValues, "_source", string(out.Source))

	for name, obj := range out.Fields {
		keysAndValues = append(keysAndValues, name, fmt.Sprintf("%v", obj))
	}

	log.L.Info(out.Message, keysAndValues...)

	if out.Diff != "" {
		fmt.Println(out.Diff)
	}
	return nil
}

func (lo *LogOutput) objectRefString(objectRef corev1.ObjectReference) string {
	return fmt.Sprintf("%s/%s,%s/%s", objectRef.APIVersion, objectRef.Kind, objectRef.Namespace, objectRef.Name)
}
