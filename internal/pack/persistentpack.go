package pack

import (
	"context"
	"errors"
	"strconv"

	packrepo "github.com/aleksandarv/pack-optimizer/internal/repo/pack"
	"github.com/aleksandarv/pack-optimizer/internal/telemetry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type persistentSvc struct {
	db     *pgxpool.Pool
	tracer telemetry.Tracer
}

func NewPersistentSvc(repo *pgxpool.Pool, tracer telemetry.Tracer) PackSvc {
	return &persistentSvc{
		db:     repo,
		tracer: tracer,
	}
}

func (s *persistentSvc) GetSizes(ctx context.Context) (_ []int, err error) {
	rows, err := s.listPackSizes(ctx)
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
	if err = s.deletePackSizes(ctx, qtx); err != nil {
		return err
	}
	if err = s.insertPackSizes(ctx, qtx, sizes); err != nil {
		return err
	}
	return nil
}

func (s *persistentSvc) listPackSizes(ctx context.Context) (_ []int16, err error) {
	ctx, end := s.tracer.Trace(ctx, "pack.db.listPackSizes")
	defer end(&err)

	return packrepo.New(s.db).ListPackSizes(ctx)
}

func (s *persistentSvc) deletePackSizes(ctx context.Context, qtx *packrepo.Queries) (err error) {
	ctx, end := s.tracer.Trace(ctx, "pack.db.deletePackSizes")
	defer end(&err)

	return qtx.DeletePackSizes(ctx)
}

func (s *persistentSvc) insertPackSizes(ctx context.Context, qtx *packrepo.Queries, sizes []int16) (err error) {
	ctx, end := s.tracer.Trace(ctx, "pack.db.insertPackSizes", telemetry.Attribute{
		Key:   "pack.sizes_count",
		Value: strconv.Itoa(len(sizes)),
	})
	defer end(&err)

	return qtx.InsertPackSizes(ctx, sizes)
}
