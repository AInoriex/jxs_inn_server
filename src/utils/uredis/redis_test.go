package uredis

import (
	"fmt"
	"testing"
)

func TestSetString(t *testing.T) {
	con := New("127.0.0.1:6379", "", 0)

	err := SetString(con, "a", "start", 20)
	if err != nil {
		t.Error("err=", err)
	}
}

func TestGetString(t *testing.T) {
	con := New("127.0.0.1:6379", "", 0)
	data, err := GetString(con, "a")
	fmt.Println(err)
	fmt.Println(string(data))
}
