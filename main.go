package main

import (
	"cdp-go/container"
	"cdp-go/utils"
	"context"
	"fmt"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"strings"
	"time"
)

var (
	pool *container.Container
	wb   *utils.WriteBuffer
)

func main() {
	pool = container.NewContainer()
	wb = utils.NewWriteBuffer()

	//incase context timed-out
	//persist already fetched data
	defer func() {
		must(wb.DumpContent("./commits"))
		must(pool.WriteCSV("stats.csv"))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	devt := devtool.New("http://127.0.0.1:9222")
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		must(err)
	}

	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	must(err)
	defer conn.Close()

	c := cdp.NewClient(conn)

	domContent, err := c.Page.DOMContentEventFired(ctx)
	must(err)
	defer domContent.Close()

	err = c.Page.Enable(ctx)
	must(err)
	FULLURL := "https://chromium.googlesource.com/chromiumos/platform/tast"
	splits := strings.Split(FULLURL, "/")
	HARDURL := strings.Join(splits[:3], "/")
	SOFTURL := "/" + strings.Join(splits[3:], "/")
	fmt.Println(HARDURL)
	fmt.Println(SOFTURL)

	defaultBranch := "main"
	//switch to a branch & check if it exists
	SOFTURL += "/+/refs/heads/" + defaultBranch
	//navArgs := page.NewNavigateArgs(HARDURL + SOFTURL)
	//_, err = c.Page.Navigate(ctx, navArgs)
	//if err != nil {
	//	panic(fmt.Errorf("bad branch name!! Default branch is main"))
	//}

	numCnt := 20

	for cnt := 0; cnt < numCnt; cnt++ {
		fmt.Println("Commit No: ", cnt+1)
		navArgs := page.NewNavigateArgs(HARDURL + SOFTURL)
		_, err := c.Page.Navigate(ctx, navArgs)
		must(err)

		//wait until dom loaded successfully
		_, err = domContent.Recv()
		must(err)
		//fmt.Printf("Page loaded with frame ID: %s\n", nav.FrameID)

		doc, err := c.DOM.GetDocument(ctx, nil)
		must(err)

		//Fetch Commit Message
		pre, err := c.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(doc.Root.NodeID, "pre"))

		preRes, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
			NodeID: &pre.NodeID,
		})

		parsedPre := utils.UnwrapPre(preRes.OuterHTML)
		revBy := utils.ExtractIdentity(parsedPre.RevBy)
		pool.AddReview(revBy)

		//fmt.Println(parsedPre.Message)

		//getCurrentCommitID & Commit Author
		tableDatas, err := c.DOM.QuerySelectorAll(ctx, &dom.QuerySelectorAllArgs{
			NodeID:   doc.Root.NodeID,
			Selector: "td",
		})

		commitHTML, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
			NodeID: &tableDatas.NodeIDs[0],
		})

		commitID := utils.UnwrapTd(commitHTML.OuterHTML)

		//fmt.Println(commitID)

		//add commitId as filename and body as message to writebuffer
		wb.AppendContent(parsedPre.Message, commitID)

		//commitAuthor
		authorHTML, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
			NodeID: &tableDatas.NodeIDs[2],
		})

		authorRaw := utils.UnwrapTd(authorHTML.OuterHTML)

		author := utils.ExtractIdentity(authorRaw)

		pool.AddCommit(author)
		//fmt.Printf("%+v", author)

		//get next url
		hyperlinks, err := c.DOM.QuerySelector(ctx, &dom.QuerySelectorArgs{
			NodeID:   tableDatas.NodeIDs[len(tableDatas.NodeIDs)-1],
			Selector: "a",
		})

		anchorTag, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
			NodeID: &hyperlinks.NodeID,
		})

		nextCommitAnchor := utils.UnwrapA(anchorTag.OuterHTML)
		//fmt.Printf("%+v", nextCommitAnchor)

		SOFTURL = nextCommitAnchor.Url
		//fmt.Println("-----------------------------------------------------\n-------------------------------------------------------")

	}

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
