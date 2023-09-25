package services

import (
	"fmt"
	"xm-task/smodels"
)

func (s *ServiceFacade) FetchAuth(authD smodels.AccessDetails) (string, error) {
	if authD.AccessUuid == "" {
		return "", fmt.Errorf("empty token details")
	}

	walletAddr, ok, err := s.dao.GetAuthToken(authD.AccessUuid)
	if err != nil || !ok {
		return "", err
	}

	return walletAddr.(string), nil
}

func (s *ServiceFacade) DeleteAuth(UUID ...string) error {
	for i := range UUID {
		err := s.dao.RemoveAuthToken(UUID[i])
		if err != nil {
			return err
		}
	}

	return nil
}
