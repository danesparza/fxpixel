/* Default config */
insert into system_config(gpio, leds) values (18, 5);

/* Initial step types */
INSERT INTO timeline_step_type (id, step_type) VALUES (1, 'effect');
INSERT INTO timeline_step_type (id, step_type) VALUES (2, 'sleep');
INSERT INTO timeline_step_type (id, step_type) VALUES (3, 'random-sleep');
INSERT INTO timeline_step_type (id, step_type) VALUES (4, 'trigger');
INSERT INTO timeline_step_type (id, step_type) VALUES (5, 'loop');

/* Initial effect types */
INSERT INTO timeline_step_effect_type (id, effect_type) VALUES (1, 'solid');
INSERT INTO timeline_step_effect_type (id, effect_type) VALUES (2, 'fade');
INSERT INTO timeline_step_effect_type (id, effect_type) VALUES (3, 'gradient');
INSERT INTO timeline_step_effect_type (id, effect_type) VALUES (4, 'sequence');
INSERT INTO timeline_step_effect_type (id, effect_type) VALUES (5, 'rainbow');
INSERT INTO timeline_step_effect_type (id, effect_type) VALUES (6, 'zip');
INSERT INTO timeline_step_effect_type (id, effect_type) VALUES (7, 'knight-rider');
INSERT INTO timeline_step_effect_type (id, effect_type) VALUES (8, 'lightning');

/* Test data */
INSERT INTO timeline (id, enabled, created, name, gpio, tags) VALUES ('tl1', 1, '2023-09-06 23:25:23', 'TL test', null, null);
INSERT INTO timeline_step (id, timeline_id, step_type_id, effect_type_id, led_range, step_time, step_meta, step_number) VALUES ('ts1', 'tl1', 1, 1, null, 3000, '{"color": {"R": 128}}', 1);
INSERT INTO timeline_step (id, timeline_id, step_type_id, effect_type_id, led_range, step_time, step_meta, step_number) VALUES ('ts2', 'tl1', 1, 1, null, 4000, '{"color": {"B": 128}}', 2);
INSERT INTO timeline_step (id, timeline_id, step_type_id, effect_type_id, led_range, step_time, step_meta, step_number) VALUES ('ts3', 'tl1', 1, 1, null, 5000, '{"color": {"G": 128}}', 3);
