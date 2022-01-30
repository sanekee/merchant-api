package repo

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/sanekee/merchant-api/backend/internal/entity"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

type TeamMemberRepo struct {
	db *pg.DB
}

func NewTeamMemberRepo(db *pg.DB) *TeamMemberRepo {
	return &TeamMemberRepo{db: db}
}

func (t *TeamMemberRepo) GetByMerchantID(merchant_id string, page model.Pagination) ([]*model.TeamMember, error) {
	var tms []*entity.TeamMember
	err := t.db.Model(&tms).
		Limit(page.Limit).
		Offset(page.Offset).
		Where("merchant_id = ?", merchant_id).
		Select()
	if err != nil {
		return nil, err
	}

	ret := make([]*model.TeamMember, len(tms))
	for i, tm := range tms {
		ret[i] = tm.ToSchema()
	}
	return ret, nil
}

func (t *TeamMemberRepo) Get(id string) (*model.TeamMember, error) {
	tm := entity.TeamMember{
		ID: id,
	}

	err := t.db.Model(&t).WherePK().Select()
	if err != nil {
		return nil, toAppError(err)
	}
	return tm.ToSchema(), nil
}

func (t *TeamMemberRepo) Insert(tm *model.TeamMember) (*model.TeamMember, error) {
	ent := entity.TeamMemberFromSchema(tm)
	now := time.Now()
	ent.UpdatedAt = &now
	_, err := t.db.Model(ent).
		Returning("*").
		Insert()
	if err != nil {
		return nil, toAppError(err)
	}
	return ent.ToSchema(), nil
}

func (t *TeamMemberRepo) Update(mc *model.TeamMember) (*model.TeamMember, error) {
	ent := entity.TeamMemberFromSchema(mc)
	now := time.Now()
	ent.UpdatedAt = &now
	_, err := t.db.Model(ent).
		WherePK().
		Returning("*").
		Update()
	if err != nil {
		return nil, toAppError(err)
	}
	return ent.ToSchema(), nil
}

func (t *TeamMemberRepo) Delete(id string) error {
	_, err := t.db.Model(&entity.TeamMember{ID: id}).WherePK().Delete()
	return toAppError(err)
}
