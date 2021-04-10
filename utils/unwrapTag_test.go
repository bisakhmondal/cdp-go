package utils

import (
	"strings"
	"testing"
)

func TestTagPre(t *testing.T) {
	pre := `<pre class="u-pre u-monospace MetadataMessage">policy: Demote policy.PromptForDownloadLocation to informational

Temporarily mark policy.PromptForDownloadLocation to informational
in order not to block the CQ.

BUG=b:183970071
TEST=Test on DUT and make it fail on purpose. Test passes.

Exempt-From-Owner-Approval: Owners are OoO and we need to get CQ green.
Change-Id: <a href="https://chromium-review.googlesource.com/#/q/I8c22214a92740cd4d0a0f9bc61908dec4533561c">I8c22214a92740cd4d0a0f9bc61908dec4533561c</a>
Reviewed-on: <a href="https://chromium-review.googlesource.com/c/chromiumos/platform/tast-tests/+/2802694">https://chromium-review.googlesource.com/c/chromiumos/platform/tast-tests/+/2802694</a>
Tested-by: Andrew Lassalle &lt;andrewlassalle@chromium.org&gt;
Commit-Queue: Andrew Lassalle &lt;andrewlassalle@chromium.org&gt;
Reviewed-by: Jorge Lucangeli Obes &lt;jorgelo@chromium.org&gt;
</pre>
`
	got, err := UnwrapPre(pre)
	if err != nil {
		t.Fatal(err)
	}
	wantMessage := `policy: Demote policy.PromptForDownloadLocation to informational

Temporarily mark policy.PromptForDownloadLocation to informational
in order not to block the CQ.

`
	if got.Message != wantMessage {
		t.Fatal("Error while parsing message")
	}

	wantRev := "Jorge Lucangeli Obes &lt;jorgelo@chromium.org&gt;"

	if len(got.RevBy) == 0 {
		t.Fatal("error while parsing reviewed by")
	}

	if strings.Trim(got.RevBy[0], " ") != wantRev {
		t.Fatal("error during regex match")
	}
}

func TestTagAnchor(t *testing.T) {
	anc := `<a href="/chromiumos/platform/tast-tests/+/7c0d166f03835f218d69f6fc6deec293512600ce%5E">564f598b77e3644a4de5f44b4046a93d2febee91</a>`
	got := UnwrapA(anc)
	want := &Anchor{
		Name: "564f598b77e3644a4de5f44b4046a93d2febee91",
		Url:  "/chromiumos/platform/tast-tests/+/7c0d166f03835f218d69f6fc6deec293512600ce%5E",
	}
	if *got != *want {
		t.Fatal("unwrap anchor failed")
	}
}

func TestTagTD(t *testing.T) {
	td := "<td>7c0d166f03835f218d69f6fc6deec293512600ce</td>"
	got := UnwrapTd(td)
	want := "7c0d166f03835f218d69f6fc6deec293512600ce"

	if got != want {
		t.Fatal("unwrap td failed")
	}
}
