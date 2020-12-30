package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	switch {
	case len(os.Args) == 1:
		fmt.Println("Error: Must pass directory to scan and output file")
		os.Exit(1)
	case len(os.Args) == 2:
		fmt.Println("Error: Must pass output file")
		os.Exit(1)
	case len(os.Args) > 3:
		fmt.Println("Error: Too many arguments, must pass root directory to scan and output file")
		os.Exit(1)
	}

	result := scanFiles(os.Args[1])

	if file, err := os.Create(os.Args[2]); err == nil {
		defer file.Close()

		writer := bufio.NewWriter(file)
		for k, _ := range result {
			_, err := writer.WriteString(k + "\n")
			if err != nil {
				panic(err)
			}
		}
		writer.Flush()
	} else {
		panic(err)
	}
}

func scanFiles(dirname string) map[string]bool {
	var files []string
	result := map[string]bool{}

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(strings.ToLower(path), ".java") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if keys, err := scanFile(file); err == nil {
			for k, _ := range keys {
				result[k] = true
			}
		} else {
			panic(err)
		}
	}

	return result
}

func scanFile(filename string) (map[string]bool, error) {
	result := map[string]bool{}
	if file, err := os.Open(filename); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			keys := scanLine(scanner.Text())
			for _, k := range keys {
				result[k] = true
			}
		}

		if err := scanner.Err(); err != nil {
			return result, err
		}

	} else {
		return result, err
	}

	return result, nil
}

func scanLine(text string) []string {
	result := []string{}
	currIndex := 0
	scanText := strings.ToLower(text)

	for {
		if currIndex >= len(text) {
			break
		} else if idx := strings.Index(scanText[currIndex:], "session.getproperty("); idx >= 0 {
			currIndex += idx + len("session.getproperty(")
			if s, err := scanForKey(scanText[currIndex:], text[currIndex:]); err == nil {
				result = append(result, s)
			} else {
				panic(err)
			}
		} else {
			break
		}
	}

	return result
}

func scanForKey(scanText string, text string) (string, error) {
	result := ""

	if idx := strings.Index(scanText, ","); idx >= 0 {
		if idx2 := strings.Index(scanText[idx+1:], ")"); idx2 >= 0 {
			return strings.Trim(text[idx+1:idx2+idx+1], " "), nil
		} else {
			return result, errors.New(fmt.Sprintf("Didn't find ')' in %s", text))
		}
	} else {
		return result, errors.New(fmt.Sprintf("Didn't find ',' in %s", text))
	}
}
