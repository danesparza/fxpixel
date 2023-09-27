/* Turn on foreign key support (I can't believe we have to do this) */
/* More information: https://www.sqlite.org/foreignkeys.html */
PRAGMA foreign_keys = ON;

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
    gpio    integer,
    tags    TEXT    /* JSON string array */
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
    timeline_id   TEXT    not null,
    step_type_id integer not null,
    effect_type_id integer,
    led_range     TEXT, /* Null or blank means use the entire strip */
    step_time integer, /* Time (in ms) to execute the step */
    step_meta    TEXT,
    step_number    integer not null,

    FOREIGN KEY(timeline_id) REFERENCES timeline(id) ON DELETE CASCADE
    FOREIGN KEY(step_type_id) REFERENCES timeline_step_type(id)
    FOREIGN KEY(effect_type_id) REFERENCES timeline_step_effect_type(id)

    /* Each step number (ordinal) in a timeline must be unique */
    UNIQUE(timeline_id,step_number)

);
