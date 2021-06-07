package storage

import (
	"context"
	"database/sql"

	"github.com/FedoseevAlex/banner-rotation/internal/types"
	"github.com/google/uuid"

	// Posgresql driver.
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	EventTypeClick string = "click"
	EventTypeShow  string = "show"
)

type Storage struct {
	db      *sqlx.DB
	connStr string
}

func New(connStr string) *Storage {
	return &Storage{connStr: connStr}
}

var ErrNoRowWasAffected = errors.New("no row was affected")

func execTxQuery(tx *sql.Tx, query string, args ...interface{}) error {
	res, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRowWasAffected
	}

	return nil
}

// This method is to be used for tests only.
func (s *Storage) CleanDB() error {
	cleanEventTimestamps := `DELETE FROM events`
	_, err := s.db.Exec(cleanEventTimestamps)
	if err != nil {
		return err
	}

	cleanRotations := `DELETE FROM rotations`
	_, err = s.db.Exec(cleanRotations)
	if err != nil {
		return err
	}

	cleanBanners := `DELETE FROM banners`
	_, err = s.db.Exec(cleanBanners)
	if err != nil {
		return err
	}

	cleanSlots := `DELETE FROM slots`
	_, err = s.db.Exec(cleanSlots)
	if err != nil {
		return err
	}

	cleanGroups := `DELETE FROM groups`
	_, err = s.db.Exec(cleanGroups)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddEventTimestamp(ctx context.Context, eventType string, rotationID int) error {
	query := ` INSERT INTO events(rotation_id, stamp, event_type)
	VALUES (:rotation_id, now(), :event_type)
	`
	eventStamp := event{
		RotationID: rotationID,
		Type:       eventType,
	}

	res, err := s.db.NamedExecContext(ctx, query, eventStamp)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRowWasAffected
	}

	return nil
}

