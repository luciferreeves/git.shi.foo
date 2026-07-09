package repos

import (
	"context"

	"git.shi.foo/account"
	"git.shi.foo/git"
	"git.shi.foo/repositories/repo"

	"github.com/gofiber/fiber/v2"
)

func GetTagsData(requestContext context.Context, currentUser *account.Response, owner string, name string) (*TagsContext, *fiber.Error) {
	record, viewError := resolveViewable(requestContext, currentUser, owner, name)
	if viewError != nil {
		return nil, viewError
	}

	tagsContext := &TagsContext{
		Title:         record.Owner + "/" + record.Name,
		Owner:         record.Owner,
		Name:          record.Name,
		Private:       record.Private,
		Description:   record.Description,
		DefaultBranch: record.DefaultBranch,
		Status:        record.Status,
		CloneURL:      cloneURL(record.Owner, record.Name),
		ActiveTab:     TabTags,
	}

	if record.Status == repo.StatusActive {
		if refs, refError := git.Tags(record.Owner, record.Name); refError == nil {
			tagsContext.Tags = toRefViews(refs)
		}
	}

	return tagsContext, nil
}
