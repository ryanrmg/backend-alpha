package service

import (
	"context"

	projectx "github.com/ryanrmg/projectx-api"
)

type AccountService struct {
	client *projectx.ProjectXClient
}

func NewAccountService(
	client *projectx.ProjectXClient,
) *AccountService {

	return &AccountService{
		client: client,
	}
}

func (s *AccountService) GetAccounts(
	ctx context.Context,
) ([]projectx.Account, error) {

	return s.client.Accounts.Search(ctx, projectx.AccountSearchRequest{OnlyActiveAccounts: true})
}
