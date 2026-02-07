//go:build integration

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func randomTaskName(t *testing.T) string {
	t.Helper()
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var b strings.Builder
	for range 32 {
		b.WriteByte(chars[r.Intn(len(chars))])
	}
	return b.String()
}

func TestIntegration(t *testing.T) {
	url := "http://localhost:8080"
	if os.Getenv("TODO_API_URL") != "" {
		url = os.Getenv("TODO_API_URL")
	}

	today := time.Now().Format("Jan/02")
	task := randomTaskName(t)
	taskId := ""

	t.Run("Add task", func(t *testing.T) {
		args := []string{task}
		var out bytes.Buffer
		expOut := fmt.Sprintf("Task: %s, added to the list\n", task)
		if err := addAction(&out, url, args); err != nil {
			t.Fatal(err)
		}
		if expOut != out.String() {
			t.Errorf("Expected out: %s, got %s", expOut, out.String())
		}
	})

	t.Run("list task", func(t *testing.T) {
		var out bytes.Buffer
		if err := listAction(&out, url); err != nil {
			t.Fatal(err)
		}
		outList := ""
		scanner := bufio.NewScanner(&out)
	l:
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), task) {
				outList = scanner.Text()
				break l
			}
		}
		if outList == "" {
			t.Errorf("Expect task: %s, got nothing, task is out of the list", task)
		}
		taskCompleteStatus := strings.Fields(outList)[0]
		if taskCompleteStatus != "-" {
			t.Errorf("Expect task uncomplete, got %s", taskCompleteStatus)
		}
		taskId = strings.Fields(outList)[1]
	})

	t.Run("View task", func(t *testing.T) {
		var out bytes.Buffer
		if err := viewAction(&out, url, taskId); err != nil {
			t.Fatal(err)
		}
		outStrings := strings.Split(out.String(), "\n")
		if !strings.Contains(outStrings[0], task) {
			t.Errorf("Expect task name %s, got %s", task, outStrings[0])
		}
		if !strings.Contains(outStrings[1], today) {
			t.Errorf("Expect task date %s, got %s", today, outStrings[1])
		}
		if !strings.Contains(outStrings[2], "No") {
			t.Errorf("Expect task Not completed, got %s", outStrings[2])
		}
	})

	t.Run("Complete task", func(t *testing.T) {
		var out bytes.Buffer
		expOut := fmt.Sprintf("Item No %s set as completed", taskId)
		args := []string{taskId}

		if err := completeAction(&out, url, args); err != nil {
			t.Fatal(err)
		}
		if expOut != out.String() {
			t.Errorf("Expected out %s, got %s", expOut, out.String())
		}
	})

	t.Run("list completed task", func(t *testing.T) {
		var out bytes.Buffer
		if err := viewAction(&out, url, taskId); err != nil {
			t.Fatal(err)
		}
		outStrings := strings.Split(out.String(), "\n")
		if !strings.Contains(outStrings[0], task) {
			t.Errorf("Expect task name %s, got %s", task, outStrings[0])
		}
		if !strings.Contains(outStrings[1], today) {
			t.Errorf("Expect task date %s, got %s", today, outStrings[1])
		}
		if !strings.Contains(outStrings[2], "Yes") {
			t.Errorf("Expect task completed, got %s", outStrings[2])
		}

	})

	t.Run("Delete task", func(t *testing.T) {
		var out bytes.Buffer
		args := []string{taskId}
		expOut := fmt.Sprintf("Item No %s removed", taskId)
		if err := removeAction(&out, url, args); err != nil {
			t.Fatal(err)
		}
		if expOut != out.String() {
			t.Errorf("Expected out %s, got %s", expOut, out.String())
		}
	})

	t.Run("View deleted task", func(t *testing.T) {
		var out bytes.Buffer
		if err := listAction(&out, url); err != nil {
			t.Fatal(err)
		}
		outList := ""
		scanner := bufio.NewScanner(&out)
	l:
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), task) {
				outList = scanner.Text()
				break l
			}
		}
		if outList != "" {
			t.Errorf("Expect task: %s, still in the list", task)
		}
	})

}
