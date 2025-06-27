package AliyunWebFC

import (
	"encoding/json"
	"github.com/aliyun/fc-runtime-go-sdk/events"
	"github.com/danvei233/fchandler2go/AliyunWebFC/config"
	"github.com/gin-gonic/gin"
	"reflect"
)

// some broken features need to fix
type AliyunWebFCProvider struct {
	config config.Config
}

func (p *AliyunWebFCProvider) GetInputAllowedList() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(events.HTTPTriggerEvent{}),

		reflect.TypeOf([]byte(nil)),
	}
}
func (p *AliyunWebFCProvider) GetOutputAllowedList() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&events.HTTPTriggerResponse{}),
	}
}

func (p *AliyunWebFCProvider) TransIn(c *gin.Context) ([]byte, error) {
	reqRaw, err := ConvertRequest(c, p.config)
	if err != nil {
		return nil, err
	}
	req, err := json.Marshal(&reqRaw)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func (p *AliyunWebFCProvider) TransOut(c *gin.Context, resRaw []byte) error {
	res := new(events.HTTPTriggerResponse)
	err := json.Unmarshal(resRaw, res)
	if err != nil {
		return err
	}
	return Recall(c, res, p.config)
}
