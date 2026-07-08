package git

func Blob(owner string, name string, ref string, path string) ([]byte, error) {
	return runGit(RepoPath(owner, name), "show", ref+":"+path)
}
