package rst

const (
	LogPrefix    = "RST"
	PandocBinary = "pandoc"
)

var PandocArgs = []string{"--from", "rst", "--to", "html", "--no-highlight"}
