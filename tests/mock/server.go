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

	// The agent-network endpoints nest one level deeper
	// (/api/agent-network/<kind>[/<id>]), so the prefix is collapsed into a
	// single resource key and the generic CRUD routing below still applies.
	if len(parts) >= 3 && parts[1] == "agent-network" {
		parts = append([]string{parts[0], parts[1] + "/" + parts[2]}, parts[3:]...)
	}

	// Settings are a per-account singleton with a bootstrap POST, an
	// echo-checked PUT and a guarded DELETE, so they bypass generic routing.
	if len(parts) == 2 && parts[1] == agentNetworkSettingsKey {
		s.serveAgentNetworkSettings(w, r)

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

		if len(parts) == 2 {
			s.list(w, parts[1])

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
	// Nested resource keys carry a slash ("agent-network/providers"), which must
	// not leak into a generated ID or the follow-up request path stops routing.
	id := fmt.Sprintf("%s-%d", strings.ReplaceAll(resource, "/", "-"), s.nextID)
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

func (s *Server) list(w http.ResponseWriter, resource string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := s.store(resource)
	out := make([]map[string]any, 0, len(items))

	for _, item := range items {
		out = append(out, item)
	}

	writeJSON(w, http.StatusOK, out)
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
		data["peers"] = groupMinimums(data["peers"])
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
	case "users":
		defaultValue(data, "is_blocked", false)
		defaultValue(data, "auto_groups", []any{})
		defaultValue(data, "pending_approval", false)

		// Mirrors real API behavior: blocking a user flips their status too.
		if blocked, _ := data["is_blocked"].(bool); blocked {
			data["status"] = "blocked"
		} else {
			defaultValue(data, "status", "active")
		}
	case "agent-network/providers":
		// The API seals the key and never returns it, and always reports both
		// identity headers on the wire — empty when unset.
		delete(data, "api_key")
		delete(data, "bootstrap_cluster")
		defaultValue(data, "identity_header_user_id", "")
		defaultValue(data, "identity_header_groups", "")
		defaultValue(data, "enabled", true)
		defaultValue(data, "skip_tls_verification", false)
		defaultValue(data, "metadata_disabled", false)
		defaultValue(data, "models", []any{})
		defaultValue(data, "created_at", mockTimestamp)
		defaultValue(data, "updated_at", mockTimestamp)
	case "agent-network/policies", "agent-network/guardrails":
		// description is always on the wire, empty when unset.
		defaultValue(data, "description", "")
		defaultValue(data, "enabled", true)
		defaultValue(data, "created_at", mockTimestamp)
		defaultValue(data, "updated_at", mockTimestamp)
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

const (
	agentNetworkSettingsKey = "agent-network/settings"
	agentNetworkSettingsRow = "singleton"
	mockTimestamp           = "2024-01-01T00:00:00Z"
)

// serveAgentNetworkSettings implements the per-account Agent Network settings
// singleton: POST bootstraps and assigns the immutable endpoint, PUT replaces
// the mutable fields while requiring the identity fields to be echoed back
// unchanged, and DELETE is refused while providers still exist.
func (s *Server) serveAgentNetworkSettings(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row, bootstrapped := s.store(agentNetworkSettingsKey)[agentNetworkSettingsRow]

	switch r.Method {
	case http.MethodGet:
		if !bootstrapped {
			writeError(w, http.StatusNotFound, "agent network settings have not been bootstrapped yet")

			return
		}

		writeJSON(w, http.StatusOK, row)
	case http.MethodPost:
		if bootstrapped {
			writeError(w, http.StatusConflict, "agent network settings already bootstrapped for account")

			return
		}

		s.bootstrapAgentNetworkSettings(w, r)
	case http.MethodPut:
		if !bootstrapped {
			writeError(w, http.StatusNotFound, "agent network settings have not been bootstrapped yet")

			return
		}

		s.updateAgentNetworkSettings(w, r, row)
	case http.MethodDelete:
		if len(s.store("agent-network/providers")) > 0 {
			writeError(w, http.StatusPreconditionFailed, "agent network settings cannot be deleted while providers exist")

			return
		}

		delete(s.store(agentNetworkSettingsKey), agentNetworkSettingsRow)
		writeJSON(w, http.StatusOK, map[string]any{})
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}

func (s *Server) bootstrapAgentNetworkSettings(w http.ResponseWriter, r *http.Request) {
	data, ok := readJSON(w, r)
	if !ok {
		return
	}

	proxyAddress, hasProxyAddress := data["proxy_address"].(string)
	endpoint, hasEndpoint := data["endpoint"].(string)

	if hasProxyAddress == hasEndpoint {
		writeError(w, http.StatusUnprocessableEntity, "exactly one of proxy_address and endpoint is required")

		return
	}

	if hasProxyAddress {
		// A labeled bootstrap: the server allocates the label and the endpoint
		// hangs one label beneath the requested cluster address.
		endpoint = "brave-otter." + proxyAddress
	} else {
		// A self-addressed bootstrap: the endpoint is claimed verbatim and the
		// proxy serving it declares exactly that address.
		proxyAddress = endpoint
	}

	row := map[string]any{
		"endpoint":                  endpoint,
		"proxy_address":             proxyAddress,
		"dedicated":                 endpoint == proxyAddress,
		"enable_log_collection":     boolOrDefault(data["enable_log_collection"], true),
		"enable_prompt_collection":  boolOrDefault(data["enable_prompt_collection"], false),
		"redact_pii":                boolOrDefault(data["redact_pii"], false),
		"access_log_retention_days": intOrDefault(data["access_log_retention_days"], 30),
		"created_at":                mockTimestamp,
		"updated_at":                mockTimestamp,
	}

	s.store(agentNetworkSettingsKey)[agentNetworkSettingsRow] = row
	writeJSON(w, http.StatusOK, row)
}

func (s *Server) updateAgentNetworkSettings(w http.ResponseWriter, r *http.Request, row map[string]any) {
	data, ok := readJSON(w, r)
	if !ok {
		return
	}

	for _, field := range []string{"endpoint", "proxy_address"} {
		if sent, _ := data[field].(string); sent != row[field] {
			writeError(w, http.StatusUnprocessableEntity, field+" is immutable once assigned")

			return
		}
	}

	row["enable_log_collection"] = boolOrDefault(data["enable_log_collection"], false)
	row["enable_prompt_collection"] = boolOrDefault(data["enable_prompt_collection"], false)
	row["redact_pii"] = boolOrDefault(data["redact_pii"], false)
	row["access_log_retention_days"] = intOrDefault(data["access_log_retention_days"], 30)

	s.store(agentNetworkSettingsKey)[agentNetworkSettingsRow] = row
	writeJSON(w, http.StatusOK, row)
}

func boolOrDefault(v any, fallback bool) bool {
	if v, ok := v.(bool); ok {
		return v
	}

	return fallback
}

func intOrDefault(v any, fallback int) int {
	if v, ok := v.(float64); ok {
		return int(v)
	}

	return fallback
}
