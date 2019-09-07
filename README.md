fls-core
========

This microservice is responsible for querying the bandsintown API, maintaining a list of shows, and serving the data through a REST API.


# JSON Objects

Show data is structured in a few different ways:

## Show Data from BandsInTown API

### Requesting shows for an artist

```json
{
    "data": "???"
}
```

## Internal Representation

### Primary Data Store

This is where fls-core stores configuration and raw data from BandsInTown.

```json
{
    "artists": [ "...list of followed artists..."],
    "config": { "???": "various other config options can be added here"},
    "query_date": 1567827883,
    "raw_data": {
        # bandsintown raw data
    }
}
```

### Regional Cache Files

These files serve as a cache for data within a region. This prevents the page data from having to be generated on every request.

The format is identical to the `/v1/shows/TX/` endpoint response below

## fls-core API Endpoints

### Getting main show list for a region

Served on `/v1/shows/TX/`

```json
{
    "query_date": 1567827883,
    "region": "TX",
    "shows": [
        {
            "show_id": 0,
            "artist": "Tame Impala",
            "date": "Sat 8/31 7:00p",
            "date_added": 1565826883,
            "venue": "White Oak Music Hall",
            "lineup": "White Denim",
            "city": "Houston",
            "region": "TX"
        },
        {
            "show_id": 1,
            "artist": "King Gizzard and the Lizard Wizard",
            "date": "Sun 9/01 7:30p",
            "date_added": 1567826883,
            "venue": "White Oak Music Hall",
            "lineup": "Mildlife, Orb",
            "city": "Houston",
            "region": "TX"
        },
        {
            "show_id": 2,
            "artist": "Unknown Mortal Orchestra",
            "date": "Mon 9/02 8:00p",
            "date_added": 1567826883,
            "venue": "White Oak Music Hall",
            "lineup": "Shakey Graves",
            "city": "Houston",
            "region": "TX"
        },
        . . .
    ]
}
```


### Getting list of followed artists

Served on `/v1/artists/`

```json
{
    artists: [
        "Tame Impala",
        "King Gizzard and the Lizard Wizard",
        "Unknown Mortal Orchestra",
        . . .
    ]
}

```





