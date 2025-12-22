package core

type Merger[DIFF any] interface {
	Merge([][]Operation[DIFF]) []Operation[DIFF]
}
