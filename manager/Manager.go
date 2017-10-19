package manager

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/trusch/aliasd/storage"
	kvstorage "github.com/trusch/storage"
	"github.com/trusch/storage/engines/meta"
)

type Manager struct {
	db     kvstorage.Storage
	scopes map[string]*storage.Storage
	mux    *mux.Router
}

func New(dbURI string) (*Manager, error) {
	db, err := meta.NewStorage(dbURI, nil)
	if err != nil {
		return nil, err
	}
	mgr := &Manager{db, make(map[string]*storage.Storage), mux.NewRouter()}
	mgr.mux.Path("/alias/{scope}").HandlerFunc(mgr.handleGetFullAliasList)
	mgr.mux.Path("/alias/{scope}/{key}").Methods("GET").HandlerFunc(mgr.handleGetAlias)
	mgr.mux.Path("/alias/{scope}/{key}").Methods("POST").HandlerFunc(mgr.handleCreateAlias)
	mgr.mux.Path("/alias/{scope}/{key}").Methods("PUT").HandlerFunc(mgr.handlePutAlias)
	mgr.mux.Path("/alias/{scope}/{key}").Methods("DELETE").HandlerFunc(mgr.handleDeleteAlias)
	return mgr, err
}

func (manager *Manager) Put(scope, key, val string) error {
	store, ok := manager.scopes[scope]
	if !ok {
		s, err := storage.NewStorage(manager.db, scope)
		if err != nil {
			return err
		}
		store = s
		manager.scopes[scope] = store
	}
	return store.Set(key, val)
}

func (manager *Manager) Get(scope, key string) (string, error) {
	store, ok := manager.scopes[scope]
	if !ok {
		s, err := storage.NewStorage(manager.db, scope)
		if err != nil {
			return "", err
		}
		store = s
		manager.scopes[scope] = store
	}
	return store.Get(key)
}

func (manager *Manager) GetAll(scope string) (map[string]string, error) {
	store, ok := manager.scopes[scope]
	if !ok {
		s, err := storage.NewStorage(manager.db, scope)
		if err != nil {
			return nil, err
		}
		store = s
		manager.scopes[scope] = store
	}
	return store.GetAll(), nil
}

func (manager *Manager) Del(scope, key string) error {
	store, ok := manager.scopes[scope]
	if !ok {
		s, err := storage.NewStorage(manager.db, scope)
		if err != nil {
			return err
		}
		store = s
		manager.scopes[scope] = store
	}
	return store.Del(key)
}

func (manager *Manager) Has(scope, key string) bool {
	store, ok := manager.scopes[scope]
	if !ok {
		s, err := storage.NewStorage(manager.db, scope)
		if err != nil {
			return false
		}
		store = s
		manager.scopes[scope] = store
	}
	return store.Has(key)
}

func (manager *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	manager.mux.ServeHTTP(w, r)
}

func (manager *Manager) Close() error {
	return manager.db.Close()
}

func getVars(r *http.Request) (scope, key string, err error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		err = errors.New("malformed URL")
		return
	}
	if len(parts) == 4 {
		scope = parts[2]
		key = parts[3]
	}
	if len(parts) == 3 {
		scope = "default"
		key = parts[2]
	}
	return
}

func checkIfValueExists(m map[string]string, val string) bool {
	for _, v := range m {
		if v == val {
			return true
		}
	}
	return false
}
