package main

import "time"

type Lookup struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Expression  string    `json:"expression,omitempty"`
	Added_by    string    `json:"added_by,omitempty"`
	Added_on    time.Time `json:"added_on,omitempty"`
}
