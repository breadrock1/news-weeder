package weeder

type Weeder struct {
	Weeder CombinedIWeeder
}

type CombinedIWeeder interface {
	IWeeder
	Connect
}

type IWeeder interface {
	CreateSchema() error
	Delete(docId string) error
	Append(doc *Document) error
	Search(params *SearchParams) ([]*Document, error)
}

type Connect interface {
	Connect() error
	Close() error
}
