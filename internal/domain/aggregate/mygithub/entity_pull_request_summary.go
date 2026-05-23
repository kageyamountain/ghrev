package mygithub

import "time"

// PullRequestSummary はレビュー情報を持たない PR の概要を表す。
// 期間／ラベルなどメタデータだけで判定できるフィルタリングはこの型のメソッドで行う。
type PullRequestSummary struct {
	Number    int
	Author    string
	CreatedAt time.Time
	Labels    []string
	Assignees []string
	HTMLURL   string
}

// IsCreatedWithin は CreatedAt が [from, to] の閉区間に含まれるかを返す。
func (p *PullRequestSummary) IsCreatedWithin(from, to time.Time) bool {
	return !p.CreatedAt.Before(from) && !p.CreatedAt.After(to)
}

// ContainsAnyLabel は targetLabels のいずれかに一致するラベルが付与されているかを返す。
func (p *PullRequestSummary) ContainsAnyLabel(targetLabels []string) bool {
	for _, prLabel := range p.Labels {
		for _, target := range targetLabels {
			if prLabel == target {
				return true
			}
		}
	}
	return false
}

// HasAnyAssignee は targetAssignees のいずれかに一致する assignee が PR に設定されているかを返す。
func (p *PullRequestSummary) HasAnyAssignee(targetAssignees []string) bool {
	for _, prAssignee := range p.Assignees {
		for _, target := range targetAssignees {
			if prAssignee == target {
				return true
			}
		}
	}
	return false
}
