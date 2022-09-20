package main

import (
	"fmt"
	"strings"
)

type Items interface {
	Use(*Player, Item)
}

type Item struct {
	Name      string
	Usable    bool
	PutOnable bool
}

func (item Item) Use(p *Player, anotherItem Item) {

	var canUse bool
	for _, forItem := range usability[item.Name] {
		if forItem == anotherItem.Name {
			canUse = true
		}
	}
	if !canUse {
		fmt.Print("не к чему применить")
		return
	}

	result := events[item.Name+"+"+anotherItem.Name]
	fmt.Print(result)

	if result == "дверь открыта" {
		for _, out := range p.Room.Outs {
			if rooms[out].Closed {
				if entry, ok := rooms[out]; ok {
					entry.Closed = false
					rooms[out] = entry
				}
			}
		}
	}
}

type Room struct {
	Name   string
	Items  map[string][]Item
	Outs   []string
	Closed bool
}

type Player struct {
	Room      Room
	Backpack  bool
	Inventory []Item
}

func (p *Player) LookAround() {

	if p.Room.Name == rooms["кухня"].Name {
		fmt.Print("ты находишься на кухне, ")
	}

	if len(p.Room.Items) != 0 {
		placeNum := 1
		for place, items := range p.Room.Items {
			fmt.Print(place, ": ")
			for i, item := range items {
				fmt.Print(item.Name)
				if i+1 < len(items) {
					fmt.Print(", ")
				} else if placeNum != len(p.Room.Items) {
					fmt.Print(", ")
				} else {
					fmt.Print("")

				}
			}
			placeNum++
		}
	} else {
		fmt.Print("пустая комната")
	}

	keys := make([]string, 0, len(tasks))
	for k, val := range tasks {
		if !val {
			keys = append(keys, k)
		}
	}
	toDo := strings.Join(keys, " и ")

	if p.Room.Name == rooms["кухня"].Name {
		fmt.Print(", надо " + toDo)
	}
	fmt.Print(". можно пройти - ")

	for i, out := range p.Room.Outs {
		fmt.Print(out)
		if i+1 < len(p.Room.Outs) {
			fmt.Print(",")
		} else {
			fmt.Print("")
		}
	}
}

func (p *Player) Go(room Room) {

	var canReach bool
	for _, out := range p.Room.Outs {
		if out == room.Name {
			canReach = true
		}
	}

	if !canReach {
		fmt.Print("нет пути в ", room.Name)
		return
	}

	if room.Closed {
		fmt.Print("дверь закрыта")
		return
	}
	p.Room = rooms[room.Name]

	switch p.Room.Name {
	case "коридор":
		fmt.Print("ничего интересного.")
	case "комната":
		fmt.Print("ты в своей комнате.")
	case "кухня":
		fmt.Print("кухня, ничего интересного.")
	case "улица":
		fmt.Print("на улице весна.")
	}

	fmt.Print(" можно пройти - ")
	for i, out := range p.Room.Outs {
		fmt.Print(out)
		if i+1 < len(p.Room.Outs) {
			fmt.Print(", ")
		} else {
			fmt.Print("")
		}
	}

}

func (p *Player) TakeItem(item Items) {
	if !p.Backpack {
		fmt.Print("некуда класть")
		return
	}

	var itemPresent bool
	var itemPlace string
	var itemIndex int
	for place, items := range p.Room.Items {
		for i, itemRoom := range items {
			if item == itemRoom {
				itemPlace = place
				itemPresent = true
				itemIndex = i
			}
		}
	}

	if !itemPresent {
		fmt.Print("нет такого")
		return
	}

	if tmpItem, ok := item.(Item); ok {
		p.Inventory = append(p.Inventory, tmpItem)
		itemSlice := rooms[p.Room.Name].Items[itemPlace]
		newItemSlice := append(itemSlice[:itemIndex], itemSlice[itemIndex+1:]...)

		if len(newItemSlice) == 0 {
			delete(rooms[p.Room.Name].Items, itemPlace)
		} else {
			rooms[p.Room.Name].Items[itemPlace] = newItemSlice
		}
		fmt.Print("предмет добавлен в инвентарь: ", tmpItem.Name)

	} else {
		fmt.Print("невозможное действие")
	}

}

func (p *Player) PutOnClothes(item Items) {

	var itemPresent bool
	var itemPlace string
	var itemIndex int
	for place, items := range p.Room.Items {
		for i, itemRoom := range items {
			if item == itemRoom {
				itemPlace = place
				itemPresent = true
				itemIndex = i
			}
		}
	}

	if !itemPresent {
		fmt.Print("нет такого")
		return
	}

	if tmpItem, ok := item.(Item); ok {
		if !tmpItem.PutOnable {
			fmt.Print("невозможное действие")
			return
		}

		if tmpItem.Name == "рюкзак" {
			p.Backpack = true
			tasks["собрать рюкзак"] = true
		}

		itemSlice := rooms[p.Room.Name].Items[itemPlace]
		newItemSlice := append(itemSlice[:itemIndex], itemSlice[itemIndex+1:]...)

		if len(newItemSlice) == 0 {
			delete(rooms[p.Room.Name].Items, itemPlace)
		} else {
			rooms[p.Room.Name].Items[itemPlace] = newItemSlice
		}

		fmt.Print("вы надели: ", tmpItem.Name)
	} else {
		fmt.Print("невозможное действие")
		return
	}
}

