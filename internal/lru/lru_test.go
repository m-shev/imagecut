package lru

import (
	"strconv"
	"testing"
)

func TestLru_Set(t *testing.T) {
	lru := NewLru(3)
	var err error
	var excluded []interface{}

	for i := 0; i < 3; i++ {
		excluded, err = lru.Set(strconv.Itoa(i), i, 1)
	}

	if err != nil {
		t.Error(err.Error())
	}

	gotLen := len(excluded)
	expectedLen := 0

	if gotLen != expectedLen {
		t.Error("expected len of excluded:", expectedLen, "got:", gotLen)
	}

	excluded, err = lru.Set("3", 3, 1)

	if err != nil {
		t.Error(err.Error())
	}

	gotLen = len(excluded)
	expectedLen = 1

	if gotLen != expectedLen {
		t.Error("expected len of excluded:", expectedLen, "got:", gotLen)
	}
}

func TestLru_Get(t *testing.T) {
	var got int
	expected := 2

	lru := NewLru(1)
	_, err := lru.Set("1", 2, 1)

	if err != nil {
		t.Error(err)
	}

	v, err := lru.Get("1")

	if err != nil {
		t.Error(err)
	}

	if v != nil {
		got = v.(int)
	}

	if expected != got {
		t.Error("expected Value from cache:", expected, "got:", got)
	}
}