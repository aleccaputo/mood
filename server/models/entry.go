package models

import (
	"github.com/google/uuid"
	"time"
)

type OverallMood int
type MoodDescriptor int

const (
	Excellent OverallMood = iota + 1
	Good
	Okay
	Bad
	Awful
)

const (
	Anger MoodDescriptor = iota + 1
	Anxiety
	Fear
	Depressed
	Sad
	Lonely
	Happy
	Motivated
	Calm
	Tired
	Excited
)

type Entry struct {
	Id          uuid.UUID        `json:"id" db:"id"`
	Overall     OverallMood      `json:"overall" db:"overall"`
	Descriptors []MoodDescriptor `json:"descriptors" db:"descriptors"`
	GoodNotes   string           `json:"goodNotes" db:"good_notes"`
	BadNotes    string           `json:"badNotes" db:"bad_notes"`
	Exercise    bool             `json:"exercise" db:"exercise"`
	Alcohol     bool             `json:"alcohol" db:"alcohol"`
	UserId      uuid.UUID        `json:"userId" db:"user_id"`
	CreatedAt   time.Time        `json:"createdAt " db:"created_at"`
}

type CreateEntryRequest struct {
	Overall     OverallMood      `json:"overall"`
	Descriptors []MoodDescriptor `json:"descriptors"`
	GoodNotes   string           `json:"goodNotes"`
	BadNotes    string           `json:"badNotes"`
	Exercise    bool             `json:"exercise"`
	Alcohol     bool             `json:"alcohol"`
	UserId      string           `json:"userId"`
}

func NewEntry(mood OverallMood, moodDescriptors []MoodDescriptor, goodNotes string, badNotes string, exercise bool, alcohol bool, userId uuid.UUID) *Entry {
	return &Entry{
		Id:          uuid.New(),
		Overall:     mood,
		Descriptors: moodDescriptors,
		GoodNotes:   goodNotes,
		BadNotes:    badNotes,
		Exercise:    exercise,
		Alcohol:     alcohol,
		CreatedAt:   time.Now().UTC(),
		UserId:      userId,
	}
}
