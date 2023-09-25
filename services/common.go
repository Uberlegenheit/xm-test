package services

func (s *ServiceFacade) CheckDBStatus() bool {
	return s.dao.CheckDBStatus()
}
