package markdown

import (
	"bytes"

	"git.shi.foo/utils/syntax"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

var converter = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle(syntax.StyleName),
			highlighting.WithFormatOptions(chromahtml.WithClasses(true)),
		),
	),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()),
)

var policy = buildPolicy()

func buildPolicy() *bluemonday.Policy {
	sanitizer := bluemonday.UGCPolicy()
	sanitizer.AllowAttrs("class").Globally()
	sanitizer.AllowAttrs("id").OnElements("h1", "h2", "h3", "h4", "h5", "h6")
	return sanitizer
}

func Render(source []byte) (string, error) {
	var buffer bytes.Buffer
	if convertError := converter.Convert(source, &buffer); convertError != nil {
		return "", convertError
	}

	return string(policy.SanitizeBytes(buffer.Bytes())), nil
}
