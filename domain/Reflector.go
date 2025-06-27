package domain

import (
	"context"
	"github.com/gin-gonic/gin"
	"reflect"
)

type Reflector interface {
	T(h interface{}) gin.HandlerFunc
}

type StandardHandler func(ctx context.Context, b []byte) (output []byte, err error)

type HandlerStandardizer interface {
	GetStandardHandler() (StandardHandler, error)
	CheckOutputValid(list []reflect.Type) (bool, reflect.Type)
	CheckInputValid(list []reflect.Type) (bool, reflect.Type)
}

type MapProvider interface {
	GetInputAllowedList() []reflect.Type
	GetOutputAllowedList() []reflect.Type
	TransIn(c *gin.Context) ([]byte, error)
	TransOut(c *gin.Context, resRaw []byte) error
}
