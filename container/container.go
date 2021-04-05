package container

import (
	"cdp-go/utils"
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

//email identity mapping
type Container struct {
	maps map[string]*Developer
}

func NewContainer() *Container {
	return &Container{map[string]*Developer{}}
}

func (c *Container) AddCommit(id *utils.Identity) {

	dev, ok := c.maps[id.Email]

	if !ok {
		dev = NewDeveloper(id)
		c.maps[id.Email] = dev
	}

	dev.NumCommit++
}

func (c *Container) AddReview(id *utils.Identity) {
	dev, ok := c.maps[id.Email]

	if !ok {
		dev = NewDeveloper(id)
		c.maps[id.Email] = dev
	}

	dev.NumReview++
}

func (c *Container) WriteCSV(csvName string) error {
	path, err := filepath.Abs(csvName)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	var data = [][]string{
		{"contributor", "created", "reviewed"},
	}

	for _, dev := range c.maps {
		data = append(data,
			[]string{dev.GetIdentityString(), strconv.Itoa(dev.NumCommit), strconv.Itoa(dev.NumReview)},
		)
	}

	return csvWriter.WriteAll(data)
}
