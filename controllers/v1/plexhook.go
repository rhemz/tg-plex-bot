package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hekmon/plexwebhooks"
	"github.com/rhemz/tg-plex-bot/config"
	"github.com/rhemz/tg-plex-bot/util"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type PlexWebhookController struct{}

var ipCache = make(map[string]util.IpInfoResponse)

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
	msgBody := ""
	locString := "local network"

	if payload.Event == "media.play" {

		ip := payload.Player.PublicAddress.String()
		if !strings.HasPrefix(ip, "192.168") {
			// look up ip info.  try cache first
			ipInfo, inCache := ipCache[ip]
			if !inCache {
				zap.S().Info(ip, " was not in the cache, looking it up")
				ipInfo, err = util.IpInfoLookup(cfg.GetString("ipinfo.apiToken"), ip)
				if err != nil {
					zap.S().Error("error looking up IpInfo.io data for ", ip, ": ", err)
					ipInfo = util.IpInfoResponse{}
				} else {
					// throw it in the cache
					ipCache[ip] = ipInfo
				}
			} else {
				zap.S().Info(ip, " was in the cache, skipping lookup")
			}
			locString = fmt.Sprintf("%s, %s", ipInfo.City, ipInfo.Region)
		}

		if payload.Metadata.LibrarySectionType == "show" { // show
			msgBody += fmt.Sprintf(`
%s started watching a TV Show from %s (%s)

<b>%s</b>
%s, Episode %d
%s
`,
				payload.Account.Title, payload.Player.PublicAddress.String(), locString, payload.Metadata.GrandparentTitle, payload.Metadata.ParentTitle, payload.Metadata.Index, payload.Metadata.Title)
		} else if payload.Metadata.LibrarySectionType == "movie" { // movie
			msgBody += fmt.Sprintf(`
%s started watching a Movie from %s (%s)

<b>%s</b>
Ⓒ%d
`,
				payload.Account.Title, payload.Player.PublicAddress.String(), locString, payload.Metadata.Title, payload.Metadata.Year)
		}
	} else if payload.Event == "media.scrobble" { // scrobble is > 90% completion of a media item
		msgBody += fmt.Sprintf("%s has finished their %s", payload.Account.Title, payload.Metadata.LibrarySectionType)
	}

	// is sending telegram messages disabled?
	if config.TelegramMessagesDisabled && len(msgBody) > 0 {
		zap.S().Info("Telegram messages are temporarily disabled.  Would have posted: ", msgBody)
		c.JSON(http.StatusOK, gin.H{"status": "no action"})
		return
	}

	// send telegram message(s) if any msg body was generated
	if len(msgBody) > 0 {
		if thumb != nil {
			err := util.SendPhotoToChats(thumb.Data, msgBody, cfg.GetIntSlice("telegram.broadcastChannels"))
			if err != nil {
				zap.S().Error("error sending photo(s) to telegram: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error sending photo to telegram channel(s)"})
			}
		} else {
			err := util.SendMessageToChats(msgBody, cfg.GetIntSlice("telegram.broadcastChannels"))
			if err != nil {
				zap.S().Error("error sending message(s) to telegram: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error sending message to telegram channel(s)"})
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "no action"})
}
