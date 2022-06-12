package config

import (
	"regexp"
	"strconv"

	"kube-snap.io/kube-snap/internal/utilities"
)

const (
	//Config properties
	EVENT_BASED_SNAPS  = "event_based_snaps"
	RESYNC_DURATION    = "resync_duration"
	REASON_REGEX       = "reason_regex"
	PRINT_WARNINGS     = "print_warnings"
	LAST_SEEN_DURATION = "last_seen_duration"

	// Other constants
	configsDir = "/etc/configs/"
)

func FetchIsEventBasedSnaps() bool {
	val, err := strconv.ParseBool(utilities.GetValueOf(configsDir, EVENT_BASED_SNAPS))
	utilities.CheckIfError(err, createErrMessage(EVENT_BASED_SNAPS))
	return val
}

func FetchResyncDuration() int {
	val, err := strconv.Atoi(utilities.GetValueOf(configsDir, RESYNC_DURATION))
	utilities.CheckIfError(err, createErrMessage(RESYNC_DURATION))
	return val
}

func FetchReasonRegex() *regexp.Regexp {
	val, err := regexp.Compile(utilities.GetValueOf(configsDir, REASON_REGEX))
	utilities.CheckIfError(err, createErrMessage(REASON_REGEX))
	return val
}

func FetchIsPrintWarnings() bool {
	val, err := strconv.ParseBool(utilities.GetValueOf(configsDir, PRINT_WARNINGS))
	utilities.CheckIfError(err, createErrMessage(PRINT_WARNINGS))
	return val
}

func FetchLastSeenThreshold() int {
	val, err := strconv.Atoi(utilities.GetValueOf(configsDir, LAST_SEEN_DURATION))
	utilities.CheckIfError(err, createErrMessage(LAST_SEEN_DURATION))
	return val
}

func createErrMessage(property string) string {
	return "Invalid value provided for " + property + " in configmap!"
}
