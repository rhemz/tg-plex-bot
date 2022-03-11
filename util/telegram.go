package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var api *tgbotapi.BotAPI

func InitTelegramAPI(botId string, token string) error {
	var err error

	api, err = tgbotapi.NewBotAPI(botId + ":" + token)
	if err != nil {
		return err
	}

	api.Debug = true

	return nil
}

func SendMessageToChats(body string, chats []int) error {
	for _, chat := range chats {
		msg := tgbotapi.NewMessage(int64(chat), body)
		msg.ParseMode = "HTML"
		_, err := api.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendPhotoToChats(photoBytes []uint8, caption string, chats []int) error {
	photoData := tgbotapi.FileBytes{
		Name:  "thumbnail",
		Bytes: photoBytes,
	}
	for _, chat := range chats {
		photoMsg := tgbotapi.NewPhoto(int64(chat), photoData)
		photoMsg.ParseMode = "HTML"
		photoMsg.Caption = caption
		_, err := api.Send(photoMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendMessageResponse(body string, update tgbotapi.Update) error {
	return SendMessageToChats(body, []int{int(update.Message.Chat.ID)})
}

func GetTelegramAPI() *tgbotapi.BotAPI {
	return api
}
