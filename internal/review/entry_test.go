package review_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/eeay0/cortex/internal/review"
)

func TestReview_NewEntryWithoutOptions(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid",
			title:       "test entry",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "empty",
			title:       "",
			wantErr:     true,
			expectedErr: review.ErrTitleEmpty,
		},
		{
			name:        "too long",
			title:       strings.Repeat("t", review.MaxTitleLength+1),
			wantErr:     true,
			expectedErr: review.ErrTitleTooLong,
		},
		{
			name:        "with spaces",
			title:       "   test entry   ",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "only spaces",
			title:       "     ",
			wantErr:     true,
			expectedErr: review.ErrTitleEmpty,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, gotErr := review.NewEntry(test.title)
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected err %v, entry nil", test.expectedErr)
				}
				if !errors.Is(gotErr, test.expectedErr) {
					t.Errorf("expected error %v, entry %v", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected no error, entry %v", gotErr)
				}
				if entry.Title != strings.TrimSpace(test.title) {
					t.Errorf(
						"expecetd title %q, entry %q",
						entry.Title,
						strings.TrimSpace(test.title),
					)
				}
				if entry.Description != "" {
					t.Errorf("expected no description, entry %q", entry.Description)
				}
				if entry.Recall != -1 {
					t.Errorf("expected recall %d, entry %q", -1, entry.Recall)
				}
				if entry.Category != "None" {
					t.Errorf("expected category %q, entry %q", "None", entry.Category)
				}
				if entry.Interval != 1 {
					t.Errorf("expected interval %d, entry %d", 1, entry.Interval)
				}
				if entry.ReviewDate.IsZero() {
					t.Errorf("expected ReviewDate to be set, entry zero value")
				}
				if !entry.LastReview.IsZero() {
					t.Errorf("expected LastReview to be zero, entry nonzero")
				}
				if entry.CreatedAt.IsZero() {
					t.Errorf("expected CreatedAt to be set, entry zero value")
				}
			}
		})
	}
}

func TestReview_NewEntryWithDescription(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid",
			title:       "test entry",
			description: "this is a test entry.",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "too long",
			title:       "test entry",
			description: strings.Repeat("d", review.MaxDescriptionLength+1),
			wantErr:     true,
			expectedErr: review.ErrDescriptionTooLong,
		},
		{
			name:        "with spaces",
			title:       "   test entry   ",
			description: "   this is a test entry.   ",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "only spaces",
			title:       "test entry",
			description: "     ",
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, gotErr := review.NewEntry(test.title, review.WithDescription(test.description))
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected error %v, entry nil", test.expectedErr)
				}
				if !errors.Is(gotErr, test.expectedErr) {
					t.Errorf("expected error %v, entry %v", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected no error, entry %v", gotErr)
				}
				if entry.Description != strings.TrimSpace(test.description) {
					t.Errorf(
						"expected description %q, entry %q",
						strings.TrimSpace(test.description),
						entry.Description,
					)
				}
			}
		})
	}
}

func TestReview_NewEntryWithCategory(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		category    string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "valid",
			title:       "test entry",
			category:    "valid category",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "with spaces",
			title:       "test entry",
			category:    "   valid category   ",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "empty",
			title:       "test entry",
			category:    "",
			wantErr:     true,
			expectedErr: review.ErrCategoryEmpty,
		},
		{
			name:        "too long",
			title:       "test entry",
			category:    strings.Repeat("c", review.MaxCategoryLength+1),
			wantErr:     true,
			expectedErr: review.ErrCategoryTooLong,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, gotErr := review.NewEntry(test.title, review.WithCategory(test.category))
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected error %q, entry nil", test.expectedErr)
				}
				if !errors.Is(test.expectedErr, gotErr) {
					t.Errorf("expected error %q, entry %q", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected err nil, entry %q", gotErr)
				}
				if entry.Category != strings.TrimSpace(test.category) {
					t.Errorf(
						"expected category %q, entry %q",
						strings.TrimSpace(test.category),
						entry.Category,
					)
				}
			}
		})
	}
}