func (p *Player) UseItem(item Items, anotherItem Items) {
	if tmpItem, ok1 := item.(Item); ok1 {
		if tmpAnotherItem, ok2 := anotherItem.(Item); ok2 {
			if !tmpItem.Usable {
				fmt.Print("Невозможно применить")
				return
			}

			var itemPresent bool
			for _, items := range p.Inventory {
				if items == tmpItem {
					itemPresent = true
				}
			}
			if !itemPresent {
				fmt.Print("нет предмета в инвентаре - ", tmpItem.Name)
				return
			}

			item.Use(p, tmpAnotherItem)
		}
	}

}

var rooms = make(map[string]Room, 3)
var startingRoom = "кухня"
var player = Player{}
var tasks = map[string]bool{
	"собрать рюкзак": false,
	"идти в универ":  false,
}
var usability = map[string][]string{
	"ключи": {"дверь"},
}
var events = map[string]string{
	"ключи+дверь": "дверь открыта",
}
var itemsGlobal = map[string]Item{
	"чай":       {"чай", true, false},
	"ключи":     {"ключи", true, false},
	"конспекты": {"конспекты", true, false},
	"рюкзак":    {"рюкзак", false, true},
	"дверь":     {"дверь", false, false},
	"шкаф":      {"шкаф", false, false},
	"телефон":   {"телефон", true, false},
}

func main() {
	initGame()

	// handleCommand("осмотреться")
	// handleCommand("идти коридор")
	// handleCommand("идти комната")
	// handleCommand("осмотреться")
	// handleCommand("надеть рюкзак")
	// handleCommand("взять ключи")
	// handleCommand("взять конспекты")
	// handleCommand("идти коридор")
	// handleCommand("применить ключи дверь")
	// handleCommand("идти улица")

	handleCommand("осмотреться")
	// handleCommand("завтракать")
	// handleCommand("идти комната")
	// handleCommand("идти коридор")
	// handleCommand("применить ключи дверь")
	// handleCommand("идти комната")
	// handleCommand("осмотреться")
	// handleCommand("взять ключи")
	// handleCommand("надеть рюкзак")
	// handleCommand("осмотреться")
	// handleCommand("взять ключи")
	// handleCommand("взять телефон")
	// handleCommand("взять ключи")
	// handleCommand("осмотреться")
	// handleCommand("взять конспекты")
	// handleCommand("осмотреться")
	// handleCommand("идти коридор")
	// handleCommand("идти кухня")
	// handleCommand("осмотреться")
	// handleCommand("идти коридор")
	// handleCommand("идти улица")
	// handleCommand("применить ключи дверь")
	// handleCommand("применить телефон шкаф")
	// handleCommand("применить ключи шкаф")
	// handleCommand("идти улица")

	// in := bufio.NewScanner(os.Stdin)
	// for in.Scan() {
	// 	txt := in.Text()
	// 	fmt.Println(txt)
	// }
}

func initGame() { // add var rooms and player emptiness check

	var kitchen Room = Room{
		Name: "кухня",
		Items: map[string][]Item{
			"на столе": {itemsGlobal["чай"]},
		},
		Outs: []string{"коридор"},
	}

	var corridor Room = Room{
		Name: "коридор",
		Outs: []string{"кухня", "комната", "улица"},
	}

	var myRoom Room = Room{
		Name: "комната",
		Items: map[string][]Item{
			"на столе": {itemsGlobal["ключи"], itemsGlobal["конспекты"]},
			"на стуле": {itemsGlobal["рюкзак"]},
		},
		Outs: []string{"коридор"},
	}

	var street Room = Room{
		Name:   "улица",
		Outs:   []string{"домой"},
		Closed: true,
	}

	rooms["кухня"] = kitchen
	rooms["коридор"] = corridor
	rooms["комната"] = myRoom
	rooms["улица"] = street
	player.Room = rooms[startingRoom]
}

func handleCommand(command string) string {

	words := strings.Split(command, " ")
	mainCommand := words[0]

	switch mainCommand {
	case "идти":
		player.Go(rooms[words[1]])
	case "взять":
		player.TakeItem(itemsGlobal[words[1]])
	case "осмотреться":
		player.LookAround()
	case "надеть":
		player.PutOnClothes(itemsGlobal[words[1]])
	case "применить":
		player.UseItem(itemsGlobal[words[1]], itemsGlobal[words[2]])
	default:
		return "неизвестная команда"
	}

	return "not implemented"
}
