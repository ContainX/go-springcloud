package model

import (
	"encoding/xml"
)

type StatusType string
type DataCenterType uint8
type Metadata map[string]string

//type PortType map[string]interface{}
type DataCenterInfoType map[string]string

// Supported statuses
const (
	UP           StatusType = "UP"
	DOWN         StatusType = "DOWN"
	STARTING     StatusType = "STARTING"
	OUTOFSERVICE StatusType = "OUT_OF_SERVICE"
	UNKNOWN      StatusType = "UNKNOWN"
)

// Datacenter names
const (
	Amazon              = "Amazon"
	MyOwn               = "MyOwn"
	DataCenterInfoClass = "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo"
)

var defaultDataCenter = DataCenterInfo{
	Name:      MyOwn,
	ClassName: DataCenterInfoClass,
}

type ApplicationsResponse struct {
	Response *Applications `json:"applications"`
}

type ApplicationResponse struct {
	Response *Application `json:"application"`
}

type Applications struct {
	Applications  []*Application `json:"application"`
	AppsHashcode  string         `json:"apps__hashcode"`
	VersionsDelta int            `json:"versions__delta"`
}

type Application struct {
	XMLName   xml.Name    `json:"application"`
	Name      string      `json:"name"`
	Instances []*Instance `json:"instance"`
}

type RegisrationRequest struct {
	Instance *Instance `json:"instance"`
}

type Instance struct {
	InstanceId     string         `json:"instanceId"`
	HostName       string         `json:"hostName"`
	AppName        string         `json:"app"`
	IpAddr         string         `json:"ipAddr"`
	VipAddr        string         `json:"vipAddress"`
	Status         StatusType     `json:"status"`
	Port           Port           `json:"port"`
	SecurePort     Port           `json:"securePort"`
	HomePageUrl    string         `json:"homePageUrl"`
	StatusPageUrl  string         `json:"statusPageUrl"`
	HealthCheckUrl string         `json:"healthCheckUrl"`
	DataCenterInfo DataCenterInfo `json:"dataCenterInfo,omitempty"`
	Metadata       Metadata       `json:"metadata,omitempty"`
}

type Registry struct {
	XMLName      xml.Name       `json:"applications"`
	VersionDelta int            `json:"versions__delta"`
	Hashcode     string         `json:"apps__hashcode"`
	Apps         []*Application `json:"application"`
}

type AmazonMetadata struct {
	Hostname         string `json:"hostname'`
	PublicHostName   string `json:"public-hostname"`
	LocalHostName    string `json:"local-hostname"`
	PublicIpv4       string `json:"public-ipv4'`
	LocalIpv4        string `json:"local-ipv4"`
	AvailabilityZone string `json:"availability-zone"`
	InstanceId       string `json:"instance-id"`
	InstanceType     string `json:"instance-type"`
	AmiId            string `json:"ami-id"`
	AmiLaunchIndex   string `json:"ami-launch-index"`
	AmiManifestPath  string `json:"ami-manifest-path"`
}

type DataCenterInfo struct {
	ClassName string          `json:"@class"`
	Name      string          `json:"name"`
	Metadata  *AmazonMetadata `json:"metadata,omitempty"`
}

type Port struct {
	Number  int  `json:"$"`
	Enabled bool `json:"@enabled"`
}

type LeaseInfo struct {
	RenewalIntervalInSecs int32 `json:"renewalIntervalInSecs"`
	DurationInSecs        int32 `json:"durationInSecs"`
	RegistrationTimestamp int64 `json:"registrationTimestamp"`
	LastRenewalTimestamp  int64 `json:"lastRenewalTimestamp"`
	EvictionTimestamp     int64 `json:"evictionTimestamp"`
	ServiceUpTimestamp    int64 `json:"serviceUpTimestamp"`
}
