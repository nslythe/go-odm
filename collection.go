package goodm

import (
	"context"
	"errors"
	"path/filepath"
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
	Delete(obj interface{}) error
}

func to_object(o interface{}) DataObject {
	return_value, ok := o.(DataObject)
	if !ok {
		return_value = Obj(o)
	}
	return return_value
}

func (coll CollectionStruct) Save(o interface{}) error {
	obj := to_object(o)

	err := coll.validate_object(obj)
	if err != nil {
		return err
	}

	id, err := obj.GetID()
	if err != nil {
		return err
	}

	if id == primitive.NilObjectID {
		result, err := coll.Collection.InsertOne(context.TODO(), obj.Interface())
		if err != nil {
			return err
		}
		obj.SetID(result.InsertedID.(primitive.ObjectID))
	} else {
		_, err := coll.Collection.UpdateOne(context.TODO(),
			primitive.M{"_id": id},
			primitive.M{"$set": obj.Interface()})
		if err != nil {
			return err
		}
	}

	return nil
}

func (coll CollectionStruct) Load(o interface{}) error {
	obj := to_object(o)

	err := coll.validate_object(obj)
	if err != nil {
		return err
	}

	id, err := obj.GetID()
	if id == primitive.NilObjectID {
		return errors.New("Id not set")
	}

	filter := primitive.M{"_id": id}
	result := coll.Collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return result.Err()
	}
	result.Decode(obj.obj_interface)

	return nil
}

func (coll CollectionStruct) Find(o interface{}, filter primitive.M) error {
	obj := to_object(o)

	err := coll.validate_object(obj)
	if err != nil {
		return err
	}

	if !obj.IsSlice() {
		result := coll.Collection.FindOne(context.TODO(), filter)
		if result.Err() != nil {
			return result.Err()
		}
		return result.Decode(obj.obj_interface)
	} else {
		cursor, err := coll.Collection.Find(context.TODO(), filter)
		if err != nil {
			return err
		}
		defer cursor.Close(context.TODO())

		obj.Clear()

		for cursor.Next(context.Background()) {
			new_obj := obj.CreateNew()

			err = cursor.Decode(new_obj.Interface())
			if err != nil {
				return err
			}

			obj.Append(new_obj)
		}
	}
	return nil
}

func (coll CollectionStruct) Delete(o interface{}) error {
	obj := to_object(o)

	err := coll.validate_object(obj)
	if err != nil {
		return err
	}

	if !obj.IsSlice() {
		filter := primitive.M{"_id": obj.Field("Id").Interface().(primitive.ObjectID)}
		result, _ := coll.Collection.DeleteOne(context.TODO(), filter)
		if result.DeletedCount == 1 {
			obj.Field("Id").Set(primitive.NilObjectID)
		}
	} else {
		ids := []primitive.ObjectID{}
		index := obj.Len()
		for i := 0; i < index; i++ {
			ids = append(ids, obj.Index(i).Field("Id").Interface().(primitive.ObjectID))
		}
		result, _ := coll.Collection.DeleteMany(context.TODO(), primitive.M{"_id": primitive.M{"$in": ids}})
		if result.DeletedCount == 0 {
			return errors.New("No delete")
		}
	}

	return nil
}

func Coll(o interface{}) Collection {
	obj := to_object(o)

	client, err := mongo.Connect(context.TODO(), config.root_options)
	if err != nil {
		return nil
	}

	collection := CollectionStruct{}
	collection.Collection = client.Database(config.connection_string.Database).Collection(GetCollectionName(obj))
	return collection
}

func GetCollectionName(obj DataObject) string {
	new_type_name := ""
	type_name := ""

	package_name := filepath.Base(obj.Package())
	if package_name != "." {
		type_name = package_name + "_" + obj.Name()
	} else {
		type_name = obj.Name()
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

func (coll CollectionStruct) validate_object(obj DataObject) error {
	if coll.Collection == nil {
		return errors.New("Colelction not initialised")
	}

	if !obj.FieldExists("BaseObject") {
		return errors.New("No BaseObject")
	}
	bson_tag := obj.FieldTag("BaseObject", "bson")
	if !strings.Contains(bson_tag, "inline") {
		return errors.New("BaseObject not inline")
	}
	return nil
}
