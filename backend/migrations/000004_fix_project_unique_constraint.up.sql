ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_name_key;
ALTER TABLE projects ADD CONSTRAINT projects_org_name_unique UNIQUE (organization_id, name);
