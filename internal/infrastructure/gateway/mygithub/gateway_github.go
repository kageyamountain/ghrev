package mygithub

import (
	"context"
	"sort"
	"time"

	"github.com/google/go-github/v80/github"

	"github.com/kageyamountain/ghrev/internal/common/config"
)

type Gateway interface {
	FindPullRequests(ctx context.Context, repositoryOwner string, repositoryName string) ([]*github.PullRequest, error)
	FindPullRequestFirstOpenTime(ctx context.Context, repositoryOwner string, repositoryName string, pullRequestNumber int) (*time.Time, error)
	FindPullRequestApproveTimes(ctx context.Context, repositoryOwner string, repositoryName string, pullRequestNumber int) ([]time.Time, error)
}

type gateway struct {
	appConfig    *config.AppConfig
	githubClient *github.Client
}

func NewGateway(
	appConfig *config.AppConfig,
	githubClient *github.Client,
) Gateway {
	return &gateway{
		appConfig:    appConfig,
		githubClient: githubClient,
	}
}

func (g *gateway) FindPullRequests(ctx context.Context, repositoryOwner string, repositoryName string) ([]*github.PullRequest, error) {
	pullRequestListOptions := &github.PullRequestListOptions{
		State:     "all",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var pullRequests []*github.PullRequest
	for {
		partialPullRequests, response, err := g.githubClient.PullRequests.List(ctx, repositoryOwner, repositoryName, pullRequestListOptions)
		if err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, partialPullRequests...)

		if response.NextPage == 0 {
			break
		}
		pullRequestListOptions.Page = response.NextPage
	}

	return pullRequests, nil
}

func (g *gateway) FindPullRequestFirstOpenTime(ctx context.Context, repositoryOwner string, repositoryName string, pullRequestNumber int) (*time.Time, error) {
	listOptions := &github.ListOptions{PerPage: 100}
	events, _, err := g.githubClient.Issues.ListIssueEvents(ctx, repositoryOwner, repositoryName, pullRequestNumber, listOptions)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		if event.GetEvent() == "ready_for_review" {
			// 初回のOpen日時を取得することが目的のため対象イベントが見つかったらreturnする（eventは時系列順に取得される）
			return event.CreatedAt.GetTime(), nil
		}
	}
	return nil, nil
}

func (g *gateway) FindPullRequestApproveTimes(ctx context.Context, repositoryOwner string, repositoryName string, pullRequestNumber int) ([]time.Time, error) {
	listOptions := &github.ListOptions{PerPage: 100}
	pullRequestReviews, _, err := g.githubClient.PullRequests.ListReviews(ctx, repositoryOwner, repositoryName, pullRequestNumber, listOptions)
	if err != nil {
		return nil, err
	}

	var approveTimes []time.Time
	approvedUsers := make(map[string]bool)
	for _, review := range pullRequestReviews {
		if review.GetState() == "APPROVED" {
			user := review.GetUser().GetLogin()
			if approvedUsers[user] { // 同一ユーザのapproveはスキップ
				continue
			}
			approveTimes = append(approveTimes, review.GetSubmittedAt().Time)
			approvedUsers[user] = true
		}
	}

	// 承認日時を昇順ソート
	sort.Slice(approveTimes, func(i int, j int) bool {
		return approveTimes[i].Before(approveTimes[j])
	})

	return approveTimes, nil
}
