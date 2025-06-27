package HttpFC

import (
	"github.com/danvei233/fchandler2go/Tencent/HttpFC/config"
	"github.com/danvei233/fchandler2go/Tencent/HttpFC/mock"
	"github.com/danvei233/fchandler2go/utills"
	"github.com/gin-gonic/gin"
)

func ConvertRequest(c *gin.Context, config config.Config) (*HTTPRequest, error) {
	var Request HTTPRequest

	Request.HTTPMethod = c.Request.Method
	Request.Path = c.Request.URL.Path
	body, err := utills.BodyReader(c.Request.Body, c.ContentType())
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return nil, err
	}
	Request.Body = body
	Request.RequestContext = mock.MakeMultipleRequestContext(config)
	Request.MultiValueHeaders = mock.MakeMultipleHeaders(c, config)
	Request.Headers = *utills.Header2map(c.Request.Header)
	Request.IsBase64Encoded = *utills.IsBin(c.ContentType())
	Request.QueryStringParameters = *utills.Param2map(c.Request.URL.Query())

	return &Request, nil

}
