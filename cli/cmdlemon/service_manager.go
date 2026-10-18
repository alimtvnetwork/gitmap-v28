package cmdlemon

import (
	"fmt"
	"strings"
)

type ServiceActionType string

const (
	ServiceActionTypeInstall   ServiceActionType = "install"
	ServiceActionTypeUninstall ServiceActionType = "uninstall"
	ServiceActionTypeStart     ServiceActionType = "start"
	ServiceActionTypeStop      ServiceActionType = "stop"
	ServiceActionTypeRestart   ServiceActionType = "restart"
	ServiceActionTypeStatus    ServiceActionType = "status"
)

type ServiceStatusType string

const (
	ServiceStatusTypeRunning ServiceStatusType = "running"
	ServiceStatusTypeStopped ServiceStatusType = "stopped"
	ServiceStatusTypeUnknown ServiceStatusType = "unknown"
)

type ServiceDetails struct {
	Name       string            `json:"name"`
	Platform   string            `json:"platform"`
	StatusType ServiceStatusType `json:"statusType"`
	IsEnabled  bool              `json:"isEnabled"`
	ConfigPath string            `json:"configPath,omitempty"`
}

type PlatformService interface {
	Install() error
	Uninstall() error
	Start() error
	Stop() error
	Restart() error
	Status() (*ServiceDetails, error)
	GenerateConfig() (string, error)
}

type ServiceManager struct {
	service PlatformService
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{service: resolvePlatformService()}
}

func (m *ServiceManager) Execute(action ServiceActionType) (*ServiceDetails, error) {
	var err error
	switch strings.ToLower(string(action)) {
	case string(ServiceActionTypeInstall):
		err = m.service.Install()
	case string(ServiceActionTypeUninstall):
		return nil, m.service.Uninstall()
	case string(ServiceActionTypeStart):
		err = m.service.Start()
	case string(ServiceActionTypeStop):
		err = m.service.Stop()
	case string(ServiceActionTypeRestart):
		err = m.service.Restart()
	case string(ServiceActionTypeStatus):
		return m.service.Status()
	default:
		return nil, fmt.Errorf("unknown service action: %s", action)
	}

	if err != nil {
		return nil, err
	}

	return m.service.Status()
}
