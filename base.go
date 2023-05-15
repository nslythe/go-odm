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
	sub_type  reflect.Type
	obj_value reflect.Value
}

func Obj(obj interface{}) DataObject {
	return_value := DataObject{}
	return_value.obj_value = reflect.ValueOf(obj)
	return_value.obj_type = reflect.TypeOf(obj)
	if return_value.obj_type.Kind() == reflect.Ptr {
		return_value.obj_value = return_value.obj_value.Elem()
		return_value.obj_type = return_value.obj_value.Type()
	}
	if return_value.obj_type.Kind() == reflect.Slice {
		return_value.sub_type = return_value.obj_type.Elem()
	}
	return return_value
}

func (obj DataObject) Name() string {
	if obj.sub_type != nil {
		return obj.sub_type.Name()
	}
	return obj.obj_type.Name()
}

func (obj DataObject) Package() string {
	return obj.obj_type.PkgPath()
}

func (obj DataObject) FieldExists(name string) bool {
	_, ok := obj.obj_type.FieldByName(name)
	return ok
}

func (obj DataObject) FieldTag(name string, tagName string) string {
	f, _ := obj.obj_type.FieldByName(name)
	return f.Tag.Get(tagName)
}

func (obj DataObject) Field(name string) DataObject {
	return_value := DataObject{}
	_, ok := obj.obj_type.FieldByName(name)
	if !ok {
		return return_value
	}
	value_field := obj.obj_value.FieldByName(name)

	return_value.obj_value = value_field
	return_value.obj_type = value_field.Type()
	return return_value
}

func (obj DataObject) Index(i int) DataObject {
	return_value := DataObject{}
	if obj.obj_type.Kind() != reflect.Slice {
		return return_value
	}
	return_value.obj_value = obj.obj_value.Index(i)
	return_value.obj_type = obj.obj_value.Index(i).Type()
	return return_value
}

func (obj DataObject) Append(val interface{}) {
	if obj.obj_type.Kind() != reflect.Slice {
		panic("No a slice")
	}
	obj.obj_value.Set(reflect.Append(obj.obj_value, reflect.ValueOf(val)))
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
	if !obj.obj_value.CanSet() {
		panic("Cant set")
	}
	obj.obj_value.Set(reflect.ValueOf(val))
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
