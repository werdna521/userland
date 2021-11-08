package slice

import "github.com/werdna521/userland/repository"

// Filter is similar to JS' Array.filter, for []*repository.Session
func FilterSession(
	slice []*repository.Session,
	f func(*repository.Session) bool,
) []*repository.Session {
	sessions := []*repository.Session{}
	for _, v := range slice {
		if f(v) {
			sessions = append(sessions, v)
		}
	}
	return sessions
}
