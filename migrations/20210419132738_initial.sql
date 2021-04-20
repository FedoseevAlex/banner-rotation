-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banners (
    id          UUID PRIMARY KEY,
    description TEXT
);

CREATE TABLE IF NOT EXISTS slots (
    id          UUID PRIMARY KEY,
    description TEXT
);

CREATE TABLE IF NOT EXISTS groups (
    id          UUID PRIMARY KEY,
    description TEXT
);

CREATE TABLE IF NOT EXISTS rotations (
    banner_id UUID ,
    slot_id   UUID,
    group_id  UUID,
    shows     INT,
    clicks    INT,

    FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE,
    FOREIGN KEY (slot_id)   REFERENCES slots(id)   ON DELETE CASCADE,
    FOREIGN KEY (group_id)  REFERENCES groups(id)  ON DELETE CASCADE,
    PRIMARY KEY (banner_id, slot_id, group_id)
);

CREATE TABLE IF NOT EXISTS total_shows (
    count BIGINT
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF NOT EXISTS total_shows;
DROP TABLE IF NOT EXISTS rotations;
DROP TABLE IF NOT EXISTS banners;
DROP TABLE IF NOT EXISTS slots;
DROP TABLE IF NOT EXISTS groups;
-- +goose StatementEnd
