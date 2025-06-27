package HttpFC

import (
	"github.com/danvei233/fchandler2go/Tencent/HttpFC/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type TencentHttpFCProvider struct {
	config config.Config
}

func (t TencentHttpFCProvider) GetInputAllowedList() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(http.Request{}),
		reflect.TypeOf([]byte(nil)),
	}
}

func (t TencentHttpFCProvider) GetOutputAllowedList() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(""),
		reflect.TypeOf([]byte(nil)),
		reflect.TypeOf(),
	}
}

func (t TencentHttpFCProvider) TransIn(c *gin.Context) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t TencentHttpFCProvider) TransOut(c *gin.Context, resRaw []byte) error {
	//TODO implement me
	panic("implement me")
}
