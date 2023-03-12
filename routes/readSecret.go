package routes

import (
	"bufio"
	"os"
)

func ReadSecret(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	return lines[0], nil
}
