package storage

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"sync"
)

type storage struct {
	mu   sync.RWMutex
	mngr manager.BaseStorageManager
}

func Create(storageManager manager.BaseStorageManager, baseLogger logger.BaseLogger) BaseStorage {
	return &storage{
		mngr: storageManager,
		log:  baseLogger,
	}
}

func (s *storage) EditNote(ctx context.Context, note *models.News) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.mngr.EditNote(ctx, note)
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	return nil
}

func (s *storage) GetNotes(ctx context.Context, page int64, perPage int64) ([]models.News, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	notes, err := s.mngr.GetNotes(ctx, page, perPage)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	return notes, nil
}
