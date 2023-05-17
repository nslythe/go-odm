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
	obj_type  reflect.Type
	obj_value reflect.Value
	sub_obj   *DataObject
	type_name string
}

func Obj(obj interface{}) DataObject {
	return *create_DataObject_from_reflect(reflect.TypeOf(obj), reflect.ValueOf(obj))
}

func create_DataObject_from_reflect(t reflect.Type, v reflect.Value) *DataObject {
	return_value := &DataObject{}
	return_value.sub_obj = nil
	return_value.obj_value = v
	return_value.obj_type = t
	return_value.type_name = t.Name()

	if return_value.obj_type.Kind() == reflect.Slice {
		return_value.sub_obj = create_DataObject_from_reflect(return_value.obj_type.Elem(), reflect.Value{})
	}
	if return_value.obj_type.Kind() == reflect.Ptr {
		return_value.sub_obj = create_DataObject_from_reflect(return_value.obj_type.Elem(), return_value.obj_value.Elem())
	}

	return return_value
}

func (obj DataObject) Name() string {
	if obj.sub_obj != nil {
		return obj.sub_obj.Name()
	}
	return obj.obj_type.Name()
}

func (obj DataObject) IsSlice() bool {
	if obj.obj_type.Kind() == reflect.Ptr {
		return obj.sub_obj.IsSlice()
	}
	return obj.obj_value.Kind() == reflect.Slice
}

func (obj DataObject) CreateNew() DataObject {
	if obj.sub_obj != nil {
		return obj.sub_obj.CreateNew()
	}
	v := reflect.New(obj.obj_type)
	return *create_DataObject_from_reflect(v.Type(), v)
}

func (obj DataObject) Package() string {
	if obj.sub_obj != nil {
		return obj.sub_obj.Package()
	}
	return obj.obj_type.PkgPath()
}

func (obj DataObject) FieldExists(name string) bool {
	if obj.sub_obj != nil {
		return obj.sub_obj.FieldExists(name)
	}
	_, ok := obj.obj_type.FieldByName(name)
	return ok
}

func (obj DataObject) FieldTag(name string, tagName string) string {
	if obj.sub_obj != nil {
		return obj.sub_obj.FieldTag(name, tagName)
	}
	f, _ := obj.obj_type.FieldByName(name)
	return f.Tag.Get(tagName)
}

func (obj DataObject) Field(name string) DataObject {
	if obj.sub_obj != nil {
		return obj.sub_obj.Field(name)
	}

	_, ok := obj.obj_type.FieldByName(name)
	if !ok {
		return DataObject{}
	}
	value_field := obj.obj_value.FieldByName(name)

	return *create_DataObject_from_reflect(value_field.Type(), value_field)
}

func (obj DataObject) Index(i int) DataObject {
	if obj.obj_type.Kind() != reflect.Slice {
		return obj.sub_obj.Index(i)
	}
	return *create_DataObject_from_reflect(obj.obj_value.Index(i).Type(), obj.obj_value.Index(i))
}

func (obj DataObject) Len() int {
	if obj.obj_type.Kind() != reflect.Slice {
		return obj.sub_obj.Len()
	}
	return obj.obj_value.Len()
}

func (obj DataObject) Clear() {
	if obj.sub_obj != nil && obj.sub_obj.obj_value.IsValid() {
		obj.sub_obj.Clear()
	} else {
		if !obj.IsSlice() {
			panic("No a slice")
		}

		obj.obj_value.Set(reflect.MakeSlice(obj.obj_type, 0, 0))
	}
}

func (obj DataObject) Append(val DataObject) {
	if obj.sub_obj != nil && obj.sub_obj.obj_value.IsValid() {
		obj.sub_obj.Append(val)
	} else {
		if !obj.IsSlice() {
			panic("No a slice")
		}

		if val.obj_type.Kind() == reflect.Ptr {
			obj.obj_value.Set(reflect.Append(obj.obj_value, val.obj_value.Elem()))
		} else {
			obj.obj_value.Set(reflect.Append(obj.obj_value, val.obj_value))
		}
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
	if obj.sub_obj != nil {
		obj.sub_obj.Set(val)
	} else {
		if !obj.obj_value.CanSet() {
			panic("Cant set")
		}
		obj.obj_value.Set(reflect.ValueOf(val))
	}
}

func GetID(obj interface{}) (primitive.ObjectID, error) {
	if !Obj(obj).FieldExists("BaseObject") {
		return primitive.ObjectID{}, errors.New("BaseObject not present")
	}
	return Obj(obj).Field("BaseObject").Field("Id").Interface().(primitive.ObjectID), nil
}

func SetID(obj interface{}, id primitive.ObjectID) error {
	if !Obj(obj).FieldExists("BaseObject") {
		return errors.New("BaseObject not present")
	}
	Obj(obj).Field("BaseObject").Field("Id").Set(id)
	return nil
}
