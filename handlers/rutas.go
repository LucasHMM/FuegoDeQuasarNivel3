package handlers

import (
	"fuegodequasar/internal/platform/calculos"
	"fuegodequasar/internal/platform/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TopSecretRequest struct {
	Satellites []SatelliteInfo `json:"satellites"`
}

type SatelliteInfo struct {
	Name     string   `json:"name"`
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

type TopSecretResponse struct {
	Position Position `json:"position"`
	Message  string   `json:"message"`
}

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type TopSecretSplitRequest struct {
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

func SetupRoutes(router *gin.Engine, repo repository.RepositoryService) {
	// Grupo de rutas para el servicio
	api := router.Group("/api")
	{
		// POST /topsecret/
		api.POST("/topsecret", handleTopSecret(repo))

		// POST /topsecret_split/{satellite_name}
		api.POST("/topsecret_split/:satellite_name", handleTopSecretSplit(repo))

		// GET /topsecret_split/
		api.GET("/topsecret_split", handleGetTopSecretSplit(repo))
	}
}

func handleTopSecret(repo repository.RepositoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request TopSecretRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Actualizar información de los satélites
		for _, sat := range request.Satellites {
			satellite := repository.Satellite{
				Name:     sat.Name,
				Distance: sat.Distance,
				Message:  sat.Message,
			}
			if err := repo.SaveSatellite(satellite); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save satellite info"})
				return
			}
		}

		// Obtener todos los satélites para el cálculo
		satellites, err := repo.GetAllSatellites()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve satellites"})
			return
		}

		// Preparar datos para la trilateración
		var positions []calculos.Point32
		var distances []float32
		var messages [][]string

		for _, sat := range satellites {
			positions = append(positions, calculos.Point32{
				X: sat.Position.X,
				Y: sat.Position.Y,
			})
			distances = append(distances, sat.Distance)
			messages = append(messages, sat.Message)
		}

		// Calcular posición
		if len(positions) < 3 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not enough satellite data"})
			return
		}

		location := calculos.GetLocation(positions[0], positions[1], positions[2],
			distances[0], distances[1], distances[2])

		// Recuperar mensaje
		message, err := calculos.GetMessage(messages[0], messages[1], messages[2])
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Could not decode message"})
			return
		}

		response := TopSecretResponse{
			Position: Position{
				X: location.X,
				Y: location.Y,
			},
			Message: message,
		}

		c.JSON(http.StatusOK, response)
	}
}

func handleTopSecretSplit(repo repository.RepositoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		satelliteName := c.Param("satellite_name")
		var request TopSecretSplitRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Obtener el satélite existente para mantener su posición
		satellite, err := repo.GetSatellite(satelliteName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Satellite not found"})
			return
		}

		// Actualizar la distancia y mensaje del satélite
		satellite.Distance = request.Distance
		satellite.Message = request.Message

		// Guardar la información actualizada
		if err := repo.SaveSatellite(satellite); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save satellite info"})
			return
		}

		c.Status(http.StatusOK)
	}
}

func handleGetTopSecretSplit(repo repository.RepositoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener todos los satélites
		satellites, err := repo.GetAllSatellites()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve satellites"})
			return
		}

		// Verificar que tengamos suficiente información
		validSatellites := 0
		for _, sat := range satellites {
			if sat.Distance > 0 && len(sat.Message) > 0 {
				validSatellites++
			}
		}

		if validSatellites < 3 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not enough satellite data"})
			return
		}

		// Preparar datos para la trilateración
		var positions []calculos.Point32
		var distances []float32
		var messages [][]string

		for _, sat := range satellites {
			if sat.Distance > 0 && len(sat.Message) > 0 {
				positions = append(positions, calculos.Point32{
					X: sat.Position.X,
					Y: sat.Position.Y,
				})
				distances = append(distances, sat.Distance)
				messages = append(messages, sat.Message)
			}
		}

		// Calcular posición
		location := calculos.GetLocation(positions[0], positions[1], positions[2],
			distances[0], distances[1], distances[2])

		// Recuperar mensaje
		message, err := calculos.GetMessage(messages[0], messages[1], messages[2])
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Could not decode message"})
			return
		}

		response := TopSecretResponse{
			Position: Position{
				X: location.X,
				Y: location.Y,
			},
			Message: message,
		}

		c.JSON(http.StatusOK, response)
	}
}
