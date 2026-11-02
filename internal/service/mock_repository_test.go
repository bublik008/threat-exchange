package service

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/bublik008/threat-exchange/internal/api/gen"
)


type MockThreatRepository struct {
	mock.Mock
}

func (m *MockThreatRepository) Create(ctx context.Context, threat *gen.ThreatCreateRequest, now time.Time) (*gen.ThreatResponse, error) {
	args := m.Called(ctx, threat, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gen.ThreatResponse), args.Error(1)
}

func (m *MockThreatRepository) GetByMD5(ctx context.Context, md5 string) (*gen.ThreatResponse, error) {
	args := m.Called(ctx, md5)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gen.ThreatResponse), args.Error(1)
}

func (m *MockThreatRepository) Update(ctx context.Context, md5 string, update *gen.ThreatUpdateRequest) (*gen.ThreatResponse, error) {
	args := m.Called(ctx, md5, update)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gen.ThreatResponse), args.Error(1)
}