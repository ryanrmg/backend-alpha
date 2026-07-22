package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	projectx "github.com/ryanrmg/projectx-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCreateTable(t *testing.T) {
	connStr := os.Getenv("TEST_DATABASE_URL")
	require.NotEmpty(t, connStr, "TEST_DATABASE_URL must be set")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	err = pool.Ping(ctx)
	require.NoError(t, err)

	repo := NewPostgresTradeRepository(pool)

	err = repo.DeleteUserFillsTable(ctx)
	require.NoError(t, err)

	err = repo.CreateUserFillsTable(ctx)
	require.NoError(t, err)

	err = repo.DeleteUserFillsTable(ctx)
	require.NoError(t, err)
}

func TestPostgresTradeRepository_GetTradesByAccount(t *testing.T) {
	connStr := os.Getenv("TEST_DATABASE_URL")
	require.NotEmpty(t, connStr, "TEST_DATABASE_URL must be set")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	err = pool.Ping(ctx)
	require.NoError(t, err)

	repo := NewPostgresTradeRepository(pool)

	err = repo.CreateUserFillsTable(ctx)
	require.NoError(t, err)

	testTrade := projectx.GatewayUserTrade{
		Id:                1,
		AccountId:         1234,
		ContractId:        "CON.F.US.MES.U26",
		CreationTimestamp: "2026-07-07T14:25:00Z",
		Price:             6321.75,
		ProfitAndLoss:     0,
		Fees:              2.04,
		Side:              0,
		Size:              1,
		Voided:            false,
		OrderId:           1002,
	}
	var tradeId int64 = 1
	err = repo.SaveUserFill(ctx, testTrade, &tradeId)
	require.NoError(t, err)

	testTrade2 := projectx.GatewayUserTrade{
		Id:                2,
		AccountId:         1234,
		ContractId:        "CON.F.US.MES.U26",
		CreationTimestamp: "2026-07-07T14:45:00Z",
		Price:             6341.75,
		ProfitAndLoss:     20.00,
		Fees:              2.04,
		Side:              1,
		Size:              1,
		Voided:            false,
		OrderId:           1002,
	}
	err = repo.SaveUserFill(ctx, testTrade2, &tradeId)
	require.NoError(t, err)

	// Choose an account that exists in your test database.
	accountID := 1234

	trades, err := repo.GetTradesByAccount(ctx, accountID)

	require.NoError(t, err)

	// We expect at least one trade for this account.
	assert.NotEmpty(t, trades)

	// Verify every returned trade belongs to the account.
	for _, trade := range trades {
		assert.Equal(t, accountID, trade.AccountId)
		assert.Equal(t, testTrade.Price, trade.EntryPrice)
		assert.Equal(t, testTrade2.Price, trade.ExitPrice)
		assert.Equal(t, testTrade2.ProfitAndLoss, trade.ProfitAndLoss)
		assert.NotZero(t, trade.TradeId)
		assert.NotEmpty(t, trade.ContractId)
	}

	err = repo.DeleteUserFillsTable(ctx)
	require.NoError(t, err)
}
