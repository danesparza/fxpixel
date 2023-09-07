create table system_config
(
    gpio integer not null,
    leds integer not null
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

create table timeline_step_type
(
    id         integer
        constraint timeline_step_type_pk
            primary key autoincrement,
    step_type TEXT not null
);

create table timeline_step_effect_type
(
    id         integer
        constraint timeline_step_effect_type_pk
            primary key autoincrement,
    effect_type TEXT not null
);

create table timeline_step
(
    id            TEXT    not null
        constraint timeline_step_pk
            primary key,
    timeline_id   TEXT    not null
        constraint timeline_step_timeline_id_fk
            references timeline,
    step_type_id integer not null
        constraint timeline_step_timeline_step_type_id_fk
            references timeline_step_type,
    effect_type_id integer
        constraint timeline_step_timeline_step_effect_type_id_fk
            references timeline_step_effect_type,
    led_range     TEXT, /* Null means use the entire strip */
    step_time integer, /* Time (in ms) to execute the step */
    step_meta    TEXT,
    step_number    integer not null,

    /* Each step number (ordinal) in a timeline must be unique */
    UNIQUE(timeline_id,step_number)

);

