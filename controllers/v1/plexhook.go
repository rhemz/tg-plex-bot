package v1


func plexhookPOST(c *gin.Context) {

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
}