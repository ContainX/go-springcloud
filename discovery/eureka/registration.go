package eureka

import (
	"fmt"
	"github.com/ContainX/go-springcloud/discovery/eureka/model"
	"github.com/ContainX/go-utils/httpclient"
	"github.com/cenkalti/backoff"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (e *eureka) Register(await bool) error {
	e.instance = model.NewRegistrationFromInstanceConfig(e.config.Instance)

	// If true - we block after retries and start heartbeat if enabled
	if await {
		if err := e.retryRegistrationFunc(e.reRegister); err != nil {
			return err
		}
		e.heartbeat()
	} else {
		go func() {
			if err := e.retryRegistrationFunc(e.reRegister); err == nil {
				e.heartbeat()
			}
		}()
	}
	return nil
}

func (e *eureka) reRegister() error {
	url := e.buildUrl(pathApps, e.instance.AppName)
	log.Infof("Registering [%s] with instance: %s", url, e.instance.InstanceId)
	resp := httpclient.Post(url, e.instance.WrapInRequest(), nil)

	if resp.Error != nil {
		return fmt.Errorf("Could not complete registration, error: %s, content: %s",
			resp.Error, resp.Content)
	}

	if resp.Status != 204 {
		return fmt.Errorf("HTTP returned %d registering Instance=%s App=%s Body=\"%s\"", resp.Status,
			e.instance.InstanceId, e.instance.AppName, resp.Content)
	}
	return nil
}

func (e *eureka) Unregister() {
	e.stopHeartbeat()
	url := e.buildUrl(pathApps, e.instance.AppName, e.instance.InstanceId)
	resp := httpclient.Delete(url, nil, nil)
	log.Info("Unregistering application: ", resp.Status)
}

func (e *eureka) heartbeat() {
	e.handleSigterm()

	if e.config.Client.RegisterWithEureka {
		go e.startHeartbeat()
	}
}

func (e *eureka) startHeartbeat() {
	log.Info("starting heartbeat....")
	throttle := time.NewTicker(time.Duration(e.config.Client.PollIntervalSeconds) * time.Second)
	stop := false
	for {
		if stop {
			log.Info("shutting down heartbeat service")
			break
		}
		select {
		case <-e.shutdown:
			stop = true
		case <-throttle.C:
			log.Info("sending heartbeat...")
			resp := e.sendHealthCheckUpdate(e.instance.AppName, e.instance.InstanceId)
			if resp.hasError() {
				if resp.status == 404 {
					log.Info("App not found, re-registering...")
					e.reRegister()
				} else {
					log.Error(resp.err)
				}
			}
		}
	}
}

func (e *eureka) stopHeartbeat() {
	e.shutdown <- true
}

func (e *eureka) retryRegistrationFunc(f func() error) error {
	return backoff.RetryNotify(f, NewMaxAttemptBackoff(2*time.Second, 3), e.notifyAttempts)
}

func (e *eureka) notifyAttempts(err error, i time.Duration) {
	log.Error(err.Error())
}

func (e *eureka) handleSigterm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		e.Unregister()
		os.Exit(1)
	}()
}
