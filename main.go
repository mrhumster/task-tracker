package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/google/uuid"
)

const (
	Todo       = "todo"
	InProgress = "in-progress"
	Done       = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"desc"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func showError(text string) {
	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fff")).
		Background(lipgloss.Color("#2e2e2e")).
		Width(50).
		Align(lipgloss.Center).
		Margin(1).
		PaddingBottom(1).
		PaddingTop(1)
	fmt.Println(warningStyle.Render(text))
}

func generateID() string {
	return uuid.New().String()
}

func (t *Task) pprint() string {
	todoStyle := lipgloss.NewStyle().Bold(true).Italic(true).Foreground(lipgloss.Color("#005aba"))
	progStyle := lipgloss.NewStyle().Bold(true).Italic(false).Foreground(lipgloss.Color("##005aba"))
	doneStyle := lipgloss.NewStyle().
		Bold(true).
		Italic(false).
		Foreground(lipgloss.Color("#3f478f")).
		Strikethrough(true)
	var style lipgloss.Style
	switch t.Status {
	case Todo:
		style = todoStyle
	case InProgress:
		style = progStyle
	case Done:
		style = doneStyle
	}
	return fmt.Sprintf("%d. %s", t.ID, style.Render(t.Description))
}

func NewTask(id int, description string) *Task {
	now := time.Now().UTC()
	return &Task{
		ID:          id,
		Description: description,
		Status:      "todo",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type TaskList []Task

func (tl *TaskList) save() error {
	data, err := json.Marshal(tl)
	if err != nil {
		return err
	}
	err = os.WriteFile("tasks.json", data, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func (tl *TaskList) load() error {
	f, err := os.Open("tasks.json")
	if os.IsNotExist(err) {
		*tl = TaskList{}
		return nil
	}
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(tl); err != nil {
		if err == io.EOF {
			*tl = TaskList{}
			return nil
		}
		return fmt.Errorf("error decoding JSON: %w", err)
	}
	return nil
}

func (tl *TaskList) add(description string) (int, error) {
	if description == "" {
		return -1, errors.New("description can't empty")
	}
	id := len(*tl) + 1
	newTask := NewTask(id, description)
	*tl = append(*tl, *newTask)
	return newTask.ID, nil
}

func help() {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#202020")).
		Margin(2, 0).
		Padding(2).
		Italic(true).
		Width(50).
		Align(lipgloss.Center)

	text := ` 📙 HELP PAGE
👉🏻 Adding a new task
$ task-cli add "Buy groceries"

👉🏻 Output: Task added successfully (ID: 1)

👉🏻 Updating and deleting tasks
$ task-cli update 1 "Buy groceries and cook dinner"
$task-cli delete 1

👉🏻 Marking a task as in progress or done
$ task-cli mark-in-progress 1
$ task-cli mark-done 1

👉🏻 Listing all tasks
$ task-cli list

👉🏻 Listing tasks by status
$ task-cli list done
$ task-cli list todo
$ task-cli list in-progress
`
	fmt.Println(style.Render(text))
}

func (tl *TaskList) list(filter string) {
	prettyList := list.New().Enumerator(list.Bullet)

	for _, v := range *tl {
		switch {
		case filter == Todo && v.Status == Todo,
			filter == InProgress && v.Status == InProgress,
			filter == Done && v.Status == Done:
			prettyList.Item(v.pprint())
		case filter == "":
			prettyList.Item(v.pprint())
		}
	}
	fmt.Println(prettyList)
}

func (tl *TaskList) deleteByIdx(idx int) {
	newTaskList := &TaskList{}
	for i, v := range *tl {
		if idx != i {
			if i > idx {
				v.ID = v.ID - 1
			}
			*newTaskList = append(*newTaskList, v)
		}
	}
	tl = newTaskList
	tl.save()
}

func notEnoughArgumentException() {
	showError("⚠️ Not enough arguments!!!")
	help()
	os.Exit(1)
}

func wrongIDException() {
	showError("⚠️ Wrong ID")
	os.Exit(1)
}

func main() {
	var command string

	if len(os.Args) == 1 {
		help()
		return
	}
	var tasks TaskList

	if err := tasks.load(); err != nil {
		log.Fatal(err)
	}

	command = os.Args[1]

	switch command {
	case "add":
		description := strings.Join(os.Args[2:], " ")
		if description == "" {
			showError("⚠️ Task description can't be empty")
			help()
			os.Exit(1)
		}
		id, err := tasks.add(description)
		if err != nil {
			panic(err)
		}
		tasks.save()
		fmt.Printf("Task added successfully. (ID: %d)\n", id)
	case "list":
		var filter string
		if len(os.Args) > 2 {
			filter = os.Args[2]
		}
		tasks.list(filter)
		os.Exit(0)
	case "update":
		if len(os.Args) < 3 {
			notEnoughArgumentException()
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil || id < 1 || id > len(tasks) {
			wrongIDException()
		}
		idx := id - 1
		description := strings.Join(os.Args[3:], " ")
		tasks[idx].Description = description
		tasks[idx].UpdatedAt = time.Now().UTC()
		tasks.save()
		fmt.Printf("Task updated succesfully. (ID: %d)\n", id)

	case "delete":
		if len(os.Args) < 2 {
			notEnoughArgumentException()
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil || id < 1 || id > len(tasks) {
			wrongIDException()
		}
		idx := id - 1
		tasks.deleteByIdx(idx)
		fmt.Printf("Task deleted succesfully. (ID: %d)\n", id)
	case "mark-in-progress":
		if len(os.Args) < 2 {
			notEnoughArgumentException()
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil || id < 1 || id > len(tasks) {
			wrongIDException()
		}
		idx := id - 1
		tasks[idx].Status = "in-progress"
		tasks.save()
		fmt.Printf("Task updated succesfully. (ID: %d)\n", id)
	case "mark-done":
		if len(os.Args) < 2 {
			notEnoughArgumentException()
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil || id < 1 || id > len(tasks) {
			wrongIDException()
		}
		idx := id - 1
		tasks[idx].Status = "done"
		tasks.save()
		fmt.Printf("Task updated succesfully. (ID: %d)\n", id)
	}
}
