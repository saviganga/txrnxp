package db_utils

import (
	"txrnxp/initialisers"
	"txrnxp/models"
)

func UpdateUserTicket(userTicket *models.UserTicket) (bool, string) {
	db := initialisers.ConnectDb().Db
	err := db.Save(&models.UserTicket{
		Id:            userTicket.Id,
		UserId:        userTicket.UserId,
		EventId:       userTicket.EventId,
		EventTicketId: userTicket.EventTicketId,
		Reference:     userTicket.Reference,
		Barcodee:      userTicket.Barcodee,
		Count:         userTicket.Count,
		IsValidated:   userTicket.IsValidated,
		CreatedAt:     userTicket.CreatedAt,
	}).Error
	if err != nil {
		return false, "unable to update user ticket"
	}
	return true, "successfully updated user ticket"

}

func UpdateEventTicket(eventTicket *models.EventTicket) (bool, string) {
	db := initialisers.ConnectDb().Db
	err := db.Save(&models.EventTicket{
		Id:             eventTicket.Id,
		EventId:        eventTicket.EventId,
		IsPaid:         eventTicket.IsPaid,
		IsInviteOnly:   eventTicket.IsInviteOnly,
		Reference:      eventTicket.Reference,
		TicketType:     eventTicket.TicketType,
		Description:    eventTicket.Description,
		PurchaseLimit:  eventTicket.PurchaseLimit,
		IsLimitedStock: eventTicket.IsLimitedStock,
		StockNumber:    eventTicket.StockNumber,
		Perks:          eventTicket.Perks,
		Price:          eventTicket.Price,
		SoldTickets:    eventTicket.SoldTickets,
		CreatedAt:      eventTicket.CreatedAt,
	}).Error
	if err != nil {
		return false, "unable to update event ticket"
	}
	return true, "successfully updated event ticket"

}
