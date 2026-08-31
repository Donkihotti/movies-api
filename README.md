# movies-api

A REST API for managing a movie database, built for a local film society that
outgrew its spreadsheets. The API stores movies, genres and actors, and the
many-to-many relationships between them.

Written in Go with the standard library only (`net/http`, `database/sql`) and
raw SQL against SQLite. 

## Project structure

```
backend/
├── main.go                  wiring: db, migrations, repos, services, handlers, server
├── api/                     HTTP layer — handlers, routes, status/error writers
├── internal/                repository + service layer, one file per entity
├── models/                  request/response structs
├── errs/                    custom error variables
└── database/
    ├── database.go          connection
    ├── migrations.go        migration runner
    └── migrations/          numbered .sql files, incl. seed data
```

The code is organised in three layers. Handlers parse and validate the HTTP
request and write the response. Services hold business logic and input
validation. Repositories own the SQL and are the only layer that touches the
database. Errors travel upward as the custom error values in `errs/`, and the
handler layer translates them into HTTP status codes.

## Setup

Requires Go 1.25 or later. SQLite itself needs no installation — the driver
embeds it, and the database file is created on first run.

```
git clone <repo-url>
cd movies-api/backend
go mod download
go run .
```
Make sure you run the program from the directory you just navigated to.

The server listens on `http://localhost:8080`.

### Database and migrations

On startup the app connects to `./movies.db`, creating the file if it does not
exist, then runs every `.sql` file in `database/migrations/` in filename order.
All schema statements use `CREATE TABLE IF NOT EXISTS` and all seed statements
use `INSERT OR IGNORE`, so a normal restart is a no-op.

`movies.db` is gitignored and is not part of the repository.

**After changing a migration, delete the database and let it rebuild:**

```
rm movies.db
go run .
```

`CREATE TABLE IF NOT EXISTS` will not add a new column to a table that already
exists, so an edited schema does not reach an existing database file.

### Sample data

`006_seed.sql` populates the database with 5 genres, 20 movies spanning
1993–2019, and 15 actors. Movies have one or more genres, actors appear in one
or more movies, and the number of actors per movie varies.


### Movies

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/movies` | Create a movie |
| `GET` | `/api/movies` | List all movies |
| `GET` | `/api/movies/{id}` | Get one movie |
| `PATCH` | `/api/movies/{id}` | Partially update a movie |
| `DELETE` | `/api/movies/{id}` | Delete a movie |
| `GET` | `/api/movies/{id}/actors` | List the actors in a movie |

Filters on `GET /api/movies`, combinable:

| Query | Example |
|---|---|
| `genre` | `/api/movies?genre=1` |
| `actor` | `/api/movies?actor=6` |
| `releaseYear` | `/api/movies?releaseYear=1999` |

### Genres

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/genres` | Create a genre |
| `GET` | `/api/genres` | List all genres |
| `GET` | `/api/genres/{id}` | List the movies in a genre |
| `PATCH` | `/api/genres/{id}` | Rename a genre |
| `DELETE` | `/api/genres/{id}` | Delete a genre |

