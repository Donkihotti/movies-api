# movies-api — curl-testikomennot
# Serveri: http://localhost:8080   (käynnistys: go run .)
# Vinkki: -i näyttää statuskoodin + headerit, jq muotoilee JSONin.
#
# Odotetut statuskoodit (api/handler.go WriteStatus):
#   GET -> 200 | POST -> 201 | PATCH, DELETE -> 204 (ei bodyä)
#   virheet: 400 BadRequest, 404 NotFound, 409 Conflict, 500 ServerError

# =====================================================================
# MOVIES
# =====================================================================

# --- GET kaikki elokuvat ---
curl -i http://localhost:8080/api/movies

# sama, muotoiltuna
curl -s http://localhost:8080/api/movies | jq

# --- GET suodattimilla (query-parametrit) ---
# elokuvat joissa näyttelijä id=1 (Keanu Reeves)
curl -s "http://localhost:8080/api/movies?actor=1" | jq

# elokuvat genrellä id=4 (Sci-Fi)
curl -s "http://localhost:8080/api/movies?genre=4" | jq

# elokuvat julkaisuvuodelta 1999
curl -s "http://localhost:8080/api/movies?releaseYear=1999" | jq

# useampi suodatin yhtä aikaa
curl -s "http://localhost:8080/api/movies?genre=1&actor=1&releaseYear=1999" | jq

# HUOM: duration-suodatin on kommentoitu pois repository/movieRepo.go:ssa,
# eli tämä palauttaa kaikki elokuvat (ei suodata):
curl -s "http://localhost:8080/api/movies?duration=136" | jq

# virhe: releaseYear ei ole numero -> pitäisi olla 400
curl -i "http://localhost:8080/api/movies?releaseYear=abc"

# --- GET yksi elokuva id:llä ---
curl -i http://localhost:8080/api/movies/1

# virhe: id ei ole numero -> 400
curl -i http://localhost:8080/api/movies/abc

# virhe: id ei ole olemassa -> 404
curl -i http://localhost:8080/api/movies/99999

# --- GET elokuvan näyttelijät ---
curl -s http://localhost:8080/api/movies/1/actors | jq

# virhe: ei näyttelijöitä / ei elokuvaa -> 404
curl -i http://localhost:8080/api/movies/99999/actors

# virhe: epäkelpo id -> 400
curl -i http://localhost:8080/api/movies/abc/actors

# --- POST uusi elokuva -> 201 ---
curl -i -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Testielokuva",
    "description": "Testikuvaus",
    "release_date": "2024-01-15",
    "duration": "120",
    "genres": [1, 4],
    "actors": [1, 2]
  }'

# POST ilman genrejä ja näyttelijöitä (nekin ovat valinnaisia)
curl -i -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Pelkkä elokuva",
    "description": "Ei genrejä eikä näyttelijöitä",
    "release_date": "2020-05-05",
    "duration": "95"
  }'

# virhe: tyhjä title -> 400
curl -i -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "", "description": "x", "release_date": "2024-01-15", "duration": "120"}'

# virhe: puuttuva kenttä (description) -> 400
curl -i -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "Vajaa", "release_date": "2024-01-15", "duration": "120"}'

# virhe: väärä päivämäärämuoto (pitää olla YYYY-MM-DD) -> 400
curl -i -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "Vaara", "description": "x", "release_date": "15/01/2024", "duration": "120"}'

# virhe: tuntematon kenttä (DisallowUnknownFields) -> 400
curl -i -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "x", "description": "x", "release_date": "2024-01-15", "duration": "120", "foo": "bar"}'

# virhe: rikkinäinen JSON -> 400
curl -i -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "x",'

# virhe: olematon genre_id (foreign key) -> 500 tai 400
curl -i -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "FK-testi", "description": "x", "release_date": "2024-01-15", "duration": "120", "genres": [9999]}'

# --- PATCH elokuva -> 204 ---
# vain otsikko
curl -i -X PATCH http://localhost:8080/api/movies/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "The Matrix (päivitetty)"}'

# vain kesto
curl -i -X PATCH http://localhost:8080/api/movies/1 \
  -H "Content-Type: application/json" \
  -d '{"duration": "150"}'

# monta kenttää kerralla
curl -i -X PATCH http://localhost:8080/api/movies/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Uusi nimi", "description": "Uusi kuvaus", "release_date": "1999-04-01", "duration": "136"}'

