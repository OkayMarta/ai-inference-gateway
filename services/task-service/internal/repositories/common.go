package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

// ensureRowsAffected уніфікує перевірку UPDATE/DELETE операцій, де відсутність
// змінених рядків зазвичай означає, що сутність не знайдена.
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
