package main

import (
	"fmt"
	"testing"
)

func TestMurmur3Hash(t *testing.T) {
	v := Murmur3Hash("test.txt")
	fmt.Println(v)
}