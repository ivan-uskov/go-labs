package repository

import (
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

func (o *orderRepository) withTx(fn func(*sql.Tx, context.Context, func(error) error) error) error {
	ctx := context.Background()
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	closeTx := func(err error) error {
		if err == nil {
			return tx.Commit()
		}

		log.Error(tx.Rollback())
		return err
	}

	return fn(tx, ctx, closeTx)
}
