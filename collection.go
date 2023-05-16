package goodm

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionStruct struct {
	Collection *mongo.Collection
}

type Collection interface {
	Save(obj interface{}) error
	Load(obj interface{}) error
	Find(obj interface{}, filter primitive.M) error
	Drop()
}

func (coll CollectionStruct) Save(obj interface{}) error {
	err := coll.validate_object(obj)
	if err != nil {
		return err
	}

	id, err := GetID(obj)
	if err != nil {
		return err
	}

	if id == primitive.NilObjectID {
		SetID(obj, primitive.NewObjectID())
	}

	_, err = coll.Collection.InsertOne(context.TODO(), obj)
	if err != nil {
		return err
	}

	return nil
}

func (coll CollectionStruct) Load(obj interface{}) error {
	err := coll.validate_object(obj)
	if err != nil {
		return err
	}

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return errors.New("Obj mus be ptr")
	}

	id, err := GetID(obj)
	if id == primitive.NilObjectID {
		return errors.New("Id not set")
	}

	filter := primitive.M{"_id": id}
	result := coll.Collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return result.Err()
	}
	result.Decode(obj)

	return nil
}

func (coll CollectionStruct) Find(obj interface{}, filter primitive.M) error {
	err := coll.validate_object(obj)
	if err != nil {
		return err
	}

	if !Obj(obj).sub_obj.IsSlice() {
		return errors.New("Obj must be slice")
	}

	cursor, err := coll.Collection.Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.Background()) {
		new_obj := Obj(obj).CreateNew()

		err = cursor.Decode(new_obj.Interface())
		if err != nil {
			return err
		}

		Obj(obj).Append(new_obj)
		//		append_value := reflect.Append(reflect.ValueOf(obj).Elem(), new_obj.obj_value.Elem())
		//		reflect.ValueOf(obj).Elem().Set(append_value)

	}
	return nil
}

func Coll(obj interface{}) Collection {
	client, err := mongo.Connect(context.TODO(), config.root_options)
	if err != nil {
		return nil
	}

	collection := CollectionStruct{}
	collection.Collection = client.Database(config.connection_string.Database).Collection(GetCollectionName(obj))
	return collection
}

func GetCollectionName(obj interface{}) string {
	new_type_name := ""
	type_name := ""

	package_name := filepath.Base(Obj(obj).Package())
	if package_name != "." {
		type_name = package_name + "_" + Obj(obj).Name()
	} else {
		type_name = Obj(obj).Name()
	}

	for _, c := range type_name {
		skip := false
		if unicode.IsUpper(c) && len(new_type_name) > 0 {
			new_type_name += "_"
		}
		if c == '-' {
			new_type_name += "_"
			skip = true
		}

		if !skip {
			new_type_name += string(c)
		}
	}
	return strings.ToLower(new_type_name)
}

func (coll CollectionStruct) Drop() {
	coll.Collection.Drop(context.Background())
}

func (coll CollectionStruct) validate_object(obj interface{}) error {
	if coll.Collection == nil {
		return errors.New("Colelction not initialised")
	}

	if !Obj(obj).FieldExists("BaseObject") {
		return errors.New("No BaseObject")
	}
	bson_tag := Obj(obj).FieldTag("BaseObject", "bson")
	if !strings.Contains(bson_tag, "inline") {
		return errors.New("BaseObject not inline")
	}
	return nil
}
