package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func UnmarshalJSON(data io.ReadCloser, item interface{}) {
	buf, _ := ioutil.ReadAll(data)
	json.Unmarshal(buf, &item)
}
