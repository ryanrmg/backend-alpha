package service

import (
    "context"
    "errors"
    "testing"

    projectx "github.com/ryanrmg/projectx-api"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type MockTradeRepository struct {
    GetTradesByAccountFunc func(ctx context.Context, accountID int) ([]projectx.GatewayUserTrade, error)
}

func (m *MockTradeRepository) GetTradesByAccount(
    ctx context.Context,
    accountID int,
) ([]projectx.GatewayUserTrade, error) {
    return m.GetTradesByAccountFunc(ctx, accountID)
}

func TestTradeService_GetTrades(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        expectedTrades := []projectx.GatewayUserTrade{
            {
                Id:                1,
                AccountId:         1234,
                ContractId:        "CON.F.US.MNQ.U26",
                CreationTimestamp: "2026-07-07T14:30:00Z",
                Price:             23567.25,
                ProfitAndLoss:     125.50,
                Fees:              2.04,
                Side:              1,
                Size:              2,
                Voided:            false,
                OrderId:           1001,
            },
            {
                Id:                2,
                AccountId:         1234,
                ContractId:        "CON.F.US.MES.U26",
                CreationTimestamp: "2026-07-07T14:45:00Z",
                Price:             6321.75,
                ProfitAndLoss:     -42.75,
                Fees:              2.04,
                Side:              2,
                Size:              1,
                Voided:            false,
                OrderId:           1002,
            },
        }

        mockRepo := &MockTradeRepository{
            GetTradesByAccountFunc: func(ctx context.Context, accountID int) ([]projectx.GatewayUserTrade, error) {

                assert.Equal(t, 1234, accountID)

                return expectedTrades, nil
            },
        }

        service := NewTradeService(mockRepo)

        trades, err := service.GetTrades(context.Background(), 1234)

        require.NoError(t, err)
        assert.Equal(t, expectedTrades, trades)
    })

    t.Run("repository error", func(t *testing.T) {
        expectedErr := errors.New("database unavailable")

        mockRepo := &MockTradeRepository{
            GetTradesByAccountFunc: func(ctx context.Context, accountID int) ([]projectx.GatewayUserTrade, error) {
                return nil, expectedErr
            },
        }

        service := NewTradeService(mockRepo)

        trades, err := service.GetTrades(context.Background(), 1234)

        require.Error(t, err)
        assert.Nil(t, trades)
        assert.Equal(t, expectedErr, err)
    })
}