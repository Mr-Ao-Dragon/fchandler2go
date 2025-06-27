package AliyunWebFC

//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"github.com/aliyun/fc-runtime-go-sdk/events"
//
//	"reflect"
//)
//
//type F struct {
//	IsFunc bool
//	In     in
//	Out    out
//	safe   bool
//	fn     interface{}
//	t      reflect.Type
//}
//
//type in struct {
//	Ctx       reflect.Type
//	Input     reflect.Type
//	inNum     int
//	inputType reflect.Type
//}
//
//type out struct {
//	OutPut     reflect.Type
//	Err        reflect.Type
//	outNum     int
//	outputType reflect.Type
//}
//
//func IsAllowed(list []reflect.Type, i interface{}) (bool, reflect.Type) {
//	if t, ok := i.(reflect.Type); ok {
//		for _, v := range list {
//			if t.ConvertibleTo(v) {
//				return true, v
//			}
//		}
//		return false, nil
//	}
//	for _, v := range list {
//		if v.ConvertibleTo(reflect.TypeOf(i)) {
//			return true, v
//		}
//	}
//	return false, nil
//}
//func NewReflector(i interface{}) *F {
//	t := reflect.TypeOf(i)
//	isFunc := t.Kind() == reflect.Func
//	if !isFunc {
//		return &F{IsFunc: false}
//	}
//	numIn := t.NumIn()
//	numOut := t.NumOut()
//	var ctx reflect.Type
//	var input reflect.Type
//	var output reflect.Type
//	var err reflect.Type
//
//	hasCtx := numIn >= 1 && t.In(0) == reflect.TypeOf((*context.Context)(nil)).Elem()
//	if hasCtx {
//		ctx = t.In(0)
//	}
//	if numIn >= 1 {
//		switch numIn {
//		case 2:
//			if !hasCtx {
//
//				break
//			}
//			input = t.In(1)
//			break
//		case 1:
//			if hasCtx {
//				break
//			}
//			input = t.In(0)
//			break
//		}
//	}
//	if errType := reflect.TypeOf((*error)(nil)).Elem(); numOut >= 1 && t.Out(numOut-1).Implements(errType) {
//		err = t.Out(numOut - 1)
//	} else if numOut >= 1 {
//
//	}
//
//	if numOut == 2 {
//		output = t.Out(0)
//	}
//
//	return &F{
//		fn:     i,
//		t:      t,
//		IsFunc: isFunc,
//		In: in{
//			inNum: numIn,
//			Ctx:   ctx,
//			Input: input,
//		},
//		Out: out{
//			outNum: numOut,
//			OutPut: output,
//			Err:    err,
//		},
//	}
//}
//
//func (res *F) CheckValid() (bool, error) {
//	if res.IsFunc == false {
//		return false, errors.New("not a function")
//	}
//	if res.In.inNum > 2 {
//		return false, errors.New("input num too much")
//	}
//	if res.Out.outNum > 2 {
//		return false, errors.New("output num too much")
//	}
//	if res.In.Input != nil {
//		Model := GetInputAllowedList()
//		var allowed bool
//		if allowed, res.In.inputType = IsAllowed(Model, res.In.Input); !allowed {
//			return false, errors.New("input type not allowed")
//		}
//
//	}
//	if res.Out.OutPut != nil {
//		Model := GetOutputAllowedList()
//		var allowed bool
//		if allowed, res.Out.outputType = IsAllowed(Model, res.Out.OutPut); !allowed {
//			return false, errors.New("output type not allowed")
//		}
//	}
//
//	if res.In.Ctx == nil && res.In.inNum == 2 {
//		return false, errors.New("bad input, ctx must be first")
//	}
//	if res.Out.Err == nil && res.Out.outNum != 0 {
//		return false, errors.New(" the last output must be error")
//	}
//
//	res.safe = true
//	return true, nil
//}
//
//func (res *F) Invoke(ctx context.Context, event *events.HTTPTriggerEvent) (interface{}, error) {
//	// 严格检查
//
//	if res.safe == false {
//		return nil, errors.New("safe check failed or not completed,please .CheckValid() first")
//	}
//	//构造In args
//	args := []reflect.Value{}
//	if res.In.Ctx != nil {
//		args = append(args, reflect.ValueOf(ctx))
//	}
//	if res.In.Input != nil {
//		switch res.In.inputType {
//		case reflect.TypeOf(events.HTTPTriggerEvent{}):
//			args = append(args, reflect.ValueOf(event).Elem().Convert(res.In.Input))
//			break
//		case reflect.TypeOf([]byte(nil)):
//			j, err := json.Marshal(event)
//			if err != nil {
//				return nil, errors.New("Json marshal error")
//			}
//			args = append(args, reflect.ValueOf(j))
//			break
//		default:
//			return nil, errors.New("Service Fatal Error crashed bad type")
//		}
//
//	}
//	// 执行
//	results := reflect.ValueOf(res.fn).Call(args)
//	// 解析结果
//	if res.Out.outNum != len(results) {
//		return nil, errors.New("Internal Server Error")
//	}
//	var err error
//	var output interface{}
//	if res.Out.Err != nil && results[res.Out.outNum-1].Interface() != nil {
//		err = results[res.Out.outNum-1].Interface().(error)
//	}
//	if res.Out.OutPut != nil {
//		target := reflect.ValueOf(events.HTTPTriggerResponse{}).Type()
//		output = results[0].Elem().Convert(target).Interface()
//	}
//
//	return output, err
//
//}
