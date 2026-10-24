package cmdlemon

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon/enums"

// ServiceActionType aliases the strongly-typed enum from the enums package.
type ServiceActionType = enums.ServiceActionType

const (
	ServiceActionTypeInstall   = enums.ServiceActionTypeInstall
	ServiceActionTypeUninstall = enums.ServiceActionTypeUninstall
	ServiceActionTypeStart     = enums.ServiceActionTypeStart
	ServiceActionTypeStop      = enums.ServiceActionTypeStop
	ServiceActionTypeRestart   = enums.ServiceActionTypeRestart
	ServiceActionTypeStatus    = enums.ServiceActionTypeStatus
)

// ServiceStatusType aliases the strongly-typed enum from the enums package.
type ServiceStatusType = enums.ServiceStatusType

const (
	ServiceStatusTypeRunning = enums.ServiceStatusTypeRunning
	ServiceStatusTypeStopped = enums.ServiceStatusTypeStopped
	ServiceStatusTypeUnknown = enums.ServiceStatusTypeUnknown
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

func NewServiceManagerWithService(service PlatformService) *ServiceManager {
	return &ServiceManager{service: service}
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
		return nil, apperror.NewValidationError(fmt.Sprintf("unknown service action: %s", action))
	}

	if err != nil {
		return nil, err
	}

	return m.service.Status()
}
