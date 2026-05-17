package mygithub

import "time"

type IssueEventType string

const (
	IssueEventTypeReadyForReview IssueEventType = "ready_for_review"
	IssueEventTypeConvertToDraft IssueEventType = "convert_to_draft"
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

// FirstOpenedAt は PR が初めてレビュー可能な状態になった日時を返す。
// events を時系列順に走査し、最初に出現した draft 状態遷移イベントで判定する：
//   - 最初が ready_for_review → draft で作成された PR → その日時
//   - 最初が convert_to_draft → 非 draft で作成された PR → prCreatedAt
//   - どちらのイベントも無い → 状態遷移なし → prCreatedAt
//
// 「最初にレビュー可能になった時点」を起点にする仕様のため、途中の draft↔ready toggle は無視する。
func (i IssueEvents) FirstOpenedAt(prCreatedAt time.Time) time.Time {
	for _, e := range i {
		switch e.Type {
		case IssueEventTypeReadyForReview:
			return e.CreatedAt
		case IssueEventTypeConvertToDraft:
			return prCreatedAt
		}
	}
	return prCreatedAt
}
