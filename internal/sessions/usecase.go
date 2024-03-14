package sessions

import (
	"time"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/sessions/entities"
)

const SessionKey = "token"

const Timeout = time.Minute

type Session = entity.Session

type SessionsUsecase struct {
	repo repository.Repository[Session]
}

func (u *SessionsUsecase) Authenticate() string {
	s := entity.NewSession()
	u.repo.Save(s)
	return s.Id
}

func (u *SessionsUsecase) Validate(token any) bool {
	if token == nil {
		return false
	}
	session, err := u.repo.FindBy(token.(string))
	if err != nil {
		return false
	}
	if isExpiredSession(session) {
		return false
	}

	return true
}

func InitSessionsUsecase(repo repository.Repository[Session]) *SessionsUsecase {
	return &SessionsUsecase{repo}
}

func isExpiredSession(s *Session) bool {
	timeNow := time.Now()
	expiredAt := s.CreatedAt.Add(Timeout)
	return timeNow.After(expiredAt)
}
