package parser

import (
	"log"
	"os"
	"os/exec"
	"testing"
)

const (
	testDataFilename = "testData.tar.gz"
)

func TestMain(m *testing.M) {
	unzip()

	code := m.Run()

	zip()

	os.Exit(code)
}

func unzip() {
	err := exec.Command("tar", "xvfz", testDataFilename).Run()
	if err != nil {
		log.Fatalf("Error unzipping test-data %s: %s", testDataFilename, err)
	}
}

func zip() {
	err := exec.Command("tar", "cvfz", testDataFilename, "enums/", "interfaces/", "operations/", "structs/").Run()
	if err != nil {
		log.Fatalf("Error zipping test-data %s: %s", testDataFilename, err)
	}
	err = exec.Command("rm", "-rf", "enums/", "interfaces/", "operations/", "structs/").Run()
	if err != nil {
		log.Fatalf("Error removing test-directories: %s", err)
	}

}
