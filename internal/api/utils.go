package api

import (
	"fmt"

	"gopkg.in/h2non/gentleman.v2"
)

func HandleError(response gentleman.Response, err error) error {
	if err != nil {
		return fmt.Errorf("http error: %e", err)
	}
	if !response.Ok {
		return fmt.Errorf("api error, response status %d", response.StatusCode)
	}
	return nil
}
