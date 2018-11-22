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

type Lexer struct {
	pos  int
	text string
}

type Interprter struct {
	lexer        Lexer
	currentToken Token
}

func IsDigit(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func (lexer *Lexer) GetCurrentChar() string {
	if lexer.pos > len(lexer.text)-1 {
		return ""
	} else {
		return string(lexer.text[lexer.pos])
	}
}

func (lexer *Lexer) GetNextChar() string {
	lexer.pos += 1
	return lexer.GetCurrentChar()
}

func (lexer *Lexer) SkipWhiteSpace() {
	currentChar := lexer.GetCurrentChar()
	for currentChar != "" && currentChar == " " {
		currentChar = lexer.GetNextChar()
	}
}

func (lexer *Lexer) GetNextToken() Token {
	currentChar := lexer.GetCurrentChar()

	for currentChar != "" {
		// skip whitespace
		if currentChar == " " {
			lexer.SkipWhiteSpace()
			currentChar = lexer.GetCurrentChar()
			continue
		}

		if IsDigit(currentChar) {
			integerChar := currentChar
			nextChar := lexer.GetNextChar()
			for IsDigit(nextChar) {
				integerChar += nextChar
				nextChar = lexer.GetNextChar()
			}
			return Token{INTEGER, integerChar}
		}

		switch currentChar {
		case "+",
			"-",
			"*",
			"/":
			lexer.pos += 1
			return Token{OPERATOR[currentChar], currentChar}
		}

	}
	return Token{EOF, ""}
}

func (interprter *Interprter) Eat(tokenType string) error {
	if interprter.currentToken._type == tokenType {
		interprter.currentToken = interprter.lexer.GetNextToken()
		return nil
	} else {
		return errors.New("Error parsing input")
	}
}

func (interprter *Interprter) Factor() (float64, error) {
	fmt.Println(interprter.currentToken)
	value, _ := strconv.ParseInt(interprter.currentToken.value, 10, 64)
	if err := interprter.Eat(INTEGER); err != nil {
		return 0, err
	} else {
		return float64(value), nil
	}
}

func (interprter *Interprter) Order() (float64, error) {
	result, err := interprter.Factor()
	if err != nil {
		return 0, err
	}

	for isMultDivOperator(interprter.currentToken.value) {
		operator := interprter.currentToken.value
		switch operator {
		case "*",
			"/":
			interprter.Eat(OPERATOR[operator])
			order, err := interprter.Factor()
			if err != nil {
				return 0, err
			}
			if operator == "*" {
				result = result * order
			} else if operator == "/" {
				result = result / order
			}
		}
	}

	return result, nil
}
func isMultDivOperator(char string) bool {
	for _, operator := range "*/" {
		if char == string(operator) {
			return true
		}
	}
	return false
}

func isPlusMinusOperator(char string) bool {
	for _, operator := range "+-" {
		if char == string(operator) {
			return true
		}
	}
	return false
}

func (interprter *Interprter) Expr() (float64, error) {
	// init token
	interprter.currentToken = interprter.lexer.GetNextToken()

	result, err := interprter.Order()
	if err != nil {
		return 0, err
	}

	for isPlusMinusOperator(interprter.currentToken.value) {
		operator := interprter.currentToken.value
		switch operator {
		case "+",
			"-":
			interprter.Eat(OPERATOR[operator])
			order, err := interprter.Order()
			if err != nil {
				return 0, err
			}
			if operator == "+" {
				result = result + order
			} else if operator == "-" {
				result = result - order
			}
		}
	}

	return result, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("\033[31mcalc>\033[0m")
		scanner.Scan()
		text := scanner.Text()
		if text != "" {
			lexer := Lexer{
				pos:  0,
				text: text,
			}
			interprter := &Interprter{
				lexer:        lexer,
				currentToken: Token{},
			}
			result, err := interprter.Expr()
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
