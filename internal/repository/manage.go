package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	projectx "github.com/ryanrmg/projectx-api"
)

// DBStore handles all database interactions
type DBStore struct {
	pool *pgxpool.Pool
}

// NewDBStore initializes our repository wrapper
func NewDBStore(pool *pgxpool.Pool) *DBStore {
	return &DBStore{pool: pool}
}

type UserTradesJournalEntry struct {
	TradeId    int    `json:"tradeId"`
	AccountId  int    `json:"accountId"`
	ContractId string `json:"contractId"`

	EntryTimestamp string `json:"entryTimestamp"`
	ExitTimestamp  string `json:"exitTimestamp"`

	EntryPrice float64 `json:"entryPrice"`
	ExitPrice  float64 `json:"exitPrice"`

	EntrySize int `json:"entrySize"`
	ExitSize  int `json:"exitSize"`

	ProfitAndLoss float64 `json:"profitAndLoss"`
	Fees          float64 `json:"fees"`

	JournalNotes string `json:"journalNotes"`
}

// CreateUserFillsTable creates the cache table if it doesn't already exist
func (store *DBStore) CreateUserFillsTable(ctx context.Context) error {
	query := `CREATE TABLE user_fills (
	    id INT PRIMARY KEY,
	    account_id INT NOT NULL,
	    contract_id VARCHAR(50) NOT NULL,
	    creation_timestamp TIMESTAMPTZ NOT NULL,  -- Best to parse string timestamps into real times
	    price NUMERIC(18, 8) NOT NULL,            -- Use NUMERIC/DECIMAL for financial accuracy
	    profit_and_loss NUMERIC(18, 8) NOT NULL,
	    fees NUMERIC(18, 8) NOT NULL,
	    side INT NOT NULL,                        -- e.g., 1 for Buy, 2 for Sell
	    size INT NOT NULL,
	    voided BOOLEAN NOT NULL DEFAULT FALSE,
	    order_id INT NOT NULL,
	    trade_id BIGINT
	);`

	_, err := store.pool.Exec(ctx, query)
	return err
}

// deletes the user table if it exists
func (store *DBStore) DeleteUserFillsTable(ctx context.Context) error {
	query := `DROP TABLE IF EXISTS user_fills;`
	_, err := store.pool.Exec(ctx, query)
	return err
}

// SaveUserTrade inserts a cleanly structured trade into the DB
func (store *DBStore) SaveUserFill(ctx context.Context, trade projectx.GatewayUserTrade, tradeId *int64) error {
	query := `
	INSERT INTO user_fills (
		id, account_id, contract_id, creation_timestamp, 
		price, profit_and_loss, fees, side, size, voided, order_id, trade_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	ON CONFLICT (id) DO NOTHING;` // Prevents crashes if the API sends duplicate logs

	// Parse your string timestamp into a real Go time.Time object for Postgres
	parsedTime, err := time.Parse(time.RFC3339, trade.CreationTimestamp)
	if err != nil {
		// Fallback to current time if the formatting fails
		parsedTime = time.Now()
	}

	_, err = store.pool.Exec(ctx, query,
		trade.Id,
		trade.AccountId,
		trade.ContractId,
		parsedTime,
		trade.Price,
		trade.ProfitAndLoss,
		trade.Fees,
		trade.Side,
		trade.Size,
		trade.Voided,
		trade.OrderId,
		tradeId,
	)

	return err
}

// DeleteResponsesOlderThan deletes cached records older than a specific duration
func (store *DBStore) DeleteResponsesOlderThan(ctx context.Context, duration time.Duration) (int64, error) {
	query := `DELETE FROM user_fills WHERE fetched_at < $1;`

	cutoff := time.Now().Add(-duration)
	result, err := store.pool.Exec(ctx, query, cutoff)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

// GetLatestResponse retrieves the newest cached JSON string for a given endpoint
func (store *DBStore) GetLatestResponse(ctx context.Context, endpoint string) (string, time.Time, error) {
	query := `
	SELECT response_json, fetched_at 
	FROM user_fills 
	WHERE endpoint = $1 
	ORDER BY fetched_at DESC 
	LIMIT 1;`

	var responseJSON string
	var fetchedAt time.Time

	err := store.pool.QueryRow(ctx, query, endpoint).Scan(&responseJSON, &fetchedAt)
	if err != nil {
		return "", time.Time{}, err // Will return pgx.ErrNoRows if empty
	}

	return responseJSON, fetchedAt, nil
}

// GetTradesByAccount retrieves all stored trades for a specific account ID ordered by newest first
func (store *DBStore) GetTradesByAccount(
	ctx context.Context,
	accountId int,
) ([]UserTradesJournalEntry, error) {

	query := `
		WITH ordered_fills AS (
			SELECT
				*,
				ROW_NUMBER() OVER (
					PARTITION BY trade_id
					ORDER BY creation_timestamp ASC, id ASC
				) AS entry_rank,
				ROW_NUMBER() OVER (
					PARTITION BY trade_id
					ORDER BY creation_timestamp DESC, id DESC
				) AS exit_rank
			FROM user_fills
			WHERE account_id = $1
			  AND trade_id IS NOT NULL
		),
		trade_summary AS (
			SELECT
				trade_id,
				account_id,
				contract_id,

				to_char(MIN(creation_timestamp), 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS entry_timestamp,
				to_char(MAX(creation_timestamp), 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS exit_timestamp,

				SUM(size) FILTER (WHERE entry_rank = 1) AS entry_size,
				SUM(size) FILTER (WHERE exit_rank = 1) AS exit_size,

				SUM(profit_and_loss) AS profit_and_loss,
				SUM(fees) AS fees,

				MAX(price) FILTER (WHERE entry_rank = 1) AS entry_price,
				MAX(price) FILTER (WHERE exit_rank = 1) AS exit_price

			FROM ordered_fills
			GROUP BY
				trade_id,
				account_id,
				contract_id
		)

		SELECT
			trade_id,
			account_id,
			contract_id,
			entry_timestamp,
			exit_timestamp,
			entry_price,
			exit_price,
			entry_size,
			exit_size,
			profit_and_loss,
			fees
		FROM trade_summary
		ORDER BY exit_timestamp DESC;
	`

	rows, err := store.pool.Query(ctx, query, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []UserTradesJournalEntry

	for rows.Next() {
		var t UserTradesJournalEntry

		err := rows.Scan(
			&t.TradeId,
			&t.AccountId,
			&t.ContractId,
			&t.EntryTimestamp,
			&t.ExitTimestamp,
			&t.EntryPrice,
			&t.ExitPrice,
			&t.EntrySize,
			&t.ExitSize,
			&t.ProfitAndLoss,
			&t.Fees,
		)
		if err != nil {
			return nil, err
		}

		trades = append(trades, t)
	}

	return trades, rows.Err()
}

// GetLatestTradeTimestamp returns the newest trade timestamp in the database.
// If there are no trades, it returns time.Time{} and nil.
func (store *DBStore) GetLatestFill(
	ctx context.Context,
) (time.Time, *int64, error) {
	query := `
		SELECT creation_timestamp, trade_id
		FROM user_fills
		ORDER BY creation_timestamp DESC
		LIMIT 1;
	`

	var (
		timestamp *time.Time
		tradeID   *int64
	)

	err := store.pool.QueryRow(ctx, query).Scan(&timestamp, &tradeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return time.Time{}, nil, nil
		}
		return time.Time{}, nil, fmt.Errorf("failed to get latest fill: %w", err)
	}

	if timestamp == nil {
		return time.Time{}, nil, nil
	}

	return *timestamp, tradeID, nil
}
