package goodm

import (
	"errors"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseObject struct {
	Id primitive.ObjectID `bson:"_id"`
}

type DataObject struct {
	obj_type      reflect.Type
	obj_value     reflect.Value
	obj_interface interface{}
	type_name     string
}

func Obj(obj interface{}) DataObject {
	data_obj := *create_DataObject_from_reflect(reflect.TypeOf(obj), reflect.ValueOf(obj))
	data_obj.obj_interface = obj
	return data_obj
}

func create_DataObject_from_reflect(t reflect.Type, v reflect.Value) *DataObject {
	return_value := &DataObject{}
	return_value.obj_value = v
	return_value.obj_type = t
	return_value.type_name = t.Name()
	return return_value
}

func (obj DataObject) Name() string {
	t := obj.obj_type

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return t.Name()
}

func (obj DataObject) IsSlice() bool {
	if obj.obj_type.Kind() == reflect.Ptr {
		return obj.obj_type.Elem().Kind() == reflect.Slice
	}
	return obj.obj_type.Kind() == reflect.Slice
}

func (obj DataObject) CreateNew() DataObject {
	var t reflect.Type = obj.obj_type

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	v := reflect.New(t)
	return *create_DataObject_from_reflect(v.Type(), v)
}

func (obj DataObject) Package() string {
	t := obj.obj_type

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return t.PkgPath()
}

func (obj DataObject) FieldExists(name string) bool {
	t := obj.obj_type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	_, ok := t.FieldByName(name)
	return ok
}

func (obj DataObject) FieldTag(name string, tagName string) string {
	if !obj.FieldExists(name) {
		panic("Field does not exists")
	}

	t := obj.obj_type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	f, _ := t.FieldByName(name)
	return f.Tag.Get(tagName)
}

func (obj DataObject) Field(name string) DataObject {
	if obj.IsSlice() {
		panic("Cant be slice")
	}

	var t reflect.Type
	if obj.obj_type.Kind() == reflect.Ptr {
		t = obj.obj_type.Elem()
	} else {
		t = obj.obj_type
	}
	_, ok := t.FieldByName(name)
	if !ok {
		return DataObject{}
	}

	var value_field reflect.Value
	if obj.obj_value.Kind() == reflect.Ptr {
		value_field = obj.obj_value.Elem().FieldByName(name)
	} else {
		value_field = obj.obj_value.FieldByName(name)
	}

	return *create_DataObject_from_reflect(value_field.Type(), value_field)
}

func (obj DataObject) Index(i int) DataObject {
	if !obj.IsSlice() {
		panic("Not a slice")
	}

	if obj.obj_value.Kind() == reflect.Ptr {
		return *create_DataObject_from_reflect(obj.obj_value.Elem().Index(i).Type(), obj.obj_value.Elem().Index(i))
	} else {
		return *create_DataObject_from_reflect(obj.obj_value.Index(i).Type(), obj.obj_value.Index(i))
	}

}

func (obj DataObject) Len() int {
	if !obj.IsSlice() {
		panic("Not a slice")
	}

	if obj.obj_value.Kind() == reflect.Ptr {
		return obj.obj_value.Elem().Len()
	} else {
		return obj.obj_value.Len()
	}
}

func (obj DataObject) Clear() {
	if !obj.IsSlice() {
		panic("No a slice")
	}

	if obj.obj_value.Kind() == reflect.Ptr {
		obj.obj_value.Elem().Set(reflect.MakeSlice(obj.obj_type.Elem(), 0, 0))
	} else {
		obj.obj_value.Set(reflect.MakeSlice(obj.obj_type, 0, 0))
	}
}

func (obj DataObject) Append(val DataObject) {
	if !obj.IsSlice() {
		panic("No a slice")
	}

	var reflect_value reflect.Value
	if val.obj_type.Kind() == reflect.Ptr {
		reflect_value = val.obj_value.Elem()
	} else {
		reflect_value = val.obj_value
	}

	if obj.obj_value.Kind() == reflect.Ptr {
		obj.obj_value.Elem().Set(reflect.Append(obj.obj_value.Elem(), reflect_value))
	} else {
		obj.obj_value.Set(reflect.Append(obj.obj_value, reflect_value))
	}
}

func (obj DataObject) String() string {
	kind := obj.obj_type.Kind()

	if kind == reflect.String {
		return obj.obj_value.String()
	}
	if kind == reflect.Int ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Int8 {
		return strconv.FormatInt(obj.obj_value.Int(), 10)
	}
	if kind == reflect.Bool {
		return strconv.FormatBool(obj.obj_value.Bool())
	}
	if kind == reflect.Float32 ||
		kind == reflect.Float64 {
		return strconv.FormatFloat(obj.obj_value.Float(), 'f', 5, 64)
	}
	return ""
}

func (obj DataObject) Interface() interface{} {
	return obj.obj_value.Interface()
}

func (obj DataObject) Set(val interface{}) {
	var v reflect.Value
	if obj.obj_value.Kind() == reflect.Ptr {
		v = obj.obj_value.Elem()
	} else {
		v = obj.obj_value
	}

	if !v.CanSet() {
		panic("Cant set")
	}
	v.Set(reflect.ValueOf(val))

}

func (obj DataObject) GetID() (primitive.ObjectID, error) {
	if !obj.FieldExists("BaseObject") {
		return primitive.ObjectID{}, errors.New("BaseObject not present")
	}
	return obj.Field("BaseObject").Field("Id").Interface().(primitive.ObjectID), nil
}

func (obj DataObject) SetID(id primitive.ObjectID) error {
	if !obj.FieldExists("BaseObject") {
		return errors.New("BaseObject not present")
	}
	obj.Field("BaseObject").Field("Id").Set(id)
	return nil
}
