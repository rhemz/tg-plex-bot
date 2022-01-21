package v1

import (
    "fmt"
    "net/http"
    "net/url"
    "os"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/hekmon/plexwebhooks"

)

type PlexWebhookController struct{}


func (p PlexWebhookController) Post(c *gin.Context) {
    // telegram things
    tgBotId := os.Getenv("TELEGRAM_BOT_ID")
    if tgBotId == "" {
        zap.S().Fatal("TELEGRAM_BOT_ID not set")
    }
    tgToken := os.Getenv("TELEGRAM_API_TOKEN")
    if tgToken == "" {
        zap.S().Fatal("TELEGRAM_API_TOKEN not set")
    }

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
        msg := ""

        // show
        if payload.Metadata.LibrarySectionType == "show" {
            msg = fmt.Sprintf(`
%s started watching a TV Show

<b>%s</b>
%s, Episode %d
%s
`, 
            payload.Account.Title, payload.Metadata.GrandparentTitle, payload.Metadata.ParentTitle, payload.Metadata.Index, payload.Metadata.Title)
        } else if payload.Metadata.LibrarySectionType == "movie" {  // movie
            msg = fmt.Sprintf(`
%s started watching a Movie

<b>%s</b>
Ⓒ%d
`, 
            payload.Account.Title, payload.Metadata.Title, payload.Metadata.Year)
        }

        v := url.Values{}
        v.Set("chat_id", "-1001623668262")
        v.Set("parse_mode", "HTML")
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
}
