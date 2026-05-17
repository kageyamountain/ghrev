package mygithub

import (
	"context"
	"time"

	"github.com/google/go-github/v80/github"

	"github.com/kageyamountain/ghrev/internal/domain/aggregate/mygithub"
)

type gateway struct {
	githubClient *github.Client
}

func NewGateway(githubClient *github.Client) mygithub.Gateway {
	return &gateway{
		githubClient: githubClient,
	}
}

// FindPullRequestSummaries は created 降順で PR 一覧を取得し、createdAtFrom より古い PR に
// 到達した時点でページングを打ち切ることで、期間外の PR を取りに行く無駄な API コールを避ける。
// createdAtTo 側の上限フィルタは降順取得では API コール削減に繋がらないため呼び出し側に委ねる。
func (g *gateway) FindPullRequestSummaries(ctx context.Context, owner, name string, createdAtFrom time.Time) ([]*mygithub.PullRequestSummary, error) {
	options := &github.PullRequestListOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var result []*mygithub.PullRequestSummary
	for {
		prs, response, err := g.githubClient.PullRequests.List(ctx, owner, name, options)
		if err != nil {
			return nil, err
		}

		reachedCutoff := false
		for _, pr := range prs {
			summary := toPullRequestSummary(pr)
			if summary.CreatedAt.Before(createdAtFrom) {
				reachedCutoff = true
				break
			}
			result = append(result, summary)
		}
		if reachedCutoff {
			break
		}

		if response.NextPage == 0 {
			break
		}
		options.Page = response.NextPage
	}

	return result, nil
}

// FindPullRequestDetail はイベントとレビューを取得し、PR が初めて open になった日時を解決した上で、
// 与えられた PullRequestSummary と合わせて PullRequestDetail を構築して返す。
func (g *gateway) FindPullRequestDetail(ctx context.Context, owner, name string, summary *mygithub.PullRequestSummary) (*mygithub.PullRequestDetail, error) {
	events, err := g.findIssueEvents(ctx, owner, name, summary.Number)
	if err != nil {
		return nil, err
	}

	reviews, err := g.findReviews(ctx, owner, name, summary.Number)
	if err != nil {
		return nil, err
	}

	pullRequest, _, err := g.githubClient.PullRequests.Get(ctx, owner, name, summary.Number)
	if err != nil {
		return nil, err
	}

	return mygithub.NewPullRequestDetail(
		summary,
		events,
		reviews,
		pullRequest.GetAdditions(),
		pullRequest.GetDeletions(),
	), nil
}

func (g *gateway) findIssueEvents(ctx context.Context, owner, name string, number int) (mygithub.IssueEvents, error) {
	listOptions := &github.ListOptions{PerPage: 100}

	var result []mygithub.IssueEvent
	for {
		events, response, err := g.githubClient.Issues.ListIssueEvents(ctx, owner, name, number, listOptions)
		if err != nil {
			return nil, err
		}
		for _, event := range events {
			result = append(result, mygithub.NewIssueEvent(
				mygithub.IssueEventType(event.GetEvent()),
				event.GetCreatedAt().Time,
			))
		}

		if response.NextPage == 0 {
			break
		}
		listOptions.Page = response.NextPage
	}
	return mygithub.NewIssueEvents(result), nil
}

func (g *gateway) findReviews(ctx context.Context, owner, name string, number int) (mygithub.Reviews, error) {
	listOptions := &github.ListOptions{PerPage: 100}

	var result []mygithub.Review
	for {
		reviews, response, err := g.githubClient.PullRequests.ListReviews(ctx, owner, name, number, listOptions)
		if err != nil {
			return nil, err
		}
		for _, review := range reviews {
			result = append(result, mygithub.NewReview(
				review.GetUser().GetLogin(),
				review.GetUser().GetType() == "Bot",
				mygithub.ReviewState(review.GetState()),
				review.GetSubmittedAt().Time,
			))
		}

		if response.NextPage == 0 {
			break
		}
		listOptions.Page = response.NextPage
	}
	return mygithub.NewReviews(result), nil
}

func toPullRequestSummary(pullRequest *github.PullRequest) *mygithub.PullRequestSummary {
	labels := make([]string, 0, len(pullRequest.Labels))
	for _, label := range pullRequest.Labels {
		labels = append(labels, label.GetName())
	}
	return &mygithub.PullRequestSummary{
		Number:    pullRequest.GetNumber(),
		Author:    pullRequest.GetUser().GetLogin(),
		CreatedAt: pullRequest.GetCreatedAt().Time,
		Labels:    labels,
		HTMLURL:   pullRequest.GetHTMLURL(),
	}
}
