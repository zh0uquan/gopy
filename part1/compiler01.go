package main

import (
        "fmt"
        "bufio"
        "os"
        "errors"
        "strconv"
)

const (
        INTEGER = "INTEGER"
        PLUS = "PLUS"
        EOF = "EOF"
)

type Token struct {
        _type string
        value string
}

func (token Token) String() string{
        return fmt.Sprintf("Token(%v, %v)", token._type, token.value)
}

type Interprter struct {
        pos int
        text string
        currentToken Token
}

func IsDigit(s string) bool {
        _, err := strconv.ParseInt(s, 10, 64)
        return err == nil
}

func (interprter *Interprter) GetNextToken() error {
        text := interprter.text

        if interprter.pos > len(text) - 1 {
                interprter.currentToken = Token{EOF, ""}
                return nil
        }

        current_char := string(text[interprter.pos])

        if IsDigit(current_char) {
                interprter.currentToken = Token{INTEGER, current_char}
                interprter.pos += 1
                return nil
        }

        if current_char == "+" {
                interprter.currentToken = Token{PLUS, current_char}
                interprter.pos += 1
                return nil
        }

        return errors.New("Error parsing input")
}

func (interprter *Interprter) Eat(token_type string) error {
        if interprter.currentToken._type == token_type {
                interprter.GetNextToken()
                return nil
        } else {
                return errors.New("Error parsing input")
        }
}

func Expr(interprter *Interprter) (string, error) {
        err := interprter.GetNextToken()
        if err != nil {
                return "", err
        }
        left, _ := strconv.ParseInt(interprter.currentToken.value, 10, 64)
        if err := interprter.Eat(INTEGER); err != nil {
                return "", err
        }

        // op := interprter.currentToken
        if err := interprter.Eat(PLUS); err != nil {
                return "", err
        }

        right, _ := strconv.ParseInt(interprter.currentToken.value, 10, 64)
        if err := interprter.Eat(INTEGER); err != nil {
                return "", err
        }

        result := left + right
        return fmt.Sprintf("%v", result), nil
}


func main() {
        scanner := bufio.NewScanner(os.Stdin)
        for {
                fmt.Printf("\033[31mcalc>\033[0m")
                scanner.Scan()
                text := scanner.Text()
                if text != "" {
                        interprter := &Interprter{
                                pos: 0,
                                text: text,
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