### Actors

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/actors` | Create an actor |
| `GET` | `/api/actors` | List all actors |
| `GET` | `/api/actors/{id}` | Get one actor |
| `PATCH` | `/api/actors/{id}` | Partially update an actor |
| `DELETE` | `/api/actors/{id}` | Delete an actor |
| `GET` | `/api/actors/name/{name}` | Search actors by name |
| `GET` | `/api/actors/birthdate/{date}` | Filter actors by birth date |

`GET /api/actors?name={name}` does the same case-insensitive partial match as
the `/name/{name}` path variant.

## Request bodies

Create a movie. `genres` and `actors` are arrays of existing IDs and may be
omitted:

```json
{
  "title": "Blade Runner",
  "description": "A blade runner hunts replicants",
  "release_date": "1982-06-25",
  "duration": "117",
  "genres": [1, 4],
  "actors": [3]
}
```

Update a movie. Every field is optional; send only what changes:

```json
{ "duration": "121" }
```

Create a genre / create an actor:

```json
{ "genre_name": "Western" }
```

```json
{ "name": "Sigourney Weaver", "birth_date": "1949-10-08" }
```

Unknown JSON fields are rejected with `400`.

## Validation

- `POST /api/movies` requires `title`, `description`, `release_date` and `duration`
- `POST /api/actors` requires `name` and `birth_date`
- Dates must be ISO 8601 (`YYYY-MM-DD`)
- `PATCH` requires at least one field in the body
- `?releaseYear=` must be numeric
- `id` is assigned by the database and cannot be changed after creation

## Status codes

| Code | When |
|---|---|
| `200 OK` | Successful `GET` |
| `201 Created` | Successful `POST` |
| `204 No Content` | Successful `PATCH` or `DELETE` |
| `400 Bad Request` | Validation failure, malformed JSON, non-numeric ID |
| `404 Not Found` | No such entity |
| `409 Conflict` | Conflicting state |
| `500 Internal Server Error` | Database or server failure |

## Example requests

```
curl http://localhost:8080/api/movies

curl http://localhost:8080/api/movies/3

curl "http://localhost:8080/api/movies?genre=4"

curl "http://localhost:8080/api/movies?releaseYear=1999"

curl http://localhost:8080/api/movies/1/actors

curl "http://localhost:8080/api/actors?name=hanks"

curl -X POST http://localhost:8080/api/genres \
  -H "Content-Type: application/json" \
  -d '{"genre_name":"Western"}'

curl -X PATCH http://localhost:8080/api/movies/3 \
  -H "Content-Type: application/json" \
  -d '{"duration":"148"}'

curl -X DELETE http://localhost:8080/api/movies/3
```

## Design decisions

**`release_date` instead of `releaseYear`.** The brief lists `releaseYear`, but
a full ISO date carries strictly more information and is what the film society
actually has on record. The `?releaseYear=` filter matches against the year
portion, so the documented filtering behaviour is unchanged.

**A `description` field on movies.** Not in the brief, but the society's
spreadsheets carry one-line summaries and dropping them on import would lose
data.

**Errors as package-level variables.** `errs` defines `NotFound`,
`BadRequest`, `ServerError` and `Conflict`. Repositories and services return
these rather than driver errors, so the SQL layer's internals never leak into
an HTTP response, and the handler layer has one place that maps an error to a
status code.

**Migrations re-run on every start.** The runner replays all `.sql` files each
time rather than tracking applied versions. Every statement is idempotent, so
this is safe, and it keeps the setup to a single `go run .`.

## Known limitations

Work in progress, listed honestly rather than left to be discovered:

- **Force delete is not implemented.** `DELETE` on an entity with existing
  relationships should return `400` with an explanatory message, and
  `?force=true` should remove the relationships first. Currently every delete
  succeeds and `ON DELETE CASCADE` removes the junction rows silently.
- **Error responses have no body.** The status code is correct but the response
  body is empty. A wrapper type carrying a message alongside the error value is
  designed but not yet built.
- **`GET /api/genres/{id}` returns the genre's movies, not the genre itself.**
- **`PATCH /api/movies/{id}` ignores `genres` and `actors`.** Only the scalar
  fields are updated. The same applies to an actor's associated movies.
- **`GET` responses omit relationships.** A movie is returned without its
  genres and actors; use `/api/movies/{id}/actors` and `/api/movies?genre=` in
  the meantime.
- **`GET /api/movies/{id}/actors` returns `404` for a movie with no actors**,
  where `200` with an empty array would be correct.
- **A foreign key violation returns `500`** rather than `400` — posting a movie
  with a non-existent genre ID is a client error.
- **Pagination and title search** (the optional requirements) are not
  implemented.
