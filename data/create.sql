CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    spotify_id VARCHAR(200) NOT NULL,
    email NVARCHAR(200) NOT NULL,
    UNIQUE KEY (spotify_id)
);

CREATE TABLE IF NOT EXISTS tokens (
    spotify_id VARCHAR(200) NOT NULL PRIMARY KEY,
    refresh VARCHAR(200) NOT NULL,
    UNIQUE KEY (refresh)
);