package main

// сюда писать код

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "5440179369:AAEPil19XVCOgmtDOE7d0J94xxGBKlpuSF0"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://d4bb-79-139-208-249.ngrok.io"
)

type User struct {
	TgUser     *tgbotapi.User
	UserChatId int64
}

type Task struct {
	Text     string
	Assignee *User
	Owner    *User
	Id       uint
	Done     bool
}

func (tsk *Task) HasAssignee() bool {

	return tsk.Assignee != nil
}

type TaskList struct {
	TaskList   []Task
	LastTaskId uint
}

type TaskListPrint struct {
	TaskLst *TaskList
	Caller  *tgbotapi.User
}

const NewTaskTenplate = `Задача "{{.Text}}" создана, id={{.Id}}`

func NewMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskText := strings.SplitAfter(update.Message.Text, "/new ")[1]

	taskList.LastTaskId += 1

	newTask := Task{
		Text:     taskText,
		Owner:    &User{update.Message.From, update.FromChat().ChatConfig().ChatID},
		Assignee: nil,
		Id:       taskList.LastTaskId,
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
{{$init_var := .}}
{{range .TaskLst.TaskList}}
	{{.Id}}. {{.Text}} by @{{.Owner.TgUser.UserName}}
	{{if .HasAssignee }}
		{{if (eq $init_var.Caller.UserName  .Assignee.TgUser.UserName)}}
			assignee: я
			/unassign_{{.Id}} /resolve_{{.Id}}
		{{else}}
			assignee: @{{.Assignee.TgUser.UserName}}
		{{end}}
	{{else}}
		/assign_{{.Id}}
	{{end}}
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

	err := tmpl.Execute(&resp, TaskListPrint{
		TaskLst: taskList,
		Caller:  update.Message.From,
	})
	if err != nil {
		fmt.Println("Error executing template in TaskMethod:", err)
	}

	fmt.Println(resp.String())

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		resp.String())
	bot.Send(msg)
}

func AssignMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskIdStr := strings.Split(update.Message.Text, "_")[1]
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		fmt.Println("error casting string to int in AssignMethod")
	}

	for i := 0; i < len(taskList.TaskList); i++ {
		if taskList.TaskList[i].Id == uint(taskId) {

			fmt.Println("task.Assignee", taskList.TaskList[i].Assignee)
			if taskList.TaskList[i].Assignee != nil {
				msg := tgbotapi.NewMessage(
					taskList.TaskList[i].Assignee.UserChatId,
					fmt.Sprintf("Задача \"%s\" назначена на @%s", taskList.TaskList[i].Text, update.Message.From.UserName))
				bot.Send(msg)
			}

			taskList.TaskList[i].Assignee = &User{
				TgUser:     update.SentFrom(),
				UserChatId: update.FromChat().ChatConfig().ChatID,
			}

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Задача \"%s\" назначена на вас", taskList.TaskList[i].Text))
			bot.Send(msg)

			break
		}
	}
}

func UnAssignMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskIdStr := strings.Split(update.Message.Text, "_")[1]
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		fmt.Println("error casting string to int in AssignMethod")
	}

	for _, task := range taskList.TaskList {
		if task.Id == uint(taskId) {
			// if task.Assignee != nil {
			// 	msg := tgbotapi.NewMessage(
			// 		task.Assignee.TgUser.ID,
			// 		fmt.Sprintf("Задача \"%s\" назначена на @%s", task.Text, update.Message.From.UserName))
			// 	bot.Send(msg)
			// }

			task.Assignee = nil
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Задача \"%s\" назначена на вас", task.Text))
			bot.Send(msg)

			break
		}
	}
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

	port := "8081"
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
			AssignMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/unassign"):
			fmt.Println("/unassing", update.Message.Text)
			AssignMethod(update, taskList, bot)

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
