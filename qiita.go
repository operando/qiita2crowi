package main

import (
	"time"
)

type Qiita struct {
	Articles []Articles    `json:"articles"`
	Groups   []interface{} `json:"groups"`
	Projects []Projects    `json:"projects"`
	Version  string        `json:"version"`
}

type Articles struct {
	RenderedBody string      `json:"rendered_body"`
	Body         string      `json:"body"`
	Coediting    bool        `json:"coediting"`
	CreatedAt    time.Time   `json:"created_at"`
	Group        interface{} `json:"group"`
	ID           string      `json:"id"`
	Private      bool        `json:"private"`
	Tags         []struct {
		Name     string        `json:"name"`
		Versions []interface{} `json:"versions"`
	} `json:"tags"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
	User      struct {
		ID              string `json:"id"`
		PermanentID     int    `json:"permanent_id"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"user"`
	Comments []interface{} `json:"comments"`
}

type Projects struct {
	RenderedBody string    `json:"rendered_body"`
	Archived     bool      `json:"archived"`
	Body         string    `json:"body"`
	CreatedAt    time.Time `json:"created_at"`
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         struct {
		ID              string `json:"id"`
		PermanentID     int    `json:"permanent_id"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"user"`
}