func TestReview_UpdateTitle(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		newTitle    string
		wantErr     bool
		expectedErr error
	}{
		{
			"valid",
			"test title",
			"new valid title",
			false,
			nil,
		},
		{
			"empty",
			"test title",
			"   ",
			true,
			review.ErrTitleEmpty,
		},
		{
			"too long",
			"test title",
			strings.Repeat("t", review.MaxTitleLength+1),
			true,
			review.ErrTitleTooLong,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, err := review.NewEntry(test.title)
			if err != nil {
				t.Fatalf("failed to crate entry: %v", err)
			}
			gotErr := entry.UpdateTitle(test.newTitle)
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected err %q, entry nil", test.expectedErr)
				}
				if !errors.Is(test.expectedErr, gotErr) {
					t.Errorf("expected error %q, entry %q", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected err nil, entry %q", gotErr)
				}
				if strings.TrimSpace(test.newTitle) != entry.Title {
					t.Errorf(
						"expected title %q, entry %q",
						strings.TrimSpace(test.newTitle),
						entry.Title,
					)
				}
				if entry.UpdatedAt.IsZero() {
					t.Errorf("expected UpdatedAt to be set, entry zero value")
				}
			}
		})
	}
}

func TestReview_UpdateDescription(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			"valid",
			"test title",
			"new valid description",
			false,
			nil,
		},
		{
			"empty",
			"test title",
			"   ",
			false,
			nil,
		},
		{
			"too long",
			"test title",
			strings.Repeat("d", review.MaxDescriptionLength+1),
			true,
			review.ErrDescriptionTooLong,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, err := review.NewEntry(test.title)
			if err != nil {
				t.Fatalf("failed to crate entry: %v", err)
			}
			gotErr := entry.UpdateDescription(test.description)
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected err %q, entry nil", test.expectedErr)
				}
				if !errors.Is(test.expectedErr, gotErr) {
					t.Errorf("expected error %q, entry %q", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected err nil, entry %q", gotErr)
				}
				if strings.TrimSpace(test.description) != entry.Description {
					t.Errorf(
						"expected description %q, entry %q",
						strings.TrimSpace(test.description),
						entry.Description,
					)
				}
				if entry.UpdatedAt.IsZero() {
					t.Errorf("expected UpdatedAt to be set, entry zero value")
				}
			}
		})
	}
}

func TestReview_UpdateRecall(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		recall      review.Recall
		wantErr     bool
		expectedErr error
	}{
		{
			"not reviewed",
			"test entry",
			review.NotReviewed,
			true,
			review.ErrInvalidRecall,
		},
		{
			"failed",
			"test entry",
			review.Failed,
			false,
			nil,
		},
		{
			"hard",
			"test entry",
			review.Hard,
			false,
			nil,
		},
		{
			"good",
			"test entry",
			review.Good,
			false,
			nil,
		},
		{
			"easy",
			"test entry",
			review.Easy,
			false,
			nil,
		},
		{
			"invalid",
			"test entry",
			20,
			true,
			review.ErrInvalidRecall,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, err := review.NewEntry(test.title)
			if err != nil {
				t.Fatalf("failed to crate entry: %v", err)
			}
			gotErr := entry.UpdateRecall(test.recall)
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected err %q, entry nil", test.expectedErr)
				}
				if !errors.Is(test.expectedErr, gotErr) {
					t.Errorf("expected error %q, entry %q", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected err nil, entry %q", gotErr)
				}
				if test.recall != entry.Recall {
					t.Errorf("expected recall %q, entry %q", test.recall, entry.Recall)
				}
				if entry.UpdatedAt.IsZero() {
					t.Errorf("expected UpdatedAt to be set, entry zero value")
				}
				if entry.LastReview.IsZero() {
					t.Errorf("expected LastReview be set, entry zero value")
				}
			}
		})
	}
}

