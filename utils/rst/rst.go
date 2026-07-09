package rst

import (
	"bytes"
	"fmt"
	"os/exec"

	"git.shi.foo/utils/sanitize"

	gorst "github.com/hhatto/gorst"
)

func Render(source []byte) (string, error) {
	if rendered, pandocError := renderPandoc(source); pandocError == nil {
		return sanitize.HTML(rendered), nil
	}

	return renderGorst(source)
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
