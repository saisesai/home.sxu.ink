package model

import (
	"fmt"
	"testing"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("home", "123456")
	if err != nil {
		panic(err)
	}
	fmt.Println(u)
}
