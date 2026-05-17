package mygithub

import (
	"context"
)

// Gateway は外部 GitHub API 経由でドメイン集約を取得するための抽象。
// 実装はインフラ層（infrastructure/gateway/mygithub）に置く。
type Gateway interface {
	FindAllPullRequestSummaries(ctx context.Context, owner, name string) ([]*PullRequestSummary, error)
	FindPullRequestDetail(ctx context.Context, owner, name string, summary *PullRequestSummary) (*PullRequestDetail, error)
}
