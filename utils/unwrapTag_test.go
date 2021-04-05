package utils

import (
	"fmt"
	"regexp"
	"testing"
)

func TestTagA(t *testing.T) {

	//td := "<td>7c0d166f03835f218d69f6fc6deec293512600ce</td>"
	//fmt.Println(unwrapTd(td))
	//anc := `<a href="/chromiumos/platform/tast-tests/+/7c0d166f03835f218d69f6fc6deec293512600ce%5E">564f598b77e3644a4de5f44b4046a93d2febee91</a>`
	//fmt.Printf("%+v", unwrapA(anc))
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
	reg := regexp.MustCompile(`<pre[ -~]*>([ -~\n]*)BUG`)

	match := reg.FindStringSubmatch(pre)
	fmt.Println(match[1])

	revBy := regexp.MustCompile(`[ -~]*Reviewed-by:([ -~]*)`)
	match2 := revBy.FindStringSubmatch(pre)
	fmt.Println(match2[1])

	fmt.Printf("%+v", ExtractIdentity(match2[1]))
}
