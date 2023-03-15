package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gprossliner/xhdl"
)

func main() {
	res, err := ExecuteLogic("https://www.google.com", "test")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(res)
}

// ExecuteLogic is a public function that performs something, and may return an error
func ExecuteLogic(serverurl, command string) (result string, err error) {

	// xhdl.Run creates the xhdl context (which implements context.Context),
	// and runs the specified callback function.
	// If any code down the stack calls ctx.Throw, the error is returned from xhdl.Run
	err = xhdl.Run(func(ctx xhdl.Context) {
		result = executeLogic(ctx, serverurl, command)
	})

	return
}

// executeLogic is the main function that expects a xhdl.Context, which
// shows that this func uses xhdl for error-handling
func executeLogic(ctx xhdl.Context, serverurl, command string) string {

	cmdUrl := getCommandUrl(ctx, serverurl, command)

	req, err := http.NewRequestWithContext(ctx, "GET", cmdUrl, nil)
	ctx.Throw(err)

	res, err := http.DefaultClient.Do(req)
	ctx.Throw(err)

	// ctx.Throw can even throw the error from defer calls
	defer ctx.Throw(res.Body.Close())

	fmt.Printf("Get response %d\n", res.StatusCode)

	// let something fail for demo here
	maybeFail(ctx)

	return res.Status
}

// maybeFail fails random for demonstration
func maybeFail(ctx xhdl.Context) {
	rand.Seed(time.Now().UTC().UnixNano())
	if (rand.Int() % 5) == 0 {
		ctx.Throw(fmt.Errorf("demonstration failure"))
	}
}

func getCommandUrl(ctx xhdl.Context, serverurl, command string) string {

	// parse url for validation
	u, err := url.Parse(serverurl)
	ctx.Throw(err) // .Throw checks for err != nil internaly

	if u.Scheme != "https" {
		ctx.Throw(fmt.Errorf("only https is supported"))
	}

	// set the command
	values := u.Query()
	values.Add("command", command)
	u.RawQuery = values.Encode()

	return u.String()
}
