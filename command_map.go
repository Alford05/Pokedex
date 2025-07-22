package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type MapList struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []result `json:"results"`
}

func commandMapf(cfg *config) error {
	var url string
	if cfg.NextURL == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = cfg.NextURL
	}
	var mapData MapList

	if entry, ok := cfg.Pokecache.Get(url); ok {
		if err := json.Unmarshal(entry, &mapData); err != nil {
			return fmt.Errorf("failed to decode data from cache: %w", err)
		}
		for _, area := range mapData.Results {
			fmt.Println(area.Name)
		}
		cfg.NextURL = mapData.Next
		cfg.PrevURL = mapData.Previous
		return nil
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch map data: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	cfg.Pokecache.Add(url, body)
	if err := json.Unmarshal(body, &mapData); err != nil {
		return fmt.Errorf("failed to decode map data: %w", err)
	}

	for _, area := range mapData.Results {
		fmt.Println(area.Name)
	}
	cfg.NextURL = mapData.Next
	cfg.PrevURL = mapData.Previous
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.PrevURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	var url string = cfg.PrevURL
	var mapData MapList

	if entry, ok := cfg.Pokecache.Get(url); ok {
		if err := json.Unmarshal(entry, &mapData); err != nil {
			return fmt.Errorf("failed to decode data from cache: %w", err)
		}
		for _, area := range mapData.Results {
			fmt.Println(area.Name)
		}
		cfg.NextURL = mapData.Next
		cfg.PrevURL = mapData.Previous
		return nil
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch map data: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	cfg.Pokecache.Add(url, body)

	if err := json.Unmarshal(body, &mapData); err != nil {
		return fmt.Errorf("failed to decode map data: %w", err)
	}
	for _, area := range mapData.Results {
		fmt.Println(area.Name)
	}
	cfg.NextURL = mapData.Next
	cfg.PrevURL = mapData.Previous
	return nil
}