func TestReview_UpdateCategory(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		category    string
		wantErr     bool
		expectedErr error
	}{
		{
			"valid",
			"test title",
			"new category",
			false,
			nil,
		},
		{
			"empty",
			"test title",
			"   ",
			true,
			review.ErrCategoryEmpty,
		},
		{
			"too long",
			"test title",
			strings.Repeat("c", review.MaxCategoryLength+1),
			true,
			review.ErrCategoryTooLong,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, err := review.NewEntry(test.title)
			if err != nil {
				t.Fatalf("failed to crate entry: %v", err)
			}
			gotErr := entry.UpdateCategory(test.category)
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected err %q, entry nil", test.expectedErr)
				}
				if !errors.Is(test.expectedErr, gotErr) {
					t.Errorf("expected error %q, entry %q", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected err nil, entry %q", gotErr)
				}
				if strings.TrimSpace(test.category) != entry.Category {
					t.Errorf(
						"expected category %q, entry %q",
						strings.TrimSpace(test.category),
						entry.Category,
					)
				}
				if entry.UpdatedAt.IsZero() {
					t.Errorf("expected UpdatedAt to be set, entry zero value")
				}
			}
		})
	}
}

func TestReview_UpdateInterval(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		recall      review.Recall
		wantErr     bool
		expectedErr error
	}{
		{
			"not reviewed",
			"test entry",
			review.NotReviewed,
			true,
			review.ErrInvalidRecall,
		},
		{
			"failed",
			"test entry",
			review.Failed,
			false,
			nil,
		},
		{
			"hard",
			"test entry",
			review.Hard,
			false,
			nil,
		},
		{
			"good",
			"test entry",
			review.Good,
			false,
			nil,
		},
		{
			"easy",
			"test entry",
			review.Easy,
			false,
			nil,
		},
		{
			"invalid",
			"test entry",
			20,
			true,
			review.ErrInvalidRecall,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, err := review.NewEntry(test.title)
			if err != nil {
				t.Fatalf("failed to crate entry: %v", err)
			}
			gotErr := entry.UpdateRecall(test.recall)
			entry.UpdateInterval()
			expectedInterval := 1 * int(test.recall)
			if test.recall == review.Failed {
				expectedInterval = 1
			}
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("expected err %q, entry nil", test.expectedErr)
				}
				if !errors.Is(test.expectedErr, gotErr) {
					t.Errorf("expected error %q, entry %q", test.expectedErr, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("expected err nil, entry %q", gotErr)
				}
				if expectedInterval != entry.Interval {
					t.Errorf("expected interval %d, entry %d", expectedInterval, entry.Interval)
				}
				if entry.UpdatedAt.IsZero() {
					t.Errorf("expected UpdatedAt to be set, entry zero value")
				}
			}
		})
	}
}

func TestReview_SecondInterval(t *testing.T) {
	tests := []struct {
		name   string
		title  string
		recall review.Recall
	}{
		{
			"failed",
			"test entry",
			review.Failed,
		},
		{
			"hard",
			"test entry",
			review.Hard,
		},
		{
			"easy",
			"test entry",
			review.Easy,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry, err := review.NewEntry(test.title)
			if err != nil {
				t.Fatalf("failed to crate entry: %v", err)
			}
			entry.UpdateRecall(review.Good)
			entry.UpdateInterval()

			entry.UpdateRecall(test.recall)
			entry.UpdateInterval()

			expectedInterval := int(1 * test.recall * review.Good)
			if test.recall == review.Failed {
				expectedInterval = 1
			}
			if expectedInterval != entry.Interval {
				t.Errorf("expected interval %d, entry %d", expectedInterval, entry.Interval)
			}
			if entry.UpdatedAt.IsZero() {
				t.Errorf("expected UpdatedAt to be set, entry zero value")
			}
		})
	}
}

func TestReview_UpdateInterval_ReviewDateCalculation(t *testing.T) {
	entry, _ := review.NewEntry("test")
	initialReviewDate := entry.ReviewDate

	entry.UpdateRecall(review.Good)
	entry.UpdateInterval()

	if entry.ReviewDate.Equal(initialReviewDate) {
		t.Error("ReviewDate should have changed after interval update")
	}

	expectedDate := time.Now().Add(time.Duration(entry.Interval) * 24 * time.Hour)

	diff := entry.ReviewDate.Sub(expectedDate).Abs()
	if diff > time.Second {
		t.Errorf("ReviewDate not calculated correctly, diff: %v", diff)
	}
}