# genret (HUOM: PATCHissa genret ovat objekteja, ei id-numeroita)
curl -i -X PATCH http://localhost:8080/api/movies/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Matrix", "genres": [{"id": 1, "genre_name": "Action"}]}'

# näyttelijät (id-lista)
curl -i -X PATCH http://localhost:8080/api/movies/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Matrix", "actors": [1, 2, 3]}'

# virhe: tyhjä body / ei yhtään päivitettävää kenttää -> 400
curl -i -X PATCH http://localhost:8080/api/movies/1 \
  -H "Content-Type: application/json" \
  -d '{}'

# virhe: väärä päivämäärämuoto -> 400
curl -i -X PATCH http://localhost:8080/api/movies/1 \
  -H "Content-Type: application/json" \
  -d '{"release_date": "1999"}'

# virhe: epäkelpo id -> 400
curl -i -X PATCH http://localhost:8080/api/movies/abc \
  -H "Content-Type: application/json" \
  -d '{"title": "x"}'

# virhe: olematon elokuva -> 404
curl -i -X PATCH http://localhost:8080/api/movies/99999 \
  -H "Content-Type: application/json" \
  -d '{"title": "x"}'

# --- DELETE elokuva -> 204 ---
curl -i -X DELETE http://localhost:8080/api/movies/20

# virhe: epäkelpo id -> 400
curl -i -X DELETE http://localhost:8080/api/movies/abc

# virhe: olematon elokuva -> 404
curl -i -X DELETE http://localhost:8080/api/movies/99999

# =====================================================================
# GENRES
# =====================================================================

# --- GET kaikki genret ---
curl -i http://localhost:8080/api/genres
curl -s http://localhost:8080/api/genres | jq

# --- GET yksi genre ---
curl -i http://localhost:8080/api/genres/1

# virhe: epäkelpo id -> 400 ("invalid id: abc")
curl -i http://localhost:8080/api/genres/abc

# virhe: olematon genre -> 404
curl -i http://localhost:8080/api/genres/99999

# --- POST uusi genre -> 201 ---
curl -i -X POST http://localhost:8080/api/genres \
  -H "Content-Type: application/json" \
  -d '{"genre_name": "Horror"}'

# virhe: tyhjä nimi -> 400
curl -i -X POST http://localhost:8080/api/genres \
  -H "Content-Type: application/json" \
  -d '{"genre_name": ""}'

# virhe: tuntematon kenttä -> 400
curl -i -X POST http://localhost:8080/api/genres \
  -H "Content-Type: application/json" \
  -d '{"name": "Horror"}'

# virhe: rikkinäinen JSON -> 400
curl -i -X POST http://localhost:8080/api/genres \
  -H "Content-Type: application/json" \
  -d '{"genre_name":'

# duplikaatti (aja sama kahdesti) -> 409 Conflict, jos UNIQUE-rajoite
curl -i -X POST http://localhost:8080/api/genres \
  -H "Content-Type: application/json" \
  -d '{"genre_name": "Action"}'

# --- PATCH genre -> 204 ---
curl -i -X PATCH http://localhost:8080/api/genres/1 \
  -H "Content-Type: application/json" \
  -d '{"genre_name": "Toiminta"}'

# virhe: tyhjä nimi -> 400
curl -i -X PATCH http://localhost:8080/api/genres/1 \
  -H "Content-Type: application/json" \
  -d '{"genre_name": ""}'

# virhe: epäkelpo id -> 400
curl -i -X PATCH http://localhost:8080/api/genres/abc \
  -H "Content-Type: application/json" \
  -d '{"genre_name": "X"}'

# virhe: olematon genre -> 404
curl -i -X PATCH http://localhost:8080/api/genres/99999 \
  -H "Content-Type: application/json" \
  -d '{"genre_name": "X"}'

# --- DELETE genre -> 204 ---
curl -i -X DELETE http://localhost:8080/api/genres/5

# virhe: epäkelpo id -> 400
curl -i -X DELETE http://localhost:8080/api/genres/abc

# virhe: olematon genre -> 404
curl -i -X DELETE http://localhost:8080/api/genres/99999

# =====================================================================
# ACTORS
# =====================================================================

# --- GET kaikki näyttelijät ---
curl -i http://localhost:8080/api/actors
curl -s http://localhost:8080/api/actors | jq

