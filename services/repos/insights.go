package repos

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"git.shi.foo/git"
	"git.shi.foo/models"
	"git.shi.foo/repositories/cache"
	"git.shi.foo/services/credentials"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
)

var languagePalette = []string{"#61d67e", "#7eeaf6", "#8066f5", "#fadd75", "#e36d6d", "#df81cf"}

type insights struct {
	Topics       []string          `json:"topics"`
	Homepage     string            `json:"homepage"`
	CreatedAt    time.Time         `json:"created_at"`
	Languages    map[string]int64  `json:"languages"`
	Contributors []ContributorView `json:"contributors"`
}

func loadInsights(requestContext context.Context, record *models.Repo) *insights {
	key := fmt.Sprintf(InsightsCacheKey, record.Owner, record.Name)

	if entry, cacheError := cache.Get(key); cacheError == nil {
		var cached insights
		if json.Unmarshal(entry.Value, &cached) == nil {
			return &cached
		}
	}

	token, tokenError := credentials.RepositoryAccessToken(requestContext, record.Owner, record.Name)
	if tokenError != nil {
		return nil
	}

	metadata, metadataError := github.FetchRepo(requestContext, token, record.Owner, record.Name)
	if metadataError != nil {
		return nil
	}

	languages, _ := github.FetchLanguages(requestContext, token, record.Owner, record.Name)

	result := &insights{
		Topics:       metadata.Topics,
		Homepage:     metadata.Homepage,
		CreatedAt:    metadata.CreatedAt,
		Languages:    languages,
		Contributors: topContributors(record.Owner, record.Name),
	}

	if payload, marshalError := json.Marshal(result); marshalError == nil {
		if setError := cache.Set(key, payload, time.Now().Add(InsightsTTL)); setError != nil {
			logger.Errorf(LogPrefix, InsightsCacheLog, setError)
		}
	}

	return result
}

func topContributors(owner string, name string) []ContributorView {
	contributors, contributorError := git.Contributors(owner, name, git.HeadRef)
	if contributorError != nil {
		return nil
	}

	views := make([]ContributorView, 0, RailContributorLimit)
	for index := range contributors {
		if index >= RailContributorLimit {
			break
		}
		views = append(views, ContributorView{
			Name:  contributors[index].Name,
			Email: contributors[index].Email,
			Count: contributors[index].Count,
		})
	}

	return views
}

func toLanguageViews(languages map[string]int64) []LanguageView {
	var total int64
	for _, count := range languages {
		total += count
	}
	if total == 0 {
		return nil
	}

	names := make([]string, 0, len(languages))
	for language := range languages {
		names = append(names, language)
	}
	sort.SliceStable(names, func(first int, second int) bool {
		return languages[names[first]] > languages[names[second]]
	})

	names = names[:min(len(names), LanguageLimit)]
	var displayedTotal int64
	for _, language := range names {
		displayedTotal += languages[language]
	}

	views := make([]LanguageView, 0, len(names))
	for index, language := range names {
		views = append(views, LanguageView{
			Name:    language,
			Percent: fmt.Sprintf("%.1f", float64(languages[language])/float64(total)*100),
			Width:   fmt.Sprintf("%.6f%%", float64(languages[language])/float64(displayedTotal)*100),
			Color:   languagePalette[index%len(languagePalette)],
		})
	}

	return views
}
