package check

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// ConfigParser parses a chk file
type ConfigParser struct {
	configMap ConfigMap
}

type ConfigMap map[string]string

// Parse parses a file
func (c *ConfigParser) Parse(filename string) (checks []Check) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	fileString := string(fileBytes)
	checkStrings := regexp.MustCompile("\\n\\s*\\n").Split(fileString, -1)
	for i, checkString := range checkStrings {
		configMap := c.getMap(checkString)
		checkType, err := configMap.GetString("type")
		if err != nil {
			fmt.Println(fmt.Errorf("Check #%d in file %s missing a readable type", i+1, filename))
			continue
		}
		// make check object (factory?)
		var check Check
		switch checkType {
		case "net_dns":
			check = &NetDNSCheck{}
		case "net_ping":
			check = &NetPingCheck{}
		case "net_tls":
			check = &NetTLSCheck{}
		case "http_status":
			check = &HTTPStatusCheck{}
		case "http_content":
			check = &HTTPContentCheck{}
		}

		if check == nil {
			fmt.Println(fmt.Errorf("Invalid check type: %s", checkType))
			continue
		}
		check.Parse(configMap)

		checks = append(checks, check)
	}
	return
}

func (c *ConfigParser) getMap(input string) ConfigMap {
	lines := strings.Split(input, "\n")
	configMap := make(map[string]string, len(lines))
	for _, line := range lines {
		lineParts := regexp.MustCompile("\\s*:\\s*").Split(line, 2)
		if len(lineParts) < 2 {
			continue
		}
		key := strings.Trim(lineParts[0], " \t")
		value := strings.Trim(lineParts[1], " \t")
		configMap[key] = value
	}
	return configMap
}

// GetString returns a string
func (c ConfigMap) GetString(key string) (string, error) {
	content, ok := c[key]
	if !ok {
		return "", fmt.Errorf("Key does not exist (%s)", key)
	}
	// format "content" -> content
	if strings.HasPrefix(content, "\"") && strings.HasSuffix(content, "\"") {
		return content[1 : len(content)-1], nil
	}

	// format content -> content
	return content, nil
}

// GetInt returns a int repesentation
func (c ConfigMap) GetInt(key string) (int, error) {
	content, ok := c[key]
	if !ok {
		return 0, fmt.Errorf("Key does not exist (%s)", key)
	}
	// format 23
	value, err := strconv.Atoi(content)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// GetInt returns a int array repesentation
func (c ConfigMap) GetIntArray(key string) ([]int, error) {
	result := make([]int, 0, 0)
	content, ok := c[key]
	if !ok {
		return result, fmt.Errorf("Key does not exist (%s)", key)
	}
	// format 23, 78
	parts := regexp.MustCompile("\\s*,\\s*").Split(content, -1)
	for _, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil {
			return make([]int, 0, 0), err
		}
		result = append(result, value)
	}
	return result, nil
}
