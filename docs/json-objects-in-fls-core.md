JSON Objects in fls-core
========================

Show data is structured in a few different ways:

## Show Data from BandsInTown API

### Requesting shows for an artist

```json
[
  {
    "id": "1016966641",
    "artist_id": "2228815",
    "url": "https://www.bandsintown.com/e/101696...",
    "on_sale_datetime": "2019-09-04T10:00:00",
    "datetime": "2019-09-11T20:00:00",
    "description": "Mac DeMarco",
    "venue": {
      "country": "United States",
      "city": "Los Angeles",
      "latitude": "34.0522222",
      "name": "John Anson Ford Theatres",
      "region": "CA",
      "longitude": "-118.2427778"
    },
    "lineup": [
      "Mac DeMarco"
    ],
    "offers": [
      {
        "type": "Tickets",
        "url": "https://www.bandsintown.com/t/101696...",
        "status": "available"
      }
    ]
  },
  // {...Event...},
  // {...Event...},
  // {...Event...},
]

```

## Internal Representation

### Primary Data Store

This is where fls-core stores configuration and raw data from BandsInTown.

```json
{
    "artists": [ "...list of followed artists..."],
    "config": { "???": "various other config options can be added here"},
    "query_date": 1567827883,
    "data": {
        "artist_name": {
            "bandsintown_events": [
                // raw event data
            ],
            "bandsintown_artist_info": {
                // raw artist info data
            },
        }
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
        ...
    ]
}
```


### Getting list of followed artists

Served on `/v1/artists/`

```json
{
    "artists": [
        "Tame Impala",
        "King Gizzard and the Lizard Wizard",
        "Unknown Mortal Orchestra",
        ...
    ]
}

```


### Getting details for a single show

Served on `/v1/show_details/?id=123`

# TODO

