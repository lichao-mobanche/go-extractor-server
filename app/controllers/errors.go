package controllers

import "fmt"

// InvalidBody TODO
type InvalidBody string

func (e InvalidBody) Error() string {
	return fmt.Sprintf("invalid %s", string(e))
}
