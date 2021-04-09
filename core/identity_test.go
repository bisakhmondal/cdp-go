package core

import "testing"

func TestExtractIdentity(t *testing.T) {
	got := " Jorge Lucangeli Obes &lt;jorgelo@chromium.org&gt;"
	getRevBy := ExtractIdentity(got)
	wantRevBy := &Identity{
		Name:  "Jorge Lucangeli Obes",
		Email: "jorgelo@chromium.org",
	}
	if *getRevBy != *wantRevBy {
		t.Fatal("extract identity parsing failed")
	}
}
