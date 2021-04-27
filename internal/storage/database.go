package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"

	// Posgresql driver.
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

func (s *Storage) AddBanner(ctx context.Context, banner Banner) error {
	insertBannerQuery := `
	INSERT INTO banners (id, description) VALUES (:id, :description);
	`
	_, err := s.db.NamedExecContext(ctx, insertBannerQuery, banner)
	return err
}

func (s *Storage) GetBanner(ctx context.Context, bannerID uuid.UUID) (Banner, error) {
	var banner Banner
	query := `
	SELECT * FROM banners WHERE id=$1
	`
	row := s.db.QueryRowxContext(ctx, query, bannerID)
	if row.Err() != nil {
		return banner, row.Err()
	}

	err := row.StructScan(&banner)
	if err != nil {
		return banner, err
	}

	return banner, nil
}

func (s *Storage) DeleteBanner(ctx context.Context, bannerID uuid.UUID) error {
	deleteBannerQuery := `
	DELETE FROM banners WHERE id=$1
	`
	_, err := s.db.ExecContext(ctx, deleteBannerQuery, bannerID)
	return err
}

func (s *Storage) AddSlot(ctx context.Context, slot Slot) error {
	insertSlotQuery := `
	INSERT INTO slots (id, description) VALUES (:id, :description);
	`
	_, err := s.db.NamedExecContext(ctx, insertSlotQuery, slot)
	return err
}

func (s *Storage) GetSlot(ctx context.Context, slotID uuid.UUID) (Slot, error) {
	var slot Slot
	query := `
	SELECT * FROM slots WHERE id=$1
	`
	row := s.db.QueryRowxContext(ctx, query, slotID)
	if row.Err() != nil {
		return slot, row.Err()
	}

	err := row.StructScan(&slot)
	if err != nil {
		return slot, err
	}

	return slot, nil
}

func (s *Storage) AddGroup(ctx context.Context, group Group) error {
	insertGroupQuery := `
	INSERT INTO groups (id, description) VALUES (:id, :description);
	`
	_, err := s.db.NamedExecContext(ctx, insertGroupQuery, group)
	return err
}

func (s *Storage) GetGroup(ctx context.Context, groupID uuid.UUID) (Group, error) {
	var group Group
	query := `
	SELECT * FROM groups WHERE id=$1
	`
	row := s.db.QueryRowxContext(ctx, query, groupID)
	if row.Err() != nil {
		return group, row.Err()
	}

	err := row.StructScan(&group)
	if err != nil {
		return group, err
	}

	return group, nil
}

func (s *Storage) AddRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	insertRotationQuery := `
	INSERT INTO rotations (banner_id, slot_id, group_id, shows, clicks)
	VALUES (:banner_id, :slot_id, :group_id, :shows, :clicks);
	`
	rotation := Rotation{
		BannerID: bannerID,
		SlotID:   slotID,
		GroupID:  groupID,
	}
	_, err := s.db.NamedExecContext(ctx, insertRotationQuery, rotation)
	return err
}

func (s *Storage) DeleteRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	deleteRotationQuery := `
	DELETE FROM rotations WHERE banner_id=$1 AND slot_id=$2 AND group_id=$3
	`
	_, err := s.db.ExecContext(ctx, deleteRotationQuery, bannerID, slotID, groupID)
	return err
}

func (s *Storage) GetRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) (Rotation, error) {
	var rotation Rotation

	selectRotationQuery := `
	SELECT * FROM rotations WHERE banner_id=$1 AND slot_id=$2 AND group_id=$3
	`
	row := s.db.QueryRowxContext(ctx, selectRotationQuery, bannerID, slotID, groupID)
	if row.Err() != nil {
		return rotation, row.Err()
	}

	err := row.StructScan(&rotation)
	return rotation, err
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
		return errors.New("no rotation was updated")
	}

	return err
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
	query := `SELECT sum(shows) FROM rotations`

	row := s.db.QueryRowxContext(ctx, query)
	if row.Err() != nil {
		return 0, row.Err()
	}

	err = row.Scan(&totalShows)
	return
}
