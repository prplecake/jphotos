CREATE EXTENSION
IF NOT EXISTS "uuid-ossp";

CREATE TABLE
IF NOT EXISTS users
(
	-- per https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)
	-- uuid_generate_v4() generates a random UUID
	id uuid 	DEFAULT	uuid_generate_v4(),
	username	TEXT UNIQUE	NOT NULL,
	hash		TEXT		NOT NULL,
	created		TIMESTAMPTZ	NOT NULL,
	last_login	TIMESTAMPTZ,
	PRIMARY KEY (id)
);

CREATE TABLE
IF NOT EXISTS sessions
(
	user_id	UUID 	REFERENCES users (id),
	token	TEXT	UNIQUE,
	expires	TIMESTAMPTZ
);

CREATE TABLE
IF NOT EXISTS groups
(
	id 	uuid 	DEFAULT uuid_generate_v4(),
	group_name	TEXT	NOT NULL,
	created		TIMESTAMPTZ	DEFAULT NOW(),
	PRIMARY KEY (id)
);

CREATE TABLE 
IF NOT EXISTS member
(
	member	uuid REFERENCES users (id),
	groups	uuid REFERENCES groups (id),
	admin	BOOLEAN
);

CREATE TABLE
IF NOT EXISTS albums
(
	id uuid DEFAULT uuid_generate_v4(),
	name	TEXT	NOT NULL,
	slug	TEXT	NOT NULL,
	created	TIMESTAMPTZ	DEFAULT NOW(),
	PRIMARY KEY (id)
);

CREATE TABLE
IF NOT EXISTS photos
(
	id 	uuid	DEFAULT uuid_generate_v4(),
	caption		TEXT,
	location 	TEXT,
	added		TIMESTAMPTZ DEFAULT NOW(),
	PRIMARY KEY (id)
);

CREATE TABLE
IF NOT EXISTS album_photos
(
	photo uuid REFERENCES photos (id) ON DELETE CASCADE,
	album uuid REFERENCES albums (id) ON DELETE CASCADE
);