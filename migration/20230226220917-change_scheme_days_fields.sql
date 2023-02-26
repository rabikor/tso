-- +migrate Up
ALTER TABLE scheme_days DROP CONSTRAINT unq_scheme_medications;
ALTER TABLE scheme_days RENAME COLUMN day_number to `order`;
ALTER TABLE scheme_days RENAME COLUMN each_hours to frequency;
ALTER TABLE scheme_days
    ADD CONSTRAINT unq_scheme_days UNIQUE (scheme_id, procedure_id, drug_id, `order`);

-- +migrate Down
ALTER TABLE scheme_days DROP CONSTRAINT unq_scheme_days;
ALTER TABLE scheme_days RENAME COLUMN `order` to day_number;
ALTER TABLE scheme_days RENAME COLUMN frequency to each_hours;
ALTER TABLE scheme_days
    ADD CONSTRAINT unq_scheme_medications UNIQUE (scheme_id, procedure_id, drug_id);
