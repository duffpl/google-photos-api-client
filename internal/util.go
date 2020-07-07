package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func UnmarshalResponse(res *http.Response, dst interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read body bytes: %w", err)
	}
	if strPtr, ok := dst.(*string); ok {
		*strPtr = string(b)
		return nil
	}
	err = json.Unmarshal(b, dst)
	if err != nil {
		return fmt.Errorf("cannot unmarshal: %w", err)
	}
	return nil
}
