ALTER TABLE environments DROP CONSTRAINT IF EXISTS environments_project_name_unique;
ALTER TABLE environments ADD CONSTRAINT environments_name_key UNIQUE (name);