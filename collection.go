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
}

func (coll CollectionStruct) Save(obj interface{}) error {
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

	id, err := GetID(obj)
	if err != nil {
		return err
	}

	ctx, client, cancel, err := CreateConnection()
	defer client.Disconnect(ctx)
	defer cancel()
	if err != nil {
		return err
	}

	if id == primitive.NilObjectID {
		SetID(obj, primitive.NewObjectID())
	}

	_, err = coll.Collection.InsertOne(ctx, obj)
	if err != nil {
		return err
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

func Find(obj interface{}, filter primitive.M) error {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return errors.New("Obj mus be ptr")
	}

	if reflect.ValueOf(obj).Elem().Kind() != reflect.Slice {
		return errors.New("Obj must be slice")
	}

	obj_type := reflect.TypeOf(obj).Elem().Elem()
	collection_name := GetCollectionName(obj_type.Name())

	ctx, client, cancel, err := CreateConnection()
	defer client.Disconnect(ctx)
	defer cancel()
	if err != nil {
		return err
	}

	cursor, err := client.Database(config.connection_string.Database).Collection(collection_name).Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		new_obj := reflect.New(obj_type)

		err = cursor.Decode(new_obj.Interface())
		if err != nil {
			return err
		}

		reflect.ValueOf(obj).Elem().Set(
			reflect.Append(reflect.ValueOf(obj).Elem(), new_obj.Elem()))

	}

	return nil
}

func Load(obj interface{}) error {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return errors.New("Obj mus be ptr")
	}

	collection_name := GetCollectionName(GetCollectionName(obj))
	id, err := GetID(obj)
	if id == primitive.NilObjectID {
		return errors.New("Id not set")
	}

	ctx, client, cancel, err := CreateConnection()
	defer client.Disconnect(ctx)
	defer cancel()
	if err != nil {
		return err
	}

	filter := primitive.M{"_id": id}
	result := client.Database(config.connection_string.Database).Collection(collection_name).FindOne(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}
	result.Decode(obj)

	return nil
}
