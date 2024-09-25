package tfimportprovider

import "encoding/json"

func GetDataFromJsonString(jsonString string) (map[string]interface{}, error) {

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