# --- GET näyttelijät nimisuodattimella (query) ---
curl -s "http://localhost:8080/api/actors?name=Keanu" | jq
curl -s "http://localhost:8080/api/actors?name=Tom%20Hanks" | jq

# --- GET yksi näyttelijä id:llä ---
curl -i http://localhost:8080/api/actors/1

# virhe: epäkelpo id -> 400 ("invalid actor id: abc")
curl -i http://localhost:8080/api/actors/abc

# virhe: olematon näyttelijä -> 404
curl -i http://localhost:8080/api/actors/99999

# --- GET näyttelijät nimellä (path-parametri) ---
curl -s http://localhost:8080/api/actors/name/Keanu%20Reeves | jq

# välilyönti pitää koodata itse muotoon %20:
curl -s http://localhost:8080/api/actors/name/Al%20Pacino | jq

# virhe: nimeä ei löydy -> 404
curl -i http://localhost:8080/api/actors/name/Eioo%20Olemassa

# --- GET näyttelijät syntymäpäivällä (muoto YYYY-MM-DD) ---
curl -s http://localhost:8080/api/actors/birthdate/1964-09-02 | jq

# virhe: väärä päivämäärämuoto -> 400 ("invalid birthdate: ...")
curl -i http://localhost:8080/api/actors/birthdate/02-09-1964

# virhe: ei osumia -> 404
curl -i http://localhost:8080/api/actors/birthdate/1900-01-01

# --- POST uusi näyttelijä -> 201 ---
curl -i -X POST http://localhost:8080/api/actors \
  -H "Content-Type: application/json" \
  -d '{"name": "Testi Näyttelijä", "birth_date": "1990-01-01"}'

# virhe: tyhjä nimi -> 400
curl -i -X POST http://localhost:8080/api/actors \
  -H "Content-Type: application/json" \
  -d '{"name": "", "birth_date": "1990-01-01"}'

# virhe: puuttuva birth_date -> 400
curl -i -X POST http://localhost:8080/api/actors \
  -H "Content-Type: application/json" \
  -d '{"name": "Ei Päivää"}'

# virhe: väärä päivämäärämuoto -> 400
curl -i -X POST http://localhost:8080/api/actors \
  -H "Content-Type: application/json" \
  -d '{"name": "Vaara", "birth_date": "01/01/1990"}'

# virhe: tuntematon kenttä -> 400
curl -i -X POST http://localhost:8080/api/actors \
  -H "Content-Type: application/json" \
  -d '{"name": "X", "birth_date": "1990-01-01", "ikä": 30}'

# --- PATCH näyttelijä -> 204 ---
# vain nimi
curl -i -X PATCH http://localhost:8080/api/actors/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Keanu Charles Reeves"}'

# vain syntymäpäivä
curl -i -X PATCH http://localhost:8080/api/actors/1 \
  -H "Content-Type: application/json" \
  -d '{"birth_date": "1964-09-03"}'

# molemmat
curl -i -X PATCH http://localhost:8080/api/actors/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Keanu Reeves", "birth_date": "1964-09-02"}'

# virhe: tyhjä body -> 400
curl -i -X PATCH http://localhost:8080/api/actors/1 \
  -H "Content-Type: application/json" \
  -d '{}'

# virhe: epäkelpo id -> 400
curl -i -X PATCH http://localhost:8080/api/actors/abc \
  -H "Content-Type: application/json" \
  -d '{"name": "X"}'

# virhe: olematon näyttelijä -> 404
curl -i -X PATCH http://localhost:8080/api/actors/99999 \
  -H "Content-Type: application/json" \
  -d '{"name": "X"}'

# --- DELETE näyttelijä -> 204 ---
curl -i -X DELETE http://localhost:8080/api/actors/15

# virhe: epäkelpo id -> 400
curl -i -X DELETE http://localhost:8080/api/actors/abc

# virhe: olematon näyttelijä -> 404
curl -i -X DELETE http://localhost:8080/api/actors/99999

# =====================================================================
# REITITYS / METODIT (ServeMux)
# =====================================================================

# olematon polku -> 404
curl -i http://localhost:8080/api/foobar

# väärä metodi olemassa olevalle polulle -> 405 Method Not Allowed
curl -i -X PUT http://localhost:8080/api/movies/1
curl -i -X POST http://localhost:8080/api/movies/1
curl -i -X DELETE http://localhost:8080/api/movies
