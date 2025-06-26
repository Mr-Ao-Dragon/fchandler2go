package mock

import (
	"fmt"
	"github.com/aliyun/fc-runtime-go-sdk/events"
	"github.com/danvei233/fchandler2go/config"
	"github.com/danvei233/fchandler2go/utills"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math/rand"
	"strings"
)

func MakeFakeID(prefix int, config config.Config) string {
	if config.Output.RequestIDFromMock {
		return fmt.Sprintf("%d-%s", prefix, uuid.New().String())
	}
	return config.Output.RequestIDFromCustom
}
func makeFakeRequestId() string {
	return fmt.Sprintf("%d-%s", 1, uuid.New().String())
}
func makeFakeAccountId() string {
	const length = 15
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = byte(rand.Intn(10) + '0') // 生成 '0'..'9'
	}
	return string(buf)
}
func makeFakeDomainName() string {
	const idLen = 32
	const hexChars = "0123456789abcdef"
	buf := make([]byte, idLen)
	regions := []string{"Beijing", "Shanghai", "Guangdong", "火星", "比邻星", "月球"}
	for i := range buf {
		buf[i] = hexChars[rand.Intn(len(hexChars))]
	}
	triggerID := string(buf)

	region := regions[rand.Intn(len(regions))]

	return fmt.Sprintf("%s.%s.fcapp.run", triggerID, region)
}

func MakeFakeContext(c *gin.Context, RequestContext *events.HTTPTriggerRequestContext, config config.Config) {
	RequestContext.RequestId = mockComplete(c,
		config.Input.RequestIDOrigin,
		"X-Fc-Request-Id",
		config.Input.RequestIDFromCustom,
		makeFakeRequestId)

	RequestContext.AccountId = mockComplete(c,
		config.Input.AccountIDOrigin,
		"X-Fc-Account-Id",
		config.Input.AccountIDFromCustom,
		makeFakeAccountId)

	RequestContext.DomainName = mockComplete(c,
		config.Input.DomainNameOrigin,
		"X-Fc-Domain-Name",
		config.Input.DomainNameFromCustom,
		makeFakeDomainName)
	RequestContext.DomainPrefix = mockComplete(c,
		config.Input.DomainPrefixOrigin,
		"X-Fc-Domain-Prefix",
		config.Input.DomainPrefixFromCustom,
		func() string {
			return strings.Split(*RequestContext.DomainName, ".")[0]
		})

}
func mockComplete(c *gin.Context, origin int, name string, custom string, mocker func() string) *string {
	switch origin {
	case config.FromCtx:
		v, ok := c.Value(name).(string)
		if !ok {
			return utills.StringPtr("")
		}
		return utills.StringPtr(v)
	case config.FromHeader:
		return utills.StringPtr(c.Request.Header.Get(name))
	case config.FromCustom:
		return utills.StringPtr(custom)
	}
	return utills.StringPtr(mocker())
}
