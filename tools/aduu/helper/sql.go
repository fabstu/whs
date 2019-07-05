package helper

import (
	"database/sql"

	"github.com/pkg/errors"
)

func CleanupWithError(tx *sql.Tx, err error) error {
	if err2 := tx.Rollback(); err2 != nil {
		return errors.Wrapf(err, "failed to cleanup after: %v", err2)
	}
	return err
}
