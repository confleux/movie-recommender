CREATE TABLE IF NOT EXISTS movies_info (
    id INT PRIMARY KEY,
    title VARCHAR NOT NULL,
    release_date timestamp NOT NULL,
    vote_average FLOAT NOT NULL,
    vote_count INT NOT NULL,
    is_adult BOOL NOT NULL,
    poster_path VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS movies_genre_map (
    movie_id INT REFERENCES movies_info(id) NOT NULL,
    genre_id INT NOT NULL
);
