package core

import (
	"context"
	"errors"
	"github.com/aliyun/fc-runtime-go-sdk/events"

	"reflect"
)

type F struct {
	IsFunc bool
	In     in
	Out    out
	safe   bool
	fn     *interface{}
}

type in struct {
	Num      int
	HasCtx   bool
	HasInput bool
}

type out struct {
	Num    int
	HasOut bool
	HasErr bool
}

func Check(i *interface{}) *F {
	t := reflect.TypeOf(i)
	isFunc := t.Kind() == reflect.Func
	if !isFunc {
		return &F{IsFunc: false}
	}

	numIn := t.NumIn()
	numOut := t.NumOut()

	hasCtx := numIn >= 1 && t.In(0) == reflect.TypeOf((*context.Context)(nil)).Elem()

	hasInput := numIn > 2

	errType := reflect.TypeOf((*error)(nil)).Elem()
	hasErr := numOut >= 1 && t.Out(numOut-1).Implements(errType)
	hasOut := numOut > 1

	return &F{
		fn:     i,
		IsFunc: isFunc,
		In: in{
			Num:      numIn,
			HasCtx:   hasCtx,
			HasInput: hasInput,
		},
		Out: out{
			Num:    numOut,
			HasOut: hasOut,
			HasErr: hasErr,
		},
	}
}
func CheckVaild(res *F) (bool, error) {
	if res.IsFunc == false {
		return false, errors.New("not a function")
	}
	//notfunc
	if res.In.Num > 2 {
		return false, errors.New("too many input")
	}
	if res.Out.Num > 2 {
		return false, errors.New("too many output")
	}
	if res.In.Num == 2 && (res.In.HasCtx == false || res.In.HasInput == false) {
		return false, errors.New("input error need (ctx , input) with correct order")
	}
	if res.Out.Num == 2 && (res.Out.HasOut == false || res.Out.HasErr == false) {
		return false, errors.New("output error, need (output , err) with correct order")

	}
	if res.Out.Num == 1 && res.Out.HasErr == false {
		return false, errors.New("output error, you must return err if you only return one value")
	}
	res.safe = true
	return true, nil
}
func (res F) Invoke(ctx context.Context, event events.HTTPTriggerEvent) (interface{}, error) {
	// 不做检查
	args := []reflect.Value{}
	if res.safe == false {
		return nil, errors.New("safe check failed or not completed")
	}
	switch res.In.Num {
	case 2:
		args = append(args, reflect.ValueOf(ctx))
		args = append(args, reflect.ValueOf(event))
		break

	case 1:
		args = append(args, reflect.ValueOf(event))
		break
	}
	results := reflect.ValueOf(res.fn).Call(args)
	switch res.Out.Num {
	case 2:
		return results[0].Interface(), results[1].Interface().(error)

	case 1:
		return results[0].Interface(), nil

	case 0:
		return nil, nil
	}
	return nil, errors.New("unknown error")
}
