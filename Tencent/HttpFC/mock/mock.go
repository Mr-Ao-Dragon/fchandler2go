package mock

import (
	"github.com/danvei233/fchandler2go/Tencent/HttpFC/config"

	"github.com/gin-gonic/gin"
)

func MakeMultipleHeaders(c *gin.Context, config config.Config) map[string][]string {
	if config.CustomHeaders == nil {
		return c.Request.Header
		//比较应付 因为我根本不知道这是啥
	}
	return config.CustomHeaders
}
func MakeMultipleRequestContext(config config.Config) map[string]string {
	return config.CustomrequestContext
}
