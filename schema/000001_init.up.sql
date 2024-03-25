CREATE TABLE IF NOT EXISTS users
(
  id            SERIAL       NOT NULL UNIQUE,
  name          VARCHAR(255) NOT NULL,
  username      VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users_sessions
(
  user_id       INT          REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  refresh_token VARCHAR(255) UNIQUE,
  expires_at    TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS todo_lists
(
  id          SERIAL       NOT NULL UNIQUE,
  title       VARCHAR(255) NOT NULL,
  description VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS users_lists
(
  id      SERIAL NOT NULL UNIQUE,
  user_id INT    REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  list_id INT    REFERENCES todo_lists (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS todo_items
(
  id          SERIAL       NOT NULL UNIQUE,
  title       VARCHAR(255) NOT NULL,
  description VARCHAR(255),
  done        BOOLEAN      NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS lists_items
(
  id      SERIAL NOT NULL UNIQUE,
  item_id INT    REFERENCES todo_items (id) ON DELETE CASCADE NOT NULL,
  list_id INT    REFERENCES todo_lists (id) ON DELETE CASCADE NOT NULL
);
