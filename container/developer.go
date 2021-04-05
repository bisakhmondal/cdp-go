package container

import (
	"cdp-go/utils"
	"fmt"
)

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

func (d *Developer) GetIdentityString() string {
	return fmt.Sprintf("%s <%s>", d.Identity.Name, d.Identity.Email)
}
