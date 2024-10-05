package admin_serializers

import (
	"time"
	"txrnxp/models"

	"github.com/google/uuid"
)

type CreateAdminCommissionConfigSerializer struct {
	Type       string `json:"type" validate:"required"`
	Commission string `json:"commission" validate:"required"`
	Cap        string `json:"cap" validate:"required"`
}

type ReadAdminCommissionConfigSerializer struct {
	Id         uuid.UUID `json:"id" validate:"required"`
	Type       string    `json:"type" validate:"required"`
	Commission string    `json:"commission" validate:"required"`
	Cap        string    `json:"cap" validate:"required"`
	CreatedAt  time.Time `json:"created_at" validate:"required"`
	UpdatedAt  time.Time `json:"updated_at" validate:"required"`
}

func SerializeCreateAdminCommissionConfig(commission_config models.AdminCommissionConfig) ReadAdminCommissionConfigSerializer {

	serialized_commission_config := new(ReadAdminCommissionConfigSerializer)

	serialized_commission_config.Id = commission_config.Id
	serialized_commission_config.Type = commission_config.Type
	serialized_commission_config.Commission = commission_config.Commission
	serialized_commission_config.Cap = commission_config.Cap
	serialized_commission_config.CreatedAt = commission_config.CreatedAt
	serialized_commission_config.UpdatedAt = commission_config.UpdatedAt

	return *serialized_commission_config

}
