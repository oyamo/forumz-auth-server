create table if not exists connection (
    user_id UUID not null,
    connected_to UUID not null,
    datetime_created timestamp not null default current_timestamp,
    primary key (user_id, connected_to),
    constraint fk_person_connection_user_id foreign key (user_id) references person(id),
    constraint fk_person_connection_connected_to foreign key (connected_to) references person(id)
);