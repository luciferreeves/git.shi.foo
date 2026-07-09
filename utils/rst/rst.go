package rst

import (
	"bytes"
	"fmt"

	"git.shi.foo/utils/sanitize"

	gorst "github.com/hhatto/gorst"
)

func Render(source []byte) (rendered string, renderError error) {
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
