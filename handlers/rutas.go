package handlers

import (
	"fuegodequasar/internal/platform/calculos"
	"fuegodequasar/internal/platform/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TopSecretRequest representa el payload para /topsecret
// @Description Datos de los satélites para decodificar mensaje y posición
type TopSecretRequest struct {
	Satellites []SatelliteInfo `json:"satellites"`
}

// SatelliteInfo representa la información de un satélite
// @Description Información individual de un satélite
type SatelliteInfo struct {
	Name     string   `json:"name" example:"kenobi"`
	Distance float32  `json:"distance" example:"927.75"`
	Message  []string `json:"message" example:"[\"este\", \"\", \"\", \"mensaje\", \"\"]"`
}

// TopSecretResponse representa la respuesta de /topsecret
// @Description Respuesta con posición y mensaje decodificado
type TopSecretResponse struct {
	Position Position `json:"position"`
	Message  string   `json:"message" example:"este es un mensaje secreto"`
}

// Position representa coordenadas X e Y
// @Description Coordenadas de la fuente
type Position struct {
	X float32 `json:"x" example:"426.4001"`
	Y float32 `json:"y" example:"-252.80016"`
}

type TopSecretSplitRequest struct {
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

// SetupRoutes configura las rutas HTTP de la API
func SetupRoutes(router *gin.Engine, repo repository.RepositoryService) {
	// POST /topsecret
	router.POST("/topsecret", handleTopSecret(repo))
	// POST /topsecret_split/{satellite_name}
	router.POST("/topsecret_split/:satellite_name", handleTopSecretSplit(repo))
	// GET /topsecret_split
	router.GET("/topsecret_split", handleGetTopSecretSplit(repo))
}

// @Summary Decodifica mensaje y posición
// @Description Recibe información de los satélites y retorna posición y mensaje
// @Tags topsecret
// @Accept json
// @Produce json
// @Param request body TopSecretRequest true "Datos de los satélites" example({"satellites":[{"name":"kenobi","distance":927.75,"message":["este","","","mensaje",""]},{"name":"skywalker","distance":360,"message":["","es","","","secreto"]},{"name":"sato","distance":360,"message":["este","","un","",""]}]})
// @Success 200 {object} TopSecretResponse "Ejemplo de respuesta" example({"position":{"x":426.4001,"y":-252.80016},"message":"este es un mensaje secreto"})
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /topsecret [post]
func handleTopSecret(repo repository.RepositoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request TopSecretRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Actualizar información de los satélites usando posición fija del repositorio
		for _, sat := range request.Satellites {
			// Obtener la posición fija desde el repositorio
			existing, err := repo.GetSatellite(sat.Name)
			var pos repository.Point
			if err == nil {
				pos = existing.Position
			} else {
				// Si no existe, usar posición por defecto (el repo ya lo hace en New())
				pos = repository.Point{}
			}
			satellite := repository.Satellite{
				Name:     sat.Name,
				Position: pos, // posición fija
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

// @Summary Guarda información parcial de un satélite
// @Description Permite guardar la distancia y mensaje de un satélite individualmente
// @Tags topsecret_split
// @Accept json
// @Produce json
// @Param satellite_name path string true "Nombre del satélite"
// @Param request body TopSecretSplitRequest true "Distancia y mensaje del satélite"
// @Success 200 "Actualización exitosa"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /topsecret_split/{satellite_name} [post]
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

// @Summary Decodifica mensaje y posición usando información parcial
// @Description Recupera la posición y mensaje usando los datos guardados de los satélites
// @Tags topsecret_split
// @Accept json
// @Produce json
// @Success 200 {object} TopSecretResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /topsecret_split [get]
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
