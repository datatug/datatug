package breadcrumbs

type breadcrumb struct {
	title string
}

func (b *breadcrumb) GetTitle() string {
	return b.title
}

func NewBreadcrumb(title string) Breadcrumb {
	return &breadcrumb{title: title}
}
