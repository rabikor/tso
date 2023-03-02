-- +migrate Up
CREATE TABLE IF NOT EXISTS illnesses
(
    id    INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) NOT NULL
) engine = InnoDB;

CREATE TABLE IF NOT EXISTS procedures
(
    id    INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) NOT NULL
) engine = InnoDB;

CREATE TABLE IF NOT EXISTS schemes
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    illness_id INT UNSIGNED NOT NULL,
    length     INT UNSIGNED NOT NULL
) engine = InnoDB;

CREATE TABLE IF NOT EXISTS scheme_days
(
    id           INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    scheme_id    INT UNSIGNED NOT NULL,
    procedure_id INT UNSIGNED NOT NULL,
    drug_id      INT UNSIGNED DEFAULT NULL,
    day_number   INT UNSIGNED NOT NULL COMMENT 'Порядковий номер дня в схемі',
    times        INT UNSIGNED NOT NULL,
    each_hours   INT UNSIGNED NOT NULL
) engine = InnoDB;

ALTER TABLE schemes
    ADD CONSTRAINT fk_schemes_illnesses FOREIGN KEY (illness_id) REFERENCES illnesses (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE scheme_days
    ADD CONSTRAINT fk_scheme_days_procedures FOREIGN KEY (procedure_id) REFERENCES procedures (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE scheme_days
    ADD CONSTRAINT unq_scheme_medications UNIQUE (scheme_id, procedure_id, drug_id);

ALTER TABLE scheme_days
    ADD CONSTRAINT fk_scheme_days_schemes FOREIGN KEY (scheme_id) REFERENCES schemes (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE scheme_days
    ADD CONSTRAINT fk_scheme_days_drugs FOREIGN KEY (drug_id) REFERENCES drugs (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +migrate Down
DROP TABLE IF EXISTS illnesses;
DROP TABLE IF EXISTS procedures;
DROP TABLE IF EXISTS schemes;
DROP TABLE IF EXISTS scheme_days;

ALTER TABLE schemes
    DROP CONSTRAINT fk_schemes_illnesses;
ALTER TABLE scheme_days
    DROP CONSTRAINT fk_scheme_days_procedures;
ALTER TABLE scheme_days
    DROP CONSTRAINT fk_scheme_days_schemes;
ALTER TABLE scheme_days
    DROP CONSTRAINT fk_scheme_days_drugs;
