package handler

type MetricsService interface {
	Update(kind, name, rawValue string) error
	Value(kind, name string) (string, error)
	GetAll() (map[string]float64, map[string]int64)
}

type Handler struct {
	metrics MetricsService
}

func NewHandler(metrics MetricsService) *Handler {
	return &Handler{
		metrics: metrics,
	}
}
