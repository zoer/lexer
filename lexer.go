// Lexer package provides an ability to tokenize text.
package lexer

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	cantMatchErrorMessage = `Can't match any existed matchers for the following text: %q`
)

// Lexer contains the input text and token matchers.
type Lexer struct {
	Input        string         // string being scanned
	Matchers     []TokenMatcher // tokens' matchers
	currentInput []byte         // current working input
	currentToken *Token         // matched token
	Error        error          // error of scanning
}

// Token represents the scanned token info.
type Token struct {
	Name interface{} // token name
	Text []byte      // token body
}

// TokenMatcher represents token's matcher function type.
type TokenMatcher func([]byte) (bool, int, interface{}, []byte)

// NewLexer creates new lexer with given input.
func NewLexer(text string) *Lexer {
	l := &Lexer{Input: text}
	l.Reset()
	return l
}

// NewLexerWithMatchers creates new lexer with given input and matchers.
//
//   text := `text which need to be tokenized`
//   l := NewLexerWithMatchers(text, []TokenMatchers{
//     TokenizeIfMatches(`\d+`, "DIGIT"),
//     SkipIfMatches(`\s+`),
//   })
func NewLexerWithMatchers(text string, matchers []TokenMatcher) *Lexer {
	l := NewLexer(text)
	for _, m := range matchers {
		l.AddMatcher(m)
	}
	return l
}

// NewToken creates new token with given name and body.
func NewToken(name interface{}, text []byte) *Token {
	return &Token{Name: name, Text: text}
}

// AddMatcher adds new matter to end of the matchers list.
//
//   l := NewLexer(`some text`)
//   l.AddMatcher(TokenizeIfMatches(`\d+`, "DIGIT"))
//   l.AddMatcher(SkipIfMatches(`\s+`))
func (l *Lexer) AddMatcher(fn TokenMatcher) {
	l.Matchers = append(l.Matchers, fn)
}

// Scan scans for a new token. It returns false if can't find any new token.
func (l *Lexer) Scan() bool {
	var matched bool
	var tokenName interface{}
	var tokenText []byte
	var shift int

	l.currentToken = nil
F:
	for _, fn := range l.Matchers {
		matched, shift, tokenName, tokenText = fn(l.currentInput)
		if shift > 0 {
			l.currentInput = l.currentInput[shift:]
		}
		if matched || shift > 0 {
			break F
		}
	}

	if matched {
		l.currentToken = NewToken(tokenName, tokenText)
		return true
	} else if shift > 0 {
		return l.Scan()
	} else {
		if len(l.currentInput) > 0 {
			l.Error = errors.New(fmt.Sprintf(cantMatchErrorMessage, string(l.currentInput)))
		}
		return false
	}
}

// Token returns current mached token.
func (l *Lexer) Token() *Token {
	return l.currentToken
}

// normalizePattern normalize regexp patterns.
func normalizePattern(pattern string) string {
	if res, err := regexp.MatchString(`^\^`, pattern); err == nil && res == false {
		pattern = "^" + pattern
	}
	return pattern
}

// SkipIfMatches skips the matches without creating a token.
// It's useful to skip space and any other charaters which don't need to
// be tokinized.
func SkipIfMatches(pattern string) TokenMatcher {
	return func(input []byte) (matched bool, shift int, name interface{}, text []byte) {
		re := regexp.MustCompile(normalizePattern(pattern))
		match := re.Find(input)
		if match == nil {
			return
		}
		return false, len(match), nil, nil
	}
}

// TokenizeIfMatches creates token with given name if pattern matches.
// Special character '^' will be insert in the beggining of pattern if it's
// missed.
//
// Usage examples:
//   TokenizeIfMatches(`\d+`, "DIGIT")
//   TokenizeIfMatches(`(?i)if`, "IF") // case-insensitive
//
//   // you can use any defined constant as token name
//   TokenizeIfMatches(`\d+`, DIGIT)
//
func TokenizeIfMatches(pattern string, tokenName interface{}) TokenMatcher {
	return func(input []byte) (matched bool, shift int, name interface{}, text []byte) {
		re := regexp.MustCompile(normalizePattern(pattern))
		match := re.Find(input)
		if match == nil {
			return
		}
		return true, len(match), tokenName, match
	}
}

// Reset resets the current scan results.
func (l *Lexer) Reset() {
	l.Error = nil
	l.currentInput = []byte(l.Input)
	l.currentToken = nil
}
