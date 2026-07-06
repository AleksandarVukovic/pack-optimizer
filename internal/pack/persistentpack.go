package pack

import (
	"context"
	"errors"

	packrepo "github.com/aleksandarv/pack-optimizer/internal/repo/pack"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type persistentSvc struct {
	db *pgxpool.Pool
}

func NewPersistentSvc(repo *pgxpool.Pool) PackSvc {
	return &persistentSvc{
		db: repo,
	}
}

func (s *persistentSvc) GetSizes(ctx context.Context) ([]int, error) {
	querier := packrepo.New(s.db)
	rows, err := querier.ListPackSizes(ctx)
	if err != nil {
		return nil, err
	}

	sizes := make([]int, len(rows))
	for i, size := range rows {
		sizes[i] = int(size)
	}
	return sizes, nil
}

func (s *persistentSvc) UpdateSizes(ctx context.Context, newSizes []int) (err error) {
	if err = validate(newSizes); err != nil {
		return err
	}
	sizes := make([]int16, len(newSizes))
	for i, size := range newSizes {
		sizes[i] = int16(size)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
				err = errors.Join(err, rbErr)
			}
			return
		}
		err = tx.Commit(ctx)
	}()

	qtx := packrepo.New(s.db).WithTx(tx)
	if err = qtx.DeletePackSizes(ctx); err != nil {
		return err
	}
	if err = qtx.InsertPackSizes(ctx, sizes); err != nil {
		return err
	}
	return nil
}
