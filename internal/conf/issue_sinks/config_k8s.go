package issue_sinks

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

const K8sType = "k8s"

type K8sSink struct {
	Namespace  string `json:"namespace"`
	SecretName string `json:"secret_name"`
}

func (f *K8sSink) GetType() string {
	return K8sType
}

func K8sBackendFromMap(args SinkConfig) (*K8sSink, error) {
	expectedTypes := map[string]reflect.Kind{
		"namespace":   reflect.String,
		"secret_name": reflect.String,
	}
	for keyword, valueType := range expectedTypes {
		t, ok := args[keyword]
		if !ok {
			return nil, fmt.Errorf("no '%s' in args", keyword)
		}

		switch v := reflect.ValueOf(t); v.Kind() {
		case valueType:
		default:
			return nil, fmt.Errorf("expected type %v for keyword '%s'", valueType, keyword)
		}
	}

	return &K8sSink{
		Namespace:  args["namespace"].(string),
		SecretName: args["secret_name"].(string),
	}, nil
}

func (b *K8sSink) Validate() (errs []error) {
	if len(b.SecretName) == 0 {
		errs = append(errs, errors.New("empty kubernetes secret name supplied"))
	}
	if len(b.Namespace) == 0 {
		errs = append(errs, errors.New("empty kubernetes namespace supplied"))
	}

	return
}

func (b *K8sSink) PrintConfig() {
	log.Info().Msgf("k8s namespace=%s", b.Namespace)
	log.Info().Msgf("k8s secret_name=%s", b.SecretName)
}
