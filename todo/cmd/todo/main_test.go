package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	// "runtime"
	"strings"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {

	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build APP: %s", err)
		os.Exit(1)
	}

	fmt.Println("Running tests ...")
	result := m.Run()
	fmt.Println("Cleaning up ...")
	dir, _ := os.Getwd()
	cmdPath := filepath.Join(dir, binName)
	os.Remove(cmdPath)
	cmdPath = filepath.Join(dir, fileName)
	os.Remove(cmdPath)
	os.Exit(result)
}

func TestTodoCli(t *testing.T) {
	task := "Task number one"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, binName)

	t.Run("Add Task Check", func(t *testing.T) {
		taskCmd := strings.Split("-add "+task, " ")
		cmd := exec.Command(cmdPath, taskCmd...)
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to run command: %s", err)
		}
	})

	t.Run("List Task Check", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		result, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run command: %s", err)
		}
		if string(result) != fmt.Sprintf("  1: %s\n", task) {
			t.Fatalf("Expected task name 1: %s, got %s", task, string(result))
		}
	})

	t.Run("Complete Task Check", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to run command: %s", err)
		}
	})
}
