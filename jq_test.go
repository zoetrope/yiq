package jiq

import (
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
