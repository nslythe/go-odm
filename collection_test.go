package goodm

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_GetColletionName_1(t *testing.T) {
	type Test struct{}
	if GetCollectionName(Obj(&Test{})) != "goodm__test" {
		t.Error("failed")
	}
}

func Test_GetColletionName_2(t *testing.T) {
	type Test_allo struct{}
	if GetCollectionName(Obj(&Test_allo{})) != "goodm__test_allo" {
		t.Error("failed")
	}
}

func Test_GetColletionName_3(t *testing.T) {
	type tests struct{}
	if GetCollectionName(Obj(&tests{})) != "goodm_tests" {
		t.Error("failed")
	}
}

func Test_GetColletionName_4(t *testing.T) {
	type TestTest struct{}
	if GetCollectionName(Obj(&TestTest{})) != "goodm__test_test" {
		t.Error("failed")
	}
}

func Test_GetColletionName_5(t *testing.T) {
	type TestTest1 struct{}
	if GetCollectionName(Obj(&TestTest1{})) != "goodm__test_test1" {
		t.Error("failed")
	}
}

func Test_GetColletionName_6(t *testing.T) {
	type TestTest1 struct{}
	if GetCollectionName(Obj(&TestTest1{})) != "goodm__test_test1" {
		t.Error("failed")
	}
	if GetCollectionName(Obj(&[]TestTest1{})) != "goodm__test_test1" {
		t.Error("failed")
	}
}

func Test_save_1(t *testing.T) {
	type TestTest1 struct{}
	obj := Obj(&TestTest1{})

	err := Coll(obj).Save(obj)
	if err == nil {
		t.Error("No BaseObject in parent")
	}
}

func Test_save_2(t *testing.T) {
	type TestTest1 struct {
		BaseObject
	}
	obj := Obj(&TestTest1{})
	err := Coll(obj).Save(obj)
	if err == nil {
		t.Errorf("BaseObject not inline %s", err)
	}
}

func Test_save_21(t *testing.T) {
	type TestTest1 struct {
		BaseObject
	}
	obj := &TestTest1{}
	err := Coll(obj).Save(obj)
	if err == nil {
		t.Errorf("BaseObject not inline %s", err)
	}
}

func Test_save_22(t *testing.T) {
	type TestTest1 struct {
		BaseObject `goodm-collection:"test_tesfdfsdfsdt11"`
	}
	if GetCollectionName(Obj(TestTest1{})) != "test_tesfdfsdfsdt11" {
		t.Error()
	}
}

func Test_save_3(t *testing.T) {
	type TestTest1 struct {
		BaseObject `bson:"inline"`
	}

	obj := Obj(&TestTest1{})
	err := Coll(obj).Save(obj)
	if err != nil {
		t.Errorf("BaseObject not inline %s", err)
	}
}

func Test_save_4(t *testing.T) {
	type TestTest1 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}
	test_obj := TestTest1{}

	obj1 := Obj(&test_obj)
	obj1.Field("TestStr").Set("TestStr")
	err := Coll(obj1).Save(obj1)
	if err != nil {
		t.Errorf("BaseObject not inline %s", err)
	}

	obj2 := Obj(&TestTest1{})
	err = Coll(obj2).Load(obj2)
	if err == nil {
		t.Error()
	}

	id, err := obj1.GetID()
	obj2.SetID(id)
	err = Coll(obj2).Load(obj2)
	if err != nil {
		t.Error()
	}

	s1 := obj1.Field("TestStr").Interface().(string)
	s2 := obj2.Field("TestStr").Interface().(string)
	if s1 != s2 {
		t.Error()
	}
}

func Test_save_5(t *testing.T) {
	type TestTest_save_5 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}

	test := TestTest_save_5{}
	test_slice := []TestTest_save_5{}

	obj := Obj(&test)
	obj_slice := Obj(&test_slice)

	Coll(obj).Drop()

	test.TestStr = "TestStr"
	err := Coll(obj).Save(obj)
	if err != nil {
		t.Error(err)
	}

	err = Coll(obj_slice).Find(obj_slice, primitive.M{})
	if err != nil {
		t.Error(err)
	}
	if len(test_slice) != 1 {
		t.Error()
	}
}

