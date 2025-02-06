package api

import (
	"database/sql"
	"fmt"
	"time"
)

/*
This package contains function that get never called directly
but always through other functions. They get to create
their own errors, but they do not manage the writer, as they are just
utilities.
*/

/*
It retrieves values from rows that should contain users.

This function passes errors without handling them.
*/

// It returns a proper date-time-formatted string.
func GetTime() string {
	currentTime := time.Now()
	datetime := fmt.Sprintf("%d-%d-%dT%d:%d:%dZ",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second())
	return datetime
}

// It retrieves users from sql's queried rows.
func UsersRowReading(res *sql.Rows) ([]string, error) {
	// Retrieving the values from rows
	var answer []string // array of actual values
	var id, username, propic *string
	for {
		if res.Next() { // there is another value to be scanned
			err := res.Scan(&id, &username, &propic)
			if err == nil {
				if propic == nil {
					tmp := "Null"
					propic = &tmp
				}
				answer = append(answer, *id, *username, *propic)
			} else {
				return nil, err // the scan has had an error
			}
		} else {
			if res.Err() == nil { // there are no more values to scan in the current set
				if res.NextResultSet() { // there are values to be scanned
					continue // in the next set
				} else {
					if res.Err() == nil { // there are no more values, and the scan can end
						break
					} else { // next set scan went unsuccessfully
						return nil, res.Err()
					}
				}
			} else { // next scan went unsuccessfully
				return nil, res.Err()
			}
		}
	}

	return answer, nil
}

// It retrieves groupmembers from sql's queried rows.
func GroupMembersRowReading(res *sql.Rows) ([]string, error) {
	// Retrieving the values from rows
	var answer []string // array of actual values
	var id, member *string
	for {
		if res.Next() { // there is another value to be scanned
			err := res.Scan(&id, &member)
			if err == nil {
				answer = append(answer, *id, *member)
			} else {
				return nil, err // the scan has had an error
			}
		} else {
			if res.Err() == nil { // there are no more values to scan in the current set
				if res.NextResultSet() { // there are values to be scanned
					continue // in the next set
				} else {
					if res.Err() == nil { // there are no more values, and the scan can end
						break
					} else { // next set scan went unsuccessfully
						return nil, res.Err()
					}
				}
			} else { // next scan went unsuccessfully
				return nil, res.Err()
			}
		}
	}

	return answer, nil
}

// It retrieves private chats from sql's queried rows.
func PrivchatsRowReading(res *sql.Rows) ([]string, error) {
	// Retrieving the values from rows
	var answer []string // array of actual values
	var id, member1, member2 *string
	for {
		if res.Next() { // there is another value to be scanned
			err := res.Scan(&id, &member1, &member2)
			if err == nil {
				answer = append(answer, *id, *member1, *member2)
			} else {
				return nil, err // the scan has had an error
			}
		} else {
			if res.Err() == nil { // there are no more values to scan in the current set
				if res.NextResultSet() { // there are values to be scanned
					continue // in the next set
				} else {
					if res.Err() == nil { // there are no more values, and the scan can end
						break
					} else { // next set scan went unsuccessfully
						return nil, res.Err()
					}
				}
			} else { // next scan went unsuccessfully
				return nil, res.Err()
			}
		}
	}

	return answer, nil
}

// It retrieves messages from sql's queried rows.
func MessageRowReading(res *sql.Rows) ([]string, error) {
	// Retrieving the values from rows
	var answer []string // array of actual values
	var id, sender, created_at, content, photo, checkmarks, replying_to *string
	for {
		if res.Next() { // there is another value to be scanned
			err := res.Scan(&id, &sender, &created_at, &content, &photo, &checkmarks, &replying_to)
			if err == nil {
				tmp := "Null"
				if content == nil {
					content = &tmp
				}
				if photo == nil {
					photo = &tmp
				}
				if replying_to == nil {
					replying_to = &tmp
				}
				answer = append(answer, *id, *sender, *created_at, *content, *photo, *checkmarks, *replying_to)
			} else {
				return nil, err // the scan has had an error
			}
		} else {
			if res.Err() == nil { // there are no more values to scan in the current set
				if res.NextResultSet() { // there are values to be scanned
					continue // in the next set
				} else {
					if res.Err() == nil { // there are no more values, and the scan can end
						break
					} else { // next set scan went unsuccessfully
						return nil, res.Err()
					}
				}
			} else { // next scan went unsuccessfully
				return nil, res.Err()
			}
		}
	}

	return answer, nil
}
