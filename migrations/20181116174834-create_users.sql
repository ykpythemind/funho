
-- +migrate Up
CREATE TABLE users (
  id integer PRIMARY KEY AUTO_INCREMENT,
  name varchar(128) NOT NULL UNIQUE,
  password varchar(128) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE chats (
  id integer PRIMARY KEY AUTO_INCREMENT,
  body varchar(256) NOT NULL,
  user_id integer NOT NULL,
  room_user_id integer NOT NULL,
  created_at datetime,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (room_user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE chats;
DROP TABLE users;
