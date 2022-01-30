package mock

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

type TeamMemberRepo struct {
	m   sync.Map
	err error
}

func NewTeamMemberRepo(retErr error, mcs []*model.TeamMember) *TeamMemberRepo {
	repo := &TeamMemberRepo{err: retErr}
	for _, mc := range mcs {
		repo.m.Store(mc.Id, mc)
	}
	return repo
}

func (m *TeamMemberRepo) GetByMerchantID(merchantID string, opt model.Pagination) ([]*model.TeamMember, error) {
	if m.err != nil {
		return nil, m.err
	}

	ret := make([]*model.TeamMember, 0)
	if opt.Limit == 0 {
		return ret, nil
	}
	m.m.Range(func(k, v interface{}) bool {
		tm := v.(*model.TeamMember)
		if tm.MerchantId == merchantID {
			ret = append(ret, v.(*model.TeamMember))
		}
		return true
	})
	start := opt.Offset
	end := start + opt.Limit - 1
	if start > len(ret)-1 {
		return []*model.TeamMember{}, nil
	}
	if end > len(ret)-1 {
		end = len(ret) - 1
	}
	return ret[start:end], nil
}

func (m *TeamMemberRepo) Get(id string) (*model.TeamMember, error) {
	if m.err != nil {
		return nil, m.err
	}
	if mc, ok := m.m.Load(id); ok {
		return mc.(*model.TeamMember), nil
	}
	return nil, model.ErrNoResults
}

func (m *TeamMemberRepo) Update(id string, umc *model.UpdateTeamMember) (*model.TeamMember, error) {
	if m.err != nil {
		return nil, m.err
	}

	tmt, ok := m.m.Load(id)
	if !ok {
		return nil, model.ErrNoResults
	}
	tm := tmt.(*model.TeamMember)
	tm.Email = umc.Email
	tm.MerchantId = umc.MerchantId
	m.m.Store(id, tm)
	return tm, nil
}

func (m *TeamMemberRepo) Insert(ntm *model.NewTeamMember) (*model.TeamMember, error) {
	if m.err != nil {
		return nil, m.err
	}

	now := time.Now().UTC()
	tm := &model.TeamMember{
		Id:         uuid.NewString(),
		Email:      ntm.Email,
		MerchantId: ntm.MerchantId,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	m.m.Store(tm.Id, tm)
	return tm, nil
}

func (m *TeamMemberRepo) Delete(id string) error {
	if m.err != nil {
		return m.err
	}
	if _, ok := m.m.Load(id); !ok {
		return model.ErrNoResults
	}
	m.m.Delete(id)
	return nil
}

func GenerateTeamMembers(num int, merchantID string) []*model.TeamMember {
	idFormat := "test-%d"
	emailFormat := "email-%d@t.est"
	ret := make([]*model.TeamMember, num)
	for i := 0; i < num; i++ {
		ret[i] = &model.TeamMember{
			Id:         fmt.Sprintf(idFormat, i),
			Email:      fmt.Sprintf(emailFormat, i),
			MerchantId: merchantID,
		}
	}
	return ret
}
