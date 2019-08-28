package graphqlgenerator

import (
	"fmt"
	"io"
)

type ModelVar struct {
	Name     string
	Var      []Token
	Arg      []GqlArg
	Lit      string
	Required bool
	List     bool
	Map      bool
}

type GqlArg struct {
	Name     string
	Type     Token
	Lit      string
	Default  string
	Required bool
}

type GqlModel struct {
	Name      string
	Variables []ModelVar
}

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func TokenCheck(token Token) bool {
	if token == STRING || token == FLOAT || token == MAP || token == BOOLEAN || token == INT || token == ID || token == IDENT {
		return true
	}
	return false
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) parseArg() (*GqlArg, error) {
	var thisArg GqlArg
	tok1, lit1 := p.scanIgnoreWhitespace()
	if tok1 != IDENT {
		return nil, fmt.Errorf("found %q, expected Identifier err 6", lit1)
	}
	thisArg.Name = lit1
	tok1, lit1 = p.scanIgnoreWhitespace()
	if tok1 != COLON {
		return nil, fmt.Errorf("found %q, expected ':' err 7", lit1)
	}
	tok1, lit1 = p.scanIgnoreWhitespace()
	if TokenCheck(tok1) == true {
		thisArg.Type = tok1
		thisArg.Lit = lit1
	} else {
		return nil, fmt.Errorf("found %q, expected Type err 8", lit1)
	}
	tok1, lit1 = p.scanIgnoreWhitespace()
	if tok1 == EQUAL {
		tok1, lit1 = p.scanIgnoreWhitespace()
		if tok1 != IDENT {
			return nil, fmt.Errorf("found %s, expected Type err 9", tok1)
		}
		thisArg.Default = lit1

		tok1, lit1 = p.scanIgnoreWhitespace()
	}
	if tok1 == EXCLAMATION {
		thisArg.Required = true
	} else {
		p.unscan()
	}

	return &thisArg, nil

}

func (p *Parser) parseInner() (*ModelVar, error) {
	var curvar ModelVar
	tok, lit := p.scanIgnoreWhitespace()
	if tok == EOF {
		return nil, nil
	}
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected Identifier err 14", lit)
	}
	curvar.Name = lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok == BRACKETOPEN {
		for {
			curArg, err := p.parseArg()
			if err != nil {
				return nil, err
			}
			curvar.Arg = append(curvar.Arg, *curArg)

			tok2, lit2 := p.scanIgnoreWhitespace()
			if tok2 == BRACKETCLOSE {
				break
			}
			if tok2 != COMMA {
				return nil, fmt.Errorf("found %q, expected , or ) err 15", lit2)
			}
		}
		tok, lit = p.scanIgnoreWhitespace()
	}

	if tok != COLON {
		return nil, fmt.Errorf("found %q, expected : err16", lit)
	}
	tok1, lit1 := p.scanIgnoreWhitespace()
	if tok1 == SQBRACKETOPEN {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected Identifier err17", lit)
		}
		curvar.Var = append(curvar.Var, tok)
		curvar.Lit = lit
		tok, lit = p.scanIgnoreWhitespace()
		if tok != SQBRACKETCLOSE {
			return nil, fmt.Errorf("found %q, expected ] err8", lit)
		}
		curvar.List = true
	} else if tok1 == MAP {
		curvar.Var = append(curvar.Var, tok1)
		curvar.Map = true
		tok, lit = p.scanIgnoreWhitespace()
		if tok != SQBRACKETOPEN {
			return nil, fmt.Errorf("found %q, expected [ err9", lit)
		}
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected Identifier err10", lit)
		}
		curvar.Var = append(curvar.Var, tok)
		tok, lit = p.scanIgnoreWhitespace()
		if tok != SQBRACKETCLOSE {
			return nil, fmt.Errorf("found %q, expected ] err11", lit)
		}
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected Identifier err12", lit)
		}
		curvar.Var = append(curvar.Var, tok)
	} else if TokenCheck(tok1) == true {
		curvar.Var = append(curvar.Var, tok1)
		curvar.Lit = lit1
	} else {
		return nil, fmt.Errorf("found %q, expected member variable declaration err13", lit1)
	}
	tok1, lit1 = p.scanIgnoreWhitespace()
	if tok1 == EXCLAMATION {
		curvar.Required = true
	} else {
		p.unscan()
	}
	return &curvar, nil
}

func (p *Parser) Parse() (*GqlModel, error) {
	gqlmodel := &GqlModel{}

	tok, lit := p.scanIgnoreWhitespace()
	if tok == SCHEMA{

	}
	if tok != TYPE {
		if lit == "" {
			return nil, fmt.Errorf("EOF reached")
		}else {
		return nil, fmt.Errorf("found %q, expected Type or schema, err1", lit)
		}
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected Identifier, err2", lit)
	} else {
		gqlmodel.Name = lit
	}

	if tok, lit = p.scanIgnoreWhitespace(); tok != CURLBRACKETOPEN {
		return nil, fmt.Errorf("found %q, expected open bracket err3", lit)
	}

	for {
		if tok, _ = p.scanIgnoreWhitespace(); tok == CURLBRACKETCLOSE {
			break
		} else {
			p.unscan()
		}
		mdlvar, err := p.parseInner()
		if err != nil {
			return nil, err
		}
		if mdlvar == nil {
			break
		}
		gqlmodel.Variables = append(gqlmodel.Variables, *mdlvar)
	}
	return gqlmodel, nil
}
func (p *Parser) ParsePackage()	(string, error){
	tok, lit :=  p.scanIgnoreWhitespace()
	if tok != PACKAGE{
		return "", fmt.Errorf("found %q, expected package err 18", lit)
	}
	tok ,lit = p.scanIgnoreWhitespace()
	if lit == ""{
		return "", fmt.Errorf("missing package name")
	}
	return lit, nil
}