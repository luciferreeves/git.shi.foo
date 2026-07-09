package repos

import (
	"context"

	"git.shi.foo/account"
	"git.shi.foo/git"
	"git.shi.foo/repositories/repo"

	"github.com/gofiber/fiber/v2"
)

func GetContributorsData(requestContext context.Context, currentUser *account.Response, owner string, name string) (*ContributorsContext, *fiber.Error) {
	record, viewError := resolveViewable(requestContext, currentUser, owner, name)
	if viewError != nil {
		return nil, viewError
	}

	contributorsContext := &ContributorsContext{
		Title:         record.Owner + "/" + record.Name,
		Owner:         record.Owner,
		Name:          record.Name,
		Private:       record.Private,
		Description:   record.Description,
		DefaultBranch: record.DefaultBranch,
		Status:        record.Status,
		CloneURL:      cloneURL(record.Owner, record.Name),
		SSHCloneURL:   sshCloneURL(record.Owner, record.Name),
		ActiveTab:     TabContributors,
	}

	if record.Status == repo.StatusActive {
		if contributors, contributorError := git.Contributors(record.Owner, record.Name, git.HeadRef); contributorError == nil {
			for index := range contributors {
				contributorsContext.Contributors = append(contributorsContext.Contributors, ContributorView{
					Name:  contributors[index].Name,
					Email: contributors[index].Email,
					Count: contributors[index].Count,
				})
			}
		}
	}

	return contributorsContext, nil
}
