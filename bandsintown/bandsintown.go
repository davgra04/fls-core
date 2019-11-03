package bandsintown

// bandsintown objects
////////////////////////////////////////////////////////////////////////////////

// BandsInTownEventData TODO TODO TODO
type BandsInTownEventData struct {
	// fields populated by BandsInTown
	ID             string                 `json:"id"`
	ArtistID       string                 `json:"artist_id"`
	URL            string                 `json:"url"`
	OnSaleDatetime string                 `json:"on_sale_datetime"` // 2017-03-01T18:00:00
	Datetime       string                 `json:"datetime"`
	Description    string                 `json:"description"`
	Venue          *BandsInTownVenueData  `json:"venue"`
	Offers         []BandsInTownOfferData `json:"offers"`
	Lineup         []string               `json:"lineup"`
	// fields populated by fls-data
	DateAdded   int64 `json:"date_added"`
	DateUpdated int64 `json:"date_updated"`
	DateRemoved int64 `json:"date_removed"`
}

// BandsInTownVenueData TODO TODO TODO
type BandsInTownVenueData struct {
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	City      string `json:"city"`
	Region    string `json:"region"`
	Country   string `json:"country"`
}

// BandsInTownOfferData TODO TODO TODO
type BandsInTownOfferData struct {
	Type   string `json:"type"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

// BandsInTownArtistData TODO TODO TODO
type BandsInTownArtistData struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	ImageURL        string `json:"image_url"`
	ThumbURL        string `json:"thumb_url"`
	FacebookPageURL string `json:"facebook_page_url"`
	MBID            string `json:"mbid"`
	TrackerCount    int    `json:"tracker_count"`
}

// BandsInTownData represents the full, raw picture from BandsInTown of an artist and their events
type BandsInTownData struct {
	QueryDate int64                  `json:"query_date"` // UNIX timestamp for time of last bandsintown API call
	Artist    BandsInTownArtistData  `json:"artist"`
	Events    []BandsInTownEventData `json:"events"`
}
