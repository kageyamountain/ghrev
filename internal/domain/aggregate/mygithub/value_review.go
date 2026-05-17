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
	IsBot       bool
	State       ReviewState
	SubmittedAt time.Time
}

func NewReview(user string, isBot bool, state ReviewState, submittedAt time.Time) Review {
	return Review{
		User:        user,
		IsBot:       isBot,
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

// ExcludingDismissed は DISMISSED ステートのレビューを除外して返す。（元の順序を保持）
// 「最初のレビュー反応」を扱う際、後で取り消された approve は反応としてカウントしない方針に対応する。
func (r Reviews) ExcludingDismissed() Reviews {
	result := make(Reviews, 0, len(r))
	for _, review := range r {
		if review.State == ReviewStateDismissed {
			continue
		}
		result = append(result, review)
	}
	return result
}

// ExcludingBots は GitHub App 等のボットによるレビューを除外して返す。（元の順序を保持）
// CodeRabbit のような自動レビュー bot が「最初のレビュー」として計測されるのを防ぐ。
func (r Reviews) ExcludingBots() Reviews {
	result := make(Reviews, 0, len(r))
	for _, review := range r {
		if review.IsBot {
			continue
		}
		result = append(result, review)
	}
	return result
}

// ExcludingUser は指定ユーザのレビューを除外して返す。（元の順序を保持）
// PR 作者自身がレビュー形式で残すコメント（CodeRabbit への返信等）を first review と誤認しないために使う。
func (r Reviews) ExcludingUser(login string) Reviews {
	result := make(Reviews, 0, len(r))
	for _, review := range r {
		if review.User == login {
			continue
		}
		result = append(result, review)
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
