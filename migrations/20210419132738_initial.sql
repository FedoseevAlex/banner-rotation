-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banners (
    id          UUID PRIMARY KEY,
    description TEXT,
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS slots (
    id          UUID PRIMARY KEY,
    description TEXT,
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS groups (
    id          UUID PRIMARY KEY,
    description TEXT,
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rotations (
    banner_id UUID,
    slot_id   UUID,
    group_id  UUID,
    shows     INT,
    clicks    INT,
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP,

    FOREIGN KEY (banner_id) REFERENCES banners(id),
    FOREIGN KEY (slot_id)   REFERENCES slots(id),
    FOREIGN KEY (group_id)  REFERENCES groups(id),
    PRIMARY KEY (banner_id, slot_id, group_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS total_shows;
DROP TABLE IF EXISTS rotations;
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS slots;
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd
