create table messages (
    id               uuid        primary key,
    organization_id  uuid        not null,
    application_id   uuid        not null,
    event_type_id    uuid        not null,
    payload          jsonb       not null,
    status           text        not null default 'pending'
                     check (status in ('pending', 'published', 'failed')),
    attempts         int         not null default 0,
    next_attempt_at  timestamptz not null default now(),
    created_at       timestamptz not null default now(),
    published_at     timestamptz,

    unique (id, organization_id),
    foreign key (application_id, organization_id) references applications (id, organization_id) on delete cascade,
    foreign key (event_type_id, organization_id)  references event_types  (id, organization_id) on delete cascade
);

create index messages_pending_idx
    on messages (next_attempt_at)
    where status = 'pending';

---- create above / drop below ----

drop table messages;
