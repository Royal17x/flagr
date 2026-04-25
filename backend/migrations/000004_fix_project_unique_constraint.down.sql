ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_org_name_unique;
ALTER TABLE projects ADD CONSTRAINT projects_name_key UNIQUE (name);