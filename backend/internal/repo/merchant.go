package repo

import (
	"github.com/go-pg/pg/v10"
	"github.com/sanekee/merchant-api/internal/entity"
	"github.com/sanekee/merchant-api/internal/model"
)

type MerchantRepo struct {
	db *pg.DB
}

func NewMerchantRepo(db *pg.DB) *MerchantRepo {
	return &MerchantRepo{db: db}
}

func (m *MerchantRepo) GetAll(pg model.Pagination) ([]*model.Merchant, error) {
	var mcs []*entity.Merchant

	err := m.db.Model(&mcs).
		Limit(pg.Limit).
		Offset(pg.Offset).
		Order("created_at").
		Select()

	if err != nil {
		return nil, err
	}

	ret := make([]*model.Merchant, len(mcs))
	for i, mc := range mcs {
		ret[i] = mc.ToSchema()
	}
	return ret, nil
}

func (m *MerchantRepo) Get(id string) (*model.Merchant, error) {
	mc := entity.Merchant{
		ID: id,
	}

	if err := m.db.Model(&mc).WherePK().Select(); err != nil {
		return nil, toAppError(err)
	}

	return mc.ToSchema(), nil
}

func (m *MerchantRepo) Insert(nmc *model.NewMerchant) (*model.Merchant, error) {
	ent := entity.CreateMerchantFromNewSchema(nmc)

	_, err := m.db.Model(ent).
		Returning("*").
		Insert()
	if err != nil {
		return nil, toAppError(err)
	}

	return ent.ToSchema(), nil
}

func (m *MerchantRepo) Update(id string, umc *model.UpdateMerchant) (*model.Merchant, error) {
	ent := entity.MerchantFromUpdateSchema(id, umc)

	_, err := m.db.Model(ent).
		WherePK().
		Returning("*").
		UpdateNotZero()
	if err != nil {
		return nil, toAppError(err)
	}

	return ent.ToSchema(), nil
}

func (m *MerchantRepo) Delete(id string) error {
	res, err := m.db.Model(&entity.Merchant{ID: id}).WherePK().Delete()
	if res != nil && res.RowsAffected() == 0 {
		return model.ErrNoResults
	}

	return toAppError(err)
}
