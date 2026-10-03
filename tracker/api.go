package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiURL = "https://groupietrackers.herokuapp.com/api"

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func fetchArtists() ([]Artist, error) {
	resp, err := client.Get(apiURL + "/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artists API returned status %d", resp.StatusCode)
	}

	var artists []Artist

	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}

	return artists, nil
}

func fetchRelation(id int) (Relation, error) {
	url := fmt.Sprintf("%s/relation/%d", apiURL, id)

	resp, err := client.Get(url)
	if err != nil {
		return Relation{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Relation{}, fmt.Errorf("relation API returned status %d", resp.StatusCode)
	}

	var relation Relation

	if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil {
		return Relation{}, err
	}

	return relation, nil
}
