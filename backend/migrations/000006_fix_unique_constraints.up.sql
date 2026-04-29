ALTER TABLE flags DROP CONSTRAINT IF EXISTS flags_key_key;
ALTER TABLE flags ADD CONSTRAINT flags_project_key_unique UNIQUE (project_id, key);

ALTER TABLE environments DROP CONSTRAINT IF EXISTS environments_name_key;
ALTER TABLE environments ADD CONSTRAINT environments_project_name_unique UNIQUE (project_id, name);

ALTER TABLE environments DROP CONSTRAINT IF EXISTS environments_slug_key;
ALTER TABLE environments ADD CONSTRAINT environments_project_slug_unique UNIQUE (project_id, slug);

ALTER TABLE sdk_keys DROP CONSTRAINT IF EXISTS sdk_keys_created_by_fkey;
ALTER TABLE sdk_keys ALTER COLUMN created_by DROP NOT NULL;