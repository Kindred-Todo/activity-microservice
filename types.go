package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
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
	TimeTaken  *ISO8601Duration   `bson:"timeTaken,omitempty" json:"timeTaken,omitempty"`
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

// ISO8601Duration represents an ISO 8601 duration (e.g., "PT0S", "PT1H30M")
type ISO8601Duration struct {
	time.Duration
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface
func (d *ISO8601Duration) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("invalid bson type for ISO8601Duration, got %v", t)
	}
	
	str, _, ok := bsoncore.ReadString(data)
	if !ok {
		return fmt.Errorf("failed to read string from bson data")
	}
	
	duration, err := parseISO8601Duration(str)
	if err != nil {
		return err
	}
	
	d.Duration = duration
	return nil
}

// parseISO8601Duration parses an ISO 8601 duration string (e.g., "PT0S", "PT1H30M45S")
func parseISO8601Duration(s string) (time.Duration, error) {
	// ISO 8601 duration format: P[n]Y[n]M[n]DT[n]H[n]M[n]S
	// We'll support the time portion (after T) for now
	re := regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)
	matches := re.FindStringSubmatch(s)
	
	if matches == nil {
		return 0, fmt.Errorf("invalid ISO 8601 duration format: %s", s)
	}
	
	var duration time.Duration
	
	// Parse hours
	if matches[4] != "" {
		hours, _ := strconv.Atoi(matches[4])
		duration += time.Duration(hours) * time.Hour
	}
	
	// Parse minutes
	if matches[5] != "" {
		minutes, _ := strconv.Atoi(matches[5])
		duration += time.Duration(minutes) * time.Minute
	}
	
	// Parse seconds
	if matches[6] != "" {
		seconds, _ := strconv.ParseFloat(matches[6], 64)
		duration += time.Duration(seconds * float64(time.Second))
	}
	
	return duration, nil
}