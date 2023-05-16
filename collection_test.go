package goodm

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_GetColletionName_1(t *testing.T) {
	type Test struct{}
	if GetCollectionName(Test{}) != "goodm__test" {
		t.Error("failed")
	}
}

func Test_GetColletionName_2(t *testing.T) {
	type Test_allo struct{}
	if GetCollectionName(Test_allo{}) != "goodm__test_allo" {
		t.Error("failed")
	}
}

func Test_GetColletionName_3(t *testing.T) {
	type tests struct{}
	if GetCollectionName(tests{}) != "goodm_tests" {
		t.Error("failed")
	}
}

func Test_GetColletionName_4(t *testing.T) {
	type TestTest struct{}
	if GetCollectionName(TestTest{}) != "goodm__test_test" {
		t.Error("failed")
	}
}

func Test_GetColletionName_5(t *testing.T) {
	type TestTest1 struct{}
	if GetCollectionName(TestTest1{}) != "goodm__test_test1" {
		t.Error("failed")
	}
}

func Test_save_1(t *testing.T) {
	type TestTest1 struct{}
	obj := TestTest1{}
	err := Coll(obj).Save(&obj)
	if err == nil {
		t.Error("No BaseObject in parent")
	}
}

func Test_save_2(t *testing.T) {
	type TestTest1 struct {
		BaseObject
	}
	obj := TestTest1{}
	err := Coll(obj).Save(&obj)
	if err == nil {
		t.Errorf("BaseObject not inline %s", err)
	}
}

func Test_save_3(t *testing.T) {
	type TestTest1 struct {
		BaseObject `bson:"inline"`
	}

	obj := TestTest1{}
	err := Coll(obj).Save(&obj)
	if err != nil {
		t.Errorf("BaseObject not inline %s", err)
	}
}

func Test_save_4(t *testing.T) {
	type TestTest1 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}

	obj1 := TestTest1{}
	obj1.TestStr = "TestStr"
	err := Coll(obj1).Save(&obj1)
	if err != nil {
		t.Errorf("BaseObject not inline %s", err)
	}

	obj2 := TestTest1{}
	err = Coll(obj2).Load(&obj2)
	if err == nil {
		t.Error()
	}

	obj2.Id = obj1.Id
	err = Coll(obj2).Load(&obj2)
	if err != nil {
		t.Error()
	}
	if obj2.TestStr != obj1.TestStr {
		t.Error()
	}
}

func Test_save_5(t *testing.T) {
	type TestTest_save_5 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}

	Coll(TestTest_save_5{}).Drop()

	obj1 := TestTest_save_5{}
	obj1.TestStr = "TestStr"
	err := Coll(obj1).Save(&obj1)
	if err != nil {
		t.Error(err)
	}

	obj2 := []TestTest_save_5{}
	err = Coll(obj2).Find(&obj2, primitive.M{})
	if err != nil {
		t.Error(err)
	}
	if len(obj2) != 1 {
		t.Error()
	}
}

func Test_save_6(t *testing.T) {
	type TestTest_save_5 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}

	Coll(TestTest_save_5{}).Drop()

	obj1 := TestTest_save_5{}
	obj1.TestStr = "TestStr"
	err := Coll(obj1).Save(&obj1)
	if err != nil {
		t.Error(err)
	}

	var obj2 TestTest_save_5
	err = Coll(obj2).Find(&obj2, primitive.M{})
	if err != nil {
		t.Error(err)
	}
	if obj2.TestStr != "TestStr" {
		t.Error()
	}
}

func Test_save_7(t *testing.T) {
	type Test_save_7 struct {
		BaseObject `bson:"inline"`
		TestStr    string
	}

	Coll(Test_save_7{}).Drop()

	obj1 := Test_save_7{}
	obj1.TestStr = "TestStr"
	err := Coll(obj1).Save(&obj1)
	if err != nil {
		t.Error(err)
	}

	obj2 := []Test_save_7{}
	obj2 = append(obj2, Test_save_7{})
	err = Coll(obj2).Find(&obj2, primitive.M{})
	if err != nil {
		t.Error(err)
	}
	l := len(obj2)
	if l != 1 {
		t.Error()
	}
}
