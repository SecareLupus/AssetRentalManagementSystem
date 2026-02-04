package fleet

import (
	"fmt"
	"sync"

	"github.com/desmond/rental-management-system/internal/domain"
)

type RemoteRegistry struct {
	mu        sync.RWMutex
	providers map[string]domain.RemoteManager
}

func NewRemoteRegistry() *RemoteRegistry {
	return &RemoteRegistry{
		providers: make(map[string]domain.RemoteManager),
	}
}

func (r *RemoteRegistry) Register(name string, manager domain.RemoteManager) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = manager
}

func (r *RemoteRegistry) Get(name string) (domain.RemoteManager, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	mgr, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("remote manager provider not found: %s", name)
	}
	return mgr, nil
}
