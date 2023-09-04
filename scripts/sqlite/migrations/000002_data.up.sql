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
