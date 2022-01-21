package server

import (
    "os"
    "time"

    "go.uber.org/zap"
    ginzap "github.com/gin-contrib/zap"
)


func Start() {
	// config := config.GetConfig()
	// r.Run(config.GetString("server.port"))

    r := NewRouter()

    r.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true)) // use zap logger
    r.Use(ginzap.RecoveryWithZap(zap.L(), true))      // log panic to error

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    zap.S().Info("Attempting to listen on port ", port)

    // Listen and serve on defined port
    zap.S().Info("Listening on port ", port)

    r.Run(":" + port)
}