package repo

import (
	"github.com/go-pg/pg/v10"
	"github.com/sanekee/merchant-api/internal/entity"
	"github.com/sanekee/merchant-api/internal/model"
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
		Order("created_at").
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

	if err := t.db.Model(&tm).WherePK().Select(); err != nil {
		return nil, toAppError(err)
	}

	return tm.ToSchema(), nil
}

func (t *TeamMemberRepo) Insert(ntm *model.NewTeamMember) (*model.TeamMember, error) {
	ent := entity.CreateTeamMemberFromNewSchema(ntm)

	_, err := t.db.Model(ent).
		Returning("*").
		Insert()
	if err != nil {
		return nil, toAppError(err)
	}

	return ent.ToSchema(), nil
}

func (t *TeamMemberRepo) Update(id string, umc *model.UpdateTeamMember) (*model.TeamMember, error) {
	ent := entity.TeamMemberFromUpdateSchema(id, umc)

	_, err := t.db.Model(ent).
		WherePK().
		Returning("*").
		UpdateNotZero()
	if err != nil {
		return nil, toAppError(err)
	}

	return ent.ToSchema(), nil
}

func (t *TeamMemberRepo) Delete(id string) error {
	res, err := t.db.Model(&entity.TeamMember{ID: id}).WherePK().Delete()
	if res != nil && res.RowsAffected() == 0 {
		return model.ErrNoResults
	}

	return toAppError(err)
}
