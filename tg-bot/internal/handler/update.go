package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-bot/internal/service"
)

type command struct {
	name        string
	description string
}

type UpdateHandler struct {
	bot          *tgbotapi.BotAPI
	movieService *service.MovieService
	commands     []command
}

func NewUpdateHandler(bot *tgbotapi.BotAPI, movieService *service.MovieService) *UpdateHandler {
	commands := []command{
		{name: "start", description: "Start the bot"},
		{name: "help", description: "Show help message"},
		{name: "random", description: "Get random message"},
	}
	return &UpdateHandler{bot: bot, movieService: movieService, commands: commands}
}

func (uh *UpdateHandler) ProcessUpdate(update *tgbotapi.Update) {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			uh.startCommand(update)
		case "help":
			uh.helpCommand(update)
		case "random":
			uh.randomCommand(update)
		default:
			uh.defaultCommand(update)
		}
	} else {
		uh.nonCommand(update)
	}
}

func (uh *UpdateHandler) startCommand(update *tgbotapi.Update) {
	//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ГООООЛ")
	url := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://i.imgur.com/unQLJIb.jpg"))
	mediaGroup := tgbotapi.NewMediaGroup(update.Message.Chat.ID, []interface{}{
		url,
		url,
		url,
		url,
		url,
	})
	uh.bot.SendMediaGroup(mediaGroup)
	//uh.bot.Send(url)
}

func (uh *UpdateHandler) helpCommand(update *tgbotapi.Update) {
	msgMessage := "Here are available commands:\n"
	for _, cmd := range uh.commands {
		msgMessage += fmt.Sprintf("/%s - %s\n", cmd.name, cmd.description)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgMessage)
	uh.bot.Send(msg)
}

func (uh *UpdateHandler) randomCommand(update *tgbotapi.Update) {
	movies, err := uh.movieService.GetRandomMovies(10)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to get movies")
		uh.bot.Send(msg)
	}

	msgString := "Here are your movies:\n"
	for _, movie := range movies {
		msgString += fmt.Sprintf("%s, %f\n", movie.Title, movie.VoteAverage)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgString)
	uh.bot.Send(msg)
}

func (uh *UpdateHandler) defaultCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Type /help to see the available commands.")
	uh.bot.Send(msg)
}

func (uh *UpdateHandler) nonCommand(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I didn't understand that command. Type /help for a list of commands.")
	uh.bot.Send(msg)
}
