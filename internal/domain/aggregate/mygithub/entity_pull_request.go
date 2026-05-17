package mygithub

import "time"

type PullRequest struct {
	Number    int
	CreatedAt time.Time
	Labels    []string
	HTMLURL   string

	// OpenedAt はレビュー可能になった日時。
	// nil の場合は「作成時から open」だったことを意味する。
	OpenedAt *time.Time

	Reviews Reviews
}

// FirstOpenedAt は PR が初めてレビュー可能になった日時を返す。
// OpenedAt が nil（明示的な ready-for-review 遷移がない）場合は CreatedAt にフォールバックする。
func (p *PullRequest) FirstOpenedAt() time.Time {
	if p.OpenedAt != nil {
		return *p.OpenedAt
	}
	return p.CreatedAt
}

// IsCreatedWithin は CreatedAt が [from, to] の閉区間に含まれるかを返す。
func (p *PullRequest) IsCreatedWithin(from, to time.Time) bool {
	return !p.CreatedAt.Before(from) && !p.CreatedAt.After(to)
}

// ContainsAnyLabel は targetLabels のいずれかに一致するラベルが付与されているかを返す。
func (p *PullRequest) ContainsAnyLabel(targetLabels []string) bool {
	for _, prLabel := range p.Labels {
		for _, target := range targetLabels {
			if prLabel == target {
				return true
			}
		}
	}
	return false
}

// TimeToSecondApproval は FirstOpenedAt から2件目の承認までの所要時間を返す。
// 承認が2件未満の場合は ok=false。
func (p *PullRequest) TimeToSecondApproval() (time.Duration, bool) {
	approved := p.Reviews.Approved().EarliestPerUser()
	if len(approved) < 2 {
		return 0, false
	}
	return approved[1].SubmittedAt.Sub(p.FirstOpenedAt()), true
}
