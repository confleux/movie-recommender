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

func (uh *UpdateHandler) ProcessUpdate(update *tgbotapi.Update) error {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			err := uh.startCommand(update)
			if err != nil {
				return fmt.Errorf("start command: %w", err)
			}

			return nil
		case "help":
			err := uh.helpCommand(update)
			if err != nil {
				return fmt.Errorf("help command: %w", err)
			}

			return nil
		case "random":
			err := uh.randomCommand(update)
			if err != nil {
				return fmt.Errorf("random command: %w", err)
			}

			return nil
		default:
			err := uh.defaultCommand(update)
			if err != nil {
				return fmt.Errorf("default command: %w", err)
			}

			return nil
		}
	} else {
		err := uh.nonCommand(update)
		if err != nil {
			return fmt.Errorf("non command: %w", err)
		}

		return nil
	}
}

func (uh *UpdateHandler) startCommand(update *tgbotapi.Update) error {
	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "ГООООЛ")

	_, err := uh.bot.Send(messageConfig)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (uh *UpdateHandler) helpCommand(update *tgbotapi.Update) error {
	messageText := "Here are available commands:\n"
	for _, cmd := range uh.commands {
		messageText += fmt.Sprintf("/%s - %s\n", cmd.name, cmd.description)
	}

	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)

	_, err := uh.bot.Send(messageConfig)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (uh *UpdateHandler) randomCommand(update *tgbotapi.Update) error {
	movies, err := uh.movieService.GetRandomMovies(10)
	if err != nil {
		messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to get movies")

		_, err := uh.bot.Send(messageConfig)
		if err != nil {
			return fmt.Errorf("send message: %w", err)
		}

		return nil
	}

	for _, movie := range movies {
		photoConfig := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileURL(fmt.Sprintf("https://image.tmdb.org/t/p/w1280/%s", movie.PosterPath)))

		messageText := fmt.Sprintf("Title: %s\nAverage Vote: %.2f\nRelease Date: %s\n", movie.Title, movie.VoteAverage, movie.ReleaseDate.Format("02/01/2006"))

		photoConfig.Caption = messageText
		photoConfig.ParseMode = tgbotapi.ModeMarkdown

		inlineKeyboardButton := tgbotapi.NewInlineKeyboardButtonURL("Open TMDB movie page", fmt.Sprintf("https://www.themoviedb.org/movie/%d", movie.Id))
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(inlineKeyboardButton),
		)

		photoConfig.ReplyMarkup = inlineKeyboard

		if _, err := uh.bot.Send(photoConfig); err != nil {
			return fmt.Errorf("send message: %w", err)
		}
	}

	return nil
}

func (uh *UpdateHandler) defaultCommand(update *tgbotapi.Update) error {
	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Type /help to see the available commands.")

	_, err := uh.bot.Send(messageConfig)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (uh *UpdateHandler) nonCommand(update *tgbotapi.Update) error {
	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "I didn't understand that command. Type /help for a list of commands.")

	_, err := uh.bot.Send(messageConfig)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
