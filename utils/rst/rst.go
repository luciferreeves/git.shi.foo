package rst

import (
	"bytes"
	"fmt"
	"html"
	"os/exec"
	"regexp"

	"git.shi.foo/utils/highlight"
	"git.shi.foo/utils/sanitize"

	gorst "github.com/hhatto/gorst"
)

var codeBlockPattern = regexp.MustCompile(`(?s)<pre class="([^"]+)"><code>(.*?)</code></pre>`)

func Render(source []byte) (string, error) {
	if rendered, pandocError := renderPandoc(source); pandocError == nil {
		return sanitize.HTML([]byte(highlightCodeBlocks(string(rendered)))), nil
	}

	return renderGorst(source)
}

func highlightCodeBlocks(rendered string) string {
	return codeBlockPattern.ReplaceAllStringFunc(rendered, func(match string) string {
		groups := codeBlockPattern.FindStringSubmatch(match)
		code := html.UnescapeString(groups[2])

		highlighted, highlightError := highlight.Classes(groups[1], code)
		if highlightError != nil {
			return match
		}
		return highlighted
	})
}

func renderPandoc(source []byte) ([]byte, error) {
	command := exec.Command(PandocBinary, PandocArgs...)
	command.Stdin = bytes.NewReader(source)

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if runError := command.Run(); runError != nil {
		return nil, fmt.Errorf(PandocFailed, runError, stderr.String())
	}

	return stdout.Bytes(), nil
}

func renderGorst(source []byte) (rendered string, renderError error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			rendered = ""
			renderError = fmt.Errorf(RenderPanic, recovered)
		}
	}()

	var buffer bytes.Buffer
	parser := gorst.NewParser(nil)
	parser.ReStructuredText(bytes.NewReader(source), gorst.ToHTML(&buffer))

	return sanitize.HTML(buffer.Bytes()), nil
}
