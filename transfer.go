package handler2gin

import (
	"errors"
	"github.com/danvei233/fchandler2go/core"
	"github.com/danvei233/fchandler2go/domain"
	"github.com/gin-gonic/gin"
)

type Reflector struct {
	h domain.HandlerStandardizer
	M domain.MapProvider
}

func NewReflector(M domain.MapProvider) *Reflector {
	return &Reflector{M: M}
}

func (r *Reflector) T(h interface{}) gin.HandlerFunc {
	//load cfg
	r.h = core.NewStandardizer(&h)
	var err error
	t, InType := r.h.CheckOutputValid(r.M.GetOutputAllowedList())
	if !t {
		err = errors.New("match output error")
	}
	t, _ = r.h.CheckInputValid(r.M.GetInputAllowedList())

	if !t {
		err = errors.New("match input error")
	}

	return func(c *gin.Context) {
		// if there is a dog call back that I failed
		if err != nil {
			c.String(500, err.Error())
			return
		}
		var res []byte
		if InType != nil {
			res, err = r.M.TransIn(c)
			if err != nil {
				c.String(500, err.Error())
				return
			}
		} else {
			res = []byte{}

		}
		// if we don't need we can skip
		// use handler
		var req []byte
		fc, err := r.h.GetStandardHandler()
		if err != nil {
			c.String(500, err.Error())
			return
		}
		req, err = fc(c, res)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		err = r.M.TransOut(c, req)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		return

	}

}

//Todo
// err X-Fc-Error-Type
// gin
// 细节对照
