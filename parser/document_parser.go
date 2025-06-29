package parser

import "github.com/rickykimani/gotex/lexer"

func (p *Parser) ParseDocument() *Document {
	nodes := []Node{}
	pos := p.curToken.Pos

	for p.curToken.Type != lexer.TokenEOF {
		node := p.parseNode()
		if node != nil {
			// Attempt to merge with the previous node if both are text nodes
			if len(nodes) > 0 {
				lastNode := nodes[len(nodes)-1]
				if currentText, ok := node.(*TextNode); ok {
					if lastText, ok := lastNode.(*TextNode); ok {
						// Merge consecutive text nodes
						lastText.Value += currentText.Value
					} else {
						nodes = append(nodes, node)
					}
				} else {
					nodes = append(nodes, node)
				}
			} else {
				nodes = append(nodes, node)
			}
		}
		p.nextToken()
	}

	return NewDocument(nodes, pos)
}
