package comm

import "reflect"

var (
	structTypes = make(map[string]reflect.Type)
)

// 注册结构体类型
func RegisterStruct(s interface{}) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	structTypes[t.Name()] = t
}

// 根据结构体名称创建实例
func CreateInstance(structName string) interface{} {
	if t, exists := structTypes[structName]; exists {
		instance := reflect.New(t).Elem()
		return instance.Interface()
	}
	return nil
}

// 根据泛型类型创建实例
func CreateGenericInstance[T any]() T {
	var zero T
	return zero
}
