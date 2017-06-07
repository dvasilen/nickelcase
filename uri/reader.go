package uri

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/giacomocariello/nickelcase/crypt"
)

type ReadStream func(io.ReadCloser) ([]byte, error)

func ReadDataFromEncryptedStream(password string) ReadStream {
	return func(stream io.ReadCloser) ([]byte, error) {
		defer stream.Close()
		encryptedData, err := ioutil.ReadAll(stream)
		if err != nil {
			return []byte{}, err
		}
		return crypt.AnsibleDecrypt(encryptedData, password)
	}
}

func ReadDataFromPlaintextStream() ReadStream {
	return func(stream io.ReadCloser) ([]byte, error) {
		defer stream.Close()
		return ioutil.ReadAll(stream)
	}
}

func ReadMapFromURI(c *cli.Context, uri string, fnStream ReadStream, ret map[string]interface{}) error {
	var err error
	var data []byte
	if uri == "" || uri == "-" {
		data, err = fnStream(os.Stdin)
	} else {
		parsedUrl, err := url.Parse(uri)
		if err != nil {
			return err
		}
		switch parsedUrl.Scheme {
		case "":
			fh, err := os.OpenFile(parsedUrl.Path, os.O_RDONLY, 0)
			if err != nil {
				return err
			}
			data, err = fnStream(fh)
		case "fd":
			fd, err := strconv.Atoi(parsedUrl.Opaque)
			if err != nil {
				return fmt.Errorf("Invalid file descriptor number: %s", err)
			}
			data, err = fnStream(os.NewFile(uintptr(fd), parsedUrl.Opaque))
		case "env":
			data, err = fnStream(ioutil.NopCloser(bytes.NewBufferString(os.Getenv(parsedUrl.Opaque))))
		case "http", "https":
			httpClient := GetHTTPClient(c)
			url := parsedUrl.String()
			response, err := httpClient.Get(url)
			if err != nil {
				return fmt.Errorf("Error while trying to retrieve url \"%s\": %s", url, err)
			}
			data, err = fnStream(response.Body)
		default:
			return fmt.Errorf("Unsupported URI scheme in parameter: %s", uri)
		}
	}
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &ret)
	if err != nil {
		return err
	}
	return nil
}

func ReadDataFromURI(c *cli.Context, uri string, fnStream ReadStream) ([]byte, error) {
	if uri == "" || uri == "-" {
		return fnStream(os.Stdin)
	} else {
		parsedUrl, err := url.Parse(uri)
		if err != nil {
			return []byte{}, err
		}
		switch parsedUrl.Scheme {
		case "":
			fh, err := os.OpenFile(parsedUrl.Path, os.O_RDONLY, 0)
			if err != nil {
				return []byte{}, err
			}
			return fnStream(fh)
		case "fd":
			fd, err := strconv.Atoi(parsedUrl.Opaque)
			if err != nil {
				return []byte{}, fmt.Errorf("Invalid file descriptor number: %s", err)
			}
			return fnStream(os.NewFile(uintptr(fd), parsedUrl.Opaque))
		case "env":
			return fnStream(ioutil.NopCloser(bytes.NewBufferString(os.Getenv(parsedUrl.Opaque))))
		case "http", "https":
			httpClient := GetHTTPClient(c)
			url := parsedUrl.String()
			response, err := httpClient.Get(url)
			if err != nil {
				return []byte{}, fmt.Errorf("Error while trying to retrieve url \"%s\": %s", url, err)
			}
			return fnStream(response.Body)
		default:
			return []byte{}, fmt.Errorf("Unsupported URI scheme in parameter: %s", uri)
		}
	}
}
