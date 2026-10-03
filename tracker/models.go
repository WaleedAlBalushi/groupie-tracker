package tracker

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type HomePageData struct {
	Artists []Artist
	Query   string
	Total   int
}

type ArtistPageData struct {
	Artist   Artist
	Concerts map[string][]string
}

type ErrorPageData struct {
	Code    int
	Title   string
	Message string
}
