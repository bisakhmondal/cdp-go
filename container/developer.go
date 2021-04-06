package container

import (
	"fmt"

	"cdp-go/utils"
)

// Struct for individual developer with commit, review count
type Developer struct {
	*utils.Identity
	NumCommit, NumReview int
}

func NewDeveloper(id *utils.Identity) *Developer {
	return &Developer{
		Identity:  id,
		NumCommit: 0,
		NumReview: 0,
	}
}

// Returns the identity of developer as string in "name <email>" format
func (d *Developer) GetIdentityString() string {
	return fmt.Sprintf("%s <%s>", d.Identity.Name, d.Identity.Email)
}
