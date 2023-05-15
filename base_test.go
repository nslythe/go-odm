package goodm

import (
	"testing"
)

func Test_Base_1(t *testing.T) {

	if Obj("11").String() != "11" {
		t.Error()
	}

	if Obj(11).String() != "11" {
		t.Error()
	}

	if Obj(true).String() != "true" {
		t.Error()
	}

	if Obj(3.1).String() != "3.10000" {
		t.Error()
	}
}

func Test_Base_2(t *testing.T) {
	type T1 struct {
		Test string
	}

	if Obj(T1{Test: "1"}).Interface().(T1).Test != "1" {
		t.Error()
	}

	if Obj("test").Interface().(string) != "test" {
		t.Error()
	}
}

func Test_Base_3(t *testing.T) {
	type T1 struct {
		TestStr string
	}
	type T2 struct {
		TestInt int
		T1      T1
	}

	if Obj(T1{TestStr: "1"}).Field("TestStr").String() != "1" {
		t.Error()
	}
	if Obj(T2{
		TestInt: 1,
		T1: T1{
			TestStr: "TestStr",
		}}).Field("TestInt").String() != "1" {
		t.Error()
	}
	if Obj(T2{
		TestInt: 1,
		T1: T1{
			TestStr: "TestStr",
		}}).Field("T1").Field("TestStr").String() != "TestStr" {
		t.Error()
	}
}

func Test_Base_4(t *testing.T) {
	type T1 struct {
		TestStr string
	}

	obj := []T1{}
	obj = append(obj, T1{TestStr: "TestStr1"})
	obj = append(obj, T1{TestStr: "TestStr2"})

	if Obj(obj).Index(0).Field("TestStr").String() != "TestStr1" {
		t.Error()
	}
	if Obj(obj).Index(1).Field("TestStr").String() != "TestStr2" {
		t.Error()
	}
}

func Test_Base_5(t *testing.T) {
	test_str := ""
	Obj(&test_str).Set("test_str")
	if test_str != test_str {
		t.Error()
	}
}

func Test_Base_6(t *testing.T) {
	type T1 struct {
		TestStr string
	}

	obj := []T1{}

	Obj(&obj).Append(T1{TestStr: "1"})
	if len(obj) != 1 {
		t.Error()
	}
}

// func Test_Obj_1(t *testing.T) {
// 	type Test1 struct {
// 	}
// 	if Obj(Test1{}).IsSlice {
// 		t.Error("Not a slice")
// 	}
// 	if Obj(Test1{}).IsPtr {
// 		t.Error("Not a ptr")
// 	}
// 	if Obj(&Test1{}).IsSlice {
// 		t.Error("Not a slice")
// 	}
// 	if !Obj(&Test1{}).IsPtr {
// 		t.Error("Is a ptr")
// 	}
// 	if !Obj([]Test1{}).IsSlice {
// 		t.Error("Is a slice")
// 	}
// 	if Obj([]Test1{}).IsPtr {
// 		t.Error("Not a ptr")
// 	}
// }

// func Test_Obj_2(t *testing.T) {
// 	type Test1 struct {
// 		Test string
// 	}
// 	if Obj(Test1{}).Field("dasdas") != nil {
// 		t.Error("Failed")
// 	}
// 	if Obj(Test1{}).Field("Test") == nil {
// 		t.Error("Failed")
// 	}

// 	if Obj(Test1{}).Field("dasdas") != nil {
// 		t.Error("Failed")
// 	}
// 	if Obj(Test1{}).Field("Test") == nil {
// 		t.Error("Failed")
// 	}

// 	if Obj([]Test1{}).Field("dasdas") != nil {
// 		t.Error("Failed")
// 	}
// 	if Obj(Test1{}).Field("Test") == nil {
// 		t.Error("Failed")
// 	}
// }

// func Test_Obj_3(t *testing.T) {
// 	type Test1 struct {
// 		Test string
// 	}

