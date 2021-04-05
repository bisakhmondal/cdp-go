package utils

import (
	"regexp"
	"strings"
)

type Identity struct {
	Name, Email string
}

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
