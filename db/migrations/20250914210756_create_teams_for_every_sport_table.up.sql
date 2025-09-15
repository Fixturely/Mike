create table teams (
    id serial primary key,
    name varchar(255) not null unique,
    sport_id int not null references sports(id),
    description text null,
    image_url varchar(500) null,
    is_active boolean not null default true,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    deleted_at timestamp null
);

create index idx_teams_name_sport_id on teams (name, sport_id);
create index idx_teams_sport_id on teams (sport_id);