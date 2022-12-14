create table if not exists videos (
    id bigserial primary key,
    title text,
    description text,
    video_path text,
    thumbnail_path text,
    status text not null,
    published_at timestamp(0) with time zone,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),
    version integer not null default 1
)