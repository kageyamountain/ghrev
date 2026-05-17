package mygithub

import (
	"sort"
	"time"
)

type ReviewState string

const (
	ReviewStateApproved         ReviewState = "APPROVED"
	ReviewStateChangesRequested ReviewState = "CHANGES_REQUESTED"
	ReviewStateCommented        ReviewState = "COMMENTED"
	ReviewStateDismissed        ReviewState = "DISMISSED"
	ReviewStatePending          ReviewState = "PENDING"
)

type Review struct {
	User        string
	State       ReviewState
	SubmittedAt time.Time
}

func NewReview(user string, state ReviewState, submittedAt time.Time) Review {
	return Review{
		User:        user,
		State:       state,
		SubmittedAt: submittedAt,
	}
}

func (r Review) IsApproved() bool {
	return r.State == ReviewStateApproved
}

type Reviews []Review

func NewReviews(rs []Review) Reviews {
	return rs
}

// Approved は APPROVED ステートのレビューのみを残して返す。（元の順序を保持）
func (r Reviews) Approved() Reviews {
	result := make(Reviews, 0, len(r))
	for _, review := range r {
		if review.IsApproved() {
			result = append(result, review)
		}
	}
	return result
}

// EarliestPerUser は SubmittedAt 昇順でソートし、同一ユーザにつき最古のレビューのみを残して返す。
// 「同一ユーザの2回目以降の承認は無視する」というドメインルールに対応する。
func (r Reviews) EarliestPerUser() Reviews {
	sortedReviews := make(Reviews, len(r))
	copy(sortedReviews, r)
	sort.Slice(sortedReviews, func(i, j int) bool {
		return sortedReviews[i].SubmittedAt.Before(sortedReviews[j].SubmittedAt)
	})

	seen := make(map[string]bool, len(sortedReviews))
	uniqueReviews := make(Reviews, 0, len(sortedReviews))
	for _, sortedReview := range sortedReviews {
		if seen[sortedReview.User] {
			continue
		}
		seen[sortedReview.User] = true
		uniqueReviews = append(uniqueReviews, sortedReview)
	}
	return uniqueReviews
}
