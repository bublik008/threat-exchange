// Package service содержит бизнес-логику приложения.
package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bublik008/threat-exchange/internal/api/gen"      
	"github.com/bublik008/threat-exchange/internal/repository"
)

// Доменные ошибки уровня бизнес-логики.
var (
	// ErrInvalidMD5 возвращается, если MD5 не соответствует формату.
	ErrInvalidMD5 = errors.New("invalid md5 format")
	// ErrInvalidInput возвращается при некорректных бизнес-данных
	ErrInvalidInput = errors.New("invalid input data")
)

// Регулярка для проверки MD5 (32 шестнадцатеричных символа).
var md5Regex = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)

// ThreatService реализует бизнес-логику работы с угрозами.
type ThreatService struct {
	repo repository.ThreatRepository
}

// NewThreatService создает новый сервис.
func NewThreatService(repo repository.ThreatRepository) *ThreatService {
	return &ThreatService{repo: repo}
}

// CreateThreat создает новую запись об угрозе.
// Валидирует MD5, приводит его к нижнему регистру (канонический вид) и вызывает репозиторий.
func (s *ThreatService) CreateThreat(ctx context.Context, req *gen.ThreatCreateRequest) (*gen.ThreatResponse, error) {
	// 1. Валидация MD5
	if !md5Regex.MatchString(req.Md5) {
		return nil, fmt.Errorf("%w: must be 32 hex characters", ErrInvalidMD5)
	}
	
	req.Md5 = strings.ToLower(req.Md5)

	// 2. Базовая валидация имени 
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("%w: name cannot be empty", ErrInvalidInput)
	}

	// 3. Передаем в репозиторий, добавляя текущее время
	now := time.Now().UTC()
	return s.repo.Create(ctx, req, now)
}

// GetThreat получает информацию об угрозе по MD5.
func (s *ThreatService) GetThreat(ctx context.Context, md5 string) (*gen.ThreatResponse, error) {
	if !md5Regex.MatchString(md5) {
		return nil, fmt.Errorf("%w: must be 32 hex characters", ErrInvalidMD5)
	}
	md5 = strings.ToLower(md5)

	return s.repo.GetByMD5(ctx, md5)
}

// UpdateThreat обновляет информацию об угрозе.
func (s *ThreatService) UpdateThreat(ctx context.Context, md5 string, req *gen.ThreatUpdateRequest) (*gen.ThreatResponse, error) {
	if !md5Regex.MatchString(md5) {
		return nil, fmt.Errorf("%w: must be 32 hex characters", ErrInvalidMD5)
	}
	md5 = strings.ToLower(md5)

	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("%w: name cannot be empty", ErrInvalidInput)
	}

	return s.repo.Update(ctx, md5, req)
}