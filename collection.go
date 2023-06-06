package goodm

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionStruct struct {
	Collection *mongo.Collection
}

type Collection interface {
	Save(obj interface{}) error
	Update(obj interface{}, filter interface{}) error
	UpdateAll(obj interface{}, filter interface{}) (int, error)
	Load(obj interface{}) error
	Find(obj interface{}, filter interface{}) error
	FindSpecificType(obj interface{}, filter interface{}) error
	Drop()
	Delete(obj interface{}) error
	Count(filter interface{}) (int64, error)
	MongoCollection() *mongo.Collection
	CreateIndex(name string, model mongo.IndexModel) error
}

func Coll(o interface{}) Collection {
	obj := to_object(o)

	collection := CollectionStruct{}
	collection.Collection = client.Database(config.connection_string.Database).Collection(GetCollectionName(obj))
	return collection
}

func GetCollectionName(obj DataObject) string {
	if obj.FieldExists("BaseObject") {
		goodm_collection := obj.FieldTag("BaseObject", "goodm-collection")
		if goodm_collection != "" {
			return goodm_collection
		}
	}

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

func to_object(o interface{}) DataObject {
	return_value, ok := o.(DataObject)
	if !ok {
		return_value = Obj(o)
	}
	return return_value
}

func (coll CollectionStruct) MongoCollection() *mongo.Collection {
	return coll.Collection
}

func (coll CollectionStruct) Count(filter interface{}) (int64, error) {
	return coll.Collection.CountDocuments(context.TODO(), filter)
}

func (coll CollectionStruct) CreateIndex(name string, model mongo.IndexModel) error {
	tempo_name := name
	index_model := mongo.IndexModel{
		Keys:    model.Keys,
		Options: &options.IndexOptions{Name: &tempo_name},
	}
	name, err := coll.Collection.Indexes().CreateOne(context.TODO(), index_model)
	return err
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

	obj.SetTypeName()
	obj.SetUpdateTime()

	if id == primitive.NilObjectID {
		obj.SetID(primitive.NewObjectID())
		obj.SetCreationTime()
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

func (coll CollectionStruct) Update(o interface{}, filter interface{}) error {
	obj := to_object(o)

	err := coll.validate_object(obj)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	obj.SetTypeName()
	obj.SetUpdateTime()

	result := coll.Collection.FindOneAndUpdate(context.TODO(), filter, primitive.M{"$set": obj.Interface()})
	if result.Err() != nil {
		return result.Err()
	}

	result.Decode(obj.obj_interface)

	return nil
}

func (coll CollectionStruct) UpdateAll(o interface{}, filter interface{}) (int, error) {
	obj := to_object(o)

	err := coll.validate_object(obj)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	obj.SetTypeName()
	obj.SetUpdateTime()

	result, err := coll.Collection.UpdateMany(context.TODO(), filter, primitive.M{"$set": obj.Interface()})
	if err != nil {
		return 0, err
	}

	return int(result.ModifiedCount), nil
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

func (coll CollectionStruct) FindSpecificType(o interface{}, filter interface{}) error {
	obj := to_object(o)

	new_filter := bson.M{"$and": bson.A{
		bson.M{"_type_name": obj.FullTypeName()},
		filter,
	}}
	return coll.Find(o, new_filter)
}

func (coll CollectionStruct) Find(o interface{}, filter interface{}) error {
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
