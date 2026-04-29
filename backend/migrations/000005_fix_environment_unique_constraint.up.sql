ALTER TABLE environments DROP CONSTRAINT IF EXISTS environments_name_key;
ALTER TABLE environments ADD CONSTRAINT environments_project_name_unique UNIQUE (project_id, name);