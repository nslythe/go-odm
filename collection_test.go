package goodm

import (
	"testing"
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
	Init(Config{
		ConnectionString: "mongodb://mongo1.home.slythe.net:27017,mongo2.home.slythe.net:27018,mongo3.home.slythe.net:27019/test_gogame?replicaSet=rs0",
	})

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
	Init(Config{
		ConnectionString: "mongodb://mongo1.home.slythe.net:27017,mongo2.home.slythe.net:27018,mongo3.home.slythe.net:27019/test_gogame?replicaSet=rs0",
	})

	type TestTest1 struct {
		BaseObject `bson:"inline"`
	}

	obj := TestTest1{}
	err := Coll(obj).Save(&obj)
	if err != nil {
		t.Errorf("BaseObject not inline %s", err)
	}
}
