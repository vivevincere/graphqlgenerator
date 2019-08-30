package graphqlgenerator

import (
	"fmt"
	"strings"
	"testing"
)

func Test_ScanWhiteSpace(t *testing.T) {
	testString := `    af`
	reader := strings.NewReader(testString)
	s := NewScanner(reader)
	tok, lit := s.scanWhitespace()
	if tok != WS {
		t.Errorf("scanWhitespace did not return WS token for whitespace")
	}
	if lit != `    ` {
		t.Errorf(`Scan did not return exact literal for " ", returned %s`, lit)
	}

	tok, lit = s.scanWhitespace()
	if tok != WS {
		t.Errorf("scanWhitespace did not return WS token for whitespace")
	}
	if lit != `a` {
		t.Errorf(`Scan did not return expected literal a, returned %s`, lit)
	}
	tok, lit = s.scanWhitespace()
	if tok != WS {
		t.Errorf("scanWhitespace did not return WS token for whitespace")
	}
	if lit != `f` {
		t.Errorf(`Scan did not return expected literal a, returned %sf`, lit)
	}
}

func Test_Scan(t *testing.T) {
	testString := `type Query {
  performance(word: int = "100"!): PerformanceSummary!
}`
	reader := strings.NewReader(testString)
	s := NewScanner(reader)

	check := scanHelper(TYPE, "type", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(WS, " ", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(IDENT, "Query", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(WS, " ", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(CURLBRACKETOPEN, "{", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(WS, "\n  ", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(IDENT, "performance", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(BRACKETOPEN, "(", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(IDENT, "word", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(COLON, ":", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(WS, " ", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(INT, "int", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(WS, " ", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(EQUAL, "=", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(WS, " ", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(IDENT, `"100"`, s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(EXCLAMATION, `!`, s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(BRACKETCLOSE, ")", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(COLON, ":", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(WS, " ", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(IDENT, "PerformanceSummary", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(EXCLAMATION, "!", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(WS, "\n", s)
	if check != nil {
		t.Errorf(check.Error())
	}
	check = scanHelper(CURLBRACKETCLOSE, "}", s)
	if check != nil {
		t.Errorf(check.Error())
	}

}

func scanHelper(tok Token, lit string, s *Scanner) error {
	tok1, lit1 := s.Scan()
	if tok1 != tok {
		stringTok := TokenToString(tok1)
		stringTok1 := TokenToString(tok)
		return fmt.Errorf("Scan did not return expected token %s, returned %s instead", stringTok1, stringTok)
	}
	if lit1 != lit {
		return fmt.Errorf("Scan did not return expected literal, returned %s instead", lit1)
	}
	return nil
}
