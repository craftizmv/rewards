package postgres

import (
	"database/sql"
	"fmt"
	"github.com/craftizmv/rewards/internal/app/repository"
	. "github.com/craftizmv/rewards/internal/domain/entities"
	"strings"
)

// PostgresRewardRepository is the concrete implementation of the RewardRepository interface for Postgres
type PostgresRewardRepository struct {
	db *sql.DB
}

// NewPostgresRewardRepository creates a new instance of PostgresRewardRepository
func NewPostgresRewardRepository(db *sql.DB) repository.RewardRepository {
	return &PostgresRewardRepository{
		db: db,
	}
}

// GetRewardByID retrieves a reward group by its ID
func (r *PostgresRewardRepository) GetRewardGroupByID(id int64) (*RewardGroup, error) {
	// Prepare the SQL query
	query := `SELECT id, name, expires_at, campaign_id 
			  FROM reward_groups WHERE id = $1`

	// Execute the query
	row := r.db.QueryRow(query, id)

	// Map the result to a RewardGroup entity
	var rewardGroup RewardGroup
	err := row.Scan(&rewardGroup.ID, &rewardGroup.Name, &rewardGroup.ExpiresAt, &rewardGroup.CampaignID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reward group with id %d not found", id)
		}
		return nil, err
	}

	// Return the result
	return &rewardGroup, nil
}