// 	obj := Test1{
// 		Test: "test",
// 	}
// 	if Obj(obj).Data().Field("Test").String() != "test" {
// 		t.Error("Failed")
// 	}
// 	if Obj(obj).Data().Field("Test1").IsValid() {
// 		t.Error("Failed")
// 	}

// 	if Obj(obj).Data().Field("Test").String() != "test" {
// 		t.Error("Failed")
// 	}
// 	if Obj(&obj).Data().Field("Test1").IsValid() {
// 		t.Error("Failed")
// 	}
// }

// func Test_Obj_4(t *testing.T) {
// 	type Test1 struct {
// 		Test string
// 	}

// 	obj := []Test1{}

// 	obj = append(obj, Test1{
// 		Test: "1",
// 	})
// 	obj = append(obj, Test1{
// 		Test: "2",
// 	})
// 	obj = append(obj, Test1{
// 		Test: "3",
// 	})
// 	obj = append(obj, Test1{
// 		Test: "4",
// 	})

// 	if Obj(obj).Data().Field("Test").String() != "1" {
// 		t.Error("Failed")
// 	}
// 	if Obj(obj).Data().Field("Test1").IsValid() {
// 		t.Error("Failed")
// 	}
// 	if Obj(obj).DataIdx(0).Field("Test").String() != "1" {
// 		t.Error("Failed")
// 	}
// 	if Obj(obj).DataIdx(1).Field("Test").String() != "2" {
// 		t.Error("Failed")
// 	}
// 	if Obj(obj).DataIdx(2).Field("Test").String() != "3" {
// 		t.Error("Failed")
// 	}
// 	if Obj(obj).DataIdx(3).Field("Test").String() != "4" {
// 		t.Error("Failed")
// 	}
// }

// func Test_Obj_5(t *testing.T) {
// 	type Test1 struct {
// 		Test string
// 	}

// 	obj1 := Test1{
// 		Test: "1",
// 	}

// 	obj2 := Test1{
// 		Test: "1",
// 	}

// 	if Obj(obj1).Data().Field("Test").String() != "1" {
// 		t.Error("Failed")
// 	}
// 	Obj(&obj1).Data().Field("Test").Set("test")
// 	if Obj(obj1).Data().Field("Test").String() != "test" {
// 		t.Error("Failed")
// 	}
// 	if obj1.Test != "test" {
// 		t.Error("Failed")
// 	}

// 	Obj(&obj2).Data().Field("Test").Set("test")
// 	if Obj(obj2).Data().Field("Test").String() != "test" {
// 		t.Error("Failed")
// 	}
// 	if obj2.Test != "test" {
// 		t.Error("Failed")
// 	}
// }

// func Test_Obj_6(t *testing.T) {
// 	type Test struct {
// 		Test string
// 	}
// 	type Test1 struct {
// 		Teststr string
// 		Test    Test
// 	}

// 	obj1 := Test1{
// 		Teststr: "1",
// 	}

// 	if Obj(obj1).Data().Field("Teststr").String() != "1" {
// 		t.Error("Failed")
// 	}
// 	if Obj(obj1).Data().Field("Test").Interface().(Test).Test != "" {
// 		t.Error("Failed")
// 	}

// 	Obj(&obj1).Data().Field("Test").Field("Test").Set("2")
// 	if Obj(&obj1).Data().Field("Test").Field("Test").String() != "2" {
// 		t.Error("Failed")
// 	}
// 	if obj1.Test.Test != "2" {
// 		t.Error("Failed")
// 	}
// }

// func Test_Obj_7(t *testing.T) {
// 	type Test struct {
// 		Test string
// 	}
// 	type Test1 struct {
// 		Teststr string
// 		Test    Test
// 	}

// 	obj1 := []Test1{}

// 	Obj(&obj1).Append(Test1{})
// 	if len(obj1) != 1 {
// 		t.Error("Failed")
// 	}
// }

