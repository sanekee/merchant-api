package repo

import (
	"github.com/go-pg/pg/v10"
	"github.com/sanekee/merchant-api/backend/internal/entity"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

type MerchantRepo struct {
	db *pg.DB
}

func NewMerchantRepo(db *pg.DB) *MerchantRepo {
	return &MerchantRepo{db: db}
}

func (b *MerchantRepo) GetAll(pg model.Pagination) ([]*model.Merchant, error) {
	var m []*entity.Merchant
	err := b.db.Model(&m).
		Limit(pg.Limit).
		Offset(pg.Offset).
		Order("created_at").
		Select()
	if err != nil {
		return nil, err
	}

	ret := make([]*model.Merchant, len(m))
	for i, mc := range m {
		ret[i] = mc.ToSchema()
	}
	return ret, nil
}

func (b *MerchantRepo) Get(id string) (*model.Merchant, error) {
	m := entity.Merchant{
		ID: id,
	}

	err := b.db.Model(&m).WherePK().Select()
	if err != nil {
		return nil, toAppError(err)
	}
	return m.ToSchema(), nil
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

func (b *MerchantRepo) Update(id string, umc *model.UpdateMerchant) (*model.Merchant, error) {
	ent := entity.MerchantFromUpdateSchema(id, umc)
	_, err := b.db.Model(ent).
		WherePK().
		Returning("*").
		UpdateNotZero()
	if err != nil {
		return nil, toAppError(err)
	}
	return ent.ToSchema(), nil
}

func (b *MerchantRepo) Delete(id string) error {
	_, err := b.db.Model(&entity.Merchant{ID: id}).WherePK().Delete()
	return toAppError(err)
}
