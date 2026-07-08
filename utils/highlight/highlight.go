package highlight

import (
	"bytes"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func Highlight(filename string, source string) (string, error) {
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(StyleName)
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithLineNumbers(true), html.TabWidth(TabWidth))

	iterator, tokeniseError := lexer.Tokenise(nil, source)
	if tokeniseError != nil {
		return "", tokeniseError
	}

	var buffer bytes.Buffer
	if formatError := formatter.Format(&buffer, style, iterator); formatError != nil {
		return "", formatError
	}

	return buffer.String(), nil
}
