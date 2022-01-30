package entity

import (
	"time"

	"github.com/sanekee/merchant-api/backend/internal/model"
)

type TeamMember struct {
	tableName  struct{}   `pg:"mc_team_member"`
	ID         string     `pg:",pk"`
	Email      string     `pg:",unique"`
	MerchantID string     `pg:"merchant_id"`
	CreatedAt  *time.Time `pg:"created_at,default:now()"`
	UpdatedAt  *time.Time `pg:"updated_at,default:now()"`
}

func TeamMemberFromSchema(from *model.TeamMember) *TeamMember {
	return &TeamMember{
		ID:         from.Id,
		MerchantID: from.MerchantId,
		Email:      from.Email,
		CreatedAt:  from.CreatedAt,
		UpdatedAt:  from.UpdatedAt,
	}
}

func (t *TeamMember) ToSchema() *model.TeamMember {
	return &model.TeamMember{
		Id:         t.ID,
		MerchantId: t.MerchantID,
		Email:      t.Email,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}
