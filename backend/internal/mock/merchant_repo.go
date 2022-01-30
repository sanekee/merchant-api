package mock

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

type MerchantRepo struct {
	m   sync.Map
	err error
}

func NewMerchantRepo(retErr error, mcs []*model.Merchant) *MerchantRepo {
	repo := &MerchantRepo{err: retErr}
	for _, mc := range mcs {
		repo.m.Store(mc.Id, mc)
	}
	return repo
}

func (m *MerchantRepo) GetAll(opt model.Pagination) ([]*model.Merchant, error) {
	if m.err != nil {
		return nil, m.err
	}

	ret := make([]*model.Merchant, 0)
	if opt.Limit == 0 {
		return ret, nil
	}
	m.m.Range(func(k, v interface{}) bool {
		ret = append(ret, v.(*model.Merchant))
		return true
	})
	start := opt.Offset
	end := start + opt.Limit - 1
	if start > len(ret)-1 {
		return []*model.Merchant{}, nil
	}
	if end > len(ret)-1 {
		end = len(ret) - 1
	}
	return ret[start:end], nil
}

func (m *MerchantRepo) Get(id string) (*model.Merchant, error) {
	if m.err != nil {
		return nil, m.err
	}

	if mc, ok := m.m.Load(id); ok {
		return mc.(*model.Merchant), nil
	}
	return nil, model.ErrNoResults
}

func (m *MerchantRepo) Update(id string, umc *model.UpdateMerchant) (*model.Merchant, error) {
	if m.err != nil {
		return nil, m.err
	}

	mct, ok := m.m.Load(id)
	if !ok {
		return nil, model.ErrNoResults
	}

	mc := mct.(*model.Merchant)
	mc.Code = umc.Code
	m.m.Store(id, mc)
	return mc, nil
}

func (m *MerchantRepo) Insert(nmc *model.NewMerchant) (*model.Merchant, error) {
	if m.err != nil {
		return nil, m.err
	}

	now := time.Now().UTC()
	mc := &model.Merchant{
		Id:        uuid.NewString(),
		Code:      nmc.Code,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	m.m.Store(mc.Id, nmc)
	return mc, nil
}

func (m *MerchantRepo) Delete(id string) error {
	if m.err != nil {
		return m.err
	}

	if _, ok := m.m.Load(id); !ok {
		return model.ErrNoResults
	}
	m.m.Delete(id)
	return nil
}

func GenerateMerchants(num int) []*model.Merchant {
	idFormat := "test-%d"
	codeFormat := "code-%d"
	ret := make([]*model.Merchant, num)
	for i := 0; i < num; i++ {
		ret[i] = &model.Merchant{
			Id:   fmt.Sprintf(idFormat, i),
			Code: fmt.Sprintf(codeFormat, i),
		}
	}
	return ret
}
