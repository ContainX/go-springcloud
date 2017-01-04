package eureka

import (
	"fmt"
	"github.com/ContainX/go-springcloud/discovery/eureka/model"
	"github.com/ContainX/go-utils/httpclient"
	"net/http"
)

func (e *eureka) GetCurrentInstance() (*model.Instance, error) {
	log.Infof("Calling current instance with: %s, %s", e.instance.AppName, e.instance.InstanceId)
	return e.GetInstance(e.instance.AppName, e.instance.InstanceId)
}

func (e *eureka) GetInstance(name, id string) (*model.Instance, error) {
	url := e.buildUrl(pathApps, name, id)

	result := &model.RegisrationRequest{}
	resp := httpclient.Get(url, result)

	if resp.Error != nil {
		return nil, resp.Error
	}
	return result.Instance, nil
}

func (e *eureka) GetApplication(name string) (*model.Application, error) {
	url := e.buildUrl(pathApps, name)
	log.Info(url)
	result := &model.ApplicationResponse{}
	resp := httpclient.Get(url, result)

	if resp.Error != nil {
		return nil, resp.Error
	}

	return result.Response, nil
}

func (e *eureka) GetApplications() (map[string]*model.Application, error) {
	url := e.buildUrl(pathApps)

	result := &model.ApplicationsResponse{}
	resp := httpclient.Get(url, result)

	if resp.Error != nil {
		return nil, resp.Error
	}

	apps := map[string]*model.Application{}

	for _, a := range result.Response.Applications {
		apps[a.Name] = a
	}

	return apps, nil
}

func (e *eureka) sendHealthCheckUpdate(name, id string) *eurekaResponse {
	url := e.buildUrl(pathApps, name, id)
	resp := httpclient.Put(url, "{}", nil)

	if resp.Error != nil {
		return &eurekaResponse{status: resp.Status, err: resp.Error}
	}

	if resp.Status != http.StatusOK {
		return &eurekaResponse{
			status: resp.Status,
			err:    fmt.Errorf("Heartbeat to Eureka for instance: %s returned status code of: %d", id, resp.Status),
		}
	}
	log.Debugf("successful heartbeat: %d", resp.Status)
	return &eurekaResponse{status: resp.Status}
}
