Notes on querying the BandsInTown API
=====================================

## Links

[BandsInTown API Reference](https://app.swaggerhub.com/apis-docs/Bandsintown/PublicAPI/3.0.0)

## BandsInTown API Key

I'm storing my API key in the bandsintown-api.secret file, which can be sourced:

```bash
$ source bandsintown-api.secret
$ echo $BANDSINTOWN_API_KEY
**********************************
```



## Calling the API

```bash
$ source bandsintown-api.secret 

# get artist information
$ curl https://rest.bandsintown.com/artists/King%20Gizzard%20and%20the%20Lizard%20Wizard?app_id=${BANDSINTOWN_API_KEY}

{"id": "2117024", "name": "King Gizzard & The Lizard Wizard", "url": "https://www.bandsintown.com/a/2117024?came_from=267&app_id=190935ba5aa9f01ff41f77c802cb0d60", "image_url": "https://s3.amazonaws.com/bit-photos/large/9288487.jpeg", "thumb_url": "https://s3.amazonaws.com/bit-photos/thumb/9288487.jpeg", "facebook_page_url": "http://www.facebook.com/168329496513295", "mbid": "", "tracker_count": 138454, "upcoming_event_count": 10}

# get list of shows for King Gizzard and the Lizard Wizard
$ curl https://rest.bandsintown.com/artists/King%20Gizzard%20and%20the%20Lizard%20Wizard/events?app_id=${BANDSINTOWN_API_KEY}


[{"id":"1014167875","artist_id":"2117024","url":"https:\/\/www.bandsintown.com\/e\/1014167875?app_id=190935ba5aa9f01ff41f77c802cb0d60&came_from=267&utm_medium=api&utm_source=public_api&utm_campaign=event","on_sale_datetime":"2019-03-15T10:00:00","datetime":"2019-09-30T19:00:00","description":"","venue":{"country":"United Kingdom","city":"Nottingham","latitude":"52.9666667","name":"Rock City","region":"","longitude":"-1.1666667"},"lineup":["King Gizzard & The Lizard Wizard"],"offers":[{"type":"Tickets","url":"https:\/\/www.bandsintown.com\/t\/1014167875?app_id=190935ba5aa9f01ff41f77c802cb0d60&came_from=267&utm_medium=api&utm_source=public_api&utm_campaign=ticket","status":"available"}]},{"id":"1014180868", # . . .

```


















