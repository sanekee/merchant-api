package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

type Merchant struct {
	tableName struct{}   `pg:"mc_merchant"`
	ID        string     `pg:",pk"`
	Code      string     `pg:"code,unique"`
	CreatedAt *time.Time `pg:"created_at,default:now()"`
	UpdatedAt *time.Time `pg:"updated_at,default:now()"`
}

func MerchantFromSchema(from *model.Merchant) *Merchant {
	return &Merchant{
		ID:        from.Id,
		Code:      from.Code,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}

func CreateMerchantFromNewSchema(from *model.NewMerchant) *Merchant {
	now := time.Now().UTC()
	return &Merchant{
		ID:        uuid.NewString(),
		Code:      from.Code,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func MerchantFromUpdateSchema(id string, from *model.UpdateMerchant) *Merchant {
	now := time.Now().UTC()
	return &Merchant{
		ID:        id,
		Code:      from.Code,
		UpdatedAt: &now,
	}
}

func (t *Merchant) ToSchema() *model.Merchant {
	return &model.Merchant{
		Id:        t.ID,
		Code:      t.Code,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
