package mygithub

import "time"

// PullRequestDetail は PullRequestSummary にレビュー情報を付与した詳細集約。
// レビュー情報を前提とするドメインロジックはこの型のメソッドとして実装する。
type PullRequestDetail struct {
	PullRequestSummary
	openedAt *time.Time
	reviews  Reviews
}

// NewPullRequestDetail は summary に openedAt と reviews を取り込んだ PullRequestDetail を返す。
// openedAt は ready-for-review の発火日時。nil の場合は「作成時から open」を意味する。
func NewPullRequestDetail(summary *PullRequestSummary, openedAt *time.Time, reviews Reviews) *PullRequestDetail {
	return &PullRequestDetail{
		PullRequestSummary: *summary,
		openedAt:           openedAt,
		reviews:            reviews,
	}
}

// FirstOpenedAt は PR が初めてレビュー可能になった日時を返す。
// openedAt が nil（明示的な ready-for-review 遷移がない）場合は CreatedAt にフォールバックする。
func (p *PullRequestDetail) FirstOpenedAt() time.Time {
	if p.openedAt != nil {
		return *p.openedAt
	}
	return p.CreatedAt
}

// TimeToSecondApproval は FirstOpenedAt から2件目の承認までの所要時間を返す。
// 承認が2件未満の場合は ok=false。
func (p *PullRequestDetail) TimeToSecondApproval() (time.Duration, bool) {
	approved := p.reviews.Approved().EarliestPerUser()
	if len(approved) < 2 {
		return 0, false
	}
	return approved[1].SubmittedAt.Sub(p.FirstOpenedAt()), true
}
