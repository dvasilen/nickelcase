package uri

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/giacomocariello/nickelcase/crypt"
)

type WriteDataToStream func(io.WriteCloser, []byte) error

func WriteDataToPlaintextStream(stream io.WriteCloser, data []byte) error {
	defer stream.Close()
	if _, err := stream.Write(data); err != nil {
		return err
	}
	return nil
}

func WriteDataToEncryptedStream(password string) WriteDataToStream {
	return func(stream io.WriteCloser, data []byte) error {
		defer stream.Close()
		encryptedData, err := crypt.AnsibleEncrypt(data, password)
		if err != nil {
			return err
		}
		_, err = stream.Write(encryptedData)
		return err
	}
}

func WriteMapToStream(stream io.WriteCloser, ret map[string]interface{}, fn WriteDataToStream) error {
	data, err := yaml.Marshal(&ret)
	if err != nil {
		return err
	}
	return fn(stream, data)
}

func WriteMapToURI(c *cli.Context, uri string, ret map[string]interface{}, fn WriteDataToStream) error {
	if uri == "" || uri == "-" {
		return WriteMapToStream(os.Stdout, ret, fn)
	}
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		return err
	}
	switch parsedUrl.Scheme {
	case "":
		fh, err := os.OpenFile(parsedUrl.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		return WriteMapToStream(fh, ret, fn)
	case "fd":
		fd, err := strconv.Atoi(parsedUrl.Opaque)
		if err != nil {
			return fmt.Errorf("Invalid file descriptor number: %s", err)
		}
		return WriteMapToStream(os.NewFile(uintptr(fd), parsedUrl.Opaque), ret, fn)
	default:
		return fmt.Errorf("Unsupported URI scheme in parameter: %s", uri)
	}
}

func GetOutputStreamFromURI(c *cli.Context, uri string) (io.WriteCloser, error) {
	if uri == "" || uri == "-" {
		return os.Stdout, nil
	}
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	switch parsedUrl.Scheme {
	case "":
		return os.OpenFile(parsedUrl.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	case "fd":
		fd, err := strconv.Atoi(parsedUrl.Opaque)
		if err != nil {
			return nil, fmt.Errorf("Invalid file descriptor number: %s", err)
		}
		return os.NewFile(uintptr(fd), parsedUrl.Opaque), nil
	default:
		return nil, fmt.Errorf("Unsupported URI scheme in parameter: %s", uri)
	}
}

func WriteDataToURI(c *cli.Context, uri string, data []byte, fn WriteDataToStream) error {
	stream, err := GetOutputStreamFromURI(c, uri)
	if err != nil {
		return err
	}
	return fn(stream, data)
}
