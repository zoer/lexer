package lexer_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoer/lexer"
)

func TestLexer_NewLexer(t *testing.T) {
	assert := assert.New(t)
	text := `foo`
	l := lexer.NewLexer(text)

	assert.Equal(l.Input, text)
}

func TestLexer_NewLexerWithMatchers(t *testing.T) {
	assert := assert.New(t)
	l := lexer.NewLexerWithMatchers(`foo`, []lexer.TokenMatcher{
		lexer.TokenizeIfMatches(`^foo`, "FOO"),
	})

	assert.Equal(len(l.Matchers), 1)
}

func TestLexer_AddMatcher(t *testing.T) {
	assert := assert.New(t)
	l := lexer.NewLexer(`foo`)
	assert.Equal(len(l.Matchers), 0, "The matchers list should be empty")
	fn := func([]byte) (bool, int, interface{}, []byte) {
		return true, 0, nil, []byte{}
	}
	l.AddMatcher(fn)
	assert.Equal(len(l.Matchers), 1, "Should increment matchers size by 1")
}

func TestLexer_Scan(t *testing.T) {
	assert := assert.New(t)

	l := lexer.NewLexer(`foo fooo  123`)
	l.AddMatcher(func(input []byte) (matched bool, shift int, tokenName interface{}, tokenText []byte) {
		re := regexp.MustCompile(`^fo+`)
		match := re.Find(input)
		if match == nil {
			return
		}
		return true, len(match), "TEST", match
	})
	l.AddMatcher(lexer.TokenizeIfMatches(`^\d+`, "DIGIT"))
	l.AddMatcher(lexer.SkipIfMatches(`\s+`))

	assert.True(l.Scan())
	assert.Equal(l.Token().Name, "TEST")
	assert.Equal(l.Token().Text, []byte("foo"))
	assert.NoError(l.Error)

	assert.True(l.Scan())
	assert.Equal(l.Token().Name, "TEST")
	assert.Equal(l.Token().Text, []byte("fooo"))
	assert.NoError(l.Error)

	assert.True(l.Scan())
	assert.Equal(l.Token().Name, "DIGIT")
	assert.Equal(l.Token().Text, []byte("123"))
	assert.NoError(l.Error)

	assert.False(l.Scan())
	assert.NoError(l.Error)
}

func TestLexer_Scan2(t *testing.T) {
	assert := assert.New(t)
	text := `price $12.4 foo`
	l := lexer.NewLexerWithMatchers(text, []lexer.TokenMatcher{
		lexer.TokenizeIfMatches(`\w+`, "WORD"),
		lexer.SkipIfMatches(`\s+`),
	})
	l.AddMatcher(func(input []byte) (matched bool, shift int, tokenName interface{}, tokenText []byte) {
		re := regexp.MustCompile(`^\$(\d+(?:\.\d+))`)
		match := re.FindSubmatch(input)
		if match == nil {
			return
		}
		return true, len(match[0]), "PRICE", match[1]
	})

	assert.True(l.Scan())
	assert.Equal(l.Token().Name, "WORD")
	assert.Equal(l.Token().Text, []byte("price"))
	assert.NoError(l.Error)

	assert.True(l.Scan())
	assert.Equal(l.Token().Name, "PRICE")
	assert.Equal(l.Token().Text, []byte("12.4"))
	assert.NoError(l.Error)

	assert.True(l.Scan())
	assert.Equal(l.Token().Name, "WORD")
	assert.Equal(l.Token().Text, []byte("foo"))
	assert.NoError(l.Error)

	assert.False(l.Scan())
	assert.Nil(l.Token())
	assert.NoError(l.Error)
}

func TestLexer_ScanWithError(t *testing.T) {
	assert := assert.New(t)

	l := lexer.NewLexer(`foo 123`)
	l.AddMatcher(lexer.TokenizeIfMatches(`^foo`, "WORD"))
	l.AddMatcher(lexer.TokenizeIfMatches(`\d+`, "DIGIT"))

	assert.True(l.Scan())
	assert.Equal(l.Token().Name, "WORD")
	assert.Equal(l.Token().Text, []byte("foo"))
	assert.NoError(l.Error, "Should not have an error")

	assert.False(l.Scan())
	assert.Nil(l.Token(), "Should not have a token")
	assert.Error(l.Error, "Should have an error")

	l.Reset()
	assert.NoError(l.Error, "Error should be reseted")
}

// Simple usage example.
func ExampleNewLexer() {
	text := `price 12`
	l := lexer.NewLexer(text)
	l.AddMatcher(lexer.TokenizeIfMatches(`\w+`, "WORD"))
	l.AddMatcher(lexer.SkipIfMatches(`\s+`))
	l.AddMatcher(lexer.TokenizeIfMatches(`\d+`, "PRICE"))

	for l.Scan() {
		fmt.Printf("%s => %s", l.Token().Name, l.Token().Text)
	}

	// Output:
	// WORD => price
	// PRICE => 12
}

// Using custom matcher to parse the cost and drop the $ sign.
func ExampleNewLexerWithMatchers() {
	text := `price $12.4`
	l := lexer.NewLexerWithMatchers(text, []lexer.TokenMatcher{
		lexer.TokenizeIfMatches(`\w+`, "WORD"),
		lexer.SkipIfMatches(`\s+`),
	})
	l.AddMatcher(func(input []byte) (matched bool, shift int, tokenName interface{}, tokenText []byte) {
		// Don't forget to add special symbol '^'
		re := regexp.MustCompile(`^\$(\d+(?:\.\d+))`)
		match := re.FindSubmatch(input)
		if match == nil {
			return
		}
		return true, len(match[0]), "PRICE", match[1]
	})

	for l.Scan() {
		fmt.Printf("%s => %s\n", l.Token().Name, l.Token().Text)
	}

	// Output:
	// WORD => price
	// PRICE => 12.4
}
