package e2e_test

import "os"

func BaseAddr() string {
	baseAddr := os.Getenv("BASE_URL")
	if baseAddr == "" {
		return "http://localhost:4000"
	}

	return baseAddr
}
