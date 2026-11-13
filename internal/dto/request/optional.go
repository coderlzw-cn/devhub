package request

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

// OptionalParam 保留 query 参数是否传入，区分未传与 false、0。
type OptionalParam[T any] struct {
	Value T
	IsSet bool
}

func (p OptionalParam[T]) Schema(r huma.Registry) *huma.Schema {
	return huma.SchemaFromType(r, reflect.TypeFor[T]())
}

func (p *OptionalParam[T]) Receiver() reflect.Value {
	return reflect.ValueOf(&p.Value).Elem()
}

func (p *OptionalParam[T]) OnParamSet(isSet bool, _ any) { p.IsSet = isSet }

func (p OptionalParam[T]) Pointer() *T {
	if !p.IsSet {
		return nil
	}
	value := p.Value
	return &value
}
