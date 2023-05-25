package goodm

import (
	"testing"
)

func Test_Base_2(t *testing.T) {
	type T1 struct {
		Test string
	}

	if Obj(&T1{Test: "1"}).Interface().(*T1).Test != "1" {
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

	if Obj(&T1{TestStr: "1"}).Field("TestStr").String() != "1" {
		t.Error()
	}
	if Obj(&T2{
		TestInt: 1,
		T1: T1{
			TestStr: "TestStr",
		}}).Field("TestInt").String() != "1" {
		t.Error()
	}
	if Obj(&T2{
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

	if Obj(&obj).Index(0).Field("TestStr").String() != "TestStr1" {
		t.Error()
	}
	if Obj(&obj).Index(1).Field("TestStr").String() != "TestStr2" {
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

	Obj(&obj).Append(Obj(&T1{TestStr: "1"}))
	if len(obj) != 1 {
		t.Error()
	}
}

func Test_Base_7(t *testing.T) {
	type T1 struct {
		TestStr string
	}

	obj := []T1{}
	obj = append(obj, T1{})
	Obj(&obj).Clear()

	if len(obj) != 0 {
		t.Error()
	}
}

func Test_Base_8(t *testing.T) {
	type T1 struct {
		TestStr string
	}

	obj := []T1{}
	obj = append(obj, T1{})
	obj = append(obj, T1{})
	obj = append(obj, T1{})

	if Obj(&obj).Len() != 3 {
		t.Error()
	}
}
