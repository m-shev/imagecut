package linkedlist

import "errors"

// RemoveError occurs when item has already been removed from the list
const RemoveError = "item has already been removed from the list"

// Item is node of doubly linked list
type Item struct {
	value interface{}
	prev  *Item
	next  *Item
	list  *List
}

// Value returns item value
func (item *Item) Value() interface{} {
	return item.value
}

// Next returns next *item
func (item *Item) Next() *Item {
	return item.next
}

// Prev returns prev *item
func (item *Item) Prev() *Item {
	return item.prev
}

// Remove item from list
func (item *Item) Remove() error {
	if item.list == nil {
		return errors.New(RemoveError)
	}

	switch {
	// item is first
	case item.prev == nil && item.next != nil:
		item.list.first = item.next
		item.list.first.prev = nil

	// item is last
	case item.prev != nil && item.next == nil:
		item.list.last = item.prev
		item.list.last.next = nil

	// item is single
	case item.list.length == 1:
		item.list.first = nil
		item.list.last = nil

	default:
		item.prev.next = item.next
		item.next.prev = item.prev
	}

	item.list.length--
	item.list = nil
	item.prev = nil
	item.next = nil

	return nil
}
