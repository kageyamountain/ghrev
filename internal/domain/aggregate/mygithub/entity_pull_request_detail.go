package mygithub

import (
	"time"

	"github.com/kageyamountain/ghrev/internal/common/mytime"
)

// PullRequestDetail は PullRequestSummary にレビュー情報を付与した詳細集約。
// レビュー情報を前提とするドメインロジックはこの型のメソッドとして実装する。
type PullRequestDetail struct {
	PullRequestSummary
	FirstOpenedAt time.Time
	reviews       Reviews
	Additions     int
	Deletions     int
}

// NewPullRequestDetail は summary に events と reviews を取り込んだ PullRequestDetail を返す。
// FirstOpenedAt（PR が初めてレビュー可能になった日時）は events と summary.CreatedAt からここで導出する。
func NewPullRequestDetail(summary *PullRequestSummary, events IssueEvents, reviews Reviews, additions, deletions int) *PullRequestDetail {
	return &PullRequestDetail{
		PullRequestSummary: *summary,
		FirstOpenedAt:      events.FirstOpenedAt(summary.CreatedAt),
		reviews:            reviews,
		Additions:          additions,
		Deletions:          deletions,
	}
}

// TimeToNthApproval は FirstOpenedAt から n 件目の承認までの所要時間を返す。
// n は 1-indexed（n=1 は最初の承認）。同一ユーザの2回目以降の承認はカウントしない。
// JST 基準の土日に該当する時間は計測から除外する。
// 承認が n 件未満の場合は ok=false。
func (p *PullRequestDetail) TimeToNthApproval(n int) (time.Duration, bool) {
	approved := p.reviews.Approved().EarliestPerUser()
	if len(approved) < n {
		return 0, false
	}
	return mytime.BusinessDuration(p.FirstOpenedAt, approved[n-1].SubmittedAt), true
}
