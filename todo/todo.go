package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

func NewList() *List {
	return &List{}
}

func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

func (l *List) Complete(i int) error {
	list := *l
	if i <= 0 || i > len(list) {
		return fmt.Errorf("item %d does not exist", i)
	}
	list[i-1].Done = true
	list[i-1].CompletedAt = time.Now()
	return nil
}

func (l *List) Delete(i int) error {
	list := *l
	if i <= 0 || i > len(*l) {
		return fmt.Errorf("item %d does not exist", i)
	}
	*l = append(list[:i-1], list[i:]...)
	return nil
}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, js, 0644)
}

func (l *List) GetFile(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || len(file) == 0 {
			return nil
		}
		return err
	}
	return json.Unmarshal(file, l)
}

func (l *List) Get(i int) (*item, error) {
	list := *l
	if i <= 0 || i > len(list) {
		return nil, fmt.Errorf("item %d does not exist", i)
	}
	return &list[i-1], nil
}

func (l *List) String() string {
	res := ""
	prefix := "  "
	for k, item := range *l {
		if item.Done {
			prefix = "X "
		} else {
			prefix = "  "
		}
		res += fmt.Sprintf("%s%d: %s\n", prefix, k+1, item.Task)
	}
	return res
}
