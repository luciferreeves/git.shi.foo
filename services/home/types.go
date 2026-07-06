package home

import "git.shi.foo/account"

type IndexContext struct {
	Title string
	User  *account.Response
}
