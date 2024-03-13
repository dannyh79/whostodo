package sessions

import (
	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/sessions/entities"
)

const SessionKey = "token"

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
	_, err := u.repo.FindBy(token.(string))
	if err != nil {
		return false
	}

	return true
}

func InitSessionsUsecase(repo repository.Repository[Session]) *SessionsUsecase {
	return &SessionsUsecase{repo}
}
