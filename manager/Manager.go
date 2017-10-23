package manager

import (
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/trusch/aliasd/storage"
	kvstorage "github.com/trusch/storage"
	"github.com/trusch/storage/engines/meta"
)

// Manager manages alias mappings and persists them to db
type Manager struct {
	db     kvstorage.Storage
	scopes map[string]*storage.Storage
	mux    *mux.Router
}

// New creates a new manager instance for a given database uri
func New(dbURI string) (*Manager, error) {
	log.Info("creating database connection to " + dbURI)
	db, err := meta.NewStorage(dbURI, nil)
	if err != nil {
		return nil, err
	}
	log.Info("successfully created database connection")
	mgr := &Manager{db, make(map[string]*storage.Storage), mux.NewRouter()}
	mgr.mux.Path("/alias/{scope}").HandlerFunc(mgr.handleGetFullAliasList)
	mgr.mux.Path("/alias/{scope}/{key}").Methods("GET").HandlerFunc(mgr.handleGetAlias)
	mgr.mux.Path("/alias/{scope}/{key}").Methods("POST").HandlerFunc(mgr.handleCreateAlias)
	mgr.mux.Path("/alias/{scope}/{key}").Methods("PUT").HandlerFunc(mgr.handlePutAlias)
	mgr.mux.Path("/alias/{scope}/{key}").Methods("DELETE").HandlerFunc(mgr.handleDeleteAlias)
	log.Info("successfully initialized alias manager")
	return mgr, nil
}

// Put saves a given alias to a given key
func (manager *Manager) Put(scope, key, val string) error {
	log.Infof("put alias '%v' ▶ '%v' (scope: %v)", key, val, scope)
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

// Get returns the alias for a given key
func (manager *Manager) Get(scope, key string) (string, error) {
	log.Infof("get alias '%v' (scope: %v)", key, scope)
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

// GetAll returns all aliases for a given scope
func (manager *Manager) GetAll(scope string) (map[string]string, error) {
	log.Infof("get all aliases (scope: %v)", scope)
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

// Del removes an alias mapping
func (manager *Manager) Del(scope, key string) error {
	log.Infof("delete alias for '%v' (scope: %v)", key, scope)
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

// Has only checks if a alias exists
func (manager *Manager) Has(scope, key string) bool {
	log.Infof("check for alias of '%v' (scope: %v)", key, scope)
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

// Close closes the database connection (and flushes before ;))
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
