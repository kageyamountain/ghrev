package mygithub

import (
	"testing"
)

func TestPullRequestSummary_ContainsAnyLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		labels       []string
		targetLabels []string
		want         bool
	}{
		{
			name:         "正常系: 大文字・小文字まで一致する場合、trueが返されること",
			labels:       []string{"Bug", "WIP"},
			targetLabels: []string{"WIP"},
			want:         true,
		},
		{
			name:         "正常系: 大文字・小文字だけが異なる場合、trueが返されること",
			labels:       []string{"Bug", "WIP"},
			targetLabels: []string{"wip"},
			want:         true,
		},
		{
			name:         "正常系: いずれにも一致しない場合、falseが返されること",
			labels:       []string{"Bug", "WIP"},
			targetLabels: []string{"release"},
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			summary := &PullRequestSummary{Labels: tt.labels}

			// Act
			got := summary.ContainsAnyLabel(tt.targetLabels)

			// Assert
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPullRequestSummary_IsAuthoredByAny(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		author        string
		targetAuthors []string
		want          bool
	}{
		{
			name:          "正常系: 大文字・小文字まで一致する場合、trueが返されること",
			author:        "Alice",
			targetAuthors: []string{"bob", "Alice"},
			want:          true,
		},
		{
			name:          "正常系: 大文字・小文字だけが異なる場合、trueが返されること",
			author:        "Alice",
			targetAuthors: []string{"alice"},
			want:          true,
		},
		{
			name:          "正常系: いずれにも一致しない場合、falseが返されること",
			author:        "Alice",
			targetAuthors: []string{"bob", "carol"},
			want:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			summary := &PullRequestSummary{Author: tt.author}

			// Act
			got := summary.IsAuthoredByAny(tt.targetAuthors)

			// Assert
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPullRequestSummary_HasAnyAssignee(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		assignees       []string
		targetAssignees []string
		want            bool
	}{
		{
			name:            "正常系: 大文字・小文字まで一致する場合、trueが返されること",
			assignees:       []string{"Alice", "Bob"},
			targetAssignees: []string{"Bob"},
			want:            true,
		},
		{
			name:            "正常系: 大文字・小文字だけが異なる場合、trueが返されること",
			assignees:       []string{"Alice", "Bob"},
			targetAssignees: []string{"bob"},
			want:            true,
		},
		{
			name:            "正常系: いずれにも一致しない場合、falseが返されること",
			assignees:       []string{"Alice", "Bob"},
			targetAssignees: []string{"carol"},
			want:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			summary := &PullRequestSummary{Assignees: tt.assignees}

			// Act
			got := summary.HasAnyAssignee(tt.targetAssignees)

			// Assert
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
