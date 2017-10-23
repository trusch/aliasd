package manager

/**
 * This file contains the HTTP handlers which maps from HTTP
 * to the public member functions of the manager
 */

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trusch/aliasd/generator"
	"github.com/trusch/aliasd/storage"
)

func (manager *Manager) handleGetFullAliasList(w http.ResponseWriter, r *http.Request) {
	scope := mux.Vars(r)["scope"]
	encoder := json.NewEncoder(w)
	obj, err := manager.GetAll(scope)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encoder.Encode(obj)
}

func (manager *Manager) handleGetAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	key := vars["key"]
	val, err := manager.Get(scope, key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(val))
}

func (manager *Manager) handleCreateAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	key := vars["key"]
	alias := ""
	store, ok := manager.scopes[scope]
	if !ok {
		s, err := storage.NewStorage(manager.db, scope)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		store = s
		manager.scopes[scope] = store
	}
	all := store.GetAll()
	for {
		alias = generator.Generate()
		if !checkIfValueExists(all, alias) {
			break
		}
	}
	if err := manager.Put(scope, key, string(alias)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(alias))
}

func (manager *Manager) handlePutAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	key := vars["key"]
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if err := manager.Put(scope, key, string(bs)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.Write(bs)
}

func (manager *Manager) handleDeleteAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scope := vars["scope"]
	key := vars["key"]
	if !manager.Has(scope, key) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("no such key"))
		return
	}
	if err := manager.Del(scope, key); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}
