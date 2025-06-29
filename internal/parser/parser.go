package parser

type NodeParser interface {
	parse(rawLink string) error
}
