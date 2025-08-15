package breadcrumbs

type Breadcrumb interface {
	GetTitle() string
	Action() error
}
