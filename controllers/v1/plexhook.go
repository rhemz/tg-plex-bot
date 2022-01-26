package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hekmon/plexwebhooks"
	"github.com/rhemz/tg-plex-bot/config"
	"github.com/rhemz/tg-plex-bot/util"
	"go.uber.org/zap"
	"net/http"
)

type PlexWebhookController struct{}

func (p PlexWebhookController) Post(c *gin.Context) {
	cfg := config.GetConfig()

	reader, err := c.Request.MultipartReader()
	if err != nil {
		// Detect error type for the http answer
		if err == http.ErrNotMultipart || err == http.ErrMissingBoundary {
			c.JSON(http.StatusBadRequest, gin.H{"status": "bad multipart, dawg"})
		} else {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		}
		return
	}

	payload, thumb, err := plexwebhooks.Extract(reader)

	zap.S().Debug(thumb)
	zap.S().Debug(err)

	// send a message to the channel
	zap.S().Info("Got plex event:", payload.Event)
	if payload.Event == "media.play" {
		msgBody := ""

		// show
		if payload.Metadata.LibrarySectionType == "show" {
			msgBody = fmt.Sprintf(`
%s started watching a TV Show from %s

<b>%s</b>
%s, Episode %d
%s
`,
				payload.Account.Title, payload.Player.PublicAddress.String(), payload.Metadata.GrandparentTitle, payload.Metadata.ParentTitle, payload.Metadata.Index, payload.Metadata.Title)
		} else if payload.Metadata.LibrarySectionType == "movie" { // movie
			msgBody = fmt.Sprintf(`
%s started watching a Movie from %s

<b>%s</b>
Ⓒ%d
`,
				payload.Account.Title, payload.Player.PublicAddress.String(), payload.Metadata.Title, payload.Metadata.Year)
		}

		err := util.SendMessageToChats(msgBody, cfg.GetIntSlice("telegram.broadcastChannels"))
		if err != nil {
			zap.S().Error("Error sending message(s) to telegram: ", err)
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})

	}
}
