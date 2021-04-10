package core

import (
	"regexp"
	"strings"
)

// Identity Struct with Name Email fields. Used for reviewer & committer
type Identity struct {
	Name, Email string
}

// Parse raw text to extract out name and email
func ExtractIdentity(str string) *Identity {
	str = strings.Replace(str, "&lt;", "<", -1)
	str = strings.Replace(str, "&gt;", ">", -1)
	str = strings.Trim(str, " ")

	reg := regexp.MustCompile(`([ -~]*)<([ -~]*)>`)

	matches := reg.FindAllStringSubmatch(str, -1)

	id := &Identity{
		Name:  strings.Trim(matches[0][1], " "),
		Email: strings.Trim(matches[0][2], " "),
	}

	return id
}

func ExtractIdentities(str []string) []*Identity {
	var ids []*Identity
	for _, id := range str {
		ids = append(ids, ExtractIdentity(id))
	}
	return ids
}
