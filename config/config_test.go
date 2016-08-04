package config

import (
	"testing"
	"github.com/ContainX/go-utils/mockrest"
	"github.com/stretchr/testify/assert"
	"strings"
)

const (
	FilePropertyResp = "testdata/GET-properties_response.txt"
	FileJsonResp = "testdata/GET-response.json"
)

type testStruct struct {
	Foo string `json:"foo"`
}

func TestConfigAsMap(t *testing.T) {
	server := mockrest.StartNewWithFile(FilePropertyResp)
	defer server.Stop()

	// Configure Bootstrap - use dynamic mockrest server host/port
	cfg, err := New(Bootstrap{ Name: "myapp", URI: server.Start() })
	assert.NoError(t, err)

	// Test fetch remote config as Map
	m, err := cfg.FetchAsMap()
	assert.NoError(t, err)
	assert.Equal(t, "bar", m["foo"])
}

func TestConfigAsStruct(t *testing.T) {
	server := mockrest.StartNewWithFile(FileJsonResp)
	defer server.Stop()

	// Configure Bootstrap - use dynamic mockrest server host/port
	cfg, err := New(Bootstrap{ Name: "myapp", URI: server.Start() })
	assert.NoError(t, err)

	ts := &testStruct{}

	err = cfg.Fetch(ts)
	assert.NoError(t, err)

	assert.Equal(t, "bar", ts.Foo)
}

func TestConfigAsJSON(t *testing.T) {
	server := mockrest.StartNewWithFile(FileJsonResp)
	defer server.Stop()

	// Configure Bootstrap - use dynamic mockrest server host/port
	cfg, err := New(Bootstrap{ Name: "myapp", URI: server.Start() })

	if assert.NoError(t, err) {
		json, _ := cfg.FetchAsJSON()
		assert.True(t, strings.Contains(json, `"foo": "bar"`))
	}
}

func TestLoadFromFile(t  *testing.T) {
	c, err := LoadFromFile("testdata/LOAD-test.json")

	if assert.NoError(t, err, "") {
		assert.Equal(t, "http://test", c.Bootstrap().URI)
		assert.Equal(t, "/test", c.Bootstrap().Context)
	}
}
