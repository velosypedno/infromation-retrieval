package core

type Normalizer interface {
	Normalize(tokens []string) []string
}

type Tokenizer interface {
	Tokenize(text string) []string
}

type DocReader interface {
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

type Producer interface {
	Produce() error
}

type Creator[T any] interface {
	Create(T) error
}

type Indexer interface {
	Index() (*[]TermToDocIds, error)
}

type Checker[T any] interface {
	Check(T) bool
}

type Reader[T any] interface {
	Read(string) (T, error)
}

type Writer[T any] interface {
	Write(T) error
}

type Merger interface {
	Merge(string, string) (string, error)
}

type Remover[T any] interface {
	Remove(T) error
}

type Worker interface {
	Work() error
}

type PipeLineStage []interface {
	Run() error
}
