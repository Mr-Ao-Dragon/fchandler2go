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
	fn     interface{}
	t      reflect.Type
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

func NewReflector(i interface{}) *F {
	t := reflect.TypeOf(i)

	isFunc := t.Kind() == reflect.Func
	if !isFunc {
		return &F{IsFunc: false}
	}

	numIn := t.NumIn()
	numOut := t.NumOut()

	hasCtx := numIn >= 1 && t.In(0) == reflect.TypeOf((*context.Context)(nil)).Elem()

	hasInput := numIn >= 1

	errType := reflect.TypeOf((*error)(nil)).Elem()
	hasErr := numOut >= 1 && t.Out(numOut-1).Implements(errType)
	hasOut := numOut > 1

	return &F{
		fn:     i,
		t:      t,
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
func CheckValid(res *F) (bool, error) {
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
	if res.In.Num == 2 {

		if res.In.HasCtx == false || res.In.HasInput == false {
			return false, errors.New("input error need (ctx , input) with correct order")
		}

		ev1 := reflect.ValueOf(events.HTTPTriggerEvent{}).Type()
		ev2 := reflect.ValueOf([]byte(nil)).Type()

		target := res.t.In(1)
		if !ev1.ConvertibleTo(target) && !ev2.ConvertibleTo(target) {
			return false, errors.New("input error: cannot match")
		}
	}
	if res.Out.Num == 2 {
		if res.Out.HasOut == false || res.Out.HasErr == false {
			return false, errors.New("output error, need (output , err) with correct order")
		}
		ev := reflect.ValueOf(events.HTTPTriggerResponse{}).Type()
		ev2 := reflect.ValueOf(&events.HTTPTriggerResponse{}).Type()

		target := res.t.Out(0)
		if !ev.ConvertibleTo(target) && !ev2.ConvertibleTo(target) {
			return false, errors.New("output error: cannot match")
		}

	}
	if res.In.Num == 1 {

		ev1 := reflect.ValueOf(events.HTTPTriggerEvent{}).Type()
		ev2 := reflect.ValueOf([]byte(nil)).Type()

		target := res.t.In(0)
		if !ev1.ConvertibleTo(target) && !ev2.ConvertibleTo(target) {
			return false, errors.New("input error: cannot match")
		}
	}
	if res.Out.Num == 1 {

		if res.Out.HasErr == false {
			return false, errors.New("output error, you must return err if you only return one value")
		}
		ev := reflect.ValueOf(events.HTTPTriggerResponse{}).Type()
		ev2 := reflect.ValueOf(&events.HTTPTriggerResponse{}).Type()
		target := res.t.Out(0)
		if !ev.ConvertibleTo(target) && !ev2.ConvertibleTo(target) {
			return false, errors.New("output error: cannot match")
		}

	}
	res.safe = true
	return true, nil
}
func (res F) Invoke(ctx context.Context, event *events.HTTPTriggerEvent) (interface{}, error) {
	// 不做检查
	args := []reflect.Value{}
	if res.safe == false {
		return nil, errors.New("safe check failed or not completed")
	}

	switch res.In.Num {
	case 2:
		args = append(args, reflect.ValueOf(ctx))
		args = append(args, reflect.ValueOf(event).Elem().Convert(res.t.In(1)))
		break

	case 1:
		args = append(args, reflect.ValueOf(event).Elem().Convert(res.t.In(0)))
		break
	}
	results := reflect.ValueOf(res.fn).Call(args)
	if res.Out.Num != len(results) {
		return nil, errors.New("Internal Server Error")
	}
	switch res.Out.Num {
	case 2:
		var o interface{}
		target := reflect.ValueOf(events.HTTPTriggerResponse{}).Type()
		if reflect.ValueOf(results[0].Interface()).Kind() == reflect.Ptr {
			o = reflect.ValueOf(results[0].Interface()).Elem().Convert(target).Interface()

		} else {
			o = reflect.ValueOf(results[0].Interface()).Convert(target).Interface()
		}
		if results[1].Interface() == nil {
			return o, nil
		}
		return o, results[1].Interface().(error)

	case 1:
		return results[0].Interface(), nil

	case 0:
		return nil, nil
	}
	return nil, errors.New("unknown error")
}
