package tfimportprovider

import "encoding/json"

func GetDataFromJsonString(jsonString string) (map[string]any, error) {

	var data map[string]any
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
