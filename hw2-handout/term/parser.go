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
	heads     []*Term
	all_atoms []*Term
}

var m map[string]*Term

func NewParser() Parser {
	inp := []*Term{}
	x := Node{inp, []*Term{}}
	m = make(map[string]*Term)
	return x
}

func (x Node) Parse(a string) (*Term, error) {
	if a == "" {
		return nil, nil
	}

	if IsValidParser(a) == false {
		return nil, fmt.Errorf("Invalid")
	}

	lex := newLexer(a)

	tok, _ := lex.next()

	if tok.typ == tokenNumber {

		return &Term{1, tok.literal, nil, nil}, nil

	}

	if tok.typ == tokenVariable {

		return &Term{2, tok.literal, nil, nil}, nil

	}

	var save *Term

	if tok.typ == tokenAtom {

		currentLiteral := tok.literal

		tok, _ = lex.next()
		if tok.typ == tokenEOF {

			return &Term{0, currentLiteral, nil, nil}, nil

		}

		m[currentLiteral] = &Term{0, currentLiteral, nil, nil}

		x.heads = append(x.heads, &Term{3, "", m[currentLiteral], []*Term{}})

		save = x.heads[0]

		for len(x.heads) > 0 {

			if tok.typ == tokenNumber {

				_, ok := m[tok.literal]
				if ok == false {
					m[tok.literal] = &Term{1, tok.literal, nil, nil}
				}

				x.heads[len(x.heads)-1].Args = append(x.heads[len(x.heads)-1].Args, m[tok.literal])

			}

			if tok.typ == tokenVariable {
				_, ok := m[tok.literal]
				if ok == false {
					m[tok.literal] = &Term{2, tok.literal, nil, nil}
				}

				x.heads[len(x.heads)-1].Args = append(x.heads[len(x.heads)-1].Args, m[tok.literal])

			}

			if tok.typ == tokenAtom {
				currentLiteral = tok.literal
				tok, _ := lex.next()
				if tok.typ == tokenComma || tok.typ == tokenRpar {
					adress, ok := m[currentLiteral]
					if ok == false {
						m[currentLiteral] = &Term{0, currentLiteral, nil, nil}
						adress = m[currentLiteral]

					}

					x.heads[len(x.heads)-1].Args = append(x.heads[len(x.heads)-1].Args, adress)

				}
				if tok.typ == tokenLpar {

					adress, ok := m[currentLiteral]
					if ok == false {
						m[currentLiteral] = &Term{0, currentLiteral, nil, nil}
						adress = m[currentLiteral]
					}

					x.heads = append(x.heads, &Term{3, "", adress, []*Term{}})

				}

			}

			if tok.typ == tokenRpar {
				add_to := x.heads[len(x.heads)-1]

				if len(x.heads) > 1 {
					for i := 0; i < len(x.all_atoms); i++ {
						if x.all_atoms[i].Typ == add_to.Typ && x.all_atoms[i].Literal == add_to.Literal && x.all_atoms[i].Functor == add_to.Functor && len(x.all_atoms[i].Args) == len(add_to.Args) {
							found := 1
							for j := 0; j < len(x.all_atoms[i].Args); j++ {
								if x.all_atoms[i].Args[j] != add_to.Args[j] {
									found = 0
									j = len(x.all_atoms[i].Args)
								}
							}
							if found == 1 {
								add_to = x.all_atoms[i]
								i = len(x.all_atoms)
							}
						}
					}
					x.all_atoms = append(x.all_atoms, add_to)
					x.heads[len(x.heads)-2].Args = append(x.heads[len(x.heads)-2].Args, add_to)
					x.heads = x.heads[0 : len(x.heads)-1]

				} else {
					return save, nil
				}

			}
			tok, _ = lex.next()
		}

	}

	return save, nil

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
