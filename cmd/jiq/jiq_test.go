package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/fiatjaf/jiq"
	"github.com/stretchr/testify/assert"
)

var called int = 0

func TestMain(m *testing.M) {
	called = 0
	code := m.Run()
	defer os.Exit(code)
}

func TestjiqRun(t *testing.T) {
	var assert = assert.New(t)

	e := &jiq.Engine{}
	result := run(e, false)
	assert.Zero(result)
	assert.Equal(2, called)

	result = run(e, true)
	assert.Equal(1, called)

	result = run(e, false)
	assert.Zero(result)
}

func TestjiqRunWithError(t *testing.T) {
	called = 0
	var assert = assert.New(t)
	e := &jiq.Engine{}
	result := run(e, false)
	assert.Equal(2, result)
	assert.Equal(0, called)
}

func (e *EngineMock) Run() jiq.EngineResultInterface {
	return &jiq.EngineResult{
		err:     fmt.Errorf(""),
		qs:      ".querystring",
		content: `{"test": "result"}`,
	}
}
