package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/sanekee/merchant-api/internal/model"
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

func CreateTeamMemberFromNewSchema(from *model.NewTeamMember) *TeamMember {
	now := time.Now().UTC()
	return &TeamMember{
		ID:         uuid.NewString(),
		Email:      from.Email,
		MerchantID: from.MerchantId,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
}

func TeamMemberFromUpdateSchema(id string, from *model.UpdateTeamMember) *TeamMember {
	now := time.Now().UTC()
	return &TeamMember{
		ID:        id,
		Email:     from.Email,
		UpdatedAt: &now,
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
