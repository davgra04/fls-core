package fls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// bandsintown objects
////////////////////////////////////////////////////////////////////////////////

// BandsInTownEvent represents an event for an artist
type BandsInTownEvent struct {
	ID             string            `db:"id" json:"id"`
	ArtistID       string            `db:"artist_id" json:"artist_id"`
	URL            string            `db:"url" json:"url"`
	OnSaleDatetime string            `db:"on_sale_datetime" json:"on_sale_datetime"` // 2017-03-01T18:00:00
	Datetime       string            `db:"datetime" json:"datetime"`
	Description    string            `db:"description" json:"description"`
	Venue          *BandsInTownVenue `db:"venue" json:"venue"`
	Lineup         []string          `db:"lineup" json:"lineup"`
	// Offers         []BandsInTownOffer `db:"offers" json:"offers"`
	Added   time.Time `db:"added" json:"added"`
	Updated time.Time `db:"updated" json:"updated"`
	Removed time.Time `db:"removed" json:"removed"`
}

// BandsInTownVenue represents a venue for an event
type BandsInTownVenue struct {
	Name      string    `db:"name" json:"name"`
	Latitude  string    `db:"latitude" json:"latitude"`
	Longitude string    `db:"longitude" json:"longitude"`
	City      string    `db:"city" json:"city"`
	Region    string    `db:"region" json:"region"`
	Country   string    `db:"country" json:"country"`
	Updated   time.Time `db:"updated" json:"updated"`
}

// // BandsInTownOffer TODO TODO TODO
// type BandsInTownOffer struct {
// 	Type   string `db:"type" json:"type"`
// 	URL    string `db:"url" json:"url"`
// 	Status string `db:"status" json:"status"`
// }

// BandsInTownArtist represents an artist
type BandsInTownArtist struct {
	ID              int64     `db:"id" json:"-"`
	Name            string    `db:"name" json:"name"`
	BITID           string    `db:"bit_id" json:"id"`
	URL             string    `db:"url" json:"url"`
	ImageURL        string    `db:"image_url" json:"image_url"`
	ThumbURL        string    `db:"thumb_url" json:"thumb_url"`
	FacebookPageURL string    `db:"facebook_page_url" json:"facebook_page_url"`
	MBID            string    `db:"mbid" json:"mbid"`
	Updated         time.Time `db:"updated" json:"updated"`
}

////////////////////////////////////////////////////////////////////////////////
// bandsintown query functions
////////////////////////////////////////////////////////////////////////////////

// BandsInTownClient is the http client used for querying bandsintown
var BandsInTownClient *http.Client

func escapeArtistName(artist string) string {
	artist = strings.ReplaceAll(artist, "/", "%252F")
	artist = strings.ReplaceAll(artist, "?", "%253F")
	artist = strings.ReplaceAll(artist, "*", "%252A")
	artist = strings.ReplaceAll(artist, "\"", "%27C")
	return artist
}

// QueryBITArtist queries BandsInTown for and returns artist data
func QueryBITArtist(artistName string) (*BandsInTownArtist, error) {
	// prepare request
	url := fmt.Sprintf("https://rest.bandsintown.com/artists/%s?app_id=%s", escapeArtistName(artistName), APIKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request object: %v", err)
	}

	// execute request
	res, err := BandsInTownClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %v", err)
	}

	// read response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}

	// parse json
	var artist BandsInTownArtist
	err = json.Unmarshal(body, &artist)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON in response: %v", err)
	}

	// set updated time
	artist.Updated = time.Now().UTC()

	return &artist, nil
}

// QueryBITEvents queries BandsInTown for and returns events data for an artist
func QueryBITEvents(artistName string) *BandsInTownArtist {
	return nil
}
