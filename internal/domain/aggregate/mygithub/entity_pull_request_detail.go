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

// TimeToFirstReview は FirstOpenedAt から最初のレビュー反応までの所要時間を返す。
// approve / changes_requested / commented のいずれかを「反応」とみなし、DISMISSED と bot のレビュー、
// および PR 作者自身のレビュー（CodeRabbit への返信等のセルフコメント）は除外する。
// 同一ユーザが複数回反応した場合は最古のものを採用する（既存の Nth approval と同じ方針）。
// JST 基準の土日に該当する時間は計測から除外する。
// 反応が1件もない場合は ok=false。
func (p *PullRequestDetail) TimeToFirstReview() (time.Duration, bool) {
	reviews := p.reviews.ExcludingBots().ExcludingUser(p.Author).ExcludingDismissed().EarliestPerUser()
	if len(reviews) == 0 {
		return 0, false
	}
	return mytime.BusinessDuration(p.FirstOpenedAt, reviews[0].SubmittedAt), true
}
