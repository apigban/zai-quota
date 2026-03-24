package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	port            = flag.Int("port", 19876, "Port to listen on")
	currentScenario = "success_full"
	scenarioMu      sync.RWMutex
)

type UsageDetail struct {
	ModelCode string `json:"modelCode"`
	Usage     int    `json:"usage"`
}

type Limit struct {
	Type          string        `json:"type"`
	Percentage    int           `json:"percentage"`
	Usage         int           `json:"usage"`
	CurrentValue  int           `json:"currentValue"`
	Total         int           `json:"total"`
	Remaining     int           `json:"remaining"`
	NextResetTime int64         `json:"nextResetTime"`
	UsageDetails  []UsageDetail `json:"usageDetails"`
}

type ResponseData struct {
	Limits []Limit `json:"limits"`
	Level  string  `json:"level"`
}

type APIResponse struct {
	Success bool          `json:"success"`
	Code    int           `json:"code"`
	Msg     string        `json:"msg"`
	Data    *ResponseData `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
}

var scenarios = map[string]func(w http.ResponseWriter, r *http.Request){
	"success_full":       successFullScenario,
	"success_partial":    successPartialScenario,
	"success_empty":      successEmptyScenario,
	"success_warning":    successWarningScenario,
	"success_critical":   successCriticalScenario,
	"auth_invalid":       authInvalidScenario,
	"auth_forbidden":     authForbiddenScenario,
	"server_error":       serverErrorScenario,
	"server_unavailable": serverUnavailableScenario,
	"timeout":            timeoutScenario,
	"malformed":          malformedScenario,
	"rate_limited":       rateLimitedScenario,
}

func getScenario() string {
	scenarioMu.RLock()
	defer scenarioMu.RUnlock()
	return currentScenario
}

func setScenario(name string) bool {
	scenarioMu.Lock()
	defer scenarioMu.Unlock()
	if _, exists := scenarios[name]; exists {
		currentScenario = name
		return true
	}
	return false
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func successFullScenario(w http.ResponseWriter, r *http.Request) {
	resp := APIResponse{
		Success: true,
		Code:    200,
		Msg:     "OK",
		Data: &ResponseData{
			Limits: []Limit{
				{
					Type:          "TOKENS_LIMIT",
					Percentage:    28,
					Usage:         1400,
					CurrentValue:  1400,
					Total:         5000,
					Remaining:     3600,
					NextResetTime: time.Now().Add(2 * time.Hour).UnixMilli(),
					UsageDetails: []UsageDetail{
						{ModelCode: "claude-3-opus", Usage: 800},
						{ModelCode: "claude-3-sonnet", Usage: 600},
					},
				},
				{
					Type:          "TIME_LIMIT",
					Percentage:    2,
					Usage:         1000,
					CurrentValue:  16,
					Total:         1000,
					Remaining:     984,
					NextResetTime: time.Now().Add(24 * time.Hour).UnixMilli(),
					UsageDetails: []UsageDetail{
						{ModelCode: "search-prime", Usage: 10},
						{ModelCode: "web-reader", Usage: 4},
						{ModelCode: "ref", Usage: 2},
					},
				},
			},
			Level: "pro",
		},
	}
	writeJSON(w, 200, resp)
}

func successPartialScenario(w http.ResponseWriter, r *http.Request) {
	resp := APIResponse{
		Success: true,
		Code:    200,
		Msg:     "OK",
		Data: &ResponseData{
			Limits: []Limit{
				{
					Type:          "TOKENS_LIMIT",
					Percentage:    75,
					Usage:         3750,
					CurrentValue:  3750,
					Total:         5000,
					Remaining:     1250,
					NextResetTime: time.Now().Add(1 * time.Hour).UnixMilli(),
					UsageDetails:  []UsageDetail{},
				},
			},
			Level: "free",
		},
	}
	writeJSON(w, 200, resp)
}

func successEmptyScenario(w http.ResponseWriter, r *http.Request) {
	resp := APIResponse{
		Success: true,
		Code:    200,
		Msg:     "OK",
		Data: &ResponseData{
			Limits: []Limit{},
			Level:  "free",
		},
	}
	writeJSON(w, 200, resp)
}

func successWarningScenario(w http.ResponseWriter, r *http.Request) {
	resp := APIResponse{
		Success: true,
		Code:    200,
		Msg:     "OK",
		Data: &ResponseData{
			Limits: []Limit{
				{
					Type:          "TOKENS_LIMIT",
					Percentage:    82,
					Usage:         4100,
					CurrentValue:  4100,
					Total:         5000,
					Remaining:     900,
					NextResetTime: time.Now().Add(30 * time.Minute).UnixMilli(),
					UsageDetails:  []UsageDetail{},
				},
			},
			Level: "pro",
		},
	}
	writeJSON(w, 200, resp)
}

func successCriticalScenario(w http.ResponseWriter, r *http.Request) {
	resp := APIResponse{
		Success: true,
		Code:    200,
		Msg:     "OK",
		Data: &ResponseData{
			Limits: []Limit{
				{
					Type:          "TOKENS_LIMIT",
					Percentage:    98,
					Usage:         4900,
					CurrentValue:  4900,
					Total:         5000,
					Remaining:     100,
					NextResetTime: time.Now().Add(5 * time.Minute).UnixMilli(),
					UsageDetails:  []UsageDetail{},
				},
			},
			Level: "pro",
		},
	}
	writeJSON(w, 200, resp)
}

func authInvalidScenario(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 401, ErrorResponse{
		Success: false,
		Code:    401,
		Msg:     "Invalid API key",
		Data:    nil,
	})
}

func authForbiddenScenario(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 403, ErrorResponse{
		Success: false,
		Code:    403,
		Msg:     "Access forbidden",
		Data:    nil,
	})
}

func serverErrorScenario(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 500, ErrorResponse{
		Success: false,
		Code:    500,
		Msg:     "Internal server error",
		Data:    nil,
	})
}

func serverUnavailableScenario(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 503, ErrorResponse{
		Success: false,
		Code:    503,
		Msg:     "Service temporarily unavailable",
		Data:    nil,
	})
}

func timeoutScenario(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	writeJSON(w, 200, APIResponse{
		Success: true,
		Code:    200,
		Msg:     "OK",
		Data:    &ResponseData{Limits: []Limit{}, Level: "free"},
	})
}

func malformedScenario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"success": true, invalid json`))
}

