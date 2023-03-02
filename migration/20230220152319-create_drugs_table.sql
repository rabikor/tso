-- +migrate Up
CREATE TABLE IF NOT EXISTS drugs
(
    id    INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) NOT NULL
) engine = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS drugs;
