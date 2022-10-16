package main

// сюда писать код

import (
	"context"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "5440179369:AAEPil19XVCOgmtDOE7d0J94xxGBKlpuSF0"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://9a92-79-139-208-249.ngrok.io"
)

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
	}

	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	port := "8080"
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	for update := range updates {
		fmt.Println(update.Message.Text)

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"/assign",
		)
		bot.Send(msg)

	}

	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
