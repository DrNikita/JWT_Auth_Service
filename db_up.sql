DROP TABLE IF EXISTS video_history;
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS role;
DROP TABLE IF EXISTS job_role;
DROP TABLE IF EXISTS settlement_type;
DROP TABLE IF EXISTS address;

CREATE TABLE IF NOT EXISTS role (
	id integer PRIMARY KEY,
	name varchar(256)
);

CREATE TABLE IF NOT EXISTS job_role (
	id integer PRIMARY KEY,
    role_id integer FOREIGN KEY(role.id),
	name varchar(256),

);

CREATE TABLE IF NOT EXISTS settlement_type (
	id integer PRIMARY KEY,
	name varchar(256)
);

CREATE TABLE IF NOT EXISTS address (
	id bigint PRIMARY KEY,
    settlement_type_id integer FOREIGN KEY(settlement_type.id),
	country varchar(256),
	region varchar(256),
	district varchar(256),
	settlement varchar(256),
	street varchar(256),
	hous_number varchar(256),
	flat_number varchar(256)
);

CREATE TABLE IF NOT EXISTS user (
	id bigint PRIMARY KEY,
	job_role_id integer FOREIGN KEY(job_role.id),
	address_id bigint FOREIGN KEY(address.id),
	name varchar(256),
	second_name varchar(256),
	surname varchar(256),
	email varchar(256),
	password varchar(256),
	birthday bigint,
	is_active boolean
);

CREATE TABLE IF NOT EXISTS video_history (
	id bigint PRIMARY KEY,
	user_id bigint FOREIGN KEY(user.id),
	video_name varchar(256),
	created_at bigint
);