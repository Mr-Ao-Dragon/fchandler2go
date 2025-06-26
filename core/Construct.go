package core

import (
	"github.com/aliyun/fc-runtime-go-sdk/events"
	"github.com/danvei233/fchandler2go/config"
	"github.com/danvei233/fchandler2go/mock"
	"github.com/danvei233/fchandler2go/utills"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func ConvertRequest(c *gin.Context, config config.Config) (*events.HTTPTriggerEvent, error) {
	var RequestContext events.HTTPTriggerRequestContext
	// some data need custom
	mock.MakeFakeContext(c, &RequestContext, config)
	//mocked
	now := time.Now()
	RequestContext.TimeEpoch = utills.StringPtr(now.UTC().Format(time.RFC3339))
	RequestContext.Time = utills.StringPtr(strconv.FormatInt(now.Unix(), 10))
	RequestContext.Http.Method = utills.StringPtr(c.Request.Method)
	RequestContext.Http.Path = utills.StringPtr(c.Request.URL.Path)
	RequestContext.Http.SourceIp = utills.StringPtr(c.ClientIP())
	RequestContext.Http.UserAgent = utills.StringPtr(c.Request.UserAgent())
	RequestContext.Http.Protocol = utills.StringPtr(c.Request.Proto)

	// please map it
	var event events.HTTPTriggerEvent
	event.Version = utills.StringPtr("v1")
	event.RawPath = utills.StringPtr(c.Request.URL.Path)
	body, err := utills.BodyReader(c.Request.Body, c.ContentType())
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return nil, err
	}
	event.Body = utills.StringPtr(body)
	event.IsBase64Encoded = utills.IsBin(c.ContentType())

	event.Headers = utills.Header2map(c.Request.Header)
	event.QueryParameters = utills.Param2map(c.Request.URL.Query())
	event.TriggerContext = RequestContext
	return &event, nil
}
