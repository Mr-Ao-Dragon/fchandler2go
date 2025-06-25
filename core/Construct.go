package core

import (
	"encoding/base64"
	"github.com/aliyun/fc-runtime-go-sdk/events"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func ConvertRequest(c *gin.Context) (*events.HTTPTriggerEvent, error) {
	var RequestContext events.HTTPTriggerRequestContext
	// some data need custom

	RequestContext.RequestId = StringPtr(c.Request.Header.Get("X-Fc-Request-Id"))
	RequestContext.AccountId = StringPtr(c.Request.Header.Get("X-Fc-Account-Id"))
	RequestContext.DomainName = StringPtr(c.Request.Header.Get("X-Fc-Domain-Name"))
	RequestContext.DomainPrefix = StringPtr(c.Request.Header.Get("X-Fc-Domain-Prefix"))

	now := time.Now()
	RequestContext.TimeEpoch = StringPtr(now.UTC().Format(time.RFC3339))
	RequestContext.Time = StringPtr(strconv.FormatInt(now.Unix(), 10))
	RequestContext.Http.Method = StringPtr(c.Request.Method)
	RequestContext.Http.Path = StringPtr(c.Request.URL.Path)
	RequestContext.Http.SourceIp = StringPtr(c.ClientIP())
	RequestContext.Http.UserAgent = StringPtr(c.Request.UserAgent())
	RequestContext.Http.Protocol = StringPtr(c.Request.Proto)

	// map the responsedata
	var event events.HTTPTriggerEvent
	event.Version = StringPtr("v1")
	event.RawPath = StringPtr(c.Request.URL.Path)
	body, err := BodyReader(c.Request.Body, c.ContentType())
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return nil, err
	}
	event.Body = StringPtr(body)
	event.IsBase64Encoded = IsBin(c.ContentType())

	event.Headers = Header2map(c.Request.Header)
	event.QueryParameters = param2map(c.Request.URL.Query())
	event.TriggerContext = RequestContext
	return &event, nil
}
func Recall(c *gin.Context, response events.HTTPTriggerResponse) {
	w := c.Writer
	for key, value := range response.Headers {
		w.Header().Set(key, strings.Replace(value, ",", ";", -1))
	}
	w.WriteHeader(response.StatusCode)
	if response.IsBase64Encoded {
		finalres, err := base64.URLEncoding.DecodeString(response.Body)
		if err != nil {
			c.String(500, err.Error())
		}
		w.Write(finalres)
	}
	w.Write([]byte(response.Body))
}
