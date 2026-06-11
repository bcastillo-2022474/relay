-- lowercase alphanumeric segments separated by single hyphens, max 63 chars.
-- Must stay in sync with the Slug value object in internal/domain/slug.go.
create function is_valid_slug(slug text) returns boolean
language sql immutable as $$
    select slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$' and char_length(slug) <= 63
$$;

create function set_updated_at() returns trigger
language plpgsql as $$
begin
    new.updated_at = now();
    return new;
end;
$$;

create table organizations (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    slug text not null unique check (is_valid_slug(slug)),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table users (
    id uuid primary key default gen_random_uuid(),
    email text not null unique,
    name text not null,
    last_name text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    organization_id uuid not null references organizations(id) on delete cascade
);

create table applications (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    slug text not null check (is_valid_slug(slug)),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    organization_id uuid not null references organizations(id) on delete cascade,
    unique (organization_id, slug),
    -- target for the composite FKs below; makes org-mismatched children unrepresentable
    unique (id, organization_id)
);

create table event_types (
    id uuid primary key default gen_random_uuid(),
    application_id uuid not null,
    name text not null, -- e.g. "invoice.paid"
    payload_schema jsonb, -- optional JSON Schema used by ingest validation
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    organization_id uuid not null,
    unique (application_id, name),
    unique (id, organization_id),
    foreign key (application_id, organization_id)
        references applications (id, organization_id) on delete cascade
);

create table endpoints (
    id uuid primary key default gen_random_uuid(),
    application_id uuid not null,
    url text not null,
    description text not null default '',
    signing_secret text not null,
    disabled boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    organization_id uuid not null,
    unique (id, organization_id),
    foreign key (application_id, organization_id)
        references applications (id, organization_id) on delete cascade
);

create table endpoint_subscriptions (
    endpoint_id uuid not null,
    event_type_id uuid not null,
    created_at timestamptz not null default now(),

    organization_id uuid not null,
    primary key (endpoint_id, event_type_id),
    foreign key (endpoint_id, organization_id)
        references endpoints (id, organization_id) on delete cascade,
    foreign key (event_type_id, organization_id)
        references event_types (id, organization_id) on delete cascade
);

-- Postgres does not index FK columns automatically. Cover tenant-scoped
-- lookups and cascade-delete scans. applications(organization_id) and
-- event_types(application_id) are already covered by unique-constraint prefixes.
create index idx_users_organization_id on users (organization_id);
create index idx_event_types_organization_id on event_types (organization_id);
create index idx_endpoints_application_id on endpoints (application_id);
create index idx_endpoints_organization_id on endpoints (organization_id);
create index idx_endpoint_subscriptions_event_type_id on endpoint_subscriptions (event_type_id);
create index idx_endpoint_subscriptions_organization_id on endpoint_subscriptions (organization_id);

create trigger set_updated_at before update on organizations
    for each row execute function set_updated_at();
create trigger set_updated_at before update on users
    for each row execute function set_updated_at();
create trigger set_updated_at before update on applications
    for each row execute function set_updated_at();
create trigger set_updated_at before update on event_types
    for each row execute function set_updated_at();
create trigger set_updated_at before update on endpoints
    for each row execute function set_updated_at();

---- create above / drop below ----

drop table endpoint_subscriptions;
drop table endpoints;
drop table event_types;
drop table applications;
drop table users;
drop table organizations;
drop function set_updated_at();
drop function is_valid_slug(text);
