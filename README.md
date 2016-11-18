# Lexer
[![Build
Status](https://travis-ci.org/zoer/lexer.svg)](https://travis-ci.org/zoer/lexer)
[![Go Report
Card](https://goreportcard.com/badge/github.com/zoer/lexer)](https://goreportcard.com/report/github.com/zoer/lexer)
[![GoDoc](https://godoc.org/github.com/zoer/lexer?status.svg)](https://godoc.org/github.com/zoer/lexer)

Lexer is a simple text tokenizer.

## Usage:
Simple example:
```go
l := lexer.NewLexer(`price 12`)
l.AddMatcher(lexer.TokenizeIfMatches(`[a-z]+`, "WORD"))
l.AddMatcher(lexer.SkipIfMatches(`\s+`))
l.AddMatcher(lexer.TokenizeIfMatches(`\d+`, "PRICE"))

for l.Scan() {
    fmt.Printf("%s => %s\n", l.Token().Name, l.Token().Text)
}
// prints:
// WORD => price
// PRICE => 12
```

Custom matcher:
```go
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
// prints:
// WORD => price
// PRICE => 12.4
// WORD => foo
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/zoer/xmlable.


## License

The gem is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
