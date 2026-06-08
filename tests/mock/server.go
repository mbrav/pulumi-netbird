// Package mock provides a small in-memory NetBird management API for provider tests.
package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Server implements the NetBird REST endpoints exercised by the provider tests.
type Server struct {
	mu     sync.Mutex
	nextID int
	items  map[string]map[string]map[string]any
}

// NewServer creates a mock NetBird management API.
func NewServer() *Server {
	return &Server{items: map[string]map[string]map[string]any{}}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") || strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")

		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "api" {
		writeError(w, http.StatusNotFound, "not found")

		return
	}

	switch r.Method {
	case http.MethodPost:
		if len(parts) == 2 {
			s.create(w, parts[1], r)

			return
		}
	case http.MethodGet:
		if len(parts) == 3 {
			s.get(w, parts[1], parts[2])

			return
		}
	case http.MethodPut:
		if len(parts) == 3 {
			s.update(w, parts[1], parts[2], r)

			return
		}
	case http.MethodDelete:
		if len(parts) == 3 {
			s.delete(w, parts[1], parts[2])

			return
		}
	}

	writeError(w, http.StatusNotFound, "not found")
}

func (s *Server) create(w http.ResponseWriter, resource string, r *http.Request) {
	data, ok := readJSON(w, r)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	id := fmt.Sprintf("%s-%d", resource, s.nextID)
	data["id"] = id
	data = apiShape(resource, data, s.nextID)

	s.store(resource)[id] = data
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) get(w http.ResponseWriter, resource, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.store(resource)[id]
	if !ok {
		writeError(w, http.StatusNotFound, "not found")

		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (s *Server) update(w http.ResponseWriter, resource, id string, r *http.Request) {
	data, ok := readJSON(w, r)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store(resource)[id]; !ok {
		writeError(w, http.StatusNotFound, "not found")

		return
	}

	data["id"] = id
	data = apiShape(resource, data, s.nextID)
	s.store(resource)[id] = data
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) delete(w http.ResponseWriter, resource, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store(resource), id)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) store(resource string) map[string]map[string]any {
	if s.items[resource] == nil {
		s.items[resource] = map[string]map[string]any{}
	}

	return s.items[resource]
}

func apiShape(resource string, data map[string]any, seq int) map[string]any {
	switch resource {
	case "groups":
		defaultValue(data, "peers", []any{})
		defaultValue(data, "resources", []any{})
		defaultValue(data, "peers_count", 0)
		defaultValue(data, "resources_count", 0)
	case "policies":
		defaultValue(data, "description", "api-generated policy description")
		for _, rule := range slice(data["rules"]) {
			rule, ok := rule.(map[string]any)
			if !ok {
				continue
			}
			defaultValue(rule, "id", fmt.Sprintf("rule-%d", seq))
			rule["sources"] = groupMinimums(rule["sources"])
			rule["destinations"] = groupMinimums(rule["destinations"])
		}
		defaultValue(data, "source_posture_checks", []any{})
	case "routes":
		defaultValue(data, "groups", []any{})
		defaultValue(data, "peer_groups", []any{})
		defaultValue(data, "network_type", "range")
	case "setup-keys":
		defaultValue(data, "key", fmt.Sprintf("mock-%s", data["id"]))
		defaultValue(data, "state", "valid")
		defaultValue(data, "valid", true)
		defaultValue(data, "revoked", false)
		defaultValue(data, "used_times", 0)
		defaultValue(data, "last_used", "0001-01-01T00:00:00Z")
		defaultValue(data, "expires", "2030-01-01T00:00:00Z")
		defaultValue(data, "updated_at", "2024-01-01T00:00:00Z")
		defaultValue(data, "ephemeral", false)
		defaultValue(data, "allow_extra_dns_labels", false)
	}

	return data
}

func defaultValue(data map[string]any, key string, value any) {
	if _, ok := data[key]; !ok {
		data[key] = value
	}
}

func groupMinimums(v any) []any {
	out := make([]any, 0, len(slice(v)))
	for _, item := range slice(v) {
		switch item := item.(type) {
		case string:
			out = append(out, map[string]any{
				"id":              item,
				"name":            item,
				"peers_count":     0,
				"resources_count": 0,
			})
		case map[string]any:
			out = append(out, item)
		}
	}

	return out
}

func slice(v any) []any {
	if v, ok := v.([]any); ok {
		return v
	}

	return nil
}

func readJSON(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	var data map[string]any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")

		return nil, false
	}

	return data, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"message": msg, "code": status})
}
