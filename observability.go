package obs

import "github.com/JoinVerse/obs/errtrack"

// Observer provides observer object
type Observer struct {
	Log *Logger
	*errtrack.ErrorTracker
}

// New returns a new observer.
func New(errConfig errtrack.Config) Observer {
	return Observer{Log: NewLogger(), ErrorTracker: errtrack.New(errConfig)}
}
