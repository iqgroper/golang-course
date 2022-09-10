package main

import (
	"fmt"

	"gitlab.com/vk-go/lectures-2022-2/01_intro/05_visibility/person"
)

func main() {
	p := person.NewPerson(1, "vasya", "secret")

	// p.secret undefined (cannot refer to unexported field or method secret)
	// fmt.Printf("main.PrintPerson: %+v\n", p.secret)

	secret := person.GetSecret(p)
	fmt.Println("GetSecret", secret)
}
