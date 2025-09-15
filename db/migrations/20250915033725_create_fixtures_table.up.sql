create table fixtures (
    id serial primary key,
    sport_id int not null references sports(id),
    team_id_1 int not null references teams(id),
    team_id_2 int null references teams(id),
    date_time timestamp not null,
    details jsonb not null,
    status VARCHAR(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    deleted_at timestamp null
);

create index idx_fixtures_sport_id on fixtures (sport_id);
create index idx_fixtures_team_id_1 on fixtures (team_id_1, team_id_2);
create index idx_fixtures_date_time on fixtures (date_time);
