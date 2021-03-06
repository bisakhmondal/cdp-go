package core

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

// Struct to keep track (as Key-Value pair) of all commits and reviews made by individual developers. Email has been considered as key.
type Container struct {
	maps map[string]*Developer
}

func NewContainer() *Container {
	return &Container{map[string]*Developer{}}
}

// Method to increment commit count for a particular identity
func (c *Container) AddCommit(id *Identity) {

	dev, ok := c.maps[id.Email]

	if !ok {
		dev = NewDeveloper(id)
		c.maps[id.Email] = dev
	}

	dev.NumCommit++
}

// Method to increment review count for a particular identity
func (c *Container) addReview(id *Identity) {
	dev, ok := c.maps[id.Email]

	if !ok {
		dev = NewDeveloper(id)
		c.maps[id.Email] = dev
	}

	dev.NumReview++
}

// Method to increment review counts for multiple identities
func (c *Container) AddReviews(ids []*Identity) {
	for _, id := range ids {
		c.addReview(id)
	}
}

// Method to write all gathered information into a csv file
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
