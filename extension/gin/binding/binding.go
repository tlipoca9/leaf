package binding

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Binding interface {
	Bind(*gin.Context, any)
}

type ginBinding struct {
	binding binding.Binding
}

func (b *ginBinding) Bind(c *gin.Context, obj any) {
	b.binding.Bind(c.Request, obj)
}

func Wrap(bd binding.Binding) Binding {
	return &ginBinding{
		binding: bd,
	}
}

type ginBindingUri struct {
	binding binding.BindingUri
}

func (b *ginBindingUri) Bind(c *gin.Context, obj any) {
	m := make(map[string][]string)
	for _, v := range c.Params {
		m[v.Key] = []string{v.Value}
	}
	b.binding.BindUri(m, obj)
}

func WrapBindingUri(bd binding.BindingUri) Binding {
	return &ginBindingUri{
		binding: bd,
	}
}

type composedBinding struct {
	bindings []Binding
}

func (b *composedBinding) Bind(c *gin.Context, obj any) {
	for _, bd := range b.bindings {
		bd.Bind(c, obj)
	}
}

func Composed(bindings ...Binding) Binding {
	return &composedBinding{
		bindings: bindings,
	}
}
