package diary_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/bakins/diary"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	l, err := diary.New(nil)
	assert.Nil(t, err)
	assert.NotNil(t, l)
}

func TestSetLevel(t *testing.T) {
	l, err := diary.New(nil, diary.SetLevel(diary.LevelDebug))
	assert.Nil(t, err)
	assert.NotNil(t, l)
}

func TestLog(t *testing.T) {
	var b bytes.Buffer
	l, err := diary.New(nil, diary.SetWriter(&b))
	assert.Nil(t, err)
	assert.NotNil(t, l)

	l.Error("this is the message")
	assert.True(t, strings.Contains(b.String(), `"message":"this is the message"`))
}

func TestDebug(t *testing.T) {
	var b bytes.Buffer
	l, err := diary.New(nil, diary.SetWriter(&b))
	assert.Nil(t, err)
	assert.NotNil(t, l)

	l.Debug("this is the message")
	assert.False(t, strings.Contains(b.String(), `"message":"this is the message"`))
}

func TestContext(t *testing.T) {
	var b bytes.Buffer
	l, err := diary.New(nil, diary.SetWriter(&b))
	assert.Nil(t, err)
	assert.NotNil(t, l)

	l.Error("this is the message", diary.Context{"foo": "bar"})
	assert.True(t, strings.Contains(b.String(), `"message":"this is the message"`))
	assert.True(t, strings.Contains(b.String(), `"foo":"bar"`))
}

func TestChildContext(t *testing.T) {
	var b bytes.Buffer
	l, err := diary.New(nil, diary.SetWriter(&b))
	assert.Nil(t, err)
	assert.NotNil(t, l)

	n, err := l.New(nil)
	assert.Nil(t, err)
	assert.NotNil(t, n)

	n.Error("this is the message", diary.Context{"foo": "bar"})
	assert.True(t, strings.Contains(b.String(), `"message":"this is the message"`))
	assert.True(t, strings.Contains(b.String(), `"foo":"bar"`))
}

func TestChildContextAdd(t *testing.T) {
	var b bytes.Buffer
	l, err := diary.New(nil, diary.SetWriter(&b))
	assert.Nil(t, err)
	assert.NotNil(t, l)

	n, err := l.New(diary.Context{"bar": "baz"})
	assert.Nil(t, err)
	assert.NotNil(t, n)

	n.Error("this is the message", diary.Context{"foo": "bar"})
	assert.True(t, strings.Contains(b.String(), `"message":"this is the message"`))
	assert.True(t, strings.Contains(b.String(), `"foo":"bar"`))
	assert.True(t, strings.Contains(b.String(), `"bar":"baz"`))
}

func TestSetMessageKey(t *testing.T) {
	var b bytes.Buffer
	l, err := diary.New(nil, diary.SetWriter(&b), diary.SetMessageKey("msg"))
	assert.Nil(t, err)
	assert.NotNil(t, l)
	l.Debug("this is the message")
	assert.False(t, strings.Contains(b.String(), `"msg":"this is the message"`))
}

func TestSetCallerKey(t *testing.T) {
	var b bytes.Buffer
	l, err := diary.New(nil, diary.SetWriter(&b), diary.SetCallerKey("caller"))
	assert.Nil(t, err)
	assert.NotNil(t, l)
	l.Info("caller stack")
	fmt.Println(b.String())
	assert.True(t, strings.Contains(b.String(), `"caller":`))
}
