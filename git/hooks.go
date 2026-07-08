package git

import (
	"fmt"
	"os"
	"path/filepath"

	"git.shi.foo/config"
)

const preReceiveHook = `#!/bin/sh
refspecs=""
while read old new ref; do
	if [ "$new" = "0000000000000000000000000000000000000000" ]; then
		refspecs="$refspecs :$ref"
	else
		refspecs="$refspecs $new:$ref"
	fi
done
[ -z "$refspecs" ] && exit 0
url=$(git remote get-url origin) || exit 1
exec git push --atomic "$url" $refspecs
`

func HooksDir() string {
	directory := filepath.Join(config.Git.ReposRoot, HooksSubdir)
	if absolute, absError := filepath.Abs(directory); absError == nil {
		return absolute
	}
	return directory
}

func InstallHooks() error {
	directory := HooksDir()
	if mkdirError := os.MkdirAll(directory, DirectoryPermission); mkdirError != nil {
		return fmt.Errorf(HooksDirFailed, mkdirError)
	}

	hookPath := filepath.Join(directory, PreReceiveHookName)
	if writeError := os.WriteFile(hookPath, []byte(preReceiveHook), HookPermission); writeError != nil {
		return fmt.Errorf(HookWriteFailed, writeError)
	}

	return nil
}
