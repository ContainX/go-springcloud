package model

const (
	nonSecureUrlPath       = "http://%s:%d/%s"
	homePageDefaultPath    = "/"
	statusPageDefaultPath  = "/info"
	healthCheckDefaultPath = "health"
)

func NewRegistrationFromInstanceConfig(instance EurekaInstanceConfig) *Instance {
	i := &Instance{}
	i.InstanceId = generateID(instance.AppName)
	i.AppName = instance.AppName
	i.IpAddr = getLocalIP(instance.IpAddress)
	i.HostName = i.IpAddr
	i.VipAddr = instance.AppName
	i.Status = UP
	i.Port = asPort(instance.Port, true)
	i.SecurePort = asPort(instance.SecurePort, false)
	i.HomePageUrl = toInstanceUrlPathToUrl(i.IpAddr, instance.Port, instance.HomePageUrlPath, homePageDefaultPath)
	i.StatusPageUrl = toInstanceUrlPathToUrl(i.IpAddr, instance.Port, instance.StatusPageUrlPath, statusPageDefaultPath)
	i.HealthCheckUrl = toInstanceUrlPathToUrl(i.IpAddr, instance.Port, instance.HealthCheckUrlPath, healthCheckDefaultPath)
	i.DataCenterInfo = defaultDataCenter
	return i
}

func (i *Instance) WrapInRequest() *RegisrationRequest {
	return &RegisrationRequest{Instance: i}
}
