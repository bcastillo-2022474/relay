-- Dev-only seed (not a tern migration: schema is migrated, data is seeded).
-- The auth middleware stubs this org until real JWT verification lands.
INSERT INTO organizations (id, name, slug)
VALUES ('00000000-0000-0000-0000-000000000001', 'Dev Organization', 'dev-org')
ON CONFLICT (id) DO NOTHING;
