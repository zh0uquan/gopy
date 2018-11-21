package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	INTEGER  = "INTEGER"
	PLUS     = "PLUS"
	EOF      = "EOF"
	MINUS    = "MINUS"
	MULTIPLY = "MULTIPLY"
	DIVIDE   = "DIVIDE"
)

var OPERATOR = map[string]string{
	"+": PLUS,
	"-": MINUS,
	"*": MULTIPLY,
	"/": DIVIDE,
}

type Token struct {
	_type string
	value string
}

func (token Token) String() string {
	return fmt.Sprintf("Token(%v, %v)", token._type, token.value)
}

type Interprter struct {
	pos          int
	text         string
        currentToken Token
}

func IsDigit(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func (interprter *Interprter) GetCurrentChar() string {
	if interprter.pos > len(interprter.text) - 1 {
		return ""
	} else {
		return string(interprter.text[interprter.pos])
	}
}

func (interprter *Interprter) GetNextChar() string {
        interprter.pos += 1
        return interprter.GetCurrentChar()
}

func (interprter *Interprter) SkipWhiteSpace() {
        currentChar := interprter.GetCurrentChar()
        for currentChar != "" && currentChar == " " {
                currentChar = interprter.GetNextChar()
        }
}


func (interprter *Interprter) GetNextToken() Token {
        currentChar := interprter.GetCurrentChar()

        for currentChar != "" {
                // skip whitespace
                if currentChar == " " {
                        interprter.SkipWhiteSpace()
                        currentChar = interprter.GetCurrentChar()
                        continue
                }

                if IsDigit(currentChar) {
                        integerChar := currentChar
                        nextChar := interprter.GetNextChar()
                        for IsDigit(nextChar) {
                                integerChar += nextChar
                                nextChar = interprter.GetNextChar()
                        }
                        interprter.pos +=1
                        return Token{INTEGER, integerChar}
                }

                switch currentChar {
        	case    "+",
        		"-",
        		"*",
        		"/":
                        interprter.pos +=1
        		return Token{OPERATOR[currentChar], currentChar}
        	}

        }

        return Token{EOF, ""}
}

func (interprter *Interprter) Eat(tokenType string) error {
	// fmt.Println(interprter.currentToken._type, token_type)
	if interprter.currentToken._type == tokenType {
		interprter.currentToken = interprter.GetNextToken()
		return nil
	} else {
		return errors.New("Error parsing input")
	}
}

func compute(left int64, right int64, operator string) (float64, error) {
	switch operator {
	case PLUS:
		return float64(left + right), nil
	case MINUS:
		return float64(left - right), nil
	case MULTIPLY:
		return float64(left) * float64(right), nil
	case DIVIDE:
		return float64(left) / float64(right), nil
	}
	return 0, errors.New("Operator not supported")
}

func (interprter *Interprter) Term() (int64, error) {
        value, _ := strconv.ParseInt(interprter.currentToken.value, 10, 64)
        if err := interprter.Eat(INTEGER); err != nil {
                return 0, err
        } else {
                return value, nil
        }
}


func Expr(interprter *Interprter) (string, error) {

        interprter.currentToken = interprter.GetNextToken()

        left, err := interprter.Term()
	if err != nil {
                return "", err
        }

	op := interprter.currentToken
	switch op.value {
	case    "+",
		"-",
		"*",
		"/":
		if err := interprter.Eat(op._type); err != nil {
                        fmt.Println(op._type)
			return "", err
		}
	default:
		return "", errors.New("Operator not supported")
	}

        right, err := interprter.Term()
        if err != nil {
                return "", err
        }

	result, err := compute(left, right, op._type)
	if err == nil {
		return fmt.Sprintf("%v", result), nil
	} else {
		return "", err
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("\033[31mcalc>\033[0m")
		scanner.Scan()
		text := scanner.Text()
		if text != "" {
			interprter := &Interprter{
				pos:          0,
				text:         text,
				currentToken: Token{},
			}
			result, err := Expr(interprter)
			if err == nil {
				fmt.Println(result)
			} else {
				fmt.Println(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
