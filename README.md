fls-core
========

This microservice acts as a proxy for the Bandsintown API maintaining information for a subset of artists. The REST API serves information on artists and their events, as well as accepts new artists to query Bandsintown for.

## Endpoints

### `GET /artists`

Returns list of tracked artists

### `POST /artists`

Adds artists to list of tracked artists

Parameters:
* artists *(POST body)* - A comma-separated list of artists

### `GET /artistinfo/{artistname}`

Returns information about an artist

Parameters:
* artistname *(query)* - The name of the artist

### `GET /events/{artistname}`

Returns upcoming events for an artist

Parameters:
* artistname *(query)* - The name of the artist
