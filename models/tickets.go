package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JSONBField map[string]interface{}

func (fld JSONBField) Value() (driver.Value, error) {
	j, err := json.Marshal(fld)
	return j, err
}

func (fld *JSONBField) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*fld, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed")
	}

	return nil
}

type EventTicket struct {
	Id             uuid.UUID  `gorm:"type:uuid;primaryKey;not null" json:"id"`
	EventId        uuid.UUID  `gorm:"type:uuid;not null" json:"event_id"`
	IsPaid         bool       `gorm:"type:boolean;default:true" json:"is_paid"`
	IsInviteOnly   bool       `gorm:"type:boolean;default:false" json:"is_invite_only"`
	Reference      string     `gorm:"type:varchar(50);unique" json:"reference"`
	TicketType     string     `gorm:"type:varchar(50);default: SINGLE" json:"ticket_type"`
	Description    string     `gorm:"type:varchar(50)" json:"description"`
	PurchaseLimit  int        `gorm:"type:int;default: 0" json:"purchase_limit"`
	IsLimitedStock bool       `gorm:"type:boolean;default:false" json:"is_limited_stock"`
	StockNumber    int        `gorm:"type:int" json:"stock_number"`
	Perks          JSONBField `gorm:"type:jsonb;default: '{}'" json:"perks"`
	Price          string     `gorm:"type:numeric(10,2);not null;default:0.00" json:"price"`
	SoldTickets    int        `gorm:"type:int;default: 0" json:"sold_tickets"`
	CreatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"type:timestamp with time zone" json:"updated_at"`
	Event          Event      `gorm:"foreignKey:EventId;constraint:OnDelete:CASCADE"`
}

func (event_ticket *EventTicket) BeforeCreate(*gorm.DB) (err error) {

	event_ticket.Id = uuid.New()

	// check is_paid
	if event_ticket.IsPaid {
		if event_ticket.Price == "" {
			return errors.New("price is compulsory for paid tickets")
		} else {
			event_str_price := event_ticket.Price
			event_float_price, err := strconv.ParseFloat(event_str_price, 64)
			if err != nil {
				return errors.New("error converting event ticket price")
			}
			if event_float_price <= float64(0) {
				return errors.New("paid ticket price must be greater than 0")
			}
		}
	}

	// stock
	if event_ticket.IsLimitedStock {
		if event_ticket.StockNumber <= 0 {
			return errors.New("limited stock ticket stock number must be greater than 0")
		}
	}

	return
}

type UserTicket struct {
	Id            uuid.UUID   `gorm:"type:uuid;primaryKey;not null" json:"id"`
	EventId       uuid.UUID   `gorm:"type:uuid;not null" json:"event_id"`
	EventTicketId uuid.UUID   `gorm:"type:uuid;not null" json:"event_ticket_id"`
	UserId        uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	Reference     string      `gorm:"type:varchar(50);unique" json:"reference"`
	Barcodee      JSONBField  `gorm:"type:jsonb;default: '{}'" json:"barcodee"`
	Count         int         `gorm:"type:int;default:0" json:"count"`
	IsValidated   bool        `gorm:"type:boolean;default:false" json:"is_validated"`
	CreatedAt     time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"type:timestamp with time zone" json:"updated_at"`
	Event         Event       `gorm:"foreignKey:EventId;constraint:OnDelete:CASCADE"`
	EventTicket   EventTicket `gorm:"foreignKey:EventTicketId;constraint:OnDelete:CASCADE"`
	User          Xuser       `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func (user_ticket *UserTicket) BeforeCreate(*gorm.DB) (err error) {

	user_ticket.Id = uuid.New()
	return
}