func (s *Storage) GetRotationID(ctx context.Context, bannerID, slotID, groupID uuid.UUID) (int, error) {
	query := `
	SELECT id FROM rotations
	WHERE
	banner_id=$1 AND slot_id=$2 AND group_id=$3 AND deleted=FALSE
	`

	row := s.db.QueryRowxContext(ctx, query, bannerID, slotID, groupID)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Storager implementation.
func (s *Storage) Connect() (err error) {
	s.db, err = sqlx.Connect("pgx", s.connStr)
	return
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddBanner(ctx context.Context, bannerInfo types.Banner) error {
	insertBannerQuery := `
	INSERT INTO banners (id, description) VALUES (:id, :description);
	`

	dbBanner := banner{
		ID:          bannerInfo.ID,
		Description: bannerInfo.Description,
	}
	_, err := s.db.NamedExecContext(ctx, insertBannerQuery, dbBanner)
	return err
}

func (s *Storage) GetBanner(ctx context.Context, bannerID uuid.UUID) (types.Banner, error) {
	query := `
	SELECT * FROM banners WHERE id=$1 AND deleted=FALSE
	`
	row := s.db.QueryRowxContext(ctx, query, bannerID)
	if row.Err() != nil {
		return types.Banner{}, row.Err()
	}

	var dbBanner banner
	err := row.StructScan(&dbBanner)
	if err != nil {
		return types.Banner{}, err
	}

	resultBanner := types.Banner{
		ID:          dbBanner.ID,
		Description: dbBanner.Description,
	}

	return resultBanner, nil
}

func (s *Storage) DeleteBanner(ctx context.Context, bannerID uuid.UUID) error {
	deleteBannerQuery := `
	UPDATE banners SET deleted=TRUE, deleted_at=now()
	WHERE id=$1
	`
	deleteRotationsQuery := `
	UPDATE rotations SET deleted=TRUE, deleted_at=now()
	where banner_id=$1
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = execTxQuery(tx, deleteBannerQuery, bannerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = execTxQuery(tx, deleteRotationsQuery, bannerID)
	switch {
	case errors.Is(err, ErrNoRowWasAffected):
	case err == nil:
	default:
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Storage) AddSlot(ctx context.Context, slotInfo types.Slot) error {
	insertSlotQuery := `
	INSERT INTO slots (id, description) VALUES (:id, :description);
	`
	dbSlot := slot{
		ID:          slotInfo.ID,
		Description: slotInfo.Description,
	}
	_, err := s.db.NamedExecContext(ctx, insertSlotQuery, dbSlot)
	return err
}

func (s *Storage) GetSlot(ctx context.Context, slotID uuid.UUID) (types.Slot, error) {
	query := `
	SELECT * FROM slots WHERE id=$1 AND deleted=FALSE
	`
	row := s.db.QueryRowxContext(ctx, query, slotID)
	if row.Err() != nil {
		return types.Slot{}, row.Err()
	}

	var dbSlot slot
	err := row.StructScan(&dbSlot)
	if err != nil {
		return types.Slot{}, err
	}

	resultSlot := types.Slot{
		ID:          dbSlot.ID,
		Description: dbSlot.Description,
	}
	return resultSlot, nil
}

func (s *Storage) DeleteSlot(ctx context.Context, slotID uuid.UUID) error {
	deleteSlotQuery := `
	UPDATE slots SET deleted=TRUE, deleted_at=now()
	WHERE id=$1
	`
	deleteRotationsQuery := `
	UPDATE rotations SET deleted=TRUE, deleted_at=now()
	where slot_id=$1
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = execTxQuery(tx, deleteSlotQuery, slotID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = execTxQuery(tx, deleteRotationsQuery, slotID)
	switch {
	case errors.Is(err, ErrNoRowWasAffected):
	case err == nil:
	default:
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Storage) AddGroup(ctx context.Context, groupInfo types.Group) error {
	insertGroupQuery := `
	INSERT INTO groups (id, description) VALUES (:id, :description);
	`
	dbGroup := group{
		ID:          groupInfo.ID,
		Description: groupInfo.Description,
	}
	_, err := s.db.NamedExecContext(ctx, insertGroupQuery, dbGroup)
	return err
}

func (s *Storage) GetGroup(ctx context.Context, groupID uuid.UUID) (types.Group, error) {
	query := `
	SELECT * FROM groups WHERE id=$1 AND deleted=FALSE
	`
	row := s.db.QueryRowxContext(ctx, query, groupID)
	if row.Err() != nil {
		return types.Group{}, row.Err()
	}

	var dbGroup group
	err := row.StructScan(&dbGroup)
	if err != nil {
		return types.Group{}, err
	}

	resultGroup := types.Group{
		ID:          dbGroup.ID,
		Description: dbGroup.Description,
	}

	return resultGroup, nil
}

func (s *Storage) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	deleteGroupQuery := `
	UPDATE groups SET deleted=TRUE, deleted_at=now()
	WHERE id=$1
	`
	deleteRotationsQuery := `
	UPDATE rotations SET deleted=TRUE, deleted_at=now()
	where group_id=$1
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = execTxQuery(tx, deleteGroupQuery, groupID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = execTxQuery(tx, deleteRotationsQuery, groupID)
	switch {
	case errors.Is(err, ErrNoRowWasAffected):
	case err == nil:
	default:
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Storage) AddRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) (types.Rotation, error) {
	insertRotationQuery := `
	INSERT INTO rotations (banner_id, slot_id, group_id, shows, clicks)
	VALUES (:banner_id, :slot_id, :group_id, :shows, :clicks);
	`
	rotation := rotation{
		BannerID: bannerID,
		SlotID:   slotID,
		GroupID:  groupID,
	}
	_, err := s.db.NamedExecContext(ctx, insertRotationQuery, rotation)

	resultRotation := types.Rotation{
		BannerID: rotation.BannerID,
		SlotID:   rotation.SlotID,
		GroupID:  rotation.GroupID,
	}

	return resultRotation, err
}

func (s *Storage) DeleteRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	deleteRotationQuery := `
	UPDATE rotations SET deleted=TRUE, deleted_at=now()
	WHERE
	banner_id=$1 AND slot_id=$2 AND group_id=$3
	`
	res, err := s.db.ExecContext(ctx, deleteRotationQuery, bannerID, slotID, groupID)
	if err != nil {
		return err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsUpdated == 0 {
		return errors.New("no rotation was deleted")
	}

	return nil
}

func (s *Storage) GetRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) (types.Rotation, error) {
	rotationID, err := s.GetRotationID(ctx, bannerID, slotID, groupID)
	if err != nil {
		return types.Rotation{}, err
	}

	selectRotationQuery := `
	SELECT * FROM rotations WHERE id=$1 AND deleted=FALSE
	`
	row := s.db.QueryRowxContext(ctx, selectRotationQuery, rotationID)
	if row.Err() != nil {
		return types.Rotation{}, row.Err()
	}

	var dbRotation rotation
	err = row.StructScan(&dbRotation)
	if err != nil {
		return types.Rotation{}, err
	}
	resultRotation := types.Rotation{
		BannerID: dbRotation.BannerID,
		SlotID:   dbRotation.SlotID,
		GroupID:  dbRotation.GroupID,
		Shows:    dbRotation.Shows,
		Clicks:   dbRotation.Clicks,
	}

	return resultRotation, nil
}

func (s *Storage) GetRotationStats(ctx context.Context, bannerID, slotID, groupID uuid.UUID) ([]types.Event, error) {
	query := `
	SELECT * FROM events WHERE rotation_id=$1
	`
	rotationID, err := s.GetRotationID(ctx, bannerID, slotID, groupID)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, rotationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var events []types.Event
	for rows.Next() {
		dbEvent := event{}
		err := rows.StructScan(&dbEvent)
		if err != nil {
			return nil, err
		}

		events = append(
			events,
			types.Event{
				Type:      dbEvent.Type,
				Timestamp: dbEvent.Timestamp,
			},
		)
	}
	return events, err
}

func (s *Storage) AddShow(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	rotationID, err := s.GetRotationID(ctx, bannerID, slotID, groupID)
	if err != nil {
		return err
	}

	query := `
	UPDATE rotations SET shows=shows+1 WHERE id=$1
	`
	res, err := s.db.ExecContext(ctx, query, rotationID)
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

	err = s.AddEventTimestamp(ctx, EventTypeShow, rotationID)

	return err
}

func (s *Storage) AddClick(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	rotationID, err := s.GetRotationID(ctx, bannerID, slotID, groupID)
	if err != nil {
		return err
	}

	query := `
	UPDATE rotations SET clicks=clicks+1 WHERE id=$1
	`
	res, err := s.db.ExecContext(ctx, query, rotationID)
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

	err = s.AddEventTimestamp(ctx, EventTypeClick, rotationID)

	return err
}

func (s *Storage) GetAllRotations(ctx context.Context) ([]types.Rotation, error) {
	query := `SELECT * FROM rotations WHERE deleted=FALSE`

	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var rotations []types.Rotation
	for rows.Next() {
		var r rotation

		err := rows.StructScan(&r)
		if err != nil {
			return nil, err
		}

		rotations = append(
			rotations,
			types.Rotation{
				BannerID: r.BannerID,
				SlotID:   r.SlotID,
				GroupID:  r.GroupID,
				Shows:    r.Shows,
				Clicks:   r.Clicks,
			},
		)
	}

	return rotations, nil
}

func (s *Storage) GetTotalShows(ctx context.Context) (int64, error) {
	query := `SELECT sum(shows) FROM rotations WHERE deleted=FALSE`

	row := s.db.QueryRowxContext(ctx, query)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var totalShows sql.NullInt64
	err := row.Scan(&totalShows)
	if err != nil {
		return 0, err
	}

	return totalShows.Int64, nil
}
