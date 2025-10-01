package main

import (
	"context"
	_ "fuegodequasar/docs"
	"fuegodequasar/handlers"
	"fuegodequasar/internal/platform/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @title Fuego de Quasar API
// @version 1.0
// @description API para el desafío de nivel 3 de Fuego de Quasar.
// @host localhost:8080
// @BasePath /api
func main() {
	// Configurar el modo de Gin basado en el ambiente
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configurar logger
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	// Determinar puerto para el servicio HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Inicializar el repositorio
	repo := repository.New()

	// Configurar el router con middleware de recuperación y logging
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithWriter(os.Stdout))

	// Añadir middleware de CORS
	router.Use(corsMiddleware())

	// Configurar las rutas
	handlers.SetupRoutes(router, repo)

	// Añadir la ruta de Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Añadir ruta de health check para Cloud Run
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Crear el servidor HTTP con timeouts
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Iniciar el servidor en una goroutine
	go func() {
		log.Printf("starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Configurar canal para señales de apagado
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	// Contexto para shutdown con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Intentar shutdown graceful
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	log.Println("server exited gracefully")
}

// corsMiddleware configura las cabeceras CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
