A tool to generate Markdown summaries of Go coverage results.

Useful for generating coverage summaries for Github Actions build summaries, PR comments.

	go install github.com/axw/gocov/gocov@latest
    go install github.com/ct1n/gocov-md@latest

## Usage

    go test -coverprofile=cov.out
    gocov convert cov.out | gocov-md >cov.md
