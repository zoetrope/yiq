package jiq

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJqQueries(t *testing.T) {
	var assert = assert.New(t)

	res, _ := jqrun(".a", `{"a": 23}`, nil)
	assert.Equal("23", res)

	res, _ = jqrun(`.a | {pli: ., plu: [12, "I have \(.) horses"]}`, `{"a": 23}`, nil)
	assert.Equal(`{
  "pli": 23,
  "plu": [
    12,
    "I have 23 horses"
  ]
}`, res)

	res, _ = jqrun(".a[].w", `{"a": [{"w": 18}, {"w": false}]}`, nil)
	assert.Equal("18\nfalse", res)

	res, _ = jqrun(".a", `{"a": [{"w": 18}, {"w": false}]}`, []string{"-c"})
	assert.Equal(`[{"w":18},{"w":false}]`, res)
}

func TestJqModules(t *testing.T) {
	var assert = assert.New(t)

	dir, err := ioutil.TempDir("", "")
	if !assert.NoError(err, "error creating tempdir") {
		return
	}

	defer os.RemoveAll(dir)

	content := []byte(`def hello(f): f ;`)

	err = ioutil.WriteFile(filepath.Join(dir, ".jq"), content, 0666)
	if !assert.NoError(err, "error creating tempfile") {
		return
	}

	defer os.Setenv("HOME", os.Getenv("HOME"))
	os.Setenv("HOME", dir)

	res, _ := jqrun("hello(.hi)", `{"hi": "world"}`, nil)
	assert.Equal(`"world"`, res)
}
