// Package review provides spaced repetition scheduling based on user performance ratings.
package review

import (
	"errors"
	"strings"
	"time"
)

// INFO: Constants and Variables

// Errors for entry validation.
var (
	ErrTitleEmpty         = errors.New("title cannot be empty")
	ErrTitleTooLong       = errors.New("title exceeds maximum length")
	ErrDescriptionTooLong = errors.New("description exceeds maximum length")
	ErrCategoryEmpty      = errors.New("category cannot be empty")
	ErrCategoryTooLong    = errors.New("category exceeds maximum length")
	ErrInvalidRecall      = errors.New("invalid recall value")
)

// Constants for length limits.
const (
	MaxTitleLength       = 255
	MaxDescriptionLength = 4000
	MaxCategoryLength    = 30
)

const day = time.Hour * 24

const intervalLimit = 90

// Recall represents user's recall performance.
// It's being used to calculate interval.
type Recall float32

// Constants for recall performance.
const (
	NotReviewed Recall = -1
	Failed      Recall = 0
	Hard        Recall = 1.2
	Good        Recall = 2
	Easy        Recall = 2.5
)

func (p Recall) String() string {
	switch p {
	case NotReviewed:
		return "Not Reviewed"
	case Failed:
		return "Forgotten"
	case Hard:
		return "Struggled"
	case Good:
		return "Remembered"
	case Easy:
		return "Mastered"
	default:
		return "Unknown"
	}
}

// Entry represents a topic to be reviewed.
type Entry struct {
	ID          int
	Title       string
	Description string
	Recall      Recall
	Category    string
	Interval    int
	ReviewDate  time.Time
	LastReview  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// INFO: Validation

func validateTitle(title string) error {
	if len(title) == 0 {
		return ErrTitleEmpty
	}
	if len(title) > MaxTitleLength {
		return ErrTitleTooLong
	}
	return nil
}

func validateDescription(description string) error {
	if len(description) > MaxDescriptionLength {
		return ErrDescriptionTooLong
	}
	return nil
}

func validateRecall(recall Recall) error {
	if recall < Failed || recall > Easy {
		return ErrInvalidRecall
	}
	return nil
}

func validateCategory(category string) error {
	if len(category) == 0 {
		return ErrCategoryEmpty
	}
	if len(category) > MaxCategoryLength {
		return ErrCategoryTooLong
	}
	return nil
}

// INFO: Options

// Option defines a functional option for configuring Entry creation.
type Option func(*Entry) error

// WithDescription sets Description of the Entry.
func WithDescription(description string) Option {
	return func(entry *Entry) error {
		description = strings.TrimSpace(description)
		err := validateDescription(description)
		if err != nil {
			return err
		}
		entry.Description = description
		return nil
	}
}

// WithCategory sets Category of the Entry.
func WithCategory(category string) Option {
	return func(entry *Entry) error {
		category = strings.TrimSpace(category)
		err := validateCategory(category)
		if err != nil {
			return err
		}
		entry.Category = category
		return nil
	}
}

// INFO: Constructor

// NewEntry creates a new entry with given title and options.
func NewEntry(title string, options ...Option) (*Entry, error) {
	title = strings.TrimSpace(title)
	if err := validateTitle(title); err != nil {
		return nil, err
	}

	entry := &Entry{
		ID:          0,
		Title:       title,
		Description: "",
		Recall:      -1,
		Category:    "None",
		Interval:    1,
		ReviewDate:  time.Now().Add(day).UTC(),
		LastReview:  time.Time{},
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Time{},
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(entry); err != nil {
			return nil, err
		}
	}

	return entry, nil
}

// INFO: Helper

func (e *Entry) calculateInterval() (int, error) {
	switch e.Recall {
	case NotReviewed:
		return 1, nil
	case Failed:
		return 1, nil
	case Hard:
		return int(float32(e.Interval) * float32(Hard)), nil
	case Good:
		return int(float32(e.Interval) * float32(Good)), nil
	case Easy:
		return int(float32(e.Interval) * float32(Easy)), nil
	default:
		return -1, ErrInvalidRecall
	}
}

// INFO: Update

// UpdateTitle updates the title of the entry after validating the new title.
func (e *Entry) UpdateTitle(title string) error {
	title = strings.TrimSpace(title)
	if err := validateTitle(title); err != nil {
		return err
	}

	e.Title = title
	e.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateDescription updates the description of the entry after validating the new description.
func (e *Entry) UpdateDescription(description string) error {
	description = strings.TrimSpace(description)
	if err := validateDescription(description); err != nil {
		return err
	}

	e.Description = description
	e.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateRecall updates the recall of the entry after validating the new recall.
func (e *Entry) UpdateRecall(recall Recall) error {
	if err := validateRecall(recall); err != nil {
		return err
	}

	e.Recall = recall
	e.UpdatedAt = time.Now().UTC()
	e.LastReview = time.Now().UTC()
	return nil
}

// UpdateCategory updates the category of the entry after validating the new category.
func (e *Entry) UpdateCategory(category string) error {
	category = strings.TrimSpace(category)
	if err := validateCategory(category); err != nil {
		return err
	}

	e.Category = category
	e.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateInterval updates the interval and next review date after calculating the interval day
// based on the users last recall.
func (e *Entry) UpdateInterval() error {
	newInterval, err := e.calculateInterval()
	if err != nil {
		return err
	}

	if newInterval > intervalLimit {
		newInterval = intervalLimit
	}

	e.Interval = newInterval
	e.ReviewDate = time.Now().Add(day * time.Duration(newInterval)).UTC()
	e.UpdatedAt = time.Now().UTC()
	return nil
}
