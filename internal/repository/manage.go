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

// CreateUserTradeTable creates the cache table if it doesn't already exist
func (store *DBStore) CreateUserTradeTable(ctx context.Context) error {
	query := `CREATE TABLE user_trades (
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
	    journal_notes TEXT                        -- Added an extra column for your journaling!
	);`

	_, err := store.pool.Exec(ctx, query)
	return err
}

// deletes the user table if it exists
func (store *DBStore) DeleteUserTable(ctx context.Context) error {
	query := `DROP TABLE IF EXISTS user_trades;`
	_, err := store.pool.Exec(ctx, query)
	return err
}

// SaveUserTrade inserts a cleanly structured trade into the DB
func (store *DBStore) SaveUserTrade(ctx context.Context, trade projectx.GatewayUserTrade) error {
	query := `
	INSERT INTO user_trades (
		id, account_id, contract_id, creation_timestamp, 
		price, profit_and_loss, fees, side, size, voided, order_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
	)

	return err
}

// DeleteResponsesOlderThan deletes cached records older than a specific duration
func (store *DBStore) DeleteResponsesOlderThan(ctx context.Context, duration time.Duration) (int64, error) {
	query := `DELETE FROM user_trades WHERE fetched_at < $1;`

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
	FROM user_trades 
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
func (store *DBStore) GetTradesByAccount(ctx context.Context, accountId int) ([]projectx.GatewayUserTrade, error) {
	// 1. Write the SQL query
	query := `
		SELECT
			id, account_id, contract_id, creation_timestamp,
			price, profit_and_loss, fees, side, size, voided, order_id
		FROM user_trades
		WHERE account_id = $1
		ORDER BY creation_timestamp DESC;
	`

	// 2. Execute the query
	rows, err := store.pool.Query(ctx, query, accountId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close() // Crucial to prevent connection leaks

	// 3. Create a slice to hold our results
	var trades []projectx.GatewayUserTrade

	// 4. Loop through the result rows
	for rows.Next() {
		var t projectx.GatewayUserTrade
		var dbTime interface{} // To temporarily hold the timestamp object

		// Scan the columns into the struct fields.
		// The order here MUST exactly match the order in your SELECT statement.
		err := rows.Scan(
			&t.Id,
			&t.AccountId,
			&t.ContractId,
			&dbTime,
			&t.Price,
			&t.ProfitAndLoss,
			&t.Fees,
			&t.Side,
			&t.Size,
			&t.Voided,
			&t.OrderId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert the database timestamp into the string format your JSON expects
		if tm, ok := dbTime.(string); ok {
			t.CreationTimestamp = tm
		} else if tm, ok := dbTime.(fmt.Stringer); ok {
			t.CreationTimestamp = tm.String()
		} else {
			// Fallback string conversion depending on how your specific driver environment maps TIMESTAMPTZ
			t.CreationTimestamp = fmt.Sprintf("%v", dbTime)
		}

		// Append the hydrated struct to our slice
		trades = append(trades, t)
	}

	// 5. Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return trades, nil
}

// GetLatestTradeTimestamp returns the newest trade timestamp in the database.
// If there are no trades, it returns time.Time{} and nil.
func (store *DBStore) GetLatestTradeTimestamp(
	ctx context.Context,
) (time.Time, error) {
	query := `
		SELECT MAX(creation_timestamp)
		FROM user_trades;
	`

	var latest *time.Time

	err := store.pool.QueryRow(ctx, query).Scan(&latest)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return time.Time{}, nil
		}
		return time.Time{}, fmt.Errorf("failed to get latest trade timestamp: %w", err)
	}

	if latest == nil {
		// Table exists but contains no rows.
		return time.Time{}, nil
	}

	return *latest, nil
}
