package main

// сюда писать код

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "5440179369:AAEPil19XVCOgmtDOE7d0J94xxGBKlpuSF0"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://f517-79-139-208-249.ngrok.io"
)

type User struct {
	TgUser *tgbotapi.User
	// Tasks  []uint
}

type Task struct {
	Text    string
	Asignee *User
	Owner   *User
	Id      uint
}

type TaskList struct {
	TaskList   []Task
	LastTaskId uint
}

const NewTaskTenplate = `Задача {{.Text}} создана, id={{.Id}}`

func NewMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskText := strings.SplitAfter(update.Message.Text, "/new ")[1]

	taskList.LastTaskId += 1

	newTask := Task{
		Text:    taskText,
		Owner:   &User{update.Message.From},
		Asignee: nil,
		Id:      taskList.LastTaskId,
	}
	taskList.TaskList = append(taskList.TaskList, newTask)

	tmpl := template.New("")
	tmpl, _ = tmpl.Parse(NewTaskTenplate)
	var resp bytes.Buffer

	tmpl.Execute(&resp, newTask)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		resp.String())
	bot.Send(msg)

}

const taskTemplate = `
{{range .TaskList}}
	{{.Id}}. {{.Text}} by @{{.Owner.TgUser.UserName}}
	asignee: {{.Asignee}}
{{end}}
`

func TaskMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {
	if taskList.LastTaskId == 0 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Нет задач")
		bot.Send(msg)
		return
	}

	tmpl := template.New("")
	tmpl, _ = tmpl.Parse(taskTemplate)
	var resp bytes.Buffer

	err := tmpl.Execute(&resp, *taskList)
	if err != nil {
		fmt.Println("Error executing template in  TaskMethod:", err)
	}

	fmt.Println(resp.String())

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		resp.String())
	bot.Send(msg)

	// command := update.Message.Text
	// switch {
	// case strings.Contains(command, "/my"):

	// case strings.Contains(command, "/owner"):
	// }

}

func AssignMethod(command string) {

}

func ResolveMethod(command string) {

}

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

	taskList := &TaskList{}

	for update := range updates {

		requestMethod := update.Message.Text

		switch {
		case strings.Contains(requestMethod, "/new"):
			fmt.Println("/new", update.Message.Text)
			NewMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/tasks"):
			fmt.Println("/tasks", update.Message.Text)
			TaskMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/assign"):
			fmt.Println("/assign", update.Message.Text)
			AssignMethod(requestMethod)

		case strings.Contains(requestMethod, "/unassign"):
			fmt.Println("/unassing", update.Message.Text)
			AssignMethod(requestMethod)

		case strings.Contains(requestMethod, "/resolve"):
			fmt.Println("/resolve", update.Message.Text)
			ResolveMethod(requestMethod)

		case strings.Contains(requestMethod, "/my"):
			fmt.Println("/MY", update.Message.Text)
			TaskMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/owner"):
			fmt.Println("/OWNER", update.Message.Text)
			TaskMethod(update, taskList, bot)

		default:
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"wrong api hit, try..")
			bot.Send(msg)
		}

	}

	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
