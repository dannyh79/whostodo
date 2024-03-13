package sessions

import "github.com/dannyh79/whostodo/internal/repository"

const SessionKey = "token"

type SessionsUsecase struct {
	repo repository.Repository[repository.Session]
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

func InitSessionsUsecase(repo repository.Repository[repository.Session]) *SessionsUsecase {
	return &SessionsUsecase{repo}
}