func rateLimitedScenario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", "60")
	writeJSON(w, 429, ErrorResponse{
		Success: false,
		Code:    429,
		Msg:     "Rate limit exceeded",
		Data:    nil,
	})
}

func quotaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		authInvalidScenario(w, r)
		return
	}

	scenario := getScenario()
	handler, exists := scenarios[scenario]
	if !exists {
		http.Error(w, "Unknown scenario", 500)
		return
	}
	handler(w, r)
}

func controlScenarioHandler(w http.ResponseWriter, r *http.Request) {
	scenario := r.PathValue("scenario")
	if scenario == "" {
		http.Error(w, "Scenario name required", 400)
		return
	}

	if setScenario(scenario) {
		fmt.Fprintf(w, "OK: switched to scenario '%s'", scenario)
	} else {
		http.Error(w, fmt.Sprintf("Unknown scenario: %s", scenario), 400)
	}
}

func controlStatusHandler(w http.ResponseWriter, r *http.Request) {
	scenario := getScenario()
	scenarioList := []string{
		"success_full", "success_partial", "success_empty",
		"success_warning", "success_critical",
		"auth_invalid", "auth_forbidden",
		"server_error", "server_unavailable",
		"timeout", "malformed", "rate_limited",
	}
	resp := map[string]any{
		"current_scenario":    scenario,
		"available_scenarios": scenarioList,
	}
	writeJSON(w, 200, resp)
}

func controlListHandler(w http.ResponseWriter, r *http.Request) {
	scenarios := []string{
		"success_full", "success_partial", "success_empty",
		"success_warning", "success_critical",
		"auth_invalid", "auth_forbidden",
		"server_error", "server_unavailable",
		"timeout", "malformed", "rate_limited",
	}
	writeJSON(w, 200, scenarios)
}

func main() {
	flag.Parse()

	http.HandleFunc("/api/monitor/usage/quota/limit", quotaHandler)
	http.HandleFunc("/control/scenario/", controlScenarioHandler)
	http.HandleFunc("/control/status", controlStatusHandler)
	http.HandleFunc("/control/list", controlListHandler)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Mock Z.ai API server starting on %s", addr)
	log.Printf("Available scenarios: %v", []string{
		"success_full", "success_partial", "success_empty",
		"success_warning", "success_critical",
		"auth_invalid", "auth_forbidden",
		"server_error", "server_unavailable",
		"timeout", "malformed", "rate_limited",
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
