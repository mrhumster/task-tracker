package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func help() {
	fmt.Println("ADD task text: create new task")
	fmt.Println("LIST: show all your tasks")
	fmt.Println("DONE id: close task")
	fmt.Println("DELETE id: remove task")
	fmt.Println("EDIT id: edit task")

}

func list() {
	f, err := os.OpenFile(".tasks", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {

	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

}

func add(task string) {
	f, err := os.OpenFile(".tasks", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(task + "\n")
	list()
}

func main() {
	var command, task string
	if len(os.Args) == 1 {
		help()
	} else if len(os.Args) == 2 {
		command = os.Args[1]
	} else {
		command = os.Args[1]
		task = strings.Join(os.Args[2:], " ")
	}

	switch command {
	case "h", "help":
		help()
	case "LIST", "L":
		list()
	case "ADD":
		add(task)
	}
}
