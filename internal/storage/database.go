package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db      *sqlx.DB
	connStr string
}

func New(connStr string) *Storage {
	return &Storage{connStr: connStr}
}

type Banner struct {
	ID          uuid.UUID
	Description string
}

type Slot struct {
	ID          uuid.UUID
	Description string
}

type Group struct {
	ID          uuid.UUID
	Description string
}

type Rotation struct {
	BannerID uuid.UUID `db:"banner_id"`
	SlotID   uuid.UUID `db:"slot_id"`
	GroupID  uuid.UUID `db:"group_id"`
	Shows    int
	Clicks   int
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(ctx, "pgx", s.connStr)
	return
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateBanner(ctx context.Context, description string) error {
	banner := Banner{ID: uuid.New(), Description: description}
	insertBannerQuery := `
	INSERT INTO banners (id, description) VALUES (:id, :description);
	`
	_, err := s.db.NamedExecContext(ctx, insertBannerQuery, banner)
	return err
}

func (s *Storage) CreateSlot(ctx context.Context, description string) error {
	slot := Slot{ID: uuid.New(), Description: description}
	insertSlotQuery := `
	INSERT INTO slots (id, description) VALUES (:id, :description);
	`
	_, err := s.db.NamedExecContext(ctx, insertSlotQuery, slot)
	return err
}

func (s *Storage) CreateGroup(ctx context.Context, description string) error {
	slot := Group{ID: uuid.New(), Description: description}
	insertGroupQuery := `
	INSERT INTO groups (id, description) VALUES (:id, :description);
	`
	_, err := s.db.NamedExecContext(ctx, insertGroupQuery, slot)
	return err
}

func (s *Storage) CreateRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	rotation := Rotation{
		BannerID: bannerID,
		SlotID:   slotID,
		GroupID:  groupID,
	}
	insertRotationQuery := `
	INSERT INTO rotations (banner_id, slot_id, group_id) VALUES (:banner_id, :slot_id, :group_id);
	`
	_, err := s.db.NamedExecContext(ctx, insertRotationQuery, rotation)
	return err
}

func (s *Storage) AddShow(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	query := `
	UPDATE rotations SET shows=shows+1
	WHERE
	banner_id=$1 AND slot_id=$2 AND group_id=$3
	`
	res, err := s.db.ExecContext(ctx, query, bannerID, slotID, groupID)
	if err != nil {
		return err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsUpdated == 0 {
		return errors.New("no row was updated")
	}

	// Now increment global shows counter
	query = `
	UPDATE total_shows SET count=count+1
	`
	res, err = s.db.ExecContext(ctx, query, bannerID, slotID, groupID)
	if err != nil {
		return err
	}

	rowsUpdated, err = res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsUpdated == 0 {
		return errors.New("no row was updated")
	}

	return nil
}

func (s *Storage) AddClick(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	query := `
	UPDATE rotations SET clicks=clicks+1
	WHERE
	banner_id=$1 AND slot_id=$2 AND group_id=$3
	`
	res, err := s.db.ExecContext(ctx, query, bannerID, slotID, groupID)
	if err != nil {
		return err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsUpdated == 0 {
		return errors.New("no row was updated")
	}

	return nil
}

func (s *Storage) GetTotalShows(ctx context.Context) (totalShows int64, err error) {
	query := `SELECT count FROM total_shows`

	row := s.db.QueryRowxContext(ctx, query)
	if row.Err() != nil {
		return 0, row.Err()
	}

	err = row.Scan(totalShows)
	return
}
