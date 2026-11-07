package service

type Storage interface {
	SetGauge(id string, value float64)
	AddCounter(id string, delta int64)
	GetGauge(id string) (float64, bool)
	GetCounter(id string) (int64, bool)
	GetAll() (map[string]float64, map[string]int64)
}
