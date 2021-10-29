package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func UnmarshalJSON(data io.ReadCloser, item interface{}) error {
	buf, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}
	json.Unmarshal(buf, &item)
	return nil
}
