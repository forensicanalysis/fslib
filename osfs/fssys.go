// +build go1.12

package osfs

import (
	"log"
	"time"

	"gopkg.in/djherbis/times.v1"
)

// Sys returns a map of item attributes.
func (i *Info) Sys() interface{} {
	attributes := map[string]interface{}{}

	t, err := times.Stat(i.syspath)
	if err != nil {
		log.Printf("could not stat times for %s: %s", err, i.syspath)
	}
	if err == nil {
		attributes["accessed"] = t.AccessTime().UTC().Format(time.RFC3339Nano)
		attributes["modified"] = t.ModTime().UTC().Format(time.RFC3339Nano)
		if t.HasChangeTime() {
			attributes["changed"] = t.ChangeTime().UTC().Format(time.RFC3339Nano)
		}
		if t.HasBirthTime() {
			attributes["created"] = t.BirthTime().UTC().Format(time.RFC3339Nano)
		}
	}
	return attributes
}