// 	bt2 := TestBase{
// 		BaseObject: &BaseObject{
// 			Id: bt1.Id,
// 		},
// 	}
// 	err = Load(&bt2)
// 	if err != nil {
// 		t.Errorf("Supose to pass %s", err)
// 	}

// 	if bt2.Test != bt1.Test {
// 		t.Errorf("Result does not match")
// 	}
// }

// func Test_Add_2(t *testing.T) {
// 	bt1 := TestBase2{
// 		BaseObject: &BaseObject{},
// 		TestBase: TestBase{
// 			Test: "test",
// 		},
// 		Test1: "test1",
// 	}

// 	err := Save(&bt1)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	bt2 := TestBase2{
// 		BaseObject: &BaseObject{
// 			Id: bt1.Id,
// 		},
// 	}

// 	err = Load(&bt2)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	if bt1.Test != bt2.Test {
// 		t.Errorf("Failed")
// 	}
// }

// func Test_Add_3(t *testing.T) {
// 	bt1 := BaseObject{}

// 	err := Save(&bt1)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}
// }

// func Test_Add_4(t *testing.T) {
// 	bt1 := TestBase3{
// 		BaseObject: &BaseObject{},
// 		TestBase: TestBase{
// 			Test: "test",
// 		},
// 		Test1: "test1",
// 	}

// 	err := Save(&bt1)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	bt2 := TestBase3{
// 		BaseObject: &BaseObject{
// 			Id: bt1.Id,
// 		},
// 	}

// 	err = Load(&bt2)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	if bt1.Test != bt2.Test {
// 		t.Errorf("Failed")
// 	}

// 	if bt1.Test1 != bt2.Test1 {
// 		t.Errorf("Failed")
// 	}
// }

// func Test_Add_5(t *testing.T) {
// 	bt1 := BaseObject{}

// 	err := Save(&bt1)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}
// }

// func Test_Add_6(t *testing.T) {
// 	bt1 := BaseObject{}
// 	err := Save(&bt1)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	bt2 := BaseObject{}
// 	bt2.Id = bt1.Id
// 	err = Save(&bt2)
// 	if err == nil {
// 		t.Errorf("Failed")
// 	}
// 	err = Update(&bt2)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}
// }

// func Test_Add_7(t *testing.T) {
// 	bt1 := TestBase2{
// 		BaseObject: &BaseObject{},
// 		TestBase: TestBase{
// 			Test: "test",
// 		},
// 		Test1: "test1",
// 	}

// 	err := Save(&bt1)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	bt1.TestBase.Test = "test_test_test"
// 	bt1.Test1 = "test_test_test"

// 	err = Update(&bt1)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	bt2 := TestBase2{
// 		BaseObject: &BaseObject{
// 			Id: bt1.Id,
// 		},
// 		TestBase: TestBase{},
// 	}
// 	err = Load(&bt2)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	if bt2.TestBase.Test != "test_test_test" {
// 		t.Errorf("Failed")
// 	}
// 	if bt2.Test1 != "test_test_test" {
// 		t.Errorf("Failed")
// 	}
// }

// func Test_Find_1(t *testing.T) {
// 	ctx, client, cancel, err := CreateConnection()
// 	defer client.Disconnect(ctx)
// 	defer cancel()
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	client.Database(connection_string.Database).Collection(GetCollectionName(GetCollectionName(&TestBase2{}))).Drop(ctx)

// 	bt2 := []TestBase2{}

// 	bt1 := TestBase2{
// 		BaseObject: &BaseObject{},
// 		TestBase: TestBase{
// 			Test: "test",
// 		},
// 		Test1: "test1",
// 	}

// 	err = Save(&bt1)
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}

// 	err = Find(&bt2, primitive.M{})
// 	if err != nil {
// 		t.Errorf("Failed")
// 	}
// 	if len(bt2) != 1 {
// 		t.Errorf("Failed %d", len(bt2))
// 	}
// }
