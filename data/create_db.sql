CREATE USER IF NOT EXISTS 'spotifyviews'@'%' IDENTIFIED BY 'spotifyviews!';
GRANT ALL PRIVILEGES ON *.* TO 'spotifyviews'@'%';
FLUSH PRIVILEGES;

CREATE DATABASE IF NOT EXISTS spotify_views;
USE spotify_views;

CREATE TABLE IF NOT EXISTS users  (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    spotify_id VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    UNIQUE(spotify_id)
);

CREATE TABLE IF NOT EXISTS tokens (
    spotify_id VARCHAR(200) NOT NULL PRIMARY KEY,
    refresh VARCHAR(200) NOT NULL,
    UNIQUE(refresh)
);


CREATE DATABASE IF NOT EXISTS spotify_views_development;
USE spotify_views_development;


CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    spotify_id VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    UNIQUE(spotify_id)
);

CREATE TABLE IF NOT EXISTS tokens (
    spotify_id VARCHAR(200) NOT NULL PRIMARY KEY,
    refresh VARCHAR(200) NOT NULL,
    UNIQUE(refresh)
);

CREATE DATABASE IF NOT EXISTS spotify_views_test;
USE spotify_views_test;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    spotify_id VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    UNIQUE(spotify_id)
);

CREATE TABLE IF NOT EXISTS tokens (
    spotify_id VARCHAR(200) NOT NULL PRIMARY KEY,
    refresh VARCHAR(200) NOT NULL,
    UNIQUE(refresh)
);