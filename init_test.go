package goodm

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//setup
	Init(Config{
		ConnectionString: "mongodb://mongo1.home.slythe.net:27017,mongo2.home.slythe.net:27018,mongo3.home.slythe.net:27019/test_gogame?replicaSet=rs0",
	})
	code := m.Run()
	os.Exit(code)
}

func Test_Ping(t *testing.T) {
	err := Ping()
	if err != nil {
		t.Error("Failed")
	}
}
