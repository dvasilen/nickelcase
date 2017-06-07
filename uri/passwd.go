package uri

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
)

func ReadPasswordStream(stream io.ReadCloser) (string, error) {
	defer stream.Close()
	scanner := bufio.NewScanner(stream)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
	}
	return scanner.Text(), nil
}

func SaveToPasswordURISource(uri string) (string, error) {
	if uri == "" || uri == "-" {
		return ReadPasswordStream(os.Stdin)
	}
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	switch parsedUrl.Scheme {
	case "":
		fh, err := os.OpenFile(parsedUrl.Path, os.O_RDONLY, 0)
		if err != nil {
			return "", err
		}
		return ReadPasswordStream(fh)
	case "fd":
		fd, err := strconv.Atoi(parsedUrl.Opaque)
		if err != nil {
			return "", fmt.Errorf("Invalid file descriptor number: %s", err)
		}
		return ReadPasswordStream(os.NewFile(uintptr(fd), parsedUrl.Opaque))
	case "env":
		return ReadPasswordStream(ioutil.NopCloser(bytes.NewBufferString(os.Getenv(parsedUrl.Opaque))))
	default:
		return "", fmt.Errorf("Unsupported URI scheme in parameter: %s", uri)
	}
}
