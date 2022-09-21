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
	for _, forItem := range world.Usability[item.Name] {
		if forItem == anotherItem.Name {
			canUse = true
		}
	}
	if !canUse {
		fmt.Fprint(answer, "не к чему применить")
		return
	}

	result := world.Events[item.Name+"+"+anotherItem.Name]
	fmt.Fprint(answer, result)

	if result == "дверь открыта" {
		for _, out := range p.Room.Outs {
			if world.Rooms[out].Closed {
				if entry, ok := world.Rooms[out]; ok {
					entry.Closed = false
					world.Rooms[out] = entry
				}
			}
		}
	}
}

type Room struct {
	Name             string
	Items            map[string][]Item
	Outs             []string
	Closed           bool
	GreetingPhrase   string
	LookAroundPhrase string
}

type Player struct {
	Room      Room
	Backpack  bool
	Inventory []Item
}

func (p *Player) LookAround() {

	fmt.Fprint(answer, p.Room.LookAroundPhrase)

	if len(p.Room.Items) != 0 {
		placeNum := 1
		for place, items := range p.Room.Items {
			fmt.Fprint(answer, place, ": ")
			for i, item := range items {
				fmt.Fprint(answer, item.Name)
				switch {
				case i+1 < len(items):
					fmt.Fprint(answer, ", ")
				case placeNum != len(p.Room.Items):
					fmt.Fprint(answer, ", ")
				default:
					fmt.Fprint(answer, "")
				}
			}
			placeNum++
		}
	} else {
		fmt.Fprint(answer, "пустая комната")
	}

	keys := make([]string, 0, len(world.Tasks))
	for k, val := range world.Tasks {
		if !val {
			keys = append(keys, k)
		}
	}
	toDo := strings.Join(keys, " и ")

	if p.Room.Name == world.Rooms[world.StartingRoom].Name {
		fmt.Fprint(answer, ", надо "+toDo)
	}
	fmt.Fprint(answer, ". можно пройти - ")

	for i, out := range p.Room.Outs {
		fmt.Fprint(answer, out)
		if i+1 < len(p.Room.Outs) {
			fmt.Fprint(answer, ",")
		} else {
			fmt.Fprint(answer, "")
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
		fmt.Fprint(answer, "нет пути в ", room.Name)
		return
	}

	if room.Closed {
		fmt.Fprint(answer, "дверь закрыта")
		return
	}
	p.Room = world.Rooms[room.Name]

	fmt.Fprint(answer, room.GreetingPhrase)

	fmt.Fprint(answer, " можно пройти - ")
	for i, out := range p.Room.Outs {
		fmt.Fprint(answer, out)
		if i+1 < len(p.Room.Outs) {
			fmt.Fprint(answer, ", ")
		} else {
			fmt.Fprint(answer, "")
		}
	}

}

func (p *Player) TakeItem(item Items) {
	if !p.Backpack {
		fmt.Fprint(answer, "некуда класть")
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
		fmt.Fprint(answer, "нет такого")
		return
	}

	if tmpItem, ok := item.(Item); ok {
		p.Inventory = append(p.Inventory, tmpItem)
		itemSlice := world.Rooms[p.Room.Name].Items[itemPlace]
		itemSlice = append(itemSlice[:itemIndex], itemSlice[itemIndex+1:]...)

		if len(itemSlice) == 0 {
			delete(world.Rooms[p.Room.Name].Items, itemPlace)
		} else {
			world.Rooms[p.Room.Name].Items[itemPlace] = itemSlice
		}
		fmt.Fprint(answer, "предмет добавлен в инвентарь: ", tmpItem.Name)

	} else {
		fmt.Fprint(answer, "невозможное действие")
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
		fmt.Fprint(answer, "нет такого")
		return
	}

	if tmpItem, ok := item.(Item); ok {
		if !tmpItem.PutOnable {
			fmt.Fprint(answer, "невозможное действие")
			return
		}

		if tmpItem.Name == "рюкзак" {
			p.Backpack = true
			world.Tasks["собрать рюкзак"] = true
		}

		itemSlice := world.Rooms[p.Room.Name].Items[itemPlace]
		itemSlice = append(itemSlice[:itemIndex], itemSlice[itemIndex+1:]...)

		if len(itemSlice) == 0 {
			delete(world.Rooms[p.Room.Name].Items, itemPlace)
		} else {
			world.Rooms[p.Room.Name].Items[itemPlace] = itemSlice
		}

		fmt.Fprint(answer, "вы надели: ", tmpItem.Name)
	} else {
		fmt.Fprint(answer, "невозможное действие")
		return
	}
}

func (p *Player) UseItem(item Items, anotherItem Items) {
	if tmpItem, ok1 := item.(Item); ok1 {
		if tmpAnotherItem, ok2 := anotherItem.(Item); ok2 {
			if !tmpItem.Usable {
				fmt.Fprint(answer, "Невозможно применить")
				return
			}

			var itemPresent bool
			for _, items := range p.Inventory {
				if items == tmpItem {
					itemPresent = true
				}
			}
			if !itemPresent {
				fmt.Fprint(answer, "нет предмета в инвентаре - ", tmpItem.Name)
				return
			}

			item.Use(p, tmpAnotherItem)
		}
	}

}

type World struct {
	Rooms        map[string]Room
	StartingRoom string
	Player       Player
	Tasks        map[string]bool
	Usability    map[string][]string
	Events       map[string]string
	ItemsGlobal  map[string]Item
}

var answer = new(strings.Builder)
var world World

func main() {
}

func initGame() {

	world.Rooms = make(map[string]Room, 3)
	world.Player = Player{}
	world.StartingRoom = "кухня"
	world.Tasks = map[string]bool{
		"собрать рюкзак": false,
		"идти в универ":  false,
	}
	world.Usability = map[string][]string{
		"ключи": {"дверь"},
	}
	world.Events = map[string]string{
		"ключи+дверь": "дверь открыта",
	}
	world.ItemsGlobal = map[string]Item{
		"чай":       {"чай", true, false},
		"ключи":     {"ключи", true, false},
		"конспекты": {"конспекты", true, false},
		"рюкзак":    {"рюкзак", false, true},
		"дверь":     {"дверь", false, false},
		"шкаф":      {"шкаф", false, false},
		"телефон":   {"телефон", true, false},
	}

	var kitchen = Room{
		Name: "кухня",
		Items: map[string][]Item{
			"на столе": {world.ItemsGlobal["чай"]},
		},
		Outs:             []string{"коридор"},
		GreetingPhrase:   "кухня, ничего интересного.",
		LookAroundPhrase: "ты находишься на кухне, ",
	}

	var corridor = Room{
		Name:           "коридор",
		Outs:           []string{"кухня", "комната", "улица"},
		GreetingPhrase: "ничего интересного.",
	}

	var myRoom = Room{
		Name: "комната",
		Items: map[string][]Item{
			"на столе": {world.ItemsGlobal["ключи"], world.ItemsGlobal["конспекты"]},
			"на стуле": {world.ItemsGlobal["рюкзак"]},
		},
		Outs:           []string{"коридор"},
		GreetingPhrase: "ты в своей комнате.",
	}

	var street = Room{
		Name:           "улица",
		Outs:           []string{"домой"},
		Closed:         true,
		GreetingPhrase: "на улице весна.",
	}

	world.Rooms["кухня"] = kitchen
	world.Rooms["коридор"] = corridor
	world.Rooms["комната"] = myRoom
	world.Rooms["улица"] = street
	world.Player.Room = world.Rooms[world.StartingRoom]
}

func handleCommand(command string) string {

	words := strings.Split(command, " ")
	mainCommand := words[0]

	switch mainCommand {
	case "идти":
		world.Player.Go(world.Rooms[words[1]])
	case "взять":
		world.Player.TakeItem(world.ItemsGlobal[words[1]])
	case "осмотреться":
		world.Player.LookAround()
	case "надеть":
		world.Player.PutOnClothes(world.ItemsGlobal[words[1]])
	case "применить":
		world.Player.UseItem(world.ItemsGlobal[words[1]], world.ItemsGlobal[words[2]])
	default:
		return "неизвестная команда"
	}

	toReturn := answer.String()
	answer.Reset()
	return toReturn
}
