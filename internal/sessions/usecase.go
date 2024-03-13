package sessions

const SessionKey = "token"

type SessionsUsecase struct{}

func (u *SessionsUsecase) Validate(token any) bool {
	if token == nil || token != "someToken" {
		return false
	}

	return true
}

func InitSessionsUsecase() *SessionsUsecase {
	return &SessionsUsecase{}
}
