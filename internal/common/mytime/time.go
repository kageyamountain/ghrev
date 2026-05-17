package mytime

import "time"

var jst = time.FixedZone("JST", 9*60*60)

// BusinessDuration は start から end までの経過時間から、JST 基準の土日に該当する分を差し引いて返す。
// JST のカレンダー日ごとに区切り、土日に重なる区間を加算しない方式で計算する。
func BusinessDuration(start, end time.Time) time.Duration {
	if !end.After(start) {
		return 0
	}
	start = start.In(jst)
	end = end.In(jst)

	var total time.Duration
	cursor := start
	for cursor.Before(end) {
		nextMidnight := time.Date(cursor.Year(), cursor.Month(), cursor.Day(), 0, 0, 0, 0, jst).AddDate(0, 0, 1)
		segEnd := nextMidnight
		if segEnd.After(end) {
			segEnd = end
		}
		wd := cursor.Weekday()
		if wd != time.Saturday && wd != time.Sunday {
			total += segEnd.Sub(cursor)
		}
		cursor = nextMidnight
	}
	return total
}
