package linkedlist

import (
	"fmt"
	"strings"
	"testing"
)

// Test Item

func TestItem_Value(t *testing.T) {
	list := List{}

	intValue := 1
	const stringValue = "some test value"

	list.PushBack(intValue)
	list.PushBack(stringValue)

	expectedInt := intValue
	got := list.First().Value()

	if expectedInt != got {
		printError(t, expectedInt, got)
	}

	expectedString := stringValue
	got = list.Last().Value()

	if expectedString != got {
		printError(t, expectedString, got)
	}
}

func TestItem_Prev(t *testing.T) {
	list := List{}

	list.PushBack(1)
	list.PushBack("some test value")

	expected := 1
	got := list.Last().Prev().Value()

	if expected != got {
		printError(t, expected, got)
	}
}

func TestItem_Next(t *testing.T) {
	list := List{}

	list.PushBack(1)
	list.PushBack("some test value")

	expected := "some test value"
	got := list.First().Next().Value()

	if expected != got {
		printError(t, expected, got)
	}
}

func TestItem_RemoveFirstItem(t *testing.T) {
	list := populateList(10)

	_ = list.first.Remove()

	expected := "2 3 4 5 6 7 8 9 10"
	got := printListToString(*list)

	if expected != got {
		printError(t, expected, got)
	}
}

func TestItem_RemoveLast(t *testing.T) {
	list := populateList(10)

	_ = list.Last().Remove()

	expected := "1 2 3 4 5 6 7 8 9"
	got := printListToString(*list)

	if expected != got {
		printError(t, expected, got)
	}
}

func TestItem_RemoveSingle(t *testing.T) {
	list := populateList(1)

	_ = list.First().Remove()

	expected := ""
	got := printListToString(*list)

	if expected != got {
		printError(t, expected, got)
	}
}

func TestItem_Remove(t *testing.T) {
	list := populateList(10)

	_ = list.First().Next().Next().Remove()

	expected := "1 2 4 5 6 7 8 9 10"
	got := printListToString(*list)

	if expected != got {
		printError(t, expected, got)
	}
}

func TestItem_RemoveError(t *testing.T) {
	list := populateList(10)
	first := list.First()

	_ = first.Remove()

	got := first.Remove()
	expected := RemoveError

	if expected != got.Error() {
		printError(t, expected, got)
	}
}

// Test List

func TestList_PushBack(t *testing.T) {
	list := List{}

	list.PushFront(1)

	expected := "1"
	got := printListToString(list)

	if expected != got {
		printError(t, expected, got)
	}

	list.PushFront(2)
	list.PushFront(3)

	expected = "3 2 1"
	got = printListToString(list)

	if expected != got {
		printError(t, expected, got)
	}
}

func TestList_PushFront(t *testing.T) {
	list := List{}

	list.PushBack(1)

	expected := "1"
	got := printListToString(list)

	if expected != got {
		printError(t, expected, got)
	}

	list.PushBack(2)
	list.PushBack(3)

	expected = "1 2 3"
	got = printListToString(list)

	if expected != got {
		printError(t, expected, got)
	}
}

func TestList_First(t *testing.T) {
	list := populateList(10)

	expected := "1"
	got := fmt.Sprint(list.First().Value())

	if expected != got {
		printError(t, expected, got)
	}
}

func TestList_Last(t *testing.T) {
	list := populateList(10)

	expected := "10"
	got := fmt.Sprint(list.Last().Value())

	if expected != got {
		printError(t, expected, got)
	}
}

func TestList_Len(t *testing.T) {
	list := populateList(5)

	expected := uint(5)
	got := list.Len()

	if got != expected {
		printError(t, expected, got)
	}
}

// Utils

func printListToString(list List) string {
	var strBuilder strings.Builder

	current := list.First()

	for current != nil {
		strBuilder.WriteString(fmt.Sprint(current.Value()))

		if current.next != nil {
			strBuilder.WriteString(" ")
		}

		current = current.next
	}

	return strBuilder.String()
}

func printError(t *testing.T, exp, got interface{}) {
	t.Error("\nExp:\n", exp, "\nGot:\n", got)
}

func populateList(count uint) *List {
	list := List{}

	for list.length < count {
		list.PushBack(list.length + 1)
	}

	return &list
}
