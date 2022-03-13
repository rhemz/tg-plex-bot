package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rhemz/tg-plex-bot/config"
	"github.com/rhemz/tg-plex-bot/util"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	disableTimeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1m", "1m"),
			tgbotapi.NewInlineKeyboardButtonData("10m", "10m"),
			tgbotapi.NewInlineKeyboardButtonData("30m", "30m"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1h", "1h"),
			tgbotapi.NewInlineKeyboardButtonData("12h", "12h"),
			tgbotapi.NewInlineKeyboardButtonData("24h", "24h"),
		),
	)
)

type TelegramWebhookController struct{}

func (t TelegramWebhookController) Post(c *gin.Context) {
	//

	zap.S().Info("Got a telegram webhook!")

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		// always 200 so tg doesn't start exponential backoff
		c.JSON(http.StatusOK, gin.H{"error binding telegram JSON payload": err.Error()})
		return
	}

	zap.S().Info("Update ID: ", update.UpdateID)
	err := handleUpdateEvent(update, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error handling request"})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func handleUpdateEvent(update tgbotapi.Update, c *gin.Context) error {
	cfg := config.GetConfig()
	authUsers := cfg.GetStringSlice("telegram.allowedUsers")

	// get ref depending on update type
	var uMsg tgbotapi.Message
	if update.Message != nil {
		uMsg = *update.Message
	} else if update.CallbackQuery != nil {
		uMsg = *update.CallbackQuery.Message
	}

	if !(uMsg.IsCommand() && util.Contains(authUsers, uMsg.From.UserName)) {
		// not a command from an authorized user
		zap.S().Info("Discarding DM from @", uMsg.From.UserName, ": ", uMsg.Text)
		return nil
	}

	zap.S().Info("Got command from authorized user: ", uMsg.From.UserName)

	switch uMsg.Command() {
	case "disable":
		disableMessages(update)
	case "enable":
		enableMessages(update)
	case "transcode":
		// TODO: implement transcode toggle
	}

	return nil
}

func enableMessages(update tgbotapi.Update) {
	// disable and notify
	if config.TelegramMessagesDisabled {
		zap.S().Info("Enabling posting updates")
		_ = util.SendMessageResponse(fmt.Sprintf("Enabling posting updates"), *update.Message)
		config.TelegramMessagesDisabled = false
	} else {
		_ = util.SendMessageResponse(fmt.Sprintf("Posting is not disabled"), *update.Message)
	}
}

func disableMessages(update tgbotapi.Update) {
	// send the keyboard
	if update.CallbackQuery == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyMarkup = disableTimeKeyboard
		if _, err := util.GetTelegramAPI().Send(msg); err != nil {
			zap.S().Error("Error sending disable time inline-keyboard: ", err)
		}
		return
	} else {
		// parse duration
		duration, err := time.ParseDuration(update.CallbackQuery.Data)
		if err != nil {
			zap.S().Error("Could not parse duration from: ", update.CallbackQuery.Data)
			_ = util.SendMessageResponse(fmt.Sprintf("Could not parse duration: %s", update.CallbackQuery.Data), *update.CallbackQuery.Message)
			return
		}

		// ack the callback query
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := util.GetTelegramAPI().Request(callback); err != nil {
			zap.S().Error("Error responding to CallbackQuery: ", err)
			return
		}

		// disable and notify
		zap.S().Info("Disabling posting updates for ", duration)
		_ = util.SendMessageResponse(fmt.Sprintf("Disabling for %s", duration), *update.CallbackQuery.Message)
		config.TelegramMessagesDisabled = true
		time.AfterFunc(duration, func() {
			msg := fmt.Sprintf("Re-enabling telegram messaging!  Sending was disabled for %s", duration)
			zap.S().Info(msg)
			_ = util.SendMessageResponse(msg, *update.CallbackQuery.Message)
			config.TelegramMessagesDisabled = false
		})
	}

}
