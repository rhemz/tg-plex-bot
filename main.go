package main

import (
    "fmt"
    "net/http"
    "net/url"
    "os"
    // "reflect"
    "time"

    ginzap "github.com/gin-contrib/zap"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/hekmon/plexwebhooks"
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

    // telegram things
    tgBotId := os.Getenv("TELEGRAM_BOT_ID")
    if tgBotId == "" {
        zap.S().Fatal("TELEGRAM_BOT_ID not set")
    }
    tgToken := os.Getenv("TELEGRAM_API_TOKEN")
    if tgToken == "" {
        zap.S().Fatal("TELEGRAM_API_TOKEN not set")
    }

    // register handlers
    g.POST("/plexhook", func(c *gin.Context) {
        // var e PlexEvent

        reader, err := c.Request.MultipartReader()
        if err != nil {
            // Detect error type for the http answer
            if err == http.ErrNotMultipart || err == http.ErrMissingBoundary {
                c.String(http.StatusBadRequest, "bad multipart, dawg")
            } else {
                c.String(http.StatusBadRequest, "some other kinda error")
            }
            zap.S().Warn(err)
            return
        }

        payload, thumb, err := plexwebhooks.Extract(reader)

        zap.S().Debug(thumb)
        zap.S().Debug(err)

        // send a message to the channel
        if payload.Event == "media.play" {
            zap.S().Info("got play event!")

            // show
            msg := ""
            if payload.Metadata.LibrarySectionType == "show" {
                msg = fmt.Sprintf("%s started playing %s, %s - %s", payload.Account.Title, payload.Metadata.GrandparentTitle, payload.Metadata.ParentTitle, payload.Metadata.Title)
            }

            v := url.Values{}
            v.Set("chat_id", "-1001623668262")
            v.Set("text", msg)

            url := url.URL{
                Scheme:   "https",
                Host:     "api.telegram.org",
                Path:     fmt.Sprintf("%s:%s/sendMessage", tgBotId, tgToken),
                RawQuery: v.Encode(),
            }

            urlString := url.String()
            zap.S().Info("Sending request: ", urlString)

            _, err := http.Get(urlString)
            if err != nil {
                zap.S().Error(err)
            }
        }
    })
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
