package mygithub

import "time"

type IssueEventType string

const (
	IssueEventTypeReadyForReview IssueEventType = "ready_for_review"
)

type IssueEvent struct {
	Type      IssueEventType
	CreatedAt time.Time
}

func NewIssueEvent(eventType IssueEventType, createdAt time.Time) IssueEvent {
	return IssueEvent{
		Type:      eventType,
		CreatedAt: createdAt,
	}
}

type IssueEvents []IssueEvent

func NewIssueEvents(events []IssueEvent) IssueEvents {
	return events
}

// FirstReadyForReviewAt は events を時系列順に走査し、最初に見つかった
// ready_for_review イベントの発火日時を返す。該当イベントが無い場合は nil。
func (i IssueEvents) FirstReadyForReviewAt() *time.Time {
	for _, e := range i {
		if e.Type == IssueEventTypeReadyForReview {
			createdAt := e.CreatedAt
			return &createdAt
		}
	}
	return nil
}