func Test_save_6(t *testing.T) {
	type TestTest_save_5 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}

	test := TestTest_save_5{}
	obj := Obj(&test)

	Coll(obj).Drop()

	test.TestStr = "TestStr"
	err := Coll(obj).Save(obj)
	if err != nil {
		t.Error(err)
	}

	var test2 TestTest_save_5
	obj2 := Obj(&test2)
	err = Coll(obj2).Find(obj2, primitive.M{})
	if err != nil {
		t.Error(err)
	}
	if test2.TestStr != "TestStr" {
		t.Error()
	}
}

func Test_save_7(t *testing.T) {
	type Test_save_7 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}

	test1 := Test_save_7{}
	test1.TestStr = "TestStr"
	obj1 := Obj(&test1)

	Coll(obj1).Drop()

	err := Coll(obj1).Save(obj1)
	if err != nil {
		t.Error(err)
	}

	test_slice := []Test_save_7{}
	test_slice = append(test_slice, Test_save_7{})
	obj_slice := Obj(&test_slice)
	err = Coll(obj_slice).Find(obj_slice, primitive.M{})
	if err != nil {
		t.Error(err)
	}
	l := len(test_slice)
	if l != 1 {
		t.Error()
	}
}

func Test_save_8(t *testing.T) {
	type Test_save_8 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}

	test := Test_save_8{}
	obj1 := Obj(&test)
	Coll(obj1).Drop()

	test2 := Test_save_8{}
	test2.TestStr = "TestStr"
	obj2 := Obj(&test2)
	err := Coll(obj2).Save(obj2)
	if err != nil {
		t.Error(err)
	}

	err = Coll(obj1).Delete(obj2)

	test_slice := []Test_save_8{}
	obj_slice := Obj(&test_slice)
	err = Coll(obj_slice).Find(obj_slice, primitive.M{})
	if err != nil {
		t.Error(err)
	}
	l := len(test_slice)
	if l != 0 {
		t.Error()
	}
}

func Test_save_9(t *testing.T) {
	type Test_save_9 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}
	test1 := Test_save_9{}
	test1.TestStr = "TestStr1"
	obj1 := Obj(&test1)

	test2 := Test_save_9{}
	test2.TestStr = "TestStr2"
	obj2 := Obj(&test2)

	test3 := Test_save_9{}
	test3.TestStr = "TestStr3"
	obj3 := Obj(&test3)

	Coll(obj1).Drop()

	err := Coll(obj1).Save(obj1)
	if err != nil {
		t.Error(err)
	}

	err = Coll(obj2).Save(obj2)
	if err != nil {
		t.Error(err)
	}

	err = Coll(obj3).Save(obj3)
	if err != nil {
		t.Error(err)
	}

	test_find := []Test_save_9{}
	obj_find := Obj(&test_find)
	err = Coll(obj_find).Find(obj_find, primitive.M{})
	l := len(test_find)
	if l != 3 {
		t.Error()
	}

	obj_to_delete := []Test_save_9{}
	obj_to_delete = append(obj_to_delete, Test_save_9{
		BaseObject: BaseObject{
			Id: test1.Id,
		},
	})
	obj_to_delete = append(obj_to_delete, Test_save_9{
		BaseObject: BaseObject{
			Id: test2.Id,
		},
	})
	obj_to_delete = append(obj_to_delete, Test_save_9{
		BaseObject: BaseObject{
			Id: test3.Id,
		},
	})

	err = Coll(obj1).Delete(Obj(&obj_to_delete))

	err = Coll(obj_find).Find(obj_find, primitive.M{})
	if err != nil {
		t.Error(err)
	}
	l = len(test_find)
	if l != 0 {
		t.Error()
	}
}
