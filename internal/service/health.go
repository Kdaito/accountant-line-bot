package service

import "net/http"

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (h *HealthService) HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello! This is accountant line bot!"))
}
