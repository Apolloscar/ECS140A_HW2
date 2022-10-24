package term

import (
	"errors"
	"fmt"
)

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrParser = errors.New("parser error")

//
// <start>    ::= <term> | \epsilon
// <term>     ::= ATOM | NUM | VAR | <compound>
// <compound> ::= <functor> LPAR <args> RPAR
// <functor>  ::= ATOM
// <args>     ::= <term> | <term> COMMA <args>
//

// Parser is the interface for the term parser.
// Do not change the definition of this interface.
type Parser interface {
	Parse(string) (*Term, error)
}

// NewParser creates a struct of a type that satisfies the Parser interface.
type Node struct {
	// Typ is the type of this term
	x int
}

var m map[int]bool

func NewParser() Parser {

	x := Node{0}
	return x
}

func (x Node) Parse(a string) (*Term, error) {
	if a == "" {
		return nil, nil
	}
	w := Term{0, "0", nil, nil}

	if IsValidParser(a) == false {
		return nil, fmt.Errorf("Invalid")
	}

	return &w, nil
}

func IsValidParser(given string) bool {

	eof := -1
	term := 1
	compund := 2
	functor := 3
	args := 4

	atom := 5
	num := 6
	vari := 7
	lp := 8
	rp := 9
	com := 10

	var stack []int
	stack = append(stack, eof)
	stack = append(stack, term)

	lex := newLexer(given)
	tok, err := lex.next()
	if err != nil {
		return false
	}

	for stack[len(stack)-1] != eof {

		if stack[len(stack)-1] == atom {
			if tok.typ != tokenAtom {
				return false
			}
			tok, err = lex.next()
			if err != nil {
				return false
			}

			stack = stack[0 : len(stack)-1]

		}
		if stack[len(stack)-1] == num {
			if tok.typ != tokenNumber {
				return false
			}
			tok, err = lex.next()
			if err != nil {
				return false
			}

			stack = stack[0 : len(stack)-1]

		}
		if stack[len(stack)-1] == vari {
			if tok.typ != tokenVariable {
				return false
			}
			tok, err = lex.next()
			if err != nil {
				return false
			}

			stack = stack[0 : len(stack)-1]

		}
		if stack[len(stack)-1] == lp {
			if tok.typ != tokenLpar {
				return false
			}
			tok, err = lex.next()
			if err != nil {
				return false
			}

			stack = stack[0 : len(stack)-1]

		}
		if stack[len(stack)-1] == rp {
			if tok.typ != tokenRpar {
				return false
			}
			tok, err = lex.next()
			if err != nil {
				return false
			}

			stack = stack[0 : len(stack)-1]

		}
		if stack[len(stack)-1] == com {
			if tok.typ != tokenComma {
				return false
			}
			tok, err = lex.next()
			if err != nil {
				return false
			}

			stack = stack[0 : len(stack)-1]

		}

		if stack[len(stack)-1] == term {

			if tok.typ != tokenAtom && tok.typ != tokenNumber && tok.typ != tokenVariable {
				return false
			}
			if tok.typ == tokenAtom {
				stack = stack[0 : len(stack)-1]
				stack = append(stack, compund, atom)
			}
			if tok.typ == tokenNumber {
				stack = stack[0 : len(stack)-1]
				stack = append(stack, num)
			}
			if tok.typ == tokenVariable {
				stack = stack[0 : len(stack)-1]
				stack = append(stack, vari)
			}

		}
		if stack[len(stack)-1] == compund {

			if tok.typ != tokenLpar && tok.typ != tokenRpar && tok.typ != tokenComma && tok.typ != tokenEOF {
				return false
			}
			if tok.typ == tokenLpar {
				stack = stack[0 : len(stack)-1]
				stack = append(stack, rp, args, lp)
			}
			if tok.typ == tokenRpar {
				stack = stack[0 : len(stack)-1]
			}
			if tok.typ == tokenComma {
				stack = stack[0 : len(stack)-1]
			}
			if tok.typ == tokenEOF {
				stack = stack[0 : len(stack)-1]
			}

		}

		if stack[len(stack)-1] == functor {

			if tok.typ != tokenRpar && tok.typ != tokenComma {
				return false
			}

			if tok.typ == tokenRpar {
				stack = stack[0 : len(stack)-1]
			}
			if tok.typ == tokenComma {
				stack = stack[0 : len(stack)-1]
				stack = append(stack, args, com)
			}

		}

		if stack[len(stack)-1] == args {

			if tok.typ != tokenAtom && tok.typ != tokenNumber && tok.typ != tokenVariable {
				return false
			}
			if tok.typ == tokenAtom {
				stack = stack[0 : len(stack)-1]
				stack = append(stack, functor, term)
			}
			if tok.typ == tokenNumber {
				stack = stack[0 : len(stack)-1]
				stack = append(stack, functor, term)

			}
			if tok.typ == tokenVariable {
				stack = stack[0 : len(stack)-1]
				stack = append(stack, functor, term)

			}
		}
	}

	if tok.typ == tokenEOF && stack[len(stack)-1] == eof {
		return true
	}
	return false
}
