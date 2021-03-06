CREATE TABLE IF NOT EXISTS users
(
    user_id      VARCHAR(50)                     NOT NULL,
    name         VARCHAR(50)                 NOT NULL CHECK ( name <> '' ),
    email        VARCHAR(64) UNIQUE          NOT NULL CHECK ( email <> '' ),
    password     VARCHAR(250)                NOT NULL CHECK ( octet_length(password) <> 0 ),
    salt         VARCHAR(100)                NOT NULL CHECK ( salt <> ''),
    avatar       VARCHAR(512),
    phone_number VARCHAR(20),
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS categories
(
    id                  serial                     NOT NULL,
    name                VARCHAR(20)                NOT NULL,
    PRIMARY KEY (id)

);


CREATE TABLE IF NOT EXISTS sosial_media 
(   
  id            Serial                       NOT NULL,
  user_id       VARCHAR(50)                  NOT NULL,
  google_id     VARCHAR(70)                  NULL,
  facebook_id   VARCHAR(70)                  NULL,
  github_id     VARCHAR(70)                  NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TABLE IF NOT EXISTS videos
(
    id                      serial                          NOT NULL,
    category_id             INT                             NOT NULL,
    name                    VARCHAR(50)                     not null,
    descriptions            VARCHAR(225)                    not NULL,
    subscribers             INT                             NOT NULL,
    price                   VARCHAR(10)                     not null,
    url                     VARCHAR(50)                     NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)


);

CREATE TABLE IF NOT EXISTS subcriptions
(
    id                  serial                     NOT NULL,
    user_id             VARCHAR(50)                NOT NULL,
    video_id            INT                        NOT NULL,
    type_subscription   VARCHAR(20)                NOT NULL,
    subscribe_at        TIMESTAMP WITH TIME ZONE   NOT NULL,
    expired_at          TIMESTAMP WITH TIME ZONE   NOT NULL,
    status              VARCHAR(10)                NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (user_id),
    FOREIGN KEY (video_id) REFERENCES videos (id)

);

 

CREATE TABLE IF NOT EXISTS admin
(
    id           SERIAL                      NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);




INSERT INTO categories (id, name) VALUES
    (1,'keluarga'),
    (2,'sports'),
    (3,'animasi'),
    (4, 'thriller'),
    (5, 'dokumentasi')
    ;

INSERT INTO videos (category_id, name, descriptions, price, url, subscribers) VALUES 
    (1, 'cars', 'mobil mobilan', '10000', 'cdn/cdn/cars.mp4', 0),
    (2, 'play your ball', 'permainan bola', '10000', 'cdn/cdn/playyourball.mp4', 0),
    (5, 'dokument1', 'dokumentasi', '1000', 'cdn/cdn/dokumentasi.mp4', 0),
    (5, 'dokument2', 'dokumentasi2', '1000', 'cdn/cdn/dokumentasi.mp4', 0),
    (5, 'dokument3', 'dokumentasi3', '1000', 'cdn/cdn/dokumentasi.mp4', 0),
    (1, 'cars2', 'mobil mobilan', '10000', 'cdn/cdn/cars.mp4', 0),
    (2, 'play your ball2', 'permainan bola', '10000', 'cdn/cdn/playyourball.mp4', 0),
    (1, 'cars3', 'mobil mobilan', '10000', 'cdn/cdn/cars.mp4', 0),
    (2, 'play your ball3', 'permainan bola', '10000', 'cdn/cdn/playyourball.mp4', 0)
    ;

