package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"

	"cdp-go/container"
	"cdp-go/utils"
)

var (
	timeout      = flag.Int("timeout", 20, "context timeout in seconds")
	repoURL      = flag.String("repo", "https://chromium.googlesource.com/chromiumos/platform/tast", "Repository URL")
	repoBranch   = flag.String("branch", "main", "branch name where the parser should run")
	saveLocation = flag.String("dir", "./commits", "folder where parsed commit messages is going to be stored")
	csvPath      = flag.String("csvpath", "stats.csv", "csv file location where the details statistics is going to be stored")
	numCommits   = flag.Int("commits", 20, "Number of commits to be parsed")
)

var (
	pool *container.Container
	wb   *utils.WriteBuffer
)

func main() {
	flag.Parse()
	PrintInfo()

	pool = container.NewContainer()
	wb = utils.NewWriteBuffer()

	//incase context timed-out
	//persist already fetched data
	defer func() {
		must(wb.DumpContent(*saveLocation))
		must(pool.WriteCSV(*csvPath))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	//graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sig)

	go func() {
		<-sig
		cancel()
	}()

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

	FULLURL := strings.Trim(*repoURL, "/")
	if !strings.HasPrefix(FULLURL, "https://") {
		FULLURL = "https://" + FULLURL
	}

	splits := strings.Split(FULLURL, "/")
	HARDURL := strings.Join(splits[:3], "/")
	SOFTURL := "/" + strings.Join(splits[3:], "/")

	defaultBranch := *repoBranch
	//switch to a branch & will be checked in UnwrapPre if it exists
	SOFTURL += "/+/refs/heads/" + defaultBranch

	for cnt := 0; cnt < *numCommits; cnt++ {
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

		parsedPre, err := utils.UnwrapPre(preRes.OuterHTML)
		must(err)

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

func PrintInfo() {
	fmt.Printf("Repo Branch\t\t: %s\n", *repoBranch)
	fmt.Printf("Commits Num\t\t: %d\n", *numCommits)
	fmt.Printf("Repo URL\t\t: %s\n", *repoURL)
	fmt.Printf("Timeout (seconds)\t: %d\n", *timeout)
	fmt.Printf("Folder Location\t\t: %s\n", *saveLocation)
	fmt.Printf("CSV Path\t\t: %s\n", *csvPath)

	fmt.Println("==================================================\n")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
