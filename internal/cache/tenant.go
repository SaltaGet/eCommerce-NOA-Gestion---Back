package cache

import (
	"sync"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
)

type TenantStore struct {
	mu      sync.RWMutex
	tenants map[string]schemas.TenantResponseSetting
}

func NewTenantStore() *TenantStore {
	return &TenantStore{
		tenants: make(map[string]schemas.TenantResponseSetting),
	}
}

func (s *TenantStore) Update(newList []schemas.TenantResponseSetting) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	newMap := make(map[string]schemas.TenantResponseSetting)
	for _, t := range newList {
		newMap[t.Identifier] = t
	}
	s.tenants = newMap
}

func (s *TenantStore) Exists(identifier string) schemas.TenantResponseSetting {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tenants[identifier]
}

func (s *TenantStore) Get(identifier string) (schemas.TenantResponseSetting, bool) {
  s.mu.RLock()
  defer s.mu.RUnlock()
  t, ok := s.tenants[identifier]
  return t, ok
}