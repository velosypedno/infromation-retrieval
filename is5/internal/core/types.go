package core

type Doc struct {
	Path string
	Id   int
}

type TermToDocIds struct {
	Term   string
	DocIds []int
}

type Pair[T any] struct {
	First  T
	Second T
}
