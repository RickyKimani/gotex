package parser

func (p *Parser) parseComment() *CommentNode {
	return &CommentNode{
		Value:    p.curToken.Value,
		Position: p.curToken.Pos,
	}
}