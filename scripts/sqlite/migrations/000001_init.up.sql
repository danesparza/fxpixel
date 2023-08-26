create table system_config
(
    gpio integer
);

create table timeline
(
    id      TEXT not null
        constraint Timeline_pk
            primary key,
    enabled INT     default 1,
    created integer default CURRENT_TIMESTAMP,
    name    TEXT,
    gpio    integer
);

create table timeline_frame_type
(
    id         integer
        constraint timeline_frame_type_pk
            primary key autoincrement,
    frame_type TEXT not null
);

create table timeline_scene_type
(
    id         integer
        constraint timeline_scene_type_pk
            primary key autoincrement,
    scene_type TEXT not null
);

create table timeline_frame
(
    id            TEXT    not null
        constraint timeline_frame_pk
            primary key,
    timeline_id   TEXT    not null
        constraint timeline_frame_timeline_id_fk
            references timeline,
    frame_type_id integer not null
        constraint timeline_frame_timeline_frame_type_id_fk
            references timeline_frame_type,
    scene_type_id integer
        constraint timeline_frame_timeline_scene_type_id_fk
            references timeline_scene_type,
    led_range     TEXT    not null,
    scene_meta    TEXT
);

