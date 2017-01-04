package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestServiceUrl         = "http://1.1.1.1:8761"
	TestPollInterval       = 20
	TestRegisterWithEureka = true
	TestRetries            = 4
)

func TestEurekaConfigAsYaml(t *testing.T) {
	cfg, err := Read("testdata/test_config.yml")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(cfg.Client.ServiceUrls))
		assert.Equal(t, TestServiceUrl, cfg.Client.ServiceUrls[0])
		assert.Equal(t, TestPollInterval, cfg.Client.PollIntervalSeconds)
		assert.Equal(t, TestRegisterWithEureka, cfg.Client.RegisterWithEureka)
		assert.Equal(t, TestRetries, cfg.Client.Retries)
		assert.Equal(t, true, cfg.Client.HealthCheckEnabled)
	}
}

func TestEurekaConfigAsJson(t *testing.T) {
	cfg, err := Read("testdata/test_config.json")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(cfg.Client.ServiceUrls))
		assert.Equal(t, TestServiceUrl, cfg.Client.ServiceUrls[0])
		assert.Equal(t, TestPollInterval, cfg.Client.PollIntervalSeconds)
		assert.Equal(t, TestRegisterWithEureka, cfg.Client.RegisterWithEureka)
		assert.Equal(t, TestRetries, cfg.Client.Retries)
	}
}
