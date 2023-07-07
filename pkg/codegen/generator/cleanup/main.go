package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	if err := os.RemoveAll("./pkg/generated"); err != nil {
		logrus.Fatal(err)
	}
}
