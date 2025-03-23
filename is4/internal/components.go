package internal

type Normalizer interface {
	Normalize(tokens []string) []string
}

type Tokenizer interface {
	Tokenize(text string) []string
}

type FileReader interface {
	Read(path string) (<-chan string, error)
}

type Validator[T any] interface {
	Validate(T) error
}

type Provider[T any, R any] interface {
	Provide(T) (R, error)
}

type Supplier[R any] interface {
	Supply() (R, error)
}

type Filter[T any, T2 any] interface {
	Filter(*[]T, T2) *[]T
}

type Mapper[T comparable, R any] interface {
	Map(T) R
}

type Writer[T any] interface {
	Write(*T) error
}

type Reader[T any] interface {
	Read() (*T, error)
}

type Searcher[T any, R any] interface {
	Search(T) (R, error)
}

type Indexer[T any] interface {
	Index(T) error
}

type TermDoc struct {
	Term string
	Doc  int
}

type DocsIndex []struct {
	Term string
	Docs []int
}

type DocIds map[int]string
