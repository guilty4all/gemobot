package main

import (
	"encoding/json"
	tba "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

var (
	bot      *tba.BotAPI
	botToken = os.Getenv("Bot_Token")
	baseURL  = os.Getenv("App_URL")
)

func initTelegram() {
	var err error

	bot, err = tba.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
		return
	}
	//gfhfghfghfgh

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	url := baseURL + bot.Token
	//_, err = bot.SetWebhook(tba.NewWebhook(url))
	bot.ListenForWebhook(url)
	if err != nil {
		log.Println(err)
	}
}

func webhookHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tba.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	// to monitor changes run: heroku logs --tail
	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// gin router
	router := gin.New()
	router.Use(gin.Logger())

	// telegram
	initTelegram()
	router.POST("/"+bot.Token, webhookHandler)

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
