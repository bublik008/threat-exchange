package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/bublik008/threat-exchange/internal/api/gen"
	"github.com/bublik008/threat-exchange/internal/repository"
)

func TestCreateThreat_InvalidMD5(t *testing.T) {
	mockRepo := new(MockThreatRepository)
	svc := NewThreatService(mockRepo)

	req := &gen.ThreatCreateRequest{
		Md5:      "invalid-md5",
		Name:     "Test",
		Severity: "high",
	}

	_, err := svc.CreateThreat(context.Background(), req)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidMD5)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateThreat_Success(t *testing.T) {
	mockRepo := new(MockThreatRepository)
	svc := NewThreatService(mockRepo)

	req := &gen.ThreatCreateRequest{
		Md5:      "D41D8CD98F00B204E9800998ECF8427E",
		Name:     "Test Trojan",
		Severity: "high",
	}

	md5Str := "d41d8cd98f00b204e9800998ecf8427e"
	nameStr := "Test Trojan"
	severityEnum := gen.ThreatResponseSeverityHigh

	expectedResp := &gen.ThreatResponse{
		Md5:      &md5Str,
		Name:     &nameStr,
		Severity: &severityEnum,
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(r *gen.ThreatCreateRequest) bool {
		return r.Md5 == "d41d8cd98f00b204e9800998ecf8427e"
	}), mock.Anything).Return(expectedResp, nil)

	resp, err := svc.CreateThreat(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", *resp.Md5)
	mockRepo.AssertExpectations(t)
}

func TestGetThreat_NotFound(t *testing.T) {
	mockRepo := new(MockThreatRepository)
	svc := NewThreatService(mockRepo)

	md5 := "d41d8cd98f00b204e9800998ecf8427e"
	mockRepo.On("GetByMD5", mock.Anything, md5).Return(nil, repository.ErrNotFound)

	_, err := svc.GetThreat(context.Background(), md5)

	require.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)
	mockRepo.AssertExpectations(t)
}

func TestUpdateThreat_EmptyName(t *testing.T) {
	mockRepo := new(MockThreatRepository)
	svc := NewThreatService(mockRepo)

	md5 := "d41d8cd98f00b204e9800998ecf8427e"
	req := &gen.ThreatUpdateRequest{
		Name: "   ",
	}

	_, err := svc.UpdateThreat(context.Background(), md5, req)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidInput)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}