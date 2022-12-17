SET TIME ZONE 'UTC';

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
);

insert into videos (title, description, video_path, thumbnail_path, status, published_at)
    values ('Video #0', 'Video Description', 'No Path', 'No Thumbnail', 'No Status', now());
