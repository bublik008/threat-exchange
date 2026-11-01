// Package repository описывает контракт для работы с хранилищем данных.
package repository

import (
	"context"
	"errors"
	"time"
	"github.com/bublik008/threat-exchange/internal/api/gen"
)


var (
	// ErrNotFound возвращается, если запись по MD5 не найдена.
	ErrNotFound = errors.New("threat not found")
	// ErrConflict возвращается при попытке создать дубликат MD5.
	ErrConflict = errors.New("threat with this MD5 already exists")
)

// ThreatRepository описывает операции с угрозами в базе данных.
type ThreatRepository interface {
	// Create сохраняет новую угрозу. Возвращает ErrConflict, если MD5 уже есть.
	Create(ctx context.Context, threat *gen.ThreatCreateRequest, now time.Time) (*gen.ThreatResponse, error)
	
	// GetByMD5 ищет угрозу по хешу. Возвращает ErrNotFound, если записи нет.
	GetByMD5(ctx context.Context, md5 string) (*gen.ThreatResponse, error)
	
	// Update обновляет существующую угрозу. Возвращает ErrNotFound, если записи нет.
	Update(ctx context.Context, md5 string, update *gen.ThreatUpdateRequest) (*gen.ThreatResponse, error)
}