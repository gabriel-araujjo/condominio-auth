package postgres

import (
	"database/sql"

	"github.com/gabriel-araujjo/base62"
	"github.com/gabriel-araujjo/condominio-auth/config"
)

const dbVersion = 1

type scheme struct {
	conf *config.Config
}

func (s *scheme) OnCreate(db *sql.DB) (err error) {

	_, err = db.Exec(`
CREATE FUNCTION check_cpf(cpf INT8) RETURNS boolean AS $$
DECLARE
  dv1 INT8;
  dv2 INT8;
  calc_dv1 INT8;
  calc_dv2 INT8;
  i INTEGER;
BEGIN
  IF cpf IS NULL THEN
	  RETURN TRUE;
  END IF;

  IF (cpf >= 99999999999 OR
      cpf = 88888888888 OR
      cpf = 77777777777 OR
      cpf = 66666666666 OR
      cpf = 55555555555 OR
      cpf = 44444444444 OR
      cpf = 33333333333 OR
      cpf = 22222222222 OR
      cpf = 11111111111 OR
      cpf = 0) THEN
    RETURN false;
  END IF;
  dv2 := cpf % 10;
  cpf := cpf / 10;
  dv1 := cpf % 10;
  cpf := cpf / 10;

  i := 2;
  calc_dv1 := 0;
  calc_dv2 := dv1 * i;

  LOOP
    EXIT WHEN i >= 11;
    calc_dv1 := calc_dv1 + i * (cpf % 10);
    i := i + 1;
    calc_dv2 := calc_dv2 + i * (cpf % 10);
    cpf := cpf / 10;
  END LOOP;

  calc_dv1 := calc_dv1 % 11;
  calc_dv2 := calc_dv2 % 11;

  IF ( calc_dv1 <= 1 ) THEN
    calc_dv1 := 0;
  ELSE
    calc_dv1 := 11 - calc_dv1;
  END IF;

  IF ( calc_dv2 <= 1 ) THEN
    calc_dv2 := 0;
  ELSE
    calc_dv2 := 11 - calc_dv2;
  END IF;

  RETURN ( dv1 = calc_dv1 AND dv2 = calc_dv2);
END;
$$ LANGUAGE plpgsql;

CREATE EXTENSION pgcrypto;

CREATE FUNCTION check_email_uniqueness(m TEXT) RETURNS boolean AS $$
DECLARE
	c INT;
BEGIN
	IF m IS NULL THEN
		RETURN TRUE;
	END IF;

	SELECT count(*) INTO c
	FROM "email_lookup" e
	WHERE e.email = m LIMIT 1;

	IF c > 0 THEN
		RETURN FALSE;
	ELSE
		RETURN TRUE;
	END IF;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION check_phone_uniqueness(p TEXT) RETURNS boolean AS $$
DECLARE
	c INT;
BEGIN
	IF p IS NULL THEN
		RETURN TRUE;
	END IF;

	SELECT count(*) INTO c
	FROM "phone_lookup" e
	WHERE e.phone = p LIMIT 1;

	IF c > 0 THEN
		RETURN FALSE;
	ELSE
		RETURN TRUE;
	END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "user" (
  user_id SERIAL PRIMARY KEY,
  name TEXT NOT NULL CHECK (name != ''),
  cpf INT8 CHECK (check_cpf(cpf)),
  fb_id TEXT UNIQUE,
  avatar TEXT,
  hash TEXT,
  phone TEXT CHECK (phone ~ '^\d+$') CHECK (check_phone_uniqueness(phone)),
  phone_verified BOOLEAN DEFAULT FALSE,
  email TEXT CHECK (email ~ '^[^@\s]+@[^@\s]+$') CHECK (check_email_uniqueness(email)),
  email_verified BOOLEAN DEFAULT FALSE
);


CREATE TABLE "user_email" (
  user_id INTEGER REFERENCES "user"(user_id) ON DELETE CASCADE,
  email TEXT NOT NULL CHECK (email ~ '^[^@\s]+@[^@\s]+$') CHECK (check_email_uniqueness(email)),
  verified BOOLEAN DEFAULT FALSE,
    CONSTRAINT ct_user_email_pk PRIMARY KEY (user_id, email)
);


CREATE TABLE "user_phone" (
  user_id INTEGER REFERENCES "user"(user_id) ON DELETE CASCADE,
  phone TEXT NOT NULL CHECK (phone ~ '^\d+$') CHECK (check_phone_uniqueness(phone)),
  verified BOOLEAN DEFAULT FALSE,
    CONSTRAINT ct_user_phone_pk PRIMARY KEY (user_id, phone)
);

CREATE TABLE "client" (
	client_id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	secret TEXT NOT NULL
);

CREATE TABLE "scope" (
  scope_id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT
);
CREATE UNIQUE INDEX scope_name_idx ON "scope" (name); 

CREATE TABLE "authorization" (
  client_id INTEGER REFERENCES "client"(client_id) ON DELETE CASCADE,
  user_id INTEGER REFERENCES "user"(user_id) ON DELETE CASCADE,
  scope_id INTEGER REFERENCES "user"(user_id) ON DELETE CASCADE,
    CONSTRAINT authorization_pk PRIMARY KEY(client_id, user_id, scope_id)
);

CREATE VIEW "email_lookup" AS
	SELECT u.user_id, u.email FROM "user" u
	UNION ALL SELECT e.user_id, e.email FROM "user_email" e;

CREATE VIEW "phone_lookup" AS
	SELECT u.user_id, u.phone FROM "user" u
	UNION ALL SELECT p.user_id, p.phone FROM "user_phone" p;

CREATE FUNCTION hash_password() RETURNS TRIGGER AS $$
BEGIN

  IF NEW.hash IS NOT NULL AND ( TG_OP = 'INSERT' OR NEW.hash != OLD.hash ) THEN
	  NEW.hash = crypt(NEW.hash, gen_salt('bf', 8));
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tg_user_hash BEFORE INSERT OR UPDATE ON "user"
FOR EACH ROW EXECUTE PROCEDURE hash_password();

`)

	if err != nil {
		return
	}

	for _, c := range s.conf.Clients {
		clientID, _ := base62.ParseUint(c.PublicID)
		c.ID = int64(clientID)
		err = db.QueryRow(`
			INSERT INTO "client"(client_id, name, secret)
			VALUES ($1, $2, $3)
			RETURNING "client".client_id
		`,
			c.ID, c.Name, c.Secret).Scan(&c.ID)
		if err != nil {
			return
		}
	}
	return
}

func (s *scheme) OnUpdate(db *sql.DB, oldVersion int) error {
	return nil
}

func (s *scheme) Version() int {
	return dbVersion
}

func (s *scheme) VersionStrategy() string {
	return s.conf.Dao.VersionStrategy
}
