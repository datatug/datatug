package breadcrumbs

type breadcrumb struct {
	title  string
	action func() error
}

func (b *breadcrumb) GetTitle() string {
	return b.title
}

func (b *breadcrumb) Action() error {
	if b.action == nil {
		return nil
	}
	return b.action()
}

func NewBreadcrumb(title string, action func() error) Breadcrumb {
	return &breadcrumb{title: title, action: action}
}
