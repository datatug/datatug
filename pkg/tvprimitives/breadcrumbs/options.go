package breadcrumbs

type option func(bc *Breadcrumbs) *Breadcrumbs

func WithSeparator(separator string) func(bc *Breadcrumbs) *Breadcrumbs {
	return func(bc *Breadcrumbs) *Breadcrumbs {
		bc.separator = separator
		return bc
	}
}
