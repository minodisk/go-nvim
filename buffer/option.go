package buffer

import (
	"reflect"
	"strings"
)

type Option struct {
	BufHidden  string
	BufListed  bool
	BufType    string
	ReadOnly   bool
	SwapFile   bool
	Modifiable bool
	Modified   bool
}

func (o Option) MapValue() map[string]interface{} {
	m := make(map[string]interface{})
	t, v := reflect.TypeOf(o), reflect.ValueOf(o)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("vim")
		if tag == "" {
			tag = strings.ToLower(f.Name)
		}
		m[tag] = v.Field(i).Interface()
	}
	return m
}

func (o Option) MapPointer() map[string]uintptr {
	m := make(map[string]uintptr)
	t, v := reflect.TypeOf(o), reflect.ValueOf(o)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("vim")
		if tag == "" {
			tag = strings.ToLower(f.Name)
		}
		m[tag] = v.Field(i).Pointer()
	}
	return m
}
