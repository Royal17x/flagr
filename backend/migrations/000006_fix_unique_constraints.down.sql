ALTER TABLE flags DROP CONSTRAINT IF EXISTS flags_project_key_unique;
ALTER TABLE flags ADD CONSTRAINT flags_key_key UNIQUE (key);

ALTER TABLE environments DROP CONSTRAINT IF EXISTS environments_project_name_unique;
ALTER TABLE environments ADD CONSTRAINT environments_name_key UNIQUE (name);

ALTER TABLE environments DROP CONSTRAINT IF EXISTS environments_project_slug_unique;
ALTER TABLE environments ADD CONSTRAINT environments_slug_key UNIQUE (slug);