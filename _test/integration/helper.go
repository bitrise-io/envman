package integration

import "os"

func binPath() string {
	pth := os.Getenv("INTEGRATION_TEST_BINARY_PATH")
	if pth == "" {
		pth = "envman"
	}
	return pth
}
