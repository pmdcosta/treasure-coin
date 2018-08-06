package database

const SessionCollection = "sessions"

// SessionService represents a service for managing session persistence.
type SessionService struct {
	client *Client
}

// Add adds the record to the database if it does not exist.
func (s *SessionService) Add(token, session string) error {
	return s.client.Create(SessionCollection, token, []byte(session))
}

// Find retrieves a session from the database.
func (s *SessionService) Find(token string) (string, error) {
	ses, err := s.client.Load(SessionCollection, token)
	if err != nil {
		return "", err
	}
	return string(ses), nil
}

// Remove deletes the session from the databse.
func (s *SessionService) Remove(token string) error {
	return s.client.Delete(SessionCollection, token)
}
