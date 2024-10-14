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

func (uh *UpdateHandler) ProcessUpdate(update *tgbotapi.Update) (*tgbotapi.Message, error) {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			message, err := uh.startCommand(update)
			if err != nil {
				return nil, fmt.Errorf("start command: %w", err)
			}

			return message, nil
		case "help":
			message, err := uh.helpCommand(update)
			if err != nil {
				return nil, fmt.Errorf("help command: %w", err)
			}

			return message, nil
		case "random":
			message, err := uh.randomCommand(update)
			if err != nil {
				return nil, fmt.Errorf("random command: %w", err)
			}

			return message, nil
		default:
			message, err := uh.defaultCommand(update)
			if err != nil {
				return nil, fmt.Errorf("default command: %w", err)
			}

			return message, nil
		}
	} else {
		message, err := uh.nonCommand(update)
		if err != nil {
			return nil, fmt.Errorf("non command: %w", err)
		}

		return message, nil
	}
}

func (uh *UpdateHandler) startCommand(update *tgbotapi.Update) (*tgbotapi.Message, error) {
	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "ГООООЛ")

	message, err := uh.bot.Send(messageConfig)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	return &message, nil
}

func (uh *UpdateHandler) helpCommand(update *tgbotapi.Update) (*tgbotapi.Message, error) {
	messageText := "Here are available commands:\n"
	for _, cmd := range uh.commands {
		messageText += fmt.Sprintf("/%s - %s\n", cmd.name, cmd.description)
	}

	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)

	message, err := uh.bot.Send(messageConfig)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	return &message, nil
}

func (uh *UpdateHandler) randomCommand(update *tgbotapi.Update) (*tgbotapi.Message, error) {
	movies, err := uh.movieService.GetRandomMovies(10)
	if err != nil {
		messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to get movies")

		message, err := uh.bot.Send(messageConfig)
		if err != nil {
			return nil, fmt.Errorf("send message: %w", err)
		}

		return &message, nil
	}

	//url := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://i.imgur.com/unQLJIb.jpg"))
	//mediaGroup := tgbotapi.NewMediaGroup(update.Message.Chat.ID, []interface{}{
	//	url,
	//	url,
	//	url,
	//	url,
	//	url,
	//})
	//uh.bot.SendMediaGroup(mediaGroup)

	messageText := "Here are your movies:\n"
	for _, movie := range movies {
		messageText += fmt.Sprintf("%s, %f\n", movie.Title, movie.VoteAverage)
	}

	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)

	message, err := uh.bot.Send(messageConfig)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	return &message, nil
}

func (uh *UpdateHandler) defaultCommand(update *tgbotapi.Update) (*tgbotapi.Message, error) {
	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Type /help to see the available commands.")

	message, err := uh.bot.Send(messageConfig)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	return &message, nil
}

func (uh *UpdateHandler) nonCommand(update *tgbotapi.Update) (*tgbotapi.Message, error) {
	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "I didn't understand that command. Type /help for a list of commands.")

	message, err := uh.bot.Send(messageConfig)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	return &message, nil
}
