package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

func ensureRowsAffected(result sql.Result, notFoundMsg string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New(notFoundMsg)
	}

	return nil
}
