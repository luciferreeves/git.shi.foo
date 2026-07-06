package home

import "git.shi.foo/account"

func GetIndexData(currentUser *account.Response) *IndexContext {
	return &IndexContext{
		Title: HomeTitle,
		User:  currentUser,
	}
}