// GetRewardItemIDsFromRewardGroup retrieves the list of reward item IDs associated with a reward group
func (r *PostgresRewardRepository) GetRewardItemIDsFromRewardGroup(rewardGroupID int64) ([]int64, error) {
	// Prepare the SQL query to retrieve RewardItemIDs directly from the associative table
	query := `SELECT reward_item_id 
			  FROM reward_group_reward_items 
			  WHERE reward_group_id = $1`

	// Execute the query
	rows, err := r.db.Query(query, rewardGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect reward item IDs
	var rewardItemIDs []int64
	for rows.Next() {
		var rewardItemID int64

		// Scan the result into the rewardItemID variable
		if err := rows.Scan(&rewardItemID); err != nil {
			return nil, err
		}

		// Append rewardItemID to the slice
		rewardItemIDs = append(rewardItemIDs, rewardItemID)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the list of reward item IDs
	return rewardItemIDs, nil
}

// GetRewardGroupIDByOrderID retrieves the RewardGroupID associated with a given OrderID from the OrderRewardGroup table
func (r *PostgresRewardRepository) GetRewardGroupIDByOrderID(orderID int64) ([]int64, error) {
	// Prepare the SQL query
	query := `
		SELECT reward_group_id
		FROM order_reward_group
		WHERE order_id = $1
	`

	// Execute the query
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RewardGroupID for OrderID %d: %v", orderID, err)
	}
	defer rows.Close()

	// Collect the RewardGroupID results
	var rewardGroupIDs []int64
	for rows.Next() {
		var rewardGroupID int64
		if err := rows.Scan(&rewardGroupID); err != nil {
			return nil, fmt.Errorf("failed to scan RewardGroupID: %v", err)
		}
		rewardGroupIDs = append(rewardGroupIDs, rewardGroupID)
	}

	// Check for errors in row iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Return the list of RewardGroupIDs
	return rewardGroupIDs, nil
}

// GetProductIDsFromRewardGroup retrieves the list of product IDs associated with a reward group
func (r *PostgresRewardRepository) GetProductIDsFromRewardGroup(rewardGroupID int64) ([]int64, error) {
	// Prepare the SQL query to retrieve ProductIDs from RewardGroupRewardProduct table
	query := `SELECT product_id 
			  FROM reward_group_reward_products 
			  WHERE reward_group_id = $1`

	// Execute the query
	rows, err := r.db.Query(query, rewardGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect product IDs
	var productIDs []int64
	for rows.Next() {
		var productID int64

		// Scan the result into the productID variable
		if err := rows.Scan(&productID); err != nil {
			return nil, err
		}

		// Append productID to the slice
		productIDs = append(productIDs, productID)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the list of product IDs
	return productIDs, nil
}

// InsertRewardGroupRewardItem inserts a new mapping between a reward group and a reward item
func (r *PostgresRewardRepository) InsertRewardGroupRewardItem(rewardGroupID, rewardItemID int64) error {
	// Prepare the SQL query for insertion
	query := `INSERT INTO reward_group_reward_items (reward_group_id, reward_item_id) 
			  VALUES ($1, $2)`

	// Execute the insert query
	_, err := r.db.Exec(query, rewardGroupID, rewardItemID)
	if err != nil {
		return fmt.Errorf("failed to insert reward group and reward item mapping: %v", err)
	}

	// Return nil if insert was successful
	return nil
}

// InsertRewardGroupRewardItemsBatch inserts RewardItemIDs in batches for a given RewardGroupID
func (r *PostgresRewardRepository) InsertRewardGroupRewardItemsBatch(rewardGroupID int64, rewardItemIDs []int64, batchSize int) error {
	// Validate input: Ensure there are items to insert
	if len(rewardItemIDs) == 0 {
		return fmt.Errorf("no reward items to insert for reward group ID %d", rewardGroupID)
	}

	// Start a new transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Helper function to perform batch insert
	insertBatch := func(itemIDs []int64) error {
		// Prepare the SQL query template
		query := `INSERT INTO reward_group_reward_items (reward_group_id, reward_item_id) VALUES `

		// Create placeholders for each item in the batch insert
		values := make([]string, len(itemIDs))
		args := make([]interface{}, 0, len(itemIDs)*2)

		for i, rewardItemID := range itemIDs {
			// Use $2, $3, $4... for reward item placeholders
			values[i] = fmt.Sprintf("($1, $%d)", i+2)
			args = append(args, rewardItemID)
		}

		// Add the rewardGroupID as the first argument for the bulk insert
		args = append([]interface{}{rewardGroupID}, args...)

		// Build the final query with placeholders
		query += strings.Join(values, ", ")

		// Execute the query within the transaction
		_, err := tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to insert batch: %v", err)
		}

		return nil
	}

	// Process the rewardItemIDs in batches
	for start := 0; start < len(rewardItemIDs); start += batchSize {
		end := start + batchSize
		if end > len(rewardItemIDs) {
			end = len(rewardItemIDs)
		}

		// Insert the current batch
		err := insertBatch(rewardItemIDs[start:end])
		if err != nil {
			tx.Rollback() // Rollback the transaction on error
			return err
		}
	}

	// Commit the transaction after all batches are inserted
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Return nil if the insert was successful
	return nil
}

// UpdateOrderRewardItemsBatch performs a bulk update of OrderRewardItem records in batches
func (r *PostgresRewardRepository) UpdateOrderRewardItemsBatch(orderRewardItems []*OrderRewardItem, batchSize int) error {
	// Validate input: Ensure there are items to update
	if len(orderRewardItems) == 0 {
		return fmt.Errorf("no order reward items to update")
	}

	// Start a new transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Helper function to perform batch update
	updateBatch := func(items []*OrderRewardItem) error {
		// Prepare the base SQL query
		queryParts := make([]string, 0, len(items))
		args := make([]interface{}, 0, len(items)*6) // Each update has 6 parameters

		// Build the query for each OrderRewardItem
		for i, item := range items {
			queryParts = append(queryParts, fmt.Sprintf(`
				UPDATE order_reward_item 
				SET 
					shipment_id = $%d, 
					allocated_date = $%d, 
					is_redeemed = $%d, 
					redeemed_date = $%d 
				WHERE 
					order_id = $%d 
					AND reward_item_id = $%d
			`, i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))

			// Add values to the args slice (account for optional fields)
			args = append(args,
				sql.NullInt64{Int64: *item.ShipmentID, Valid: item.ShipmentID != nil},
				item.AllocatedDate,
				sql.NullBool{Bool: *item.IsRedeemed, Valid: item.IsRedeemed != nil},
				sql.NullTime{Time: *item.RedeemedDate, Valid: item.RedeemedDate != nil},
				item.OrderID,
				item.RewardItemID,
			)
		}

		// Combine the query parts into a single SQL statement
		query := strings.Join(queryParts, "; ")

		// Execute the batch update within the transaction
		_, err := tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to update batch: %v", err)
		}

		return nil
	}

	// Process the OrderRewardItems in batches
	for start := 0; start < len(orderRewardItems); start += batchSize {
		end := start + batchSize
		if end > len(orderRewardItems) {
			end = len(orderRewardItems)
		}

		// Update the current batch
		err := updateBatch(orderRewardItems[start:end])
		if err != nil {
			tx.Rollback() // Rollback the transaction on error
			return err
		}
	}

	// Commit the transaction after all batches are updated
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Return nil if the update was successful
	return nil
}

// DeleteRewardGroupByOrderID deletes the association between an OrderID and RewardGroupID from the order_reward_group table
func (r *PostgresRewardRepository) DeleteRewardGroupByOrderID(orderID int64, rewardGroupID int64) error {
	// Prepare the SQL delete query
	query := `
		DELETE FROM order_reward_group
		WHERE order_id = $1 AND reward_group_id = $2
	`

	// Execute the delete query
	result, err := r.db.Exec(query, orderID, rewardGroupID)
	if err != nil {
		return fmt.Errorf("failed to delete RewardGroupID %d for OrderID %d: %v", rewardGroupID, orderID, err)
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted for OrderID %d and RewardGroupID %d", orderID, rewardGroupID)
	}

	return nil
}

// DeleteRewardItemsByOrderID deletes all associations between an OrderID and RewardItems from the order_reward_item table
func (r *PostgresRewardRepository) DeleteRewardItemsByOrderID(orderID int64) error {
	// Prepare the SQL delete query
	query := `
		DELETE FROM order_reward_item
		WHERE order_id = $1
	`

	// Execute the delete query
	result, err := r.db.Exec(query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete RewardItems for OrderID %d: %v", orderID, err)
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted for OrderID %d", orderID)
	}

	return nil
}

// DeleteRewardItemsByRewardGroupID deletes all associations between a RewardGroupID and RewardItems from the reward_group_reward_item table
func (r *PostgresRewardRepository) DeleteRewardItemsByRewardGroupID(rewardGroupID int64) error {
	// Prepare the SQL delete query
	query := `
		DELETE FROM reward_group_reward_item
		WHERE reward_group_id = $1
	`

	// Execute the delete query
	result, err := r.db.Exec(query, rewardGroupID)
	if err != nil {
		return fmt.Errorf("failed to delete RewardItems for RewardGroupID %d: %v", rewardGroupID, err)
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted for RewardGroupID %d", rewardGroupID)
	}

	return nil
}
