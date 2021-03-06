package fls

import (
	"strings"
	"time"
)

func stringArrToMap(arr []string) map[string]int {
	m := map[string]int{}
	for _, s := range arr {
		if _, ok := m[s]; !ok {
			m[s] = 1
		}
	}
	return m
}

func getArtistList() ([]string, error) {
	return DAO.getArtistList()
}

func getAliasList() ([]string, error) {
	return DAO.getAliasList()
}

func addArtists(artists []string) error {

	// get existing artist and alias map to check against
	existingArtistsList, err := getArtistList()
	if err != nil {
		return err
	}
	existingArtists := stringArrToMap(existingArtistsList)

	existingAliasesList, err := getAliasList()
	if err != nil {
		return err
	}
	existingAliases := stringArrToMap(existingAliasesList)

	// determine which artists are new
	//   this first pass minimizes bandsintown queries, but will need to check
	//   again once the real name is determined
	var newAliases []string
	for _, aName := range artists {
		aNameLower := strings.ToLower(aName)
		if _, ok := existingAliases[aNameLower]; !ok {
			newAliases = append(newAliases, aName)
		}
	}

	// Info.Printf("DEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUG\n")
	// Info.Printf("existingAliases: %#v\n", existingAliases)
	// Info.Printf("newAliases: %#v\n", newAliases)
	// Info.Printf("DEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUGDEBUG\n")

	if len(newAliases) == 0 {
		Info.Printf("No new artists to add\n")
		return nil
	}

	// obtain artist data and add aliases as necessary
	artistsAdded := []string{}
	aliasesAdded := []string{}
	limiter := time.Tick(Cfg.RPSPeriod)

	for _, aName := range newAliases {
		<-limiter
		aNameLower := strings.ToLower(aName)

		a, err := QueryBITArtist(aNameLower)
		if err != nil {
			Error.Printf("failed to query artist %v: %v\n", a, err)
			continue
		}
		// Info.Printf("a: %#v\n", a)

		if _, ok := existingArtists[a.Name]; !ok {
			// truly haven't seen this artist before, add to db
			err := DAO.addArtist(a)
			if err != nil {
				Error.Printf("failed to add artist data: %v\n", err)
			}

			DAO.addAlias(strings.ToLower(a.Name), a.BITID)
			existingAliases[strings.ToLower(a.Name)] = 1
			artistsAdded = append(artistsAdded, a.Name)
			aliasesAdded = append(aliasesAdded, strings.ToLower(a.Name))

		}

		if _, ok := existingAliases[aNameLower]; !ok {
			// havent seen this alias before, add to db
			DAO.addAlias(aNameLower, a.BITID)
			existingAliases[aNameLower] = 1
			aliasesAdded = append(aliasesAdded, aNameLower)
		}

	}

	Info.Printf("    Added %d artists: %v\n", len(artistsAdded), strings.Join(artistsAdded, ", "))
	Info.Printf("    Added %d aliases: %v\n", len(aliasesAdded), strings.Join(aliasesAdded, ", "))

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// fls objects
////////////////////////////////////////////////////////////////////////////////

// // FLSData represents all of the non-cache data in fls-core
// type FLSData struct {
// 	BandsInTownData map[string]*BandsInTown `json:"bandsintown_data"` // maps artist_name -> bandsintown artist info
// 	ShowIDsSeen     map[string]int64        // keeps track of which ShowIDs have been seen before, maps ShowID -> time first seen (UNIX timestamp)
// }

// // FLSEventData represents an artist's show data for a particular region
// type FLSEventData struct {
// 	ShowID    string `json:"show_id"`
// 	Artist    string `json:"artist"`
// 	Date      string `json:"date"`
// 	DateAdded int64  `json:"date_added"`
// 	Venue     string `json:"venue"`
// 	Lineup    string `json:"lineup"`
// 	City      string `json:"city"`
// 	Region    string `json:"region"`
// }

// // GetShowsResponse represents an artist's show data for a particular region, used by RouteGetShows endpoint
// type GetShowsResponse struct {
// 	QueryDate int64          `json:"query_date"`
// 	Region    string         `json:"region"`
// 	Shows     []FLSEventData `json:"shows"`
// }

////////////////////////////////////////////////////////////////////////////////
// fls functions
////////////////////////////////////////////////////////////////////////////////

// // GetCachedShowsResponse TODO TODO TODO
// func GetCachedShowsResponse(region string) string {

// 	// try to read cached json value, generate if none is there
// 	// TODO

// 	// just generating every time for now...
// 	flsdata := ReadFLSData("data.json")
// 	Info.Printf("flsdata: %v", flsdata)

// 	showData := GetShowsResponse{QueryDate: time.Now().Unix(), Region: region}

// 	filterEanbled := true
// 	if region == "all" {
// 		filterEanbled = false
// 	}

// 	for _, artist := range Cfg.Artists {
// 		for _, event := range flsdata.BandsInTownData[artist].Events {
// 			// Info.Printf("venue: %v", event.Venue)
// 			// Info.Printf("region: %v", event.Venue.Region)

// 			if (event.Venue.Region == region) || !filterEanbled {
// 				showData.Shows = append(showData.Shows, FLSEventData{
// 					ShowID:    event.ID,
// 					Artist:    artist,
// 					Date:      event.Datetime,
// 					DateAdded: event.DateAdded,
// 					Venue:     event.Venue.Name,
// 					Lineup:    strings.Join(event.Lineup, ", "),
// 					City:      event.Venue.City,
// 					Region:    event.Venue.Region,
// 				})
// 				// Info.Printf("Found a %v show: %v", region, showData.Shows[len(showData.Shows)-1])
// 			}
// 		}
// 	}

// 	// sort shows by date
// 	sort.SliceStable(showData.Shows, func(i, j int) bool {
// 		return showData.Shows[i].Date < showData.Shows[j].Date
// 	})

// 	// marshal to json
// 	showDataJSON, err := json.Marshal(showData)
// 	if err != nil {
// 		Error.Printf("Failed to marshal show data: %v", err)
// 	}

// 	return string(showDataJSON)

// }

// // GetArtistsResponse represents the data returned by the RouteGetArtists endpoint
// type GetArtistsResponse struct {
// }

////////////////////////////////////////////////////////////////////////////////
// bandsintown query goroutine
////////////////////////////////////////////////////////////////////////////////

// // ArtistRequest TODO TODO TODO
// type ArtistRequest struct {
// 	artist, url string
// }

// // PollBandsInTown periodically polls BandsInTown for show data, saves to FLSData, and initiates cache rebuild goroutine
// func PollBandsInTown() {
// 	apiKey := os.Getenv("BANDSINTOWN_API_KEY")

// 	// channel for limiting requests
// 	limiterDuration := time.Duration(int64(Cfg.RateLimitMillis) * time.Millisecond.Nanoseconds())
// 	limiter := time.Tick(limiterDuration)
// 	Info.Printf("limiter duration: %v [%T]", limiterDuration, limiterDuration)

// 	bandsInTownClient := http.Client{
// 		Timeout: time.Second * 10,
// 	}

// 	// TODO: load from file
// 	// flsdata := FLSData{}

// 	// initialize fresh flsdata

// 	// flsdata := FLSData{BandsInTownData: make(map[string]*BandsInTownData)}
// 	flsdata := ReadFLSData("data.json")

// 	// flsdata := make(map[string]*BandsInTownData)
// 	// for _, artist := range Cfg.Artists {
// 	// 	flsdata[artist] = &BandsInTownData{}
// 	// 	// flsdata[artist] = make(map[string]T)
// 	// }

// 	for {

// 		// set times for this polling period
// 		startTime := time.Now()
// 		nextPollTime := startTime.Add(time.Duration(Cfg.RefreshPeriodSeconds) * time.Second)

// 		// get data from api
// 		//////////////////////

// 		// load requests into requests channel
// 		requests := make(chan ArtistRequest, 1000)
// 		for _, artist := range Cfg.Artists {
// 			url := fmt.Sprintf("https://rest.bandsintown.com/artists/%s/events?app_id=%s", url.PathEscape(artist), apiKey)
// 			requests <- ArtistRequest{url: url, artist: artist}
// 			// Info.Printf("    preparing request for %-40v [%v]", artist, url)
// 		}
// 		close(requests)

// 		// make requests using limiter
// 		for ar := range requests {
// 			<-limiter // wait for limiter

// 			// build and make request
// 			Info.Printf("    requesting %-40v [%v]", ar.artist, ar.url)
// 			req, err := http.NewRequest(http.MethodGet, ar.url, nil)
// 			if err != nil {
// 				Error.Printf("Could not create request object: %v", err)
// 				continue
// 			}

// 			res, err := bandsInTownClient.Do(req)
// 			if err != nil {
// 				Error.Printf("Could not make request: %v", err)
// 				continue
// 			}

// 			// read response
// 			body, err := ioutil.ReadAll(res.Body)
// 			if err != nil {
// 				Error.Printf("Could not read response body: %v", err)
// 				continue
// 			}

// 			// Info.Printf("    response: %v", string(body))

// 			// parse response
// 			var events []BandsInTownEventData
// 			err = json.Unmarshal(body, &events)
// 			if err != nil {
// 				Error.Printf("Could not parse JSON in response: %v", err)
// 			}

// 			// check for new ShowIDs
// 			for eIdx, e := range events {
// 				if _, ok := flsdata.ShowIDsSeen[e.ID]; !ok {
// 					// haven't seen this ShowID before, set DateAdded and add to ShowIDsSeen
// 					Info.Printf("        New show! (ShowID: %v)", e.ID)
// 					flsdata.ShowIDsSeen[e.ID] = startTime.Unix()
// 					events[eIdx].DateAdded = startTime.Unix()
// 				} else {
// 					events[eIdx].DateAdded = flsdata.ShowIDsSeen[e.ID]
// 				}
// 			}

// 			// write events to flsdata
// 			if _, ok := flsdata.BandsInTownData[ar.artist]; !ok {
// 				// initialize BandsInTownData object if not in map
// 				flsdata.BandsInTownData[ar.artist] = &BandsInTownData{QueryDate: time.Now().Unix()}
// 			}

// 			flsdata.BandsInTownData[ar.artist].Events = events

// 		}

// 		// save data to data.json
// 		///////////////////////////
// 		WriteFLSData("data.json", flsdata)

// 		// TODO: update the event id index to track a number of things about shows:
// 		//     * which shows are new
// 		//     * which shows have disappeared from bandsintown
// 		//     * which shows have changed information

// 		// flsdata.BandsInTownData.Events = ""

// 		// trigger cache rebuild goroutine
// 		////////////////////////////////////

// 		// sleep until next query
// 		///////////////////////////
// 		time.Sleep(time.Until(nextPollTime))

// 	}

// }
