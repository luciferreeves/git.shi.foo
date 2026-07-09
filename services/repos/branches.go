package repos

import (
	"context"

	"git.shi.foo/account"
	"git.shi.foo/git"
	"git.shi.foo/repositories/repo"

	"github.com/gofiber/fiber/v2"
)

func GetBranchesData(requestContext context.Context, currentUser *account.Response, owner string, name string) (*BranchesContext, *fiber.Error) {
	record, viewError := resolveViewable(requestContext, currentUser, owner, name)
	if viewError != nil {
		return nil, viewError
	}

	branchesContext := &BranchesContext{
		Title:         record.Owner + "/" + record.Name,
		Owner:         record.Owner,
		Name:          record.Name,
		Private:       record.Private,
		Description:   record.Description,
		DefaultBranch: record.DefaultBranch,
		Status:        record.Status,
		CloneURL:      cloneURL(record.Owner, record.Name),
		ActiveTab:     TabBranches,
	}

	if record.Status == repo.StatusActive {
		if refs, refError := git.Branches(record.Owner, record.Name); refError == nil {
			branchesContext.Branches = toRefViews(record.Owner, record.Name, refs)
		}
	}

	return branchesContext, nil
}

func toRefViews(owner string, name string, refs []git.Ref) []RefView {
	views := make([]RefView, 0, len(refs))
	for index := range refs {
		views = append(views, RefView{
			Name:    refs[index].Name,
			Short:   refs[index].Short,
			Message: refs[index].Message,
			Author:  refs[index].Author,
			URL:     treeURL(owner, name, refs[index].Name, ""),
			When:    refs[index].When,
		})
	}
	return views
}
