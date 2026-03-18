package opt

func Ptr[T any](v T) *T {
	return &v
}

type Optional[V any] struct {
	V       V
	defined bool
}

func IsDefinedAndNotNil[T any](o Optional[*T]) bool {
	return o.defined && o.V != nil
}

func IsDefinedButNil[T any](o Optional[*T]) bool {
	return o.defined && o.V == nil
}

func IsUndefined[T any](o Optional[T]) bool {
	return !o.defined
}

func IsDefined[T any](o Optional[T]) bool {
	return o.defined
}

func Defined[T any](v T) Optional[T] {
	return Optional[T]{
		V:       v,
		defined: true,
	}
}

func Undefined[T any]() Optional[T] {
	return Optional[T]{
		defined: false,
	}
}
