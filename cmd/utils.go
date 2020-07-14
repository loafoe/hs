package cmd

import (
	"bytes"
	"encoding/json"
)

func pretty(data []byte) string {
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, data, "", "    ")
	return prettyJSON.String()
}
