package main

import (
	"github.com/her0elt/florida-man-bot/bot"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	bot.Run()

}
