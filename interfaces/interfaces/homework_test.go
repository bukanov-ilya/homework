package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	NotEmptyStruct bool
}
type MessageService struct {
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]func() interface{}
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]func() interface{}),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	constructorFunc, ok := constructor.(func() interface{})
	if !ok {
		panic(fmt.Sprintf("constructor for %s is not a function of type func() interface{}", name))
	}

	c.constructors[name] = constructorFunc
}

func (c *Container) Resolve(name string) (interface{}, error) {
	constructor, ok := c.constructors[name]
	if !ok {
		return nil, fmt.Errorf("constructor for %s not registered", name)
	}

	return constructor(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
