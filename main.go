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
	repoURL      = flag.String("repo", "/chromium.googlesource.com/chromiumos/platform/tast", "Repository URL")
	repoBranch   = flag.String("branch", "main", "branch name where the parser should run")
	saveLocation = flag.String("dir", "./commits", "folder where parsed commit messages is going to be stored")
	csvPath      = flag.String("csvpath", "stats.csv", "csv file location where the details statistics is going to be stored")
	numCommits   = flag.Int("commits", 20, "Number of commits to be parsed")
)

func main() {
	flag.Parse()
	PrintInfo()

	if err := scraper(); err != nil {
		panic(err)
	}

	fmt.Println("Execution Complete")
}

func scraper() error {

	// A container which keeps track of
	pool := container.NewContainer()
	// A temporary buffer to store commit messages with commitID as filename, which is going to be flushed into files later
	wb := utils.NewWriteBuffer()

	var err error

	// In case context timed-out, persist already fetched data
	defer func() error {
		err = wb.DumpContent(*saveLocation)
		if err != nil {
			return err
		}
		err = pool.WriteCSV(*csvPath)
		return err
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	//graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sig)

	go func() {
		<-sig
		fmt.Println("Exiting...")
		cancel()
	}()

	devt := devtool.New("http://127.0.0.1:9222")
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		return err
	}

	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := cdp.NewClient(conn)

	domContent, err := client.Page.DOMContentEventFired(ctx)
	if err != nil {
		return err
	}
	defer domContent.Close()

	err = client.Page.Enable(ctx)
	if err != nil {
		return err
	}

	hardURL, softURL := formatURL()

	for cnt := 0; cnt < *numCommits; cnt++ {

		fmt.Println("Commit No: ", cnt+1)
		// Navigate to next commit
		navArgs := page.NewNavigateArgs(hardURL + softURL)
		_, err = client.Page.Navigate(ctx, navArgs)
		if err != nil {
			return err
		}

		//wait until dom loaded successfully
		_, err = domContent.Recv()
		if err != nil {
			return err
		}

		nextCommitAnchor, err := parseCommitPage(ctx, client, pool, wb)
		if err != nil {
			return err
		}

		// Load next commit's relative url
		softURL = nextCommitAnchor.Url

		// bypass rate limiting quota for multiple request, please give larger timeout
		if (cnt+1)%50 == 0 {
			time.Sleep(5 * time.Second)
			if cnt >= 150 {
				time.Sleep(5 * time.Second)
			}
		}
	}
	return err
}

func parseCommitPage(ctx context.Context, client *cdp.Client, pool *container.Container, wb *utils.WriteBuffer) (*utils.Anchor, error) {
	// Get the whole DOM
	doc, err := client.DOM.GetDocument(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Fetch Commit Message
	pre, err := client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(doc.Root.NodeID, "pre"))
	if err != nil {
		return nil, &utils.DOMError{Err: err}
	}
	preRes, err := client.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &pre.NodeID,
	})
	if err != nil {
		return nil, &utils.DOMError{Err: err}
	}

	// Extract Reviewed-By and Commit message
	parsedPre, err := utils.UnwrapPre(preRes.OuterHTML)
	if err != nil {
		return nil, err
	}

	// Parsed email, name of Reviewed-BY
	revBy := utils.ExtractIdentity(parsedPre.RevBy)
	// Increment his/her review count
	pool.AddReview(revBy)

	// Current commitID & commit author & next url
	tableDatas, err := client.DOM.QuerySelectorAll(ctx, &dom.QuerySelectorAllArgs{
		NodeID:   doc.Root.NodeID,
		Selector: "td",
	})
	if err != nil {
		return nil, &utils.DOMError{Err: err}
	}
	commitHTML, err := client.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &tableDatas.NodeIDs[0],
	})
	if err != nil {
		return nil, &utils.DOMError{Err: err}
	}

	// Parsed commit ID as unique filename
	commitID := utils.UnwrapTd(commitHTML.OuterHTML)

	// add commitId as filename and body as message to the temporary buffer
	wb.AppendContent(parsedPre.Message, commitID)

	// Get commit author
	authorHTML, err := client.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &tableDatas.NodeIDs[2],
	})
	if err != nil {
		return nil, &utils.DOMError{Err: err}
	}

	// Unwrap author from HTML tag
	authorRaw := utils.UnwrapTd(authorHTML.OuterHTML)
	// Parse email and name from string
	author := utils.ExtractIdentity(authorRaw)

	// Increment his/her commit count
	pool.AddCommit(author)

	// Get next commit url
	hyperlinks, err := client.DOM.QuerySelector(ctx, &dom.QuerySelectorArgs{
		NodeID:   tableDatas.NodeIDs[len(tableDatas.NodeIDs)-1],
		Selector: "a",
	})
	if err != nil {
		return nil, &utils.DOMError{Err: err}
	}
	anchorTag, err := client.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &hyperlinks.NodeID,
	})
	if err != nil {
		return nil, &utils.DOMError{Err: err}
	}

	// Parse url from href field in anchor tag
	nextCommitAnchor := utils.UnwrapA(anchorTag.OuterHTML)
	return nextCommitAnchor, nil
}

func formatURL() (hardURL string, softURL string) {
	trimmedURL := strings.Trim(*repoURL, "/")
	if !strings.HasPrefix(trimmedURL, "https://") {
		trimmedURL = "https://" + trimmedURL
	}

	splits := strings.Split(trimmedURL, "/")
	hardURL = strings.Join(splits[:3], "/")
	softURL = "/" + strings.Join(splits[3:], "/")

	defaultBranch := *repoBranch
	// For switching to a branch & UnwrapPre will passively check if it exists
	softURL += "/+/refs/heads/" + defaultBranch
	return
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
