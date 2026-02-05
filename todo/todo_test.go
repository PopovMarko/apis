package todo_test

import (
	"os"
	"pragprog.com/rggo/interacting/todo"
	"testing"
)

func TestAdd(t *testing.T) {
	list := todo.List{}
	task := "New Task"

	list.Add(task)

	if len(list) != 1 {
		t.Errorf("Expected list length of 1, got %d", len(list))
	}

	if list[0].Task != task {
		t.Errorf("Expected task name %s, got %s", task, list[0].Task)
	}
}

func TestComplete(t *testing.T) {
	list := todo.List{}
	list.Add("Check todo to complete task")
	if list[0].Done {
		t.Errorf("Expected task to be incomplete at addition")
	}
	list.Complete(1)
	if !list[0].Done {
		t.Errorf("Expected task to be marked as complete")
	}
	if list[0].CompletedAt.IsZero() {
		t.Errorf("Expected CompletedAt time to be set")
	}
}

func TestDelete(t *testing.T) {
	list := todo.List{}
	list.Add("Check todo to Delete task")
	if len(list) == 1 {
		list.Delete(1)
		if len(list) != 0 {
			t.Errorf("Expected list length of 0 got %d, not deleted item", len(list))
		}
	} else {
		t.Errorf("Add method failed to add item")
	}
}

func TestSaveGet(t *testing.T) {
	list := todo.List{}
	task := "Task to check Save and Get"
	list.Add(task)
	list.Save("test.json")
	list.Delete(1)
	list.GetFile("test.json")
	if list[0].Task != task {
		t.Errorf("Expected task name %s, got %s", task, list[0].Task)
	}
	if err := os.Remove("test.json"); err != nil {
		t.Errorf("Failed to remove test file %s", err)
	}
}
