-- +goose Up

PRAGMA foreign_keys = ON;

CREATE TABLE sales (
    sales_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

INSERT INTO sales (name)
VALUES
    ('Nguyễn Văn An'),
    ('Trần Thị Bình'),
    ('Lê Văn Cường');

CREATE TABLE outlets (
    outlet_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    channel TEXT NOT NULL,
    tier TEXT NOT NULL,
    sales_id INTEGER NOT NULL,
    stage TEXT NOT NULL,
    note TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_outlets_sales
        FOREIGN KEY (sales_id)
        REFERENCES sales(sales_id)
);

CREATE TABLE working_schedules (
    schedule_id INTEGER PRIMARY KEY AUTOINCREMENT,
    outlet_id INTEGER NOT NULL,
    sales_id INTEGER NOT NULL,
    address TEXT NOT NULL,
    schedule_date DATETIME NOT NULL,
    current_stage TEXT NOT NULL,
    expected_stage TEXT,
    note TEXT,
    sync_status TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_working_schedules_outlet
        FOREIGN KEY (outlet_id)
        REFERENCES outlets(outlet_id),

    CONSTRAINT fk_working_schedules_sales
        FOREIGN KEY (sales_id)
        REFERENCES sales(sales_id),

    CONSTRAINT uq_working_schedules_sales_outlet_date
        UNIQUE (sales_id, outlet_id, schedule_date)
);

CREATE TABLE files (
    file_id INTEGER PRIMARY KEY AUTOINCREMENT,
    object_key TEXT NOT NULL,
    file_name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE evidences (
    evidence_id INTEGER PRIMARY KEY AUTOINCREMENT,
    schedule_id INTEGER NOT NULL,
    file_id INTEGER NOT NULL,

    CONSTRAINT fk_evidences_schedule
        FOREIGN KEY (schedule_id)
        REFERENCES working_schedules(schedule_id),

    CONSTRAINT fk_evidences_file
        FOREIGN KEY (file_id)
        REFERENCES files(file_id)
);

-- Indexes
CREATE INDEX idx_outlets_sales_id
ON outlets(sales_id);

CREATE INDEX idx_working_schedules_outlet_id
ON working_schedules(outlet_id);

CREATE INDEX idx_working_schedules_sales_id
ON working_schedules(sales_id);

CREATE INDEX idx_working_schedules_schedule_date
ON working_schedules(schedule_date);

CREATE INDEX idx_files_object_key
ON files(object_key);

CREATE INDEX idx_evidences_schedule_id
ON evidences(schedule_id);

CREATE INDEX idx_evidences_file_id
ON evidences(file_id);

-- +goose Down

DROP INDEX IF EXISTS idx_evidences_file_id;
DROP INDEX IF EXISTS idx_evidences_schedule_id;
DROP INDEX IF EXISTS idx_files_object_key;
DROP INDEX IF EXISTS idx_working_schedules_schedule_date;
DROP INDEX IF EXISTS idx_working_schedules_sales_id;
DROP INDEX IF EXISTS idx_working_schedules_outlet_id;
DROP INDEX IF EXISTS idx_outlets_sales_id;

DROP TABLE IF EXISTS evidences;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS working_schedules;
DROP TABLE IF EXISTS outlets;
DROP TABLE IF EXISTS sales;