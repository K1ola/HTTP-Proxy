CREATE TABLE IF NOT EXISTS `requests` (
                            `id` INTEGER PRIMARY KEY AUTOINCREMENT,
                            `method` VARCHAR(64) NOT NULL,
                            `uri` VARCHAR(64) NOT NULL,
                            `proto` VARCHAR(64) NOT NULL,
                            `created` DATE NULL
);

CREATE TABLE IF NOT EXISTS  `headers` (
                           `id` INTEGER PRIMARY KEY AUTOINCREMENT,
                           `request_id` INTEGER NOT NULL REFERENCES requests(id),
                           `key` VARCHAR(64) NOT NULL,
                           `value` VARCHAR(64) NOT NULL
);
