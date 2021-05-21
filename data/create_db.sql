CREATE USER 'spotifyviews'@'localhost' IDENTIFIED BY 'spotifyviews!';
GRANT ALL PRIVILEGES ON *.* TO 'spotifyviews'@'localhost';
FLUSH PRIVILEGES;

CREATE DATABASE spotify_views;
USE spotify_views;

CREATE TABLE users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    spotify_id VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    UNIQUE(spotify_id)
);

CREATE TABLE tokens (
    spotify_id VARCHAR(200) NOT NULL PRIMARY KEY,
    refresh VARCHAR(200) NOT NULL,
    UNIQUE(refresh)
);


CREATE DATABASE spotify_views_development;
USE spotify_views_development;


CREATE TABLE users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    spotify_id VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    UNIQUE(spotify_id)
);

CREATE TABLE tokens (
    spotify_id VARCHAR(200) NOT NULL PRIMARY KEY,
    refresh VARCHAR(200) NOT NULL,
    UNIQUE(refresh)
);

CREATE DATABASE spotify_views_test;
USE spotify_views_test;

CREATE TABLE users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    spotify_id VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    UNIQUE(spotify_id)
);

CREATE TABLE tokens (
    spotify_id VARCHAR(200) NOT NULL PRIMARY KEY,
    refresh VARCHAR(200) NOT NULL,
    UNIQUE(refresh)
);