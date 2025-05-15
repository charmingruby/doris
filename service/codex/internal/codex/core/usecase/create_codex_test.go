package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

func (s *Suite) Test_CreateCodex() {
	s.Run("it should be able to create a codex", func() {
		correlationID := "123"
		name := "test"
		description := "test"

		id, err := s.uc.CreateCodex(context.Background(), &CreateCodexInput{
			CorrelationID: correlationID,
			Name:          name,
			Description:   description,
		})
		s.NoError(err)

		storedCodex := s.codexRepo.Items[0]

		s.Equal(storedCodex.ID, id)
		s.Equal(storedCodex.CorrelationID, correlationID)
		s.Equal(storedCodex.Name, name)
		s.Equal(storedCodex.Description, description)
	})

	s.Run("it should be not able to create a quota if datasource fails", func() {
		correlationID := "123"
		name := "test"
		description := "test"

		s.codexRepo.IsHealthy = false

		id, err := s.uc.CreateCodex(context.Background(), &CreateCodexInput{
			CorrelationID: correlationID,
			Name:          name,
			Description:   description,
		})
		s.Empty(id)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be not able to create a codex if already exists", func() {
		ctx := context.Background()

		correlationID := "123"
		name := "test"
		description := "test"

		codex := model.NewCodex(model.CodexInput{
			CorrelationID: correlationID,
			Name:          name,
			Description:   description,
		})

		err := s.codexRepo.Create(ctx, *codex)
		s.NoError(err)

		id, err := s.uc.CreateCodex(ctx, &CreateCodexInput{
			CorrelationID: correlationID,
			Name:          name,
			Description:   description,
		})
		s.Empty(id)
		s.Error(err)

		var resourceAlreadyExistsErr *custom_err.ErrResourceAlreadyExists
		s.True(errors.As(err, &resourceAlreadyExistsErr), "error should be of type ErrResourceAlreadyExists")
	})
}
