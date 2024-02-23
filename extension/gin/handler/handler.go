package handler

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/tlipoca9/errors"
)

type HandlerPhase string

const (
	HandlerPhaseBind HandlerPhase = "bind"
	HandlerPhaseExec HandlerPhase = "exec"
	HandlerPhaseResp HandlerPhase = "resp"
)

type HandlerManagerBuilder struct {
	data HandlerManager
}

func NewHandlerManagerBuilder() *HandlerManagerBuilder {
	return &HandlerManagerBuilder{
		data: HandlerManager{
			errorHandler: func(c *gin.Context, phase HandlerPhase, err error) {
				switch phase {
				case HandlerPhaseBind:
					c.AbortWithError(http.StatusBadRequest, err)
				default:
					c.AbortWithError(http.StatusInternalServerError, err)
				}
			},
		},
	}
}

func (h *HandlerManagerBuilder) ErrorHandler(eh func(c *gin.Context, phase HandlerPhase, err error)) {
	h.data.errorHandler = eh
}

func (h *HandlerManagerBuilder) Build() *HandlerManager {
	return &h.data
}

type HandlerRequestType interface {
	Bind(c *gin.Context) error
}

type HandlerResponseType interface {
	Resp(c *gin.Context) error
}

type HandlerManager struct {
	errorHandler func(c *gin.Context, phase HandlerPhase, err error)
}

func (h *HandlerManager) CallErrorHandler(c *gin.Context, phase HandlerPhase, err error) {
	if h.errorHandler == nil {
		return
	}

	h.errorHandler(c, phase, err)
}

func (h *HandlerManager) Handle(handler any) gin.HandlerFunc {
	assertHandler(handler)

	v := reflect.ValueOf(handler)
	t := v.Type()

	// input returns handler's input arguments
	input := func(c *gin.Context) ([]reflect.Value, error) {
		in := make([]reflect.Value, 0, t.NumIn())
		in = append(in, reflect.ValueOf(c))

		// if handler has request argument then bind it
		if t.NumIn() == 2 {
			reqV := reflect.New(t.In(1).Elem())
			req := reqV.Interface().(HandlerRequestType)

			if err := req.Bind(c); err != nil {
				return nil, errors.Wrap(err, "failed to bind request")
			}
		}

		return in, nil
	}

	return func(c *gin.Context) {
		var (
			resp any
			err  error
		)

		in, err := input(c)
		if err != nil {
			h.CallErrorHandler(c, HandlerPhaseBind, err)
			return
		}

		out := v.Call(in)

		switch len(out) {
		case 0:
			return
		case 1:
			if errV := out[0].Interface(); errV != nil {
				err = errV.(error)
			}
		case 2:
			resp = out[0].Interface()
			if errV := out[1].Interface(); errV != nil {
				err = errV.(error)
			}
		default:
			panic("handler return values count must be 2 or less")
		}

		if err != nil {
			h.CallErrorHandler(c, HandlerPhaseExec, err)
			return
		}

		if resp != nil {
			if r, ok := resp.(HandlerResponseType); ok {
				if err := r.Resp(c); err != nil {
					h.CallErrorHandler(c, HandlerPhaseResp, err)
				}
			}
		}
	}
}

// assertHandler checks if handler is valid
// handler must be a function
// handler's first argument must be *gin.Context
// handler's second argument must be a struct
// handler's last return value must be error
// handler's first return value must be a pointer
// examples:
//   - func(c *gin.Context)
//   - func(c *gin.Context) error
//   - func(c *gin.Context, *req) error
//   - func(c *gin.Context) (*resp, error)
//   - func(c *gin.Context, *req) (*resp, error)
func assertHandler(handler any) {
	v := reflect.ValueOf(handler)
	t := v.Type()

	if t.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	if t.NumIn() == 0 || t.NumIn() > 2 {
		panic("handler must have 1 or 2 arguments")
	}

	if t.In(0) != reflect.TypeOf(&gin.Context{}) {
		panic("handler's first argument must be *gin.Context")
	}

	if t.NumIn() == 2 && !t.In(1).Implements(reflect.TypeOf((*HandlerRequestType)(nil)).Elem()) {
		panic("handler's second argument must be have Bind(*gin.Context) error method")
	}

	if t.NumOut() > 2 {
		panic("handler return values count must be 2 or less")
	}

	if t.NumOut() != 0 && !t.Out(t.NumOut()-1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic("handler's last return value must be error")
	}

	if t.NumOut() == 2 && t.Out(0).Kind() != reflect.Ptr {
		panic("handler's first return value must be a pointer")
	}
}
