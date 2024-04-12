package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseBoolean(data interface{}) interface{} {
	if dataStr, ok := data.(string); ok {
		dataStr := strings.TrimSpace(dataStr)
		if dataStr == "1" || strings.ToLower(dataStr) == "t" || strings.ToLower(dataStr) == "true" {
			return true
		}
		if dataStr == "0" || strings.ToLower(dataStr) == "f" || strings.ToLower(dataStr) == "false" {
			return false
		}
	}
	return nil
}

func parseNull(data interface{}) interface{} {
	boolResult := parseBoolean(data)
	if boolResult == true {
		return "null"
	}
	return nil
}

func parseNumber(data interface{}) interface{} {
	if dataStr, ok := data.(string); ok {
		dataStr := strings.TrimSpace(dataStr)
		dataStr = strings.TrimLeft(dataStr, "0")
		r := regexp.MustCompile(`^-?(\d|[1-9]\d+)(\.\d+)?([Ee][+-]?\d+)?$`)
		if r.MatchString(dataStr) {
			if strings.Contains(dataStr, ".") {
				f, err := strconv.ParseFloat(dataStr, 64)
				if err == nil {
					return f
				}
			} else {
				i, err := strconv.Atoi(dataStr)
				if err == nil {
					return i
				}
			}
		}
	}
	return nil
}

func parseString(data interface{}) interface{} {
	if dataStr, ok := data.(string); ok {
		dataStr := strings.TrimSpace(dataStr)
		if len(dataStr) > 0 && strings.ToLower(dataStr[len(dataStr)-1:]) == "z" {
			dataStr = dataStr[:len(dataStr)-1] + "+00:00"
		}
		r := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}[tT]\d{2}:\d{2}:\d{2}(\.\d+)?([+-]\d{2}:\d{2})$`)
		if r.MatchString(dataStr) {
			t, err := time.Parse(time.RFC3339, dataStr)
			if err == nil {
				return t.Unix()
			}
		}
		return dataStr
	}
	return nil
}

func parseScalar(data map[string]interface{}) interface{} {
	for key, value := range data {
		key = strings.TrimSpace(key)
		switch key {
		case "S":
			return parseString(value)
		case "N":
			return parseNumber(value)
		case "BOOL":
			return parseBoolean(value)
		default:
			return nil
		}
	}
	return nil
}

func parseValue(data map[string]interface{}) interface{} {
	for key, value := range data {
		key = strings.TrimSpace(key)
		switch key {
		case "NULL":
			return parseNull(value)
		case "L":
			return parseList(value)
		case "M":
			return parseMap(value)
		default:
			return parseScalar(data)
		}
	}
	return nil
}

func parseList(data interface{}) []interface{} {
	var result []interface{}
	if dataList, ok := data.([]interface{}); ok {
		for _, value := range dataList {
			if valueMap, ok := value.(map[string]interface{}); ok {
				if parsedValue := parseScalar(valueMap); parsedValue != nil {
					result = append(result, parsedValue)
				}
			}
		}
	}
	return result
}

func parseMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if dataMap, ok := data.(map[string]interface{}); ok {
		for key, value := range dataMap {
			key = strings.TrimSpace(key)
			if key != "" && value != nil { // Check if the value is not nil
				if valueMap, ok := value.(map[string]interface{}); ok {
					if parsedValue := parseValue(valueMap); parsedValue != nil {
						result[key] = parsedValue
					}
				}
			}
		}
	}
	return result
}

func transforms(data interface{}) string {
	parsedData := parseMap(data)
	transformedData, _ := json.MarshalIndent(parsedData, "", "  ")
	return string(transformedData)
}

func transformJSON(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON: %v", err)
	}
	return transforms(jsonData)
}

func main() {
	inputFile := "input.json"
	transformedData := transformJSON(inputFile)
	fmt.Println(transformedData)
}
