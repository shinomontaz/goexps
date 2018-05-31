CREATE SEQUENCE IF NOT EXISTS area_serial START 1;

DROP TABLE IF EXISTS area;

CREATE TABLE area (
id integer NOT NULL DEFAULT nextval('area_serial'::regclass),
name varchar(128) NOT NULL,
CONSTRAINT area_pkey PRIMARY KEY(id)
);

CREATE SEQUENCE IF NOT EXISTS zone_serial START 1;

DROP TABLE IF EXISTS zone;

CREATE TABLE zone (
id integer NOT NULL DEFAULT nextval('zone_serial'::regclass),
fk_area integer,
name varchar(128) NOT NULL,
geom GEOMETRY,
CONSTRAINT zone_pkey PRIMARY KEY(id),
CONSTRAINT fk_area FOREIGN KEY (fk_area)
      REFERENCES area (id) MATCH SIMPLE
      ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE SEQUENCE IF NOT EXISTS fleet_serial START 1;

DROP TABLE IF EXISTS fleet;

CREATE TABLE fleet (
id integer NOT NULL DEFAULT nextval('fleet_serial'::regclass),
fk_zone integer,
name varchar(128) NOT NULL,
weight numeric(10,2) NOT NULL,
volume numeric(10, 2) NOT NULL,
CONSTRAINT fleet_pkey PRIMARY KEY(id),
CONSTRAINT fk_zone FOREIGN KEY (fk_zone)
      REFERENCES zone (id) MATCH SIMPLE
      ON UPDATE CASCADE ON DELETE CASCADE
);

INSERT INTO area (id, name) VALUES (1, 'Москва и область');

COPY zone(geom, fk_area, name) FROM '/Users/drybakov/go/src/github.com/shinomontaz/goexps/FleetManager/database/boundaries.csv' DELIMITER ',' CSV;

INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (1, 'asdasd', 1500, 4.5);
INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (1, 'asdasd 2', 1500, 4.5);

INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (2, 'aaa-2', 1000, 3.5);
INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (3, 'aaa-3', 1500, 4.5);

INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (4, 'aaa-4', 1500, 4.5);
INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (4, 'bbb-4', 500, 2.0);

INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (5, 'asd-5', 1500, 4.5);
INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (6, 'fff-6', 1000, 5.5);

INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (7, 'oiru', 1500, 4.5);
INSERT INTO fleet (fk_zone, name, weight, volume) VALUES (10, 'nss', 1000, 4.5);

