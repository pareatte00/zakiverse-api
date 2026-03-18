package sqlxtool

import (
	"maps"

	opt "github.com/zakiverse/zakiverse-api/util/optional"
)

type Builder struct {
	allowed map[string]Spec
}

func New(allowed map[string]Spec) *Builder {
	m := make(map[string]Spec, len(allowed))
	maps.Copy(m, allowed)
	return &Builder{allowed: m}
}

// Clean out only what external allowed.
func (b *Builder) Clean(updates map[string]any) map[string]any {
	filtered := make(map[string]any)

	for k, spec := range b.allowed {
		if !spec.external {
			continue
		}

		v := spec.extractor(updates, k)
		if opt.IsDefined(v) {
			filtered[k] = v.V
		}
	}

	return filtered
}

// Clean out only what external allowed and is content only.
func (b *Builder) CleanContent(updates map[string]any) map[string]any {
	filtered := make(map[string]any)

	for k, spec := range b.allowed {
		if !spec.external || !spec.content {
			continue
		}

		v := spec.extractor(updates, k)
		if opt.IsDefined(v) {
			filtered[k] = v.V
		}
	}

	return filtered
}

// Clean out only what internal allowed for create.
func (b *Builder) BuildCreate(updates map[string]any) (cols []string, vals []string, args map[string]any) {
	cols = []string{}
	vals = []string{}
	args = map[string]any{}

	for k, spec := range b.allowed {
		if !spec.internal {
			continue
		}

		v := spec.extractor(updates, k)
		if !opt.IsDefined(v) {
			continue
		}

		cols = append(cols, k)
		vals = append(vals, ":"+k)
		args[k] = v.V
	}

	return
}

// Clean out only what internal allowed for update.
func (b *Builder) BuildUpdate(updates map[string]any) (set []string, args map[string]any) {
	set = []string{}
	args = map[string]any{}

	for k, spec := range b.allowed {
		if !spec.internal {
			continue
		}

		v := spec.extractor(updates, k)
		if !opt.IsDefined(v) {
			continue
		}

		set = append(set, k+" = :"+k)
		args[k] = v.V
	}

	return
}
