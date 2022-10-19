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
	WebhookURL = "https://35be-79-139-208-249.ngrok.io"
)

type User struct {
	TgUser     *tgbotapi.User
	UserChatID int64
}

type Task struct {
	Text     string
	Assignee *User
	Owner    *User
	ID       uint
	Done     bool
}

func (tsk *Task) HasAssignee() bool {
	return tsk.Assignee != nil
}

type TaskList struct {
	TaskList   []Task
	LastTaskID uint
}

type TaskListPrint struct {
	TaskLst *TaskList
	Caller  *tgbotapi.User
}

const NewTaskTenplate = `Задача "{{.Text}}" создана, id={{.ID}}`

func NewMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskText := strings.SplitAfter(update.Message.Text, "/new ")[1]

	taskList.LastTaskID += 1

	newTask := Task{
		Text:     taskText,
		Owner:    &User{update.Message.From, update.FromChat().ChatConfig().ChatID},
		Assignee: nil,
		ID:       taskList.LastTaskID,
	}
	taskList.TaskList = append(taskList.TaskList, newTask)

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(NewTaskTenplate)
	if errParse != nil {
		fmt.Println("Error parsing New method", errParse)
	}
	var resp bytes.Buffer

	err := tmpl.Execute(&resp, newTask)
	if err != nil {
		fmt.Println("Error executing template in MyMethod:", err)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		resp.String())
	_, errorBot := bot.Send(msg)
	if errorBot != nil {
		fmt.Println("Error sending massege", errorBot)
	}

}

const taskTemplate = `{{$init_var := .}}{{$first_occ := 1}}{{range $index, $value := .TaskLst.TaskList}}{{if $first_occ}}{{$first_occ = 0}}{{else}}

{{end}}{{$value.ID}}. {{$value.Text}} by @{{$value.Owner.TgUser.UserName}}
{{if $value.HasAssignee }}{{if (eq $init_var.Caller.UserName  $value.Assignee.TgUser.UserName)}}assignee: я
/unassign_{{$value.ID}} /resolve_{{$value.ID}}{{else}}assignee: @{{$value.Assignee.TgUser.UserName}}{{end}}{{else}}/assign_{{$value.ID}}{{end}}{{end}}`

func TaskMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {
	if len(taskList.TaskList) == 0 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Нет задач")
		_, errorBot := bot.Send(msg)
		if errorBot != nil {
			fmt.Println("Error sending massege", errorBot)
		}
		return
	}

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(taskTemplate)
	if errParse != nil {
		fmt.Println("Reeor parsing task method", errParse)
	}
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
	_, errorBot := bot.Send(msg)
	if errorBot != nil {
		fmt.Println("Error sending massege", errorBot)
	}
}

func UnAssignMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskIDStr := strings.Split(update.Message.Text, "_")[1]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		fmt.Println("error casting string to int in UnAssignMethod")
	}

	for i := 0; i < len(taskList.TaskList); i++ {
		if taskList.TaskList[i].ID == uint(taskID) {

			if taskList.TaskList[i].Assignee != nil {

				if taskList.TaskList[i].Assignee.TgUser.ID != update.Message.From.ID {
					msg := tgbotapi.NewMessage(
						update.FromChat().ID,
						"Задача не на вас")

					_, errorBot := bot.Send(msg)
					if errorBot != nil {
						fmt.Println("Error sending massege", errorBot)
					}
					return
				}

				msg := tgbotapi.NewMessage(
					taskList.TaskList[i].Assignee.UserChatID,
					"Принято")

				_, errorBot := bot.Send(msg)
				if errorBot != nil {
					fmt.Println("Error sending massege", errorBot)
				}

				msgOwner := tgbotapi.NewMessage(
					taskList.TaskList[i].Owner.UserChatID,
					fmt.Sprintf("Задача \"%s\" осталась без исполнителя", taskList.TaskList[i].Text))

				_, errorBot2 := bot.Send(msgOwner)
				if errorBot != nil {
					fmt.Println("Error sending massege", errorBot2)
				}
				taskList.TaskList[i].Assignee = nil
				break
			}

		}
	}
}

func AssignMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskIDStr := strings.Split(update.Message.Text, "_")[1]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		fmt.Println("error casting string to int in AssignMethod")
	}

	for i := 0; i < len(taskList.TaskList); i++ {
		if taskList.TaskList[i].ID == uint(taskID) {

			if taskList.TaskList[i].Assignee != nil {
				msg := tgbotapi.NewMessage(
					taskList.TaskList[i].Assignee.UserChatID,
					fmt.Sprintf("Задача \"%s\" назначена на @%s", taskList.TaskList[i].Text, update.Message.From.UserName))
				_, errorBot := bot.Send(msg)
				if errorBot != nil {
					fmt.Println("Error sending massege", errorBot)
				}
			} else {
				msg := tgbotapi.NewMessage(
					taskList.TaskList[i].Owner.UserChatID,
					fmt.Sprintf("Задача \"%s\" назначена на @%s", taskList.TaskList[i].Text, update.Message.From.UserName))
				_, errorBot := bot.Send(msg)
				if errorBot != nil {
					fmt.Println("Error sending massege", errorBot)
				}
			}

			taskList.TaskList[i].Assignee = &User{
				TgUser:     update.SentFrom(),
				UserChatID: update.FromChat().ChatConfig().ChatID,
			}

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Задача \"%s\" назначена на вас", taskList.TaskList[i].Text))
			_, errorBot := bot.Send(msg)
			if errorBot != nil {
				fmt.Println("Error sending massege", errorBot)
			}

			break
		}
	}
}

func ResolveMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	taskIDStr := strings.Split(update.Message.Text, "_")[1]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		fmt.Println("error casting string to int in ResolveMethod")
		return
	}

	for i := 0; i < len(taskList.TaskList); i++ {
		if taskList.TaskList[i].ID == uint(taskID) {

			msg := tgbotapi.NewMessage(
				taskList.TaskList[i].Assignee.UserChatID,
				fmt.Sprintf("Задача \"%s\" выполнена", taskList.TaskList[i].Text))

			_, errorBot := bot.Send(msg)
			if errorBot != nil {
				fmt.Println("Error sending massege", errorBot)
			}

			msgOwner := tgbotapi.NewMessage(
				taskList.TaskList[i].Owner.UserChatID,
				fmt.Sprintf("Задача \"%s\" выполнена @%s", taskList.TaskList[i].Text, taskList.TaskList[i].Assignee.TgUser.UserName))

			_, errorBot2 := bot.Send(msgOwner)
			if errorBot2 != nil {
				fmt.Println("Error sending massege", errorBot2.Error())
			}

			taskList.TaskList = append(taskList.TaskList[:i], taskList.TaskList[i+1:]...)

			break
		}
	}
}

const MyMethodTemplate = `{{$init_var := .}}{{$first_occ := 1}}{{range $index, $value := .TaskLst.TaskList}}{{if $value.HasAssignee }}{{if (eq $init_var.Caller.UserName  $value.Assignee.TgUser.UserName)}}{{if $first_occ}}{{$first_occ = 0}}{{else}}

{{end}}{{$value.ID}}. {{$value.Text}} by @{{$value.Owner.TgUser.UserName}}
/unassign_{{$value.ID}} /resolve_{{$value.ID}}{{end}}{{end}}{{end}}`

func MyMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	tmpl := template.New("")
	tmpl, errParse := tmpl.Parse(MyMethodTemplate)
	if errParse != nil {
		fmt.Println("Error parsing Owner tempalte", errParse)
	}
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
	_, errorBot := bot.Send(msg)
	if errorBot != nil {
		fmt.Println("Error sending massege", errorBot)
	}
}

const OwnerTemplate = `{{$init_var := .}}{{$first_occ := 1}}{{range $index, $value := .TaskLst.TaskList}}{{if (eq $init_var.Caller.UserName  $value.Owner.TgUser.UserName)}}{{if $first_occ}}{{$first_occ = 0}}{{else}}

{{end}}{{$value.ID}}. {{$value.Text}} by @{{$value.Owner.TgUser.UserName}}
/assign_{{$value.ID}}{{end}}{{end}}`

func OwnerMethod(update tgbotapi.Update, taskList *TaskList, bot *tgbotapi.BotAPI) {

	tmplt := template.New("")
	_, errParse := tmplt.Parse(OwnerTemplate)
	if errParse != nil {
		fmt.Println("Error parsing Owner tempalte", errParse)
	}
	var resp bytes.Buffer

	err := tmplt.Execute(&resp, TaskListPrint{
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
	_, errorBot := bot.Send(msg)
	if errorBot != nil {
		fmt.Println("Error sending massege", errorBot)
	}
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
			NewMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/tasks"):
			TaskMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/assign"):
			AssignMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/unassign"):
			UnAssignMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/resolve"):
			ResolveMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/my"):
			MyMethod(update, taskList, bot)

		case strings.Contains(requestMethod, "/owner"):
			OwnerMethod(update, taskList, bot)
		case strings.Contains(requestMethod, "/start"):
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Privet")
			_, errorBot := bot.Send(msg)
			if errorBot != nil {
				fmt.Println("Error sending massege", errorBot)
			}
		default:
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"wrong api hit, try..")
			_, errorBot := bot.Send(msg)
			if errorBot != nil {
				fmt.Println("Error sending massege", errorBot)
			}
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
