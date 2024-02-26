package ginleaf

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/tlipoca9/errors"
)

var defaultStructsBinding binding.Binding = NewStructsBindingBuilder().Build()

type StructsBindingBuilder struct {
	data StructsBinding
}

func NewStructsBindingBuilder() *StructsBindingBuilder {
	return &StructsBindingBuilder{
		data: StructsBinding{
			binding: ComposeBinding(
				NewDefaultBindingBuilder().Build(),
				binding.Header,
				binding.Query,
				binding.JSON,
			),
			validator: validator.New(),
		},
	}
}

func (s *StructsBindingBuilder) Validator(validator *validator.Validate) {
	s.data.validator = validator
}

func (s *StructsBindingBuilder) Binding(binding binding.Binding) {
	s.data.binding = binding
}

func (s *StructsBindingBuilder) Build() *StructsBinding {
	return &s.data
}

type StructsBinding struct {
	binding   binding.Binding
	validator *validator.Validate
}

func (s *StructsBinding) Name() string {
	return "structs"
}

func (s *StructsBinding) Bind(request *http.Request, obj any) error {
	if err := s.binding.Bind(request, obj); err != nil {
		return errors.Wrapf(err, "structs binding '%s' bind failed", s.binding.Name())
	}

	if err := s.validator.Struct(obj); err != nil {
		return errors.Wrapf(err, "structs binding '%s' validate failed", s.binding.Name())
	}

	return nil
}
