package parser

func (p *Parser) parseText() *TextNode {
	return &TextNode{
		Value:    p.curToken.Value,
		Position: p.curToken.Pos,
	}
}
