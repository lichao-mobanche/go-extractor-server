package global

import "fmt"

// QueueUnavailableError TODO
type QueueUnavailableError string

func (e QueueUnavailableError) Error() string {
	return fmt.Sprintf("queue unavailable %s", string(e))
}

// QueueFullError TODO
type QueueFullError string

func (e QueueFullError) Error() string {
	return fmt.Sprintf("queue full %s", string(e))
}