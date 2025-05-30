package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompletedTaskDocument struct {
	ID           primitive.ObjectID     `bson:"_id" json:"id"`
	Priority     int                    `bson:"priority" json:"priority"`
	Content      string                 `bson:"content" json:"content"`
	Value        float64                `bson:"value" json:"value"`
	Recurring    bool                   `bson:"recurring" json:"recurring"`
	RecurType    string                 `bson:"recurType,omitempty" json:"recurType,omitempty"`
	RecurDetails *RecurDetails           `bson:"recurDetails,omitempty" json:"recurDetails,omitempty"`
	Public       bool                   `bson:"public" json:"public"`
	Active       bool                   `bson:"active" json:"active"`
	Timestamp    time.Time              `bson:"timestamp" json:"timestamp"`
	LastEdited   time.Time              `bson:"lastEdited" json:"lastEdited"`
	TemplateID    primitive.ObjectID     `bson:"templateID,omitempty" json:"templateID,omitempty"`

	

	Deadline *time.Time `bson:"deadline,omitempty" json:"deadline,omitempty"`
	StartTime *time.Time `bson:"startTime,omitempty" json:"startTime,omitempty"`
	StartDate *time.Time `bson:"startDate" json:"startDate"` // Defaults to today

	Notes        string                 `bson:"notes,omitempty" json:"notes,omitempty"`
	Checklist    []ChecklistItem        `bson:"checklist,omitempty" json:"checklist,omitempty"`

	CategoryID primitive.ObjectID `bson:"category,omitempty" json:"category,omitempty"`
	UserID     primitive.ObjectID `bson:"user,omitempty" json:"user,omitempty"`
	TimeTaken  *time.Time         `bson:"timeTaken,omitempty" json:"timeTaken,omitempty"`
	TimeCompleted *time.Time `bson:"timeCompleted,omitempty" json:"timeCompleted,omitempty"`
}

type RecurDetails struct {
	Every int `validate:"required,min=1" bson:"every,omitempty" json:"every,omitempty"`
	DaysOfWeek []int `validate:"omitempty,min=7,max=7" bson:"daysOfWeek,omitempty" json:"daysOfWeek,omitempty"`
	DaysOfMonth []int `validate:"omitempty,min=1,max=31,unique" bson:"daysOfMonth,omitempty" json:"daysOfMonth,omitempty"`
	Months []int `validate:"omitempty,min=1,max=12,unique" bson:"months,omitempty" json:"months,omitempty"`
	Behavior string `validate:"required,oneof=BUILDUP ROLLING" bson:"behavior,omitempty" json:"behavior,omitempty"` // Buildup, Rolling
}

type ChecklistItem struct {
	Content   string `bson:"content" json:"content"`
	Completed bool   `bson:"completed" json:"completed"`
	Order     int    `bson:"order" json:"order"`
}