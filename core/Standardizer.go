package core

import (
	"context"
	"errors"

	"github.com/danvei233/fchandler2go/domain"
	"github.com/goccy/go-json"

	"reflect"
)

type Standardizer struct {
	RawFn *interface{}
	Fn    Fn
}
type Fn struct {
	isAcceptableFunc bool
	In               In
	Out              Out
}
type In struct {
	Ctx     reflect.Type
	Input   reflect.Type
	numIn   int
	isValid bool
}
type Out struct {
	Output  reflect.Type
	Err     reflect.Type
	numOut  int
	isValid bool
}

func NewStandardizer(rawFn *interface{}) *Standardizer {
	return &Standardizer{RawFn: rawFn}
}
func (s *Standardizer) initInArgs(t reflect.Type) error {
	hasCtx := t.In(0).Kind() == reflect.TypeOf((*context.Context)(nil)).Elem().Kind()
	switch s.Fn.In.numIn {
	case 1:
		if hasCtx {
			s.Fn.In.Ctx = t.In(0)
			break
		}
		s.Fn.In.Input = t.In(0)
		break
	case 2:
		if !hasCtx {
			return errors.New("stand error: ctx must be first")
		}
		s.Fn.In.Ctx = t.In(0)
		s.Fn.In.Input = t.In(1)
	}
	return nil
}
func (s *Standardizer) initOutArgs(t reflect.Type) error {
	hasErr := t.Out(s.Fn.Out.numOut-1).Kind() == reflect.TypeOf((*error)(nil)).Elem().Kind()
	switch s.Fn.Out.numOut {
	case 1:
		if hasErr {
			s.Fn.Out.Err = t.In(0)
			break
		}
		s.Fn.Out.Output = t.In(0)
		break
	case 2:
		if !hasErr {
			return errors.New("stand error: err must be last")
		}
		s.Fn.Out.Err = t.Out(1)
		s.Fn.Out.Output = t.Out(0)
		break
	}
	return nil
}
func (s *Standardizer) init() error {
	t := reflect.TypeOf(*s.RawFn)
	if t.Kind() != reflect.Func {
		return errors.New("stand error: not a function")
	}
	s.Fn.In.numIn = t.NumIn()
	s.Fn.Out.numOut = t.NumOut()
	// reflect arguments
	if s.Fn.In.numIn > 2 {
		return errors.New("stand error: too many input")
	}
	if s.Fn.Out.numOut > 2 {
		return errors.New("stand error: too many output")
	}
	if s.Fn.In.numIn != 0 {
		if err := s.initInArgs(t); err != nil {
			return err
		}

	}
	if s.Fn.Out.numOut != 0 {
		if err := s.initOutArgs(t); err != nil {
			return err
		}
	}
	s.Fn.isAcceptableFunc = true
	return nil

}
func (s *Standardizer) isAllowed(list []reflect.Type, i interface{}) (bool, reflect.Type) {
	if t, ok := i.(reflect.Type); ok {
		for _, v := range list {
			if t.ConvertibleTo(v) {
				return true, v
			}
		}
		return false, nil
	}
	for _, v := range list {
		if v.ConvertibleTo(reflect.TypeOf(i)) {
			return true, v
		}
	}
	return false, nil
}
func (s *Standardizer) CheckInputValid(list []reflect.Type) (bool, reflect.Type) {
	if s.Fn.isAcceptableFunc == false {
		if err := s.init(); err != nil {
			return false, nil
		}

	}
	if s.Fn.In.Input != nil {
		return s.isAllowed(list, s.Fn.In.Input)

	}
	return true, reflect.ValueOf(nil).Type()
}
func (s *Standardizer) CheckOutputValid(list []reflect.Type) (bool, reflect.Type) {
	if s.Fn.isAcceptableFunc == false {
		if err := s.init(); err != nil {
			return false, nil
		}

	}
	if s.Fn.Out.Output != nil {
		return s.isAllowed(list, s.Fn.Out.Output)
	}
	return true, nil
}
func (s *Standardizer) GetStandardHandler() (domain.StandardHandler, error) {

	if s.Fn.isAcceptableFunc == false {
		if err := s.init(); err != nil {
			return nil, err
		}

	}

	return func(ctx context.Context, b []byte) (output []byte, err error) {
		args := []reflect.Value{}
		if s.Fn.In.Ctx != nil {
			args = append(args, reflect.ValueOf(ctx))
		}

		if s.Fn.In.Input != nil {
			target := reflect.New(s.Fn.In.Input)
			if s.Fn.In.Input.Kind() == reflect.TypeOf([]byte(nil)).Kind() {
				args = append(args, reflect.ValueOf(b))
			} else {
				var argVal reflect.Value
				argVal = target.Elem()
				err := json.Unmarshal(b, target.Interface())
				if err != nil {
					return nil, err
				}
				args = append(args, argVal)
			}

		}
		Result := reflect.ValueOf(*s.RawFn).Call(args)

		if s.Fn.Out.Output != nil {
			//ptr?
			if Result[0].Kind() == reflect.Ptr {
				Result[0] = Result[0].Elem()
			}
			//byte[]?
			if Result[0].Kind() == reflect.TypeOf([]byte(nil)).Elem().Kind() {
				output = Result[0].Interface().([]byte)
			}
			//String?
			if Result[0].Kind() == reflect.TypeOf("").Kind() {
				output = []byte(Result[0].Interface().(string))
			} else {
				temp := Result[0].Interface()
				Standout, err := json.Marshal(&temp)
				if err != nil {
					return nil, err
				}
				output = Standout
			}
		}
		if s.Fn.Out.Err != nil {
			if Result[len(Result)-1].Interface() != nil {
				err = Result[len(Result)-1].Interface().(error)

			}

		}
		return

	}, nil

}
