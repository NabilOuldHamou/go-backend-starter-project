package initializers

import (
	"log"
	"os"
)

func CreateAssetsFolders() {

	if _, err := os.Stat("assets"); os.IsNotExist(err) {
		if err := os.MkdirAll("assets", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}
