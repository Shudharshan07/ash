package parser

import (
	"strings"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(command []rune) ([]string, error) {
	var sb strings.Builder

	var res []string

	for _, ch := range command {
		if ch == ' ' {
			res = append(res, sb.String())
			sb.Reset()
		} else {
			_, err := sb.WriteRune(ch)
			if err != nil {
				return nil, err
			}
		}
	}

	if sb.Len() != 0 {
		res = append(res, sb.String())
	}

	return res, nil
}
