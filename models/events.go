package models

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var event_category = []string{"ART & CULTURE", "PARTY", "EDUCATIONAL"}
var event_duration = []string{"ONE-TIME", "RECURRING"}
var event_types = []string{"ONLINE", "LIVE"}

type Event struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	OrganiserId  string    `gorm:"type:varchar(50)" json:"organiser_id"`
	OrganisaerId string    `gorm:"type:varchar(50)" json:"organisaer_id"`
	IsBusiness   bool      `gorm:"type:boolean;default:false" json:"is_business"`
	Name         string    `gorm:"type:varchar(50);not null" json:"name"`
	Image        string    `gorm:"type:varchar(150)" json:"image"`
	EventType    string    `gorm:"type:varchar(50);default: ONLINE" json:"event_type"`
	Reference    string    `gorm:"type:varchar(50);unique" json:"reference"`
	Country      string    `gorm:"type:varchar(50);default: NGA" json:"country"`
	Description  string    `gorm:"type:varchar(50)" json:"description"`
	Addresss	JSONBField `gorm:"type:jsonb;default: '{}'" json:"addresss"`
	Address      string    `gorm:"type:varchar(50)" json:"address"`
	Category     string    `gorm:"type:varchar(50);default: ART-AND-CULTURE;" json:"category"`
	Duration     string    `gorm:"type:varchar(50);default: ONE-TIME;" json:"duration"`
	StartTime    time.Time `gorm:"type:timestamp with time zone;not null" json:"start_time"`
	EndTime      time.Time `gorm:"type:timestamp with time zone;not null" json:"end_time"`
	CreatedAt    time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

func (event *Event) BeforeCreate(*gorm.DB) (err error) {

	if len(event.Category) > 0 {
		containsCategory := slices.Contains(event_category, strings.ToUpper(event.Category))
		if !containsCategory {
			event.Category = "ART-AND-CULTURE"
		} else {
			event.Category = strings.ToUpper(event.Category)
		}
	}

	if len(event.Duration) > 0 {
		containsDuration := slices.Contains(event_duration, strings.ToUpper(event.Duration))
		if !containsDuration {
			event.Duration = "ONE-TIME"
		} else {
			event.Duration = strings.ToUpper(event.Duration)
		}
	}

	if len(event.EventType) > 0 {
		containsType := slices.Contains(event_types, strings.ToUpper(event.EventType))
		if !containsType {
			event.EventType = "ONLINE"
		} else {
			event.EventType = strings.ToUpper(event.EventType)
		}
	} else {
		event.EventType = "ONLINE"
	}

	if len(event.OrganiserId) == 0 {
		return errors.New("organiser id cannot be empty")
	}

	if len(event.Reference) == 0 {
		return errors.New("reference cannot be empty")
	}
	if strings.ToUpper(event.EventType) != "ONLINE" {
		if len(event.Address) == 0 {
			return errors.New("address cannot be empty")
		}
	}

	event.Id = uuid.New()
	return
}
