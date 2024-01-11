CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    nickname      VARCHAR(32)  NOT NULL,
    email         VARCHAR(64)  NOT NULL,
    password_hash varchar(255) NOT NULL
);

CREATE TABLE pastes
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(64) NOT NULL,
    description VARCHAR(128),
    data TEXT NOT NULL
);

CREATE TABLE  users_pastes
(
  id SERIAL PRIMARY KEY,
  user_id int references users(id) ON DELETE CASCADE NOT NULL,
  paste_id int references pastes(id) ON DELETE CASCADE NOT NULL
);