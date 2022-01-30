package entity

import (
	"time"

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

func (t *Merchant) ToSchema() *model.Merchant {
	return &model.Merchant{
		Id:        t.ID,
		Code:      t.Code,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
