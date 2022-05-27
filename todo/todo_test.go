package todo_test

import (
	"testing"

	"github.com/bazsup/todoapi/todo"
)

func TestNewTodoNotAllowSleep(t *testing.T) {
	handler := todo.NewTodoHandler(&TestDB{})
	c := &TestContext{}

	handler.NewTask(c)

	want := "not allowed"

	if want != c.v["error"] {
		t.Errorf("want %s but get %s", want, c.v["error"])
	}
}

type TestDB struct{}

func (TestDB) New(t *todo.Todo) error {
	return nil
}

type TestContext struct {
	v map[string]interface{}
}

func (TestContext) Bind(v interface{}) error {
	*v.(*todo.Todo) = todo.Todo{Title: "sleep"}
	return nil
}
func (TestContext) TransactionID() string {
	return "TestTransactionID"
}
func (TestContext) Audience() string {
	return "Unit Test"
}
func (c *TestContext) JSON(code int, v interface{}) {
	c.v = v.(map[string]interface{})
}
