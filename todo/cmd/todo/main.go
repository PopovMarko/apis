package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"pragprog.com/rggo/interacting/todo"
)

var todoFile = "todo.json"

func main() {
	if os.Getenv("TODO_FILE_NAME") != "" {
		todoFile = os.Getenv("TODO_FILE_NAME")
	}
	add := flag.Bool("add", false, "Add task to the todo list by StdIn")
	list := flag.Bool("list", false, "List all todo item's names")
	complete := flag.Int("complete", 0, "Mark a task as completed by task number")
	del := flag.Int("delete", 0, "Delete a task by task number")
	get := flag.Int("get", 0, "Get a particular task by task number")
	flag.Parse()

	todolist := todo.NewList()

	if err := todolist.GetFile(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, "Warning: Could not open todo file", err)
		os.Exit(1)
	}
	switch {
	case *list:
		fmt.Print(todolist)
	case *add:
		taskText, err := GetTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Warning: Could not read the task", err)
			os.Exit(1)
		}
		todolist.Add(taskText)
		if err := todolist.Save(todoFile); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: Could not save the file", err)
			os.Exit(1)
		}
	case *complete > 0:
		if err := todolist.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: Could not mark as a completed", err)
			os.Exit(1)
		}
		if err := todolist.Save(todoFile); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: Could not save the file", err)
		}
	case *del > 0:
		if err := todolist.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: Could not delete the task", err)
			os.Exit(1)
		}
		if err := todolist.Save(todoFile); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: Could not save the file", err)
		}
	case *get > 0:
		item, err := todolist.Get(*get)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Warning: Could not get the task from the list", err)
			os.Exit(1)
		}
		bytes, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Warning: Could not marshal the item", err)
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	default:
		fmt.Fprintln(os.Stderr, "Indalid option see todo -h tor help")
		os.Exit(1)
	}
}

func GetTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}
	return s.Text(), nil
}
