package repository_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/sessions/entities"
	util "github.com/dannyh79/whostodo/internal/testutil"
)

type Session = entity.Session

func Test_InMemorySessionRepositorySave(t *testing.T) {
	repo := repository.InitInMemorySessionRepository()
	session := entity.NewSession()
	got := repo.Save(session)

	util.AssertEqual(t)(*session, got)
}

func Test_InMemorySessionRepositoryFindBy(t *testing.T) {
	t.Run("returns a task", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemorySessionRepository()
		session := entity.NewSession()
		repo.Save(session)

		got, err := repo.FindBy(session.Id)

		util.AssertEqual(t)(got, session)
		util.AssertErrorEqual(t)(err, nil)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemorySessionRepository()

		got, err := repo.FindBy("nonexistent_token")

		if got != nil {
			util.AssertEqual(t)(got, nil)
		}
		util.AssertErrorEqual(t)(err, repository.ErrorNotFound)
	})
}
