package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
)

func DecodeInput(input, param any) error {
	return mapstructure.Decode(input, param)
}

func EncodeOutput(output any) (any, error) {
	d, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	var v any
	return v, json.Unmarshal(d, &v)
}

type Box struct {
	value any
}

func NewBox(v any) *Box {
	return &Box{value: v}
}

func (vb *Box) MapString() (map[string]string, bool) {
	if vb == nil || vb.value == nil {
		return nil, false
	}
	switch m := vb.value.(type) {
	case map[string]string:
		return m, true
	case map[string]any:
		newMap := make(map[string]string, len(m))
		for k, v := range m {
			if s, ok := v.(string); !ok {
				return nil, false
			} else {
				newMap[k] = s
			}
		}
		return newMap, true
	default:
		return nil, false
	}
}

func (vb *Box) String() (string, bool) {
	if vb == nil || vb.value == nil {
		return "", false
	}
	s, ok := vb.value.(string)
	return s, ok
}

func (vb *Box) Int() (int, bool) {
	if vb == nil || vb.value == nil {
		return 0, false
	}
	switch i := vb.value.(type) {
	case int:
		return int(i), true
	case int8:
		return int(i), true
	case int16:
		return int(i), true
	case int32:
		return int(i), true
	case int64:
		return int(i), true
	case uint:
		return int(i), true
	case uint8:
		return int(i), true
	case uint16:
		return int(i), true
	case uint32:
		return int(i), true
	case uint64:
		return int(i), true
	}
	return 0, false
}

func (vb *Box) Value() any {
	if vb == nil {
		return nil
	}
	return vb.value
}

func (vb *Box) Get(name string) *Box {
	if vb == nil || vb.value == nil {
		return vb
	}
	switch m := vb.value.(type) {
	case map[string]any:
		return &Box{value: m[name]}
	case map[string]string:
		return &Box{value: m[name]}
	case map[string]int:
		return &Box{value: m[name]}
	case map[string]int8:
		return &Box{value: m[name]}
	case map[string]int16:
		return &Box{value: m[name]}
	case map[string]int32:
		return &Box{value: m[name]}
	case map[string]int64:
		return &Box{value: m[name]}
	case map[string]uint:
		return &Box{value: m[name]}
	case map[string]uint8:
		return &Box{value: m[name]}
	case map[string]uint16:
		return &Box{value: m[name]}
	case map[string]uint32:
		return &Box{value: m[name]}
	case map[string]uint64:
		return &Box{value: m[name]}
	case map[string]float32:
		return &Box{value: m[name]}
	case map[string]float64:
		return &Box{value: m[name]}
	}

	v := reflect.ValueOf(vb.value)

	// 解引用多层指针
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return &Box{}
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		// 只处理 map[string]*
		if v.Type().Key().Kind() != reflect.String {
			return &Box{}
		}
		key := reflect.ValueOf(name)
		val := v.MapIndex(key)
		if !val.IsValid() {
			return &Box{}
		}
		return &Box{value: val.Interface()}

	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Name == name {
				if !v.Field(i).CanInterface() {
					return &Box{}
				}
				return &Box{value: v.Field(i).Interface()}
			}
		}
		return &Box{}
	default:
		return &Box{}
	}

}

func NewError(pluginName string) *Error {
	return &Error{pluginName: pluginName}
}

type Error struct {
	pluginName string
}

func (e *Error) Throw(v ...any) error {
	var errs []error
	var tags []string
	var others []any
	for _, t := range v {
		switch target := t.(type) {
		case error:
			if target != nil {
				errs = append(errs, target)
			}
		case string:
			if target != "" {
				tags = append(tags, target)
			}
		default:
			if target != nil {
				others = append(others, target)
			}
		}
	}
	errorStr := "plugin"
	if e.pluginName != "" {
		errorStr += " " + e.pluginName
	}
	errorStr += " error:"
	if len(tags) > 0 {
		errorStr += " " + strings.Join(tags, " ")
	}
	if len(others) > 0 {
		errorStr += " " + fmt.Sprintf("%v", others...)
	}
	return errors.New(errorStr)
}
