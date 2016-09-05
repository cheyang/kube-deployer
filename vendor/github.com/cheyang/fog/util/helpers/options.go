package helpers

import "github.com/docker/machine/libmachine/drivers"

type ConfigFlagger struct {
	Data map[string]interface{}
}

func NewConfigFlagger(data map[string]interface{}) drivers.DriverOptions {
	return ConfigFlagger{Data: data}
}

func (this ConfigFlagger) String(key string) string {
	if value, ok := this.Data[key]; ok {
		return value.(string)
	}
	return ""
}

func (this ConfigFlagger) StringSlice(key string) []string {
	if value, ok := this.Data[key]; ok {
		// return value.([]string)
		switch value.(type) {
		case []string:
			return value.([]string)
		case []interface{}:
			// interface slice
			is := value.([]interface{})
			// string slice
			ss := []string{}
			for _, v := range is {
				ss = append(ss, v.(string))
			}
		}
	}
	return []string{}
}

func (this ConfigFlagger) Int(key string) int {
	if value, ok := this.Data[key]; ok {
		// return value.(int)
		switch value.(type) {
		case int:
			return value.(int)
		case float64:
			return int(value.(float64))
		}
	}
	return 0
}

func (this ConfigFlagger) Bool(key string) bool {
	if value, ok := this.Data[key]; ok {
		return value.(bool)
	}
	return false
}
