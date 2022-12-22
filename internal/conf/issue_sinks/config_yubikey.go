package issue_sinks

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

const YubiType = "yubi"

type YubikeySink struct {
	YubikeyPin  string `json:"pin"`
	YubikeySlot uint32 `json:"slot"`
}

func YubiSinkFromMap(args SinkConfig) (*YubikeySink, error) {
	expectedTypes := map[string]reflect.Kind{
		"pin":  reflect.String,
		"slot": reflect.Uint32,
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

	return &YubikeySink{
		YubikeyPin:  args["pin"].(string),
		YubikeySlot: args["slot"].(uint32),
	}, nil
}

func (f *YubikeySink) GetType() string {
	return YubiType
}

func (b *YubikeySink) Validate() (errs []error) {
	// TODO
	return
}

func (b *YubikeySink) PrintConfig() {
	log.Info().Msgf("yubi slot=%s", b.YubikeySlot)
	log.Info().Msg("yubi pin=*** (redacted)")
}
