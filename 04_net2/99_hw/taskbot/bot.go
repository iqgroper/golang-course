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
	WebhookURL = "https://a3ec-79-139-208-249.ngrok.io"
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
	Length     int
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
	taskList.Length += 1

	tmpl := template.New("")
	tmpl, _ = tmpl.Parse(NewTaskTenplate)
	var resp bytes.Buffer

	tmpl.Execute(&resp, newTask)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		resp.String())
	bot.Send(msg)

}

// const taskTemplate = `
// {{$init_var := .}}
// {{range $index, $value := .TaskLst.TaskList}}
// 	{{$value.Id}}. {{$value.Text}} by @{{$value.Owner.TgUser.UserName}}
// 	{{if $value.HasAssignee }}
// 		{{if (eq $init_var.Caller.UserName  $value.Assignee.TgUser.UserName)}}
// 			assignee: я
// 			/unassign_{{$value.Id}} /resolve_{{$value.Id}}
// 		{{else}}
// 			assignee: @{{$value.Assignee.TgUser.UserName}}
// 		{{end}}
// 	{{else}}
// 		/assign_{{$value.Id}}
// 	{{end}}
// {{end}}`

const taskTemplate = `{{$init_var := .}}{{range $index, $value := .TaskLst.TaskList}}{{$value.Id}}. {{$value.Text}} by @{{$value.Owner.TgUser.UserName}}
{{if $value.HasAssignee }}{{if (eq $init_var.Caller.UserName  $value.Assignee.TgUser.UserName)}}assignee: я
/unassign_{{$value.Id}} /resolve_{{$value.Id}}{{else}}assignee: @{{$value.Assignee.TgUser.UserName}}{{end}}{{else}}/assign_{{$value.Id}}{{end}}{{if (ne $index $init_var.TaskLst.Length)}}

{{end}}{{end}}`

func TaskMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {
	if len(taskList.TaskList) == 0 {
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
		return
	}

	fmt.Println(resp.String())

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		resp.String())
	bot.Send(msg)
}

func UnAssignMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskIdStr := strings.Split(update.Message.Text, "_")[1]
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		fmt.Println("error casting string to int in UnAssignMethod")
	}

	for i := 0; i < len(taskList.TaskList); i++ {
		if taskList.TaskList[i].Id == uint(taskId) {

			if taskList.TaskList[i].Assignee != nil {

				if taskList.TaskList[i].Assignee.TgUser.ID != update.Message.From.ID {
					msg := tgbotapi.NewMessage(
						update.FromChat().ID,
						"Задача не на вас")

					bot.Send(msg)
					return
				}

				msg := tgbotapi.NewMessage(
					taskList.TaskList[i].Assignee.UserChatId,
					"Принято")

				bot.Send(msg)

				msgOwner := tgbotapi.NewMessage(
					taskList.TaskList[i].Owner.UserChatId,
					fmt.Sprintf("Задача \"%s\" осталась без исполнителя", taskList.TaskList[i].Text))

				bot.Send(msgOwner)
				taskList.TaskList[i].Assignee = nil
				break
			}

		}
	}
}

func AssignMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskIdStr := strings.Split(update.Message.Text, "_")[1]
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		fmt.Println("error casting string to int in AssignMethod")
	}

	for i := 0; i < len(taskList.TaskList); i++ {
		if taskList.TaskList[i].Id == uint(taskId) {

			if taskList.TaskList[i].Assignee != nil {
				msg := tgbotapi.NewMessage(
					taskList.TaskList[i].Assignee.UserChatId,
					fmt.Sprintf("Задача \"%s\" назначена на @%s", taskList.TaskList[i].Text, update.Message.From.UserName))
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(
					taskList.TaskList[i].Owner.UserChatId,
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

func ResolveMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskIdStr := strings.Split(update.Message.Text, "_")[1]
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		fmt.Println("error casting string to int in ResolveMethod")
		return
	}

	for i := 0; i < len(taskList.TaskList); i++ {
		if taskList.TaskList[i].Id == uint(taskId) {

			msg := tgbotapi.NewMessage(
				taskList.TaskList[i].Assignee.UserChatId,
				fmt.Sprintf("Задача \"%s\" выполнена", taskList.TaskList[i].Text))

			bot.Send(msg)

			msgOwner := tgbotapi.NewMessage(
				taskList.TaskList[i].Owner.UserChatId,
				fmt.Sprintf("Задача \"%s\" выполнена @%s", taskList.TaskList[i].Text, taskList.TaskList[i].Assignee.TgUser.UserName))

			bot.Send(msgOwner)

			newSlice := append(taskList.TaskList[:i], taskList.TaskList[i+1:]...)
			taskList.TaskList = newSlice

			taskList.Length -= 1

			break
		}
	}
}

// const MyMethodTemplate = `
// {{$init_var := .}}
// {{range .TaskLst.TaskList}}
// 	{{if .HasAssignee }}
// 		{{if (eq $init_var.Caller.UserName  .Assignee.TgUser.UserName)}}
// 			{{.Id}}. {{.Text}} by @{{.Owner.TgUser.UserName}}
// 			/unassign_{{.Id}} /resolve_{{.Id}}
// 		{{end}}
// 	{{end}}
// {{end}}
// `

const MyMethodTemplate = `{{$init_var := .}}{{range $index, $value := .TaskLst.TaskList}}{{if $value.HasAssignee }}{{if (eq $init_var.Caller.UserName  $value.Assignee.TgUser.UserName)}}{{$value.Id}}. {{$value.Text}} by @{{$value.Owner.TgUser.UserName}}
/unassign_{{$value.Id}} /resolve_{{$value.Id}}{{end}}{{end}}{{if (ne $index $init_var.TaskLst.Length)}}

{{end}}{{end}}`

func MyMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	tmpl := template.New("")
	tmpl, _ = tmpl.Parse(MyMethodTemplate)
	var resp bytes.Buffer

	err := tmpl.Execute(&resp, TaskListPrint{
		TaskLst: taskList,
		Caller:  update.Message.From,
	})
	if err != nil {
		fmt.Println("Error executing template in MyMethod:", err)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		resp.String())
	bot.Send(msg)
}

// const OwnerTemplate = `
// {{$init_var := .}}
// {{range .TaskLst.TaskList}}
// 	{{if (eq $init_var.Caller.UserName  .Owner.TgUser.UserName)}}
// 		{{.Id}}. {{.Text}} by @{{.Owner.TgUser.UserName}}
// 		/assign_{{.Id}}
// 	{{end}}
// {{end}}
// `

const OwnerTemplate = `{{$init_var := .}}{{range .TaskLst.TaskList}}{{if (eq $init_var.Caller.UserName  .Owner.TgUser.UserName)}}{{.Id}}. {{.Text}} by @{{.Owner.TgUser.UserName}}
/assign_{{.Id}}{{end}}{{end}}
`

func OwnerMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	tmpl := template.New("")
	tmpl, _ = tmpl.Parse(OwnerTemplate)
	var resp bytes.Buffer

	err := tmpl.Execute(&resp, TaskListPrint{
		TaskLst: taskList,
		Caller:  update.Message.From,
	})
	if err != nil {
		fmt.Println("Error executing template in OwnerMethod:", err)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		resp.String())
	bot.Send(msg)
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

	taskList := &TaskList{Length: -1}

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
			UnAssignMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/resolve"):
			fmt.Println("/resolve", update.Message.Text)
			ResolveMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/my"):
			fmt.Println("/MY", update.Message.Text)
			MyMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/owner"):
			fmt.Println("/OWNER", update.Message.Text)
			OwnerMethod(update, taskList, bot)
		case strings.Contains(requestMethod, "/start"):
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Privet")
			bot.Send(msg)
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
