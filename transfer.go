package fchandler2go

import (
	"github.com/aliyun/fc-runtime-go-sdk/events"
	"github.com/danvei233/fchandler2go/core"
	"github.com/gin-gonic/gin"
)

func T(h *interface{}) gin.HandlerFunc {
	// first check if it is ok
	result := core.Check(h)
	IsValid, reason := core.CheckVaild(result)
	if IsValid == false {
		goto rtn
	}
	// then map the data
rtn:
	return func(c *gin.Context) {
		// if there is a dog call back that I failed
		if IsValid != true {
			c.String(400, reason.Error())
			return
		}

		// map the context
		// if we don't need we can skip
		var request *events.HTTPTriggerEvent
		var err error
		if result.In.HasInput == true {
			request, err = core.ConvertRequest(c)
			if err != nil {
				c.String(500, err.Error())
			}

		}

		// use handler
		res, err := result.Invoke(c, *request)
		response, ok := res.(events.HTTPTriggerResponse)
		if !ok {

			if err != nil {
				c.JSON(500, gin.H{"msg": err.Error()})
			}
			//c.JSON(200, gin.H{"msg": "No response"})
			// 如果您的函数返回有效的JSON但是没有包含statusCode字段，或者返回的不是有效的JSON，函数计算会做出以下假设，构造响应结构体。
			core.Recall(c, response)
		}
		//post out data

		core.Recall(c, response)

		return

	}

}

//Todo
// err X-Fc-Error-Type
// gin
// 细节对照
