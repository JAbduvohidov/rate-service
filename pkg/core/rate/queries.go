package rate

const getAllRatesDML = `SELECT movie_id,
       title,
       description,
       image,
       year,
       country,
       actors,
       genres,
       creators,
       studio,
       extlink,
       avg(rate)
FROM rates
WHERE removed = FALSE GROUP BY movie_id, title, description, image, year, country, actors, genres, creators, studio, extlink;`

const getRateDML = `SELECT movie_id,
       title,
       description,
       image,
       year,
       country,
       actors,
       genres,
       creators,
       studio,
       extlink,
       avg(rate)
FROM rates
WHERE removed = FALSE AND movie_id = $1 GROUP BY movie_id, title, description, image, year, country, actors, genres, creators, studio, extlink;`

const addRateDML = `INSERT INTO rates (movie_id, title, description, image, year, country, actors, genres, creators, studio, extlink, user_id, user_name, rate)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

const deleteRateDML = `UPDATE rates SET removed = TRUE WHERE movie_id = $1 AND user_id = $2;`
