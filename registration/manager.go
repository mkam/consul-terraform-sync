package registration

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul-terraform-sync/client"
	"github.com/hashicorp/consul-terraform-sync/config"
	"github.com/hashicorp/consul-terraform-sync/logging"
	consulapi "github.com/hashicorp/consul/api"
)

const (
	// Service defaults
	defaultServiceName = "Consul-Terraform-Sync"

	// Check defaults
	defaultCheckName                      = "CTS Health Status"
	defaultCheckNotes                     = "Check created by Consul-Terraform-Sync"
	defaultDeregisterCriticalServiceAfter = "30m"
	defaultCheckStatus                    = consulapi.HealthCritical

	logSystemName = "registration"
)

var defaultServiceTags = []string{"cts"}

type RegistrationManager struct {
	client  client.ConsulClientInterface
	service *Service

	logger logging.Logger
}

type Service struct {
	Name      string
	ID        string
	Tags      []string
	Port      int
	Namespace *string

	Checks []*consulapi.AgentServiceCheck
}

func NewRegistrationManager(conf config.Config, client client.ConsulClientInterface) *RegistrationManager {
	logger := logging.Global().Named(logSystemName)

	name := defaultServiceName

	var ns *string
	if conf.Consul != nil && conf.Consul.SelfRegistration != nil {
		ns = conf.Consul.SelfRegistration.Namespace
	}

	var checks []*consulapi.AgentServiceCheck
	checks = append(checks, defaultHTTPCheck(conf))

	return &RegistrationManager{
		client: client,
		logger: logger,
		service: &Service{
			Name:      name,
			ID:        *conf.ID,
			Tags:      defaultServiceTags,
			Port:      *conf.Port,
			Namespace: ns,
			Checks:    checks,
		},
	}
}

func (m *RegistrationManager) ServiceRegister(ctx context.Context) error {
	s := m.service
	r := &consulapi.AgentServiceRegistration{
		ID:     s.ID,
		Name:   s.Name,
		Tags:   s.Tags,
		Port:   s.Port,
		Checks: s.Checks,
	}

	if s.Namespace != nil {
		r.Namespace = *s.Namespace
	}

	m.logger.Debug("registering service with Consul", "name", s.Name, "id", s.ID)
	err := m.client.RegisterService(ctx, r)
	if err != nil {
		m.logger.Error("error registering service with Consul", "name", s.Name, "id", s.ID)
		return err
	}
	return nil
}

func defaultHTTPCheck(conf config.Config) *consulapi.AgentServiceCheck {
	id := *conf.ID
	port := *conf.Port
	var protocol string
	if conf.TLS != nil && *conf.TLS.Enabled {
		protocol = "https"
	} else {
		protocol = "http"
	}
	// address := fmt.Sprintf("%s://localhost:%d/v1/health", protocol, port)
	address := fmt.Sprintf("%s://localhost:%d/v1/status", protocol, port) // TODO: temporary until /health implemented
	return &consulapi.AgentServiceCheck{
		Name:                           defaultCheckName,
		CheckID:                        fmt.Sprintf("%s-health", id),
		Notes:                          defaultCheckNotes,
		DeregisterCriticalServiceAfter: defaultDeregisterCriticalServiceAfter,
		Status:                         defaultCheckStatus,
		HTTP:                           address,
		Method:                         "GET",
		Interval:                       "10s",
		Timeout:                        "2s",
		TLSSkipVerify:                  true,
	}
}
