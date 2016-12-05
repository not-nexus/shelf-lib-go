//usr/bin/env go run $0 $@; exit;
package main

import (
    "fmt"
	"github.com/docopt/docopt-go"
	"math/rand"
	"time"
    "github.com/not-nexus/shelf-lib-go"
    "log"
    "os"
    "path/filepath"
)

func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, strlen)

	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

func checkError(message string, err *shelflib.ShelfError) {
    if err != nil {
        os.Stderr.WriteString(message + "\n")
        os.Stderr.WriteString(err.Error() + "\n")
        os.Exit(1)
    }
}

func main() {
    doc := `example.go

        Example usage of shelflib. This is a hacky example that serves the dual
        purpose of manual testing of shelflib.

        Usage:
            ./example/example.go <host> <refName> <path> [options]

        Options:
            -h --help               Show this message.

            -v --verbose            Verbose logging.

            --shelf-token <token>   Shelf token. Required if token not set in
                                    then environment as SHELF_AUTH_TOKEN.

        Arguments:
            <host>                  Host of Shelf to point to.

            <refName>               Bucket to use.

            <path>                  Path to use.
    `
    arguments, _ := docopt.Parse(doc, nil, true, "shelflib-example 1.0.0", false)
    refName := arguments["<refName>"].(string)
    path := arguments["<path>"].(string)
    host := arguments["<host>"].(string)
    shelfToken := os.Getenv("SHELF_AUTH_TOKEN")

    if shelfToken == "" {
        if val, ok := arguments["--shelf-token"]; ok {
            shelfToken = val.(string)
        } else {
            fmt.Println("Must supply a shelf token via the command line or environment variable.")
            os.Exit(1)
        }
    }

    logger := log.New(os.Stderr, "", 0)
    shelfLib := shelflib.New(shelfToken, logger)
    wd, _ := os.Getwd()
    dir, _ := filepath.Abs(filepath.Dir(wd))

    // Let's just use the repositories README.md for our test.
    filePath := filepath.Join(dir, "shelf-lib-go", "README.md")

    basePath := host + refName + "/artifact/" + path
    artifactPath := basePath + "/" + randomString(32)
    fmt.Println("Creating artifact with path: " + artifactPath + " with " + filePath)
    //err := shelfLib.UploadArtifactFromFile(artifactPath, filePath)
    //checkError("Error creating artifact.", err)

    fmt.Println("Perfoming a HEAD request on " + basePath)
    links, err := shelfLib.ListArtifact(basePath)
    checkError("Error listing artifacts.", err)
    fmt.Println("Links:")
    fmt.Println(links)

    fmt.Println("Perfoming a POST _search request on " + basePath)
    criteria := &shelflib.SearchCriteria{}
    links, err = shelfLib.Search(basePath, criteria)
    checkError("Error listing artifacts.", err)
    fmt.Println("Links:")
    fmt.Println(links)
}
