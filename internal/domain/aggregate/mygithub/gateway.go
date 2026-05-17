package mygithub

import (
	"context"
	"time"
)

// Gateway は外部 GitHub API 経由でドメイン集約を取得するための抽象。
// 実装はインフラ層（infrastructure/gateway/mygithub）に置く。
type Gateway interface {
	FindAllPullRequests(ctx context.Context, owner, name string) ([]*PullRequest, error)
	FindPullRequestFirstOpenedAt(ctx context.Context, owner, name string, number int) (*time.Time, error)
	FindPullRequestReviews(ctx context.Context, owner, name string, number int) (Reviews, error)
}
