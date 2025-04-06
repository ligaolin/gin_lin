package database

import (
	"fmt"
	"testing"
)

func TestToWhere(t *testing.T) {
	v := 1
	where, err := ToWhere([]Where{
		{Name: "in", Op: "in", Value: &v, Nullable: false},
		{Name: "name", Op: "like", Value: "name", Nullable: false},
		{Name: "isnull", Op: "null", Value: "data", Nullable: false},
		{Name: "set", Op: "set", Value: "set", Nullable: false},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(where)
}
