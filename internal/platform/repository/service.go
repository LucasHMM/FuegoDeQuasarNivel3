package repository

import (
	"errors"
	"sync"
)

// Errores del repositorio
var (
	ErrSatelliteNotFound = errors.New("satellite not found")
)

// Satellite representa la información de un satélite
type Satellite struct {
	Name     string   `json:"name"`
	Position Point    `json:"position"`
	Message  []string `json:"message"`
	Distance float32  `json:"distance"`
}

// Point representa una posición en coordenadas x,y
type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// INTERFAZ que debe implementar el repositorio y el mock
type RepositoryService interface {
	// GetSatellite obtiene la información de un satélite por su nombre
	GetSatellite(name string) (Satellite, error)
	// SaveSatellite guarda o actualiza la información de un satélite
	SaveSatellite(satellite Satellite) error
	// GetAllSatellites obtiene la información de todos los satélites
	GetAllSatellites() ([]Satellite, error)
}

// Estructura que implementa RepositoryService
type Service struct {
	satellites map[string]Satellite
	mutex      sync.RWMutex
}

func New() *Service {
	// Inicializamos con las posiciones conocidas de los satélites
	initialSatellites := map[string]Satellite{
		"kenobi": {
			Name: "kenobi",
			Position: Point{
				X: -500,
				Y: -200,
			},
		},
		"skywalker": {
			Name: "skywalker",
			Position: Point{
				X: 100,
				Y: -100,
			},
		},
		"sato": {
			Name: "sato",
			Position: Point{
				X: 500,
				Y: 100,
			},
		},
	}

	return &Service{
		satellites: initialSatellites,
		mutex:      sync.RWMutex{},
	}
}

func (s *Service) GetSatellite(name string) (Satellite, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if satellite, exists := s.satellites[name]; exists {
		return satellite, nil
	}
	return Satellite{}, ErrSatelliteNotFound
}

func (s *Service) SaveSatellite(satellite Satellite) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.satellites[satellite.Name] = satellite
	return nil
}

func (s *Service) GetAllSatellites() ([]Satellite, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	satellites := make([]Satellite, 0, len(s.satellites))
	for _, satellite := range s.satellites {
		satellites = append(satellites, satellite)
	}
	return satellites, nil
}
