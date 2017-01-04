package model

import (
	"github.com/ContainX/go-utils/encoding"
	"github.com/ContainX/go-utils/logger"
)

type EurekaConfig struct {
	// Instance - Application Instance Information
	Instance EurekaInstanceConfig `json:"instance"`
	// Client - Eureka Client Configuration
	Client EurekaClientConfig `json:"client"`
}

type EurekaInstanceConfig struct {
	AppName            string `json:"app"`
	IpAddress          string `json:"ipAddress"`
	HostName           string `json:"hostName"`
	Port               int    `json:"port"`
	SecurePort         int    `json:"securePort"`
	PreferIpAddress    bool   `json:"preferIpAddress"`
	HomePageUrlPath    string `json:"homePageUrlPath"`
	StatusPageUrlPath  string `json:"statusPageUrlPath"`
	HealthCheckUrlPath string `json:"healthCheckUrlPath"`
}

type EurekaClientConfig struct {
	ServiceUrls         []string `json:"serviceUrls"`
	PollIntervalSeconds int      `json:"pollIntervalSeconds"`
	RegisterWithEureka  bool     `json:"registerWithEureka"`
	HealthCheckEnabled  bool     `json:"healthCheckEnabled"`
	Retries             int      `json:"retries"`
}

func (c *EurekaConfig) PopulateDefaults() {
	c.Client.PollIntervalSeconds = 30
	c.Client.Retries = 3
	c.Client.RegisterWithEureka = true
	c.Client.HealthCheckEnabled = true
}

var log = logger.GetLogger("discovery.model")

func NewConfigFromArgs(appName, host string, port int, serviceUrls ...string) *EurekaConfig {
	ec := &EurekaConfig{
		Instance: EurekaInstanceConfig{
			AppName:   appName,
			IpAddress: host,
			Port:      port,
		},
		Client: EurekaClientConfig{
			ServiceUrls: serviceUrls,
		},
	}
	ec.PopulateDefaults()
	return ec
}

func Read(configFile string) (config EurekaConfig, err error) {

	config.PopulateDefaults()

	encoder, err := encoding.NewEncoderFromFileExt(configFile)
	if err != nil {
		log.Fatalf("Error reading config: %s - %s", configFile, err.Error())
		return config, err
	}

	err = encoder.UnMarshalFile(configFile, &config)
	if err != nil {
		log.Fatalf("Error reading config: %s - %s", configFile, err.Error())
		return config, err
	}
	return config, nil
}
