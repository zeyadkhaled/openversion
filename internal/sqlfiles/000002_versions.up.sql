CREATE TABLE IF NOT EXISTS backend.versions (
    "id" text PRIMARY KEY,
    "min_version" text NOT NULL,
    "package" text NOT NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL
);


INSERT INTO backend.versions ("id", "min_version","package", "created_at", "updated_at")
    VALUES ('demo_1', '999','me.zdkd.app', '2019-12-05 12:10:21.553', '2019-12-05 12:10:21.553')
ON CONFLICT ("id")
    DO NOTHING;

