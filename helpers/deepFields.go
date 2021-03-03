package helpers

import "reflect"

func GetDeepFields(iface interface{}) (map[string]interface{}) {
	//fields := make([]reflect.Value, 0)
	fields := make(map[string]interface{})
	irv := reflect.ValueOf(iface)
	irt := reflect.TypeOf(iface)

	for i := 0; i < irt.NumField(); i++ {
		t := irt.Field(i)
		v := irv.Field(i)

		switch v.Kind() {
		case reflect.Struct:
			fields[t.Name] = GetDeepFields(v.Interface())
		default:
			fields[t.Name] = v
		}
	}

	return fields
}