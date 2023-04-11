package service

import "net/http"

func (s *Service) getVersion(w http.ResponseWriter, _ *http.Request) error {
	version, err := s.db.Version.GetVersion()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]string{"version": version})
}
