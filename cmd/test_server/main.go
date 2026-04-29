package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type config struct {
	port         string
	clientID     string
	clientSecret string
}

func defaultConfig() *config {
	return &config{
		port:         "8080",
		clientID:     "test-client-id",
		clientSecret: "test-client-secret",
	}
}

type positionResponse struct {
	Code               string    `json:"code"`
	Title              string    `json:"title"`
	SupervisorPosition bool      `json:"supervisorPosition"`
	EffectiveDate      time.Time `json:"effectiveDate"`
	ClosedDate         time.Time `json:"closedDate"`
}

type InfoPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"personalEmail"`
	JobTitle  string `json:"jobTitle"`
}

type PositionPayload struct {
	PositionCode string `json:"positionCode"`
	EmployeeType string `json:"employeeType"`
	Department   string `json:"costCenter1"`
}

type employeeResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Status      string `json:"statusType"`

	Info     InfoPayload     `json:"info"`
	Position PositionPayload `json:"position"`
}

var (
	tokenCache   = make(map[string]time.Time)
	tokenCacheMu sync.RWMutex
)

type server struct {
	config    *config
	mu        sync.RWMutex
	positions map[string]positionResponse
	employees map[string]employeeResponse
}

func newServer(config *config) *server {
	return &server{
		config:    config,
		positions: make(map[string]positionResponse),
		employees: make(map[string]employeeResponse),
	}
}

// maxAuthBodyBytes caps the request body for the test server's auth endpoint
// so that ParseForm cannot be used to exhaust memory (gosec G120).
const maxAuthBodyBytes = 1 << 20 // 1 MiB

func (s *server) handleAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	if r.FormValue("grant_type") != "client_credentials" || r.FormValue("client_id") != s.config.clientID || r.FormValue("client_secret") != s.config.clientSecret {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := uuid.New().String()
	tokenCacheMu.Lock()
	tokenCache[token] = time.Now().Add(1 * time.Hour)
	tokenCacheMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   3600,
	}); err != nil {
		log.Printf("Error writing json response: %v", err)
	}
}

func (s *server) handlePositions(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	limit, offset := parseOffsetPagination(r)
	allPositions := make([]positionResponse, 0, len(s.positions))
	for _, p := range s.positions {
		allPositions = append(allPositions, p)
	}
	sort.Slice(allPositions, func(i, j int) bool { return allPositions[i].Code < allPositions[j].Code })
	total := len(allPositions)
	w.Header().Set("X-Pcty-Total-Count", strconv.Itoa(total))
	w.Header().Set("Content-Type", "application/json")
	if offset >= total {
		if err := json.NewEncoder(w).Encode([]positionResponse{}); err != nil {
			log.Printf("Error writing empty json response: %v", err)
		}
		return
	}
	end := offset + limit
	if end > total {
		end = total
	}
	paginatedPositions := allPositions[offset:end]
	if err := json.NewEncoder(w).Encode(paginatedPositions); err != nil {
		log.Printf("Error writing json response: %v", err)
	}
}

func (s *server) handleEmployees(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	limit, offset := parseTokenPaginationAsOffset(r)
	employees := make([]employeeResponse, 0, len(s.employees))
	for _, e := range s.employees {
		employees = append(employees, e)
	}
	sort.Slice(employees, func(i, j int) bool { return employees[i].ID < employees[j].ID })
	if offset > len(employees) {
		offset = len(employees)
	}
	end := offset + limit
	if end > len(employees) {
		end = len(employees)
	}
	paginated := employees[offset:end]
	var nextToken string
	if end < len(employees) {
		nextToken = strconv.Itoa(end)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"employees": paginated,
		"nextToken": nextToken,
	}); err != nil {
		log.Printf("Error writing json response: %v", err)
	}
}

func (s *server) handleSingleEmployee(w http.ResponseWriter, r *http.Request, employeeID string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	employee, ok := s.employees[employeeID]
	if !ok {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employee); err != nil {
		log.Printf("Error writing json response: %v", err)
	}
}

func (s *server) employeesRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/coreHr/v1/companies/%s/employees", "test-company"))
	if path == "" || path == "/" {
		s.handleEmployees(w, r)
		return
	}
	employeeID := strings.TrimPrefix(path, "/")
	s.handleSingleEmployee(w, r, employeeID)
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		tokenCacheMu.RLock()
		expiresAt, ok := tokenCache[token]
		tokenCacheMu.RUnlock()
		if !ok || time.Now().After(expiresAt) {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (s *server) addTestData() {
	s.positions["DEV-01"] = positionResponse{Code: "DEV-01", Title: "Software Developer"}
	s.positions["DEV-02"] = positionResponse{Code: "DEV-02", Title: "Senior Software Developer"}
	s.positions["MGR-01"] = positionResponse{Code: "MGR-01", Title: "Engineering Manager", SupervisorPosition: true}

	s.employees["101"] = employeeResponse{
		ID: "101", DisplayName: "Ana Gomez", Status: "Active",
		Info:     InfoPayload{FirstName: "Ana", LastName: "Gomez", Email: "ana.gomez@example.com", JobTitle: "Software Developer"},
		Position: PositionPayload{PositionCode: "DEV-01", EmployeeType: "Regular", Department: "Technology"},
	}
	s.employees["102"] = employeeResponse{
		ID: "102", DisplayName: "Juan Perez", Status: "Active",
		Info:     InfoPayload{FirstName: "Juan", LastName: "Perez", Email: "juan.perez@example.com", JobTitle: "Senior Software Developer"},
		Position: PositionPayload{PositionCode: "DEV-02", EmployeeType: "Regular", Department: "Technology"},
	}
	s.employees["103"] = employeeResponse{
		ID: "103", DisplayName: "Maria Rodriguez", Status: "Active",
		Info:     InfoPayload{FirstName: "Maria", LastName: "Rodriguez", Email: "maria.r@example.com", JobTitle: "Engineering Manager"},
		Position: PositionPayload{PositionCode: "MGR-01", EmployeeType: "Regular", Department: "Management"},
	}
	s.employees["104"] = employeeResponse{
		ID: "104", DisplayName: "Carlos Lopez", Status: "Terminated",
		Info:     InfoPayload{FirstName: "Carlos", LastName: "Lopez", Email: "carlos.lopez@example.com", JobTitle: "Software Developer"},
		Position: PositionPayload{PositionCode: "DEV-01", EmployeeType: "Regular", Department: "Technology"},
	}
}

const defaultPageLimit = 100

func parseOffsetPagination(r *http.Request) (int, int) {
	limit := defaultPageLimit
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}
	return limit, offset
}

func parseTokenPaginationAsOffset(r *http.Request) (int, int) {
	limit := defaultPageLimit
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("nextToken"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}
	return limit, offset
}

func main() {
	config := defaultConfig()
	server := newServer(config)
	server.addTestData()
	mux := http.NewServeMux()
	companyID := "test-company"
	positionsPath := fmt.Sprintf("/apiHub/positionManagement/v1/companies/%s/positions", companyID)
	employeesPath := fmt.Sprintf("/coreHr/v1/companies/%s/employees/", companyID)
	mux.HandleFunc("/public/security/v1/token", server.handleAuth)
	mux.Handle(positionsPath, authMiddleware(server.handlePositions))
	mux.Handle(employeesPath, authMiddleware(server.employeesRouter))

	srv := &http.Server{
		Addr:         ":" + config.port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting Paylocity test server on port %s", config.port)
		log.Printf("Server is configured for companyId: %s", companyID)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("could not listen on %s: %w", config.port, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Println("Shutting down server...")
	case err := <-serverErr:
		log.Printf("Server error: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("ERROR: Server graceful shutdown failed: %v", err)
	}

	log.Println("Server exiting.")
}
