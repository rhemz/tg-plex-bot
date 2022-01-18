package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "net/url"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    ginzap "github.com/gin-contrib/zap"
    "go.uber.org/zap"
)

func main() {
    // initialize gin, loggers
    g := gin.New()
    logger, _ := zap.NewProduction()
    sugar := logger.Sugar()
    g.Use(ginzap.Ginzap(logger, time.RFC3339, true)) // use zap logger
    g.Use(ginzap.RecoveryWithZap(logger, true))  // log panic to error

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    sugar.Info("Attempting to listen on port ", port)

    // telegram things
    tgBotId := os.Getenv("TELEGRAM_BOT_ID")
    if tgBotId == "" {
        sugar.Fatal("TELEGRAM_BOT_ID not set")
    }
    tgToken := os.Getenv("TELEGRAM_API_TOKEN")
    if tgToken == "" {
        sugar.Fatal("TELEGRAM_API_TOKEN not set")
    }
    

    // register handlers
    g.POST("/plexhook", func(c *gin.Context) {
        var e PlexEvent
        if c.BindJSON(&e) == nil {
            // for now lets log it
            body, _ := json.Marshal(e)
            sugar.Info(string(body))
        }

        // send a message to the channel
        if e.Event == "media.play" {
            sugar.Info("got play event!")

            v := url.Values{}
            v.Set("chat_id", "-1001623668262")
            v.Set("text", fmt.Sprintf("%s started playing %s", e.Account.Title, e.Metadata.Title))

            url := url.URL{
                Scheme:     "https",
                Host:       "api.telegram.org",
                Path:       fmt.Sprintf("%s:%s/sendMessage", tgBotId, tgToken),
                RawQuery:   v.Encode(),
            }

            urlString := url.String()
            sugar.Info("Sending request: ", urlString)
            

            _, err := http.Get(urlString)
            if err != nil {
               sugar.Error(err)
            }
        }
    })
    g.GET("/health", func(c *gin.Context) {
        c.String(http.StatusOK, "ok")
    })

    // Listen and serve on defined port
    sugar.Info("Listening on port ", port)
    g.Run(":" + port)
}
