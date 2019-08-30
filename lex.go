package graphqlgenerator

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // main

	// Misc characters
	ASTERISK         // *
	COMMA            //,
	EXCLAMATION      //!
	SQBRACKETOPEN    //[
	SQBRACKETCLOSE   //]
	CURLBRACKETOPEN  // {
	CURLBRACKETCLOSE //}
	BRACKETOPEN      // (
	BRACKETCLOSE     // )
	COLON            //:
	EQUAL            // =
	// Keywords
	TYPE
	STRING
	FLOAT
	BOOLEAN
	INT
	ID
	PACKAGE
)

func TokenToString(tok Token) string {
	switch tok {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case IDENT:
		return "IDENT"
	case ASTERISK:
		return "ASTERISK"
	case COMMA:
		return "COMMA"
	case EXCLAMATION:
		return "EXCLAMATION"
	case SQBRACKETOPEN:
		return "SQBRACKETOPEN"
	case SQBRACKETCLOSE:
		return "SQBRACKETCLOSE"
	case BRACKETOPEN:
		return "BRACKETOPEN"
	case BRACKETCLOSE:
		return "BRACKETCLOSE"
	case COLON:
		return "COLON"
	case EQUAL:
		return "EQUAL"
	case TYPE:
		return "TYPE"
	case STRING:
		return "STRING"
	case FLOAT:
		return "FLOAT"
	case BOOLEAN:
		return "BOOLEAN"
	case INT:
		return "INT"
	case ID:
		return "ID"
	case PACKAGE:
		return "PACKAGE"
	default:
		return "Token not found in function TokenToString"
	}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

var eof = rune(0)

type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) || isDigit(ch) || ch == '"' {
		s.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	case '[':
		return SQBRACKETOPEN, string(ch)
	case ']':
		return SQBRACKETCLOSE, string(ch)
	case '{':
		return CURLBRACKETOPEN, string(ch)
	case '}':
		return CURLBRACKETCLOSE, string(ch)
	case '!':
		return EXCLAMATION, string(ch)
	case ':':
		return COLON, string(ch)
	case '(':
		return BRACKETOPEN, string(ch)
	case ')':
		return BRACKETCLOSE, string(ch)
	case '=':
		return EQUAL, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' && ch != '"' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "TYPE":
		return TYPE, buf.String()
	case "STRING":
		return STRING, buf.String()
	case "FLOAT":
		return FLOAT, buf.String()
	case "BOOLEAN":
		return BOOLEAN, buf.String()
	case "INT":
		return INT, buf.String()
	case "ID":
		return ID, buf.String()
	case "PACKAGE":
		return PACKAGE, buf.String()

	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() } // puts rune back on reader
