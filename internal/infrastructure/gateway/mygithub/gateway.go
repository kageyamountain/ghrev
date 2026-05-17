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

func (g *gateway) FindAllPullRequests(ctx context.Context, owner, name string) ([]*mygithub.PullRequest, error) {
	options := &github.PullRequestListOptions{
		State:     "all",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var result []*mygithub.PullRequest
	for {
		prs, response, err := g.githubClient.PullRequests.List(ctx, owner, name, options)
		if err != nil {
			return nil, err
		}
		for _, pr := range prs {
			result = append(result, toPullRequest(pr))
		}

		if response.NextPage == 0 {
			break
		}
		options.Page = response.NextPage
	}

	return result, nil
}

func (g *gateway) FindPullRequestFirstOpenedAt(ctx context.Context, owner, name string, number int) (*time.Time, error) {
	listOptions := &github.ListOptions{PerPage: 100}
	events, _, err := g.githubClient.Issues.ListIssueEvents(ctx, owner, name, number, listOptions)
	if err != nil {
		return nil, err
	}

	// events は時系列順に返ってくるため最初に見つかった ready_for_review が初回 Open 日時。
	for _, event := range events {
		if event.GetEvent() == "ready_for_review" {
			return event.CreatedAt.GetTime(), nil
		}
	}
	return nil, nil
}

func (g *gateway) FindPullRequestReviews(ctx context.Context, owner, name string, number int) (mygithub.Reviews, error) {
	listOptions := &github.ListOptions{PerPage: 100}
	reviews, _, err := g.githubClient.PullRequests.ListReviews(ctx, owner, name, number, listOptions)
	if err != nil {
		return nil, err
	}

	result := make([]mygithub.Review, 0, len(reviews))
	for _, review := range reviews {
		result = append(result, mygithub.NewReview(
			review.GetUser().GetLogin(),
			mygithub.ReviewState(review.GetState()),
			review.GetSubmittedAt().Time,
		))
	}
	return mygithub.NewReviews(result), nil
}

func toPullRequest(pr *github.PullRequest) *mygithub.PullRequest {
	labels := make([]string, 0, len(pr.Labels))
	for _, label := range pr.Labels {
		labels = append(labels, label.GetName())
	}
	return &mygithub.PullRequest{
		Number:    pr.GetNumber(),
		CreatedAt: pr.GetCreatedAt().Time,
		Labels:    labels,
		HTMLURL:   pr.GetHTMLURL(),
	}
}
