DROP TABLE IF EXISTS video_history;
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS role;
DROP TABLE IF EXISTS job_role;
DROP TABLE IF EXISTS settlement_type;
DROP TABLE IF EXISTS address;

CREATE TABLE IF NOT EXISTS role (
	id serial PRIMARY KEY,
	name varchar(256) UNIQUE
);

CREATE TABLE IF NOT EXISTS job_role (
	id serial PRIMARY KEY,
    role_id integer REFERENCES role(id),
	name varchar(256) UNIQUE
);

CREATE TABLE IF NOT EXISTS settlement_type (
	id serial PRIMARY KEY,
	name varchar(256) UNIQUE
);

CREATE TABLE IF NOT EXISTS address (
	id serial PRIMARY KEY,
    settlement_type_id integer REFERENCES settlement_type(id) NOT NULL,
	country varchar(256),
	region varchar(256),
	district varchar(256),
	settlement varchar(256),
	street varchar(256),
	house_number varchar(256),
	flat_number varchar(256)
);

CREATE TABLE IF NOT EXISTS "user" (
	id serial PRIMARY KEY,
	job_role_id integer REFERENCES job_role(id) DEFAULT 1,
	address_id bigint REFERENCES address(id) DEFAULT null,
	name varchar(256),
	second_name varchar(256),
	surname varchar(256),
	email varchar(256) NOT NULL UNIQUE,
	password varchar(256) NOT NULL,
	birthday bigint NOT NULL,
	is_active boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS video_history (
	id serial PRIMARY KEY,
	user_id bigint REFERENCES "user"(id) NOT NULL,
	video_name varchar(256),
	created_at bigint
);

INSERT INTO role (name) VALUES ('client');
INSERT INTO role (name) VALUES ('admin');

INSERT INTO job_role (role_id, name) VALUES ((SELECT id FROM role WHERE name='admin'),'dev-ops');
INSERT INTO job_role (role_id, name) VALUES ((SELECT id FROM role WHERE name='admin'),'GO-developer');
INSERT INTO job_role (role_id, name) VALUES ((SELECT id FROM role WHERE name='client'),'unknown');
INSERT INTO job_role (role_id, name) VALUES ((SELECT id FROM role WHERE name='client'),'business analyst');
INSERT INTO job_role (role_id, name) VALUES ((SELECT id FROM role WHERE name='client'),'qa');
INSERT INTO job_role (role_id, name) VALUES ((SELECT id FROM role WHERE name='client'),'aqa');

INSERT INTO settlement_type (name) VALUES ('Поселок');
INSERT INTO settlement_type (name) VALUES ('Поселок городского типа');
INSERT INTO settlement_type (name) VALUES ('Деревня');
INSERT INTO settlement_type (name) VALUES ('Агрогородок');
INSERT INTO settlement_type (name) VALUES ('Город');