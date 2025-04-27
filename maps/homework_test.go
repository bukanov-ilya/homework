package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type node struct {
	key, value  int
	left, right *node
}

type OrderedMap struct {
	root *node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{root: nil}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = &node{key: key, value: value}
		m.size++

		return
	}

	current := m.root
	for {
		if key == m.root.key {
			current.value = value
			m.size++

			return
		} else if key < current.key {
			if current.left == nil {
				current.left = &node{key: key, value: value}
				m.size++
				return
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = &node{key: key, value: value}
				m.size++
				return
			}
			current = current.right
		}
	}

}

func (m *OrderedMap) Erase(key int) {
	m.root = deleteNode(m.root, key, &m.size)
}

func deleteNode(n *node, key int, sz *int) *node {
	if n == nil {
		return nil
	}
	switch {
	case key < n.key:
		n.left = deleteNode(n.left, key, sz)
	case key > n.key:
		n.right = deleteNode(n.right, key, sz)
	default:
		*sz--
		if n.left == nil {
			return n.right
		}
		if n.right == nil {
			return n.left
		}

		succ := n.right
		for succ.left != nil {
			succ = succ.left
		}

		n.key, n.value = succ.key, succ.value

		n.right = deleteNode(n.right, succ.key, sz)
	}
	return n
}

func (m *OrderedMap) Contains(key int) bool {
	current := m.root

	for current != nil {
		if key == current.key {
			return true
		} else if key < current.key {
			current = current.left
		} else {
			current = current.right
		}
	}

	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	current := m.root
	var traverse func(n *node)
	traverse = func(n *node) {
		if n == nil {
			return
		}
		traverse(n.left)
		action(n.key, n.value)
		traverse(n.right)
	}
	traverse(current)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
