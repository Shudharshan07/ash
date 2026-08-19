package parser

type Parser struct {
}

func (p *Parser) Parse(cmd []rune) string {
	return string(cmd)
}
