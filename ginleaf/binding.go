package ginleaf

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/tlipoca9/errors"
)

func Bind(c *gin.Context, obj any, bindings ...binding.Binding) error {
	var b binding.Binding
	switch len(bindings) {
	case 0:
		b = defaultStructsBinding
	case 1:
		b = bindings[0]
	default:
		b = ComposeBinding(bindings...)
	}

	if err := c.ShouldBindUri(obj); err != nil {
		return errors.Wrapf(err, "bind uri failed")
	}

	if err := c.ShouldBindWith(obj, b); err != nil {
		return errors.Wrapf(err, "bind failed")
	}

	return nil
}

var _ binding.Binding = (*ComposedBinding)(nil)

type ComposedBinding struct {
	name string

	bindings []binding.Binding
}

func ComposeBinding(bindings ...binding.Binding) binding.Binding {
	if len(bindings) == 0 {
		panic("at least one binding is required")
	}

	var buf bytes.Buffer
	buf.WriteString("compose")
	for _, b := range bindings {
		buf.WriteByte('_')
		buf.WriteString(b.Name())
	}
	return &ComposedBinding{
		name:     buf.String(),
		bindings: bindings,
	}
}

func (c *ComposedBinding) Name() string {
	return c.name
}

func (c *ComposedBinding) Bind(request *http.Request, obj any) error {
	for _, b := range c.bindings {
		if err := b.Bind(request, obj); err != nil {
			return errors.Wrapf(err, "compose binding '%s'-'%s' bind failed", c.Name(), b.Name())
		}
	}
	return nil
}
