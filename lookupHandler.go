package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// templ represents a single template
type lookupHandler struct {
}

func (l *lookupHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error
	tenantGuid, err := l.validateRequest(req)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	switch req.Method {
	case "GET":
		lookups := l.GetLookups(tenantGuid)
		json.NewEncoder(w).Encode(lookups)
	case "POST":
		lookup, err := l.PostLookup(tenantGuid, req)
		if err != nil {
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(lookup)
		}
	case "DELETE":
		err := l.DeleteLookup(tenantGuid, req)
		if err != nil {
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode("Success")
		}
	}
}

func (l *lookupHandler) GetLookups(tenantGuid string) []Lookup {
	var lookups []Lookup

	// read it in
	raw, err := ioutil.ReadFile(getFilename(tenantGuid))
	if err == nil {
		json.Unmarshal(raw, &lookups)
	}
	return lookups
}

func (l *lookupHandler) PostLookup(tenantGuid string, req *http.Request) (lookup Lookup, err error) {
	//var lookup Lookup

	json.NewDecoder(req.Body).Decode(&lookup)

	lookup.Added_on = time.Now()

	ls := l.GetLookups(tenantGuid)

	// append the new one
	ls = append(ls, lookup)

	// write it back
	bytes, err := json.Marshal(ls)
	err = ioutil.WriteFile(getFilename(tenantGuid), bytes, 0644)

	return
}

func (l *lookupHandler) DeleteLookup(tenantGuid string, req *http.Request) error {
	var lookup Lookup
	json.NewDecoder(req.Body).Decode(&lookup)

	ls := l.GetLookups(tenantGuid)

	// find the index of the type
	index := -1

	for i, l := range ls {
		if l.Title == lookup.Title {
			index = i
			break
		}
	}

	// remove it from the array
	if index > -1 {
		ls = append(ls[:index], ls[index+1:]...)
	}

	// write it back
	bytes, err := json.Marshal(ls)
	err = ioutil.WriteFile(getFilename(tenantGuid), bytes, 0644)

	return err
}

func (l *lookupHandler) validateRequest(req *http.Request) (tenantGuid string, err error) {
	if tenantGuid = req.Header.Get("tenant-guid"); tenantGuid == "" {
		return "", errors.New("Please specify a tenant guid")
	}

	return
}

func getFilename(tenantGuid string) string {
	return tenantGuid + "-lookups.json"
}
