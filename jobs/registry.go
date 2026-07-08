package jobs

var runners = make(map[string]Runner)

func Register(kind string, runner Runner) {
	runners[kind] = runner
}

func lookup(kind string) (Runner, bool) {
	runner, found := runners[kind]
	return runner, found
}
