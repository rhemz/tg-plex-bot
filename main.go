package main

import (
    "net/http"
    "os"
    // "reflect"
    "time"

    ginzap "github.com/gin-contrib/zap"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/rhemz/tg-plex-bot/controllers/v1"
)


func main() {
    // initialize gin, loggers
    g := gin.New()
    logger, _ := zap.NewProduction()
    undoLoggerReplace := zap.ReplaceGlobals(logger)
    defer undoLoggerReplace()

    g.Use(ginzap.Ginzap(logger, time.RFC3339, true)) // use zap logger
    g.Use(ginzap.RecoveryWithZap(logger, true))      // log panic to error

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    zap.S().Info("Attempting to listen on port ", port)

    // register handlers
    g.POST("/plexhook", v1.PlexHookPOST)
    g.GET("/health", func(c *gin.Context) {
        c.String(http.StatusOK, "ok")
    })

    // Listen and serve on defined port
    zap.S().Info("Listening on port ", port)
    g.Run(":" + port)
}


func contains(s []string, str string) bool {
    for _, v := range s {
        if v == str {
            return true
        }
    }

    return false
}
