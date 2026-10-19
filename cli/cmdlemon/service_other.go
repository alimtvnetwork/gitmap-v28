//go:build !windows && !linux && !darwin

package cmdlemon

import (
	"errors"
)

type OtherService struct{}

func resolvePlatformService() PlatformService {
	return &OtherService{}
}

func (s *OtherService) GenerateConfig() (string, error) {
	return "", errors.New("service config generation not supported on this platform")
}

func (s *OtherService) Install() error {
	return errors.New("service install not supported on this platform")
}

func (s *OtherService) Uninstall() error {
	return errors.New("service uninstall not supported on this platform")
}

func (s *OtherService) Start() error {
	return errors.New("service start not supported on this platform")
}

func (s *OtherService) Stop() error {
	return errors.New("service stop not supported on this platform")
}

func (s *OtherService) Restart() error {
	return errors.New("service restart not supported on this platform")
}

func (s *OtherService) Status() (*ServiceDetails, error) {
	return &ServiceDetails{
		Name:       "gitmap-daemon",
		Platform:   "unsupported",
		StatusType: ServiceStatusTypeUnknown,
		IsEnabled:  false,
	}, nil
}
