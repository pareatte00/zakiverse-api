package sqlxtool

import opt "github.com/zakiverse/zakiverse-api/util/optional"

type Spec struct {
	extractor func(map[string]any, string) opt.Optional[any]
	internal  bool
	external  bool
	content   bool
}

func NewSpec[T any](extractor func(map[string]any, string) opt.Optional[T]) Spec {
	return Spec{
		extractor: func(m map[string]any, k string) opt.Optional[any] {
			v := extractor(m, k)
			if opt.IsDefined(v) {
				return opt.Defined(any(v.V))
			}

			return opt.Undefined[any]()
		},
	}
}

func (s Spec) External() Spec {
	s.external = true
	return s
}

func (s Spec) Internal() Spec {
	s.internal = true
	return s
}

func (s Spec) Content() Spec {
	s.content = true
	return s
}
