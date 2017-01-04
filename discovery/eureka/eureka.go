package eureka

import (
	"github.com/ContainX/go-springcloud/discovery/eureka/model"
	"github.com/ContainX/go-utils/logger"
)

type EurekaClient interface {
	// Register the current application with Eureka and setup client side health checks
	// if enabled (default).
	//
	// If await is true then this call will block until either retries have been exhausted or an successful
	// registration.
	Register(await bool) error

	// Unregister the current application from Eureka.  This is auto triggered during normal exiting or sigterm
	// but in some rare cases may be handled manually.
	Unregister()

	// GetInstance fetches the current application instance for the specified application name
	// and id
	GetInstance(name, id string) (*model.Instance, error)

	// GetCurrentInstance fetches the current application instance that has registered within
	// the current lifecycle
	GetCurrentInstance() (*model.Instance, error)

	// GetApplication returns the information about an application by its name.  This also includes
	// information about all the available instances.
	GetApplication(name string) (*model.Application, error)

	// GetApplications fetches all applications and returns a map keyed by the application
	// name and a value of the application containing all available instances
	GetApplications() (map[string]*model.Application, error)
}

type ShutdownChan chan bool

type eureka struct {
	config   *model.EurekaConfig
	instance *model.Instance
	shutdown ShutdownChan
}

type eurekaResponse struct {
	status int
	err    error
}

var log = logger.GetLogger("discovery")

func NewClient(cfg *model.EurekaConfig) EurekaClient {
	return &eureka{config: cfg, shutdown: make(ShutdownChan, 2)}
}

func (e *eurekaResponse) hasError() bool {
	return e.err != nil
}

func isStatus2XX(status int) bool {
	return status > 199 && status < 300
}
