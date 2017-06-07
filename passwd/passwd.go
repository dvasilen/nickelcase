package passwd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/nbutton23/zxcvbn-go"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/giacomocariello/nickelcase/uri"
)


// GetPassword : read current password from TTY or URI
func GetPassword(c *cli.Context, src string) (password string, err error) {
	fdout := getOutputTTY()
	if src == "" {
		password, err = readPasswordFromTTY("Password: ", fdout)
	} else {
		password, err = uri.ReadFromPasswordURI(c, src)
	}
	return
}

// GetPassword : read new password from TTY or URI
func GetNewPassword(c *cli.Context, src string) (password string, err error) {
	fdout := getOutputTTY()
	if src == "" {
		password, err = readNewPasswordFromTTY(fdout)
	} else {
		password, err = uri.ReadFromPasswordURI(c, src)
	}
	return
}

// GetPassword : read password change from TTY or URI
func GetPasswordChange(c *cli.Context, oldSrc, newSrc string) (oldPassword string, newPassword string, err error) {
	fdout := getOutputTTY()
	if oldSrc == "" {
		oldPassword, err = readPasswordFromTTY("Old password: ", fdout)
	} else {
		oldPassword, err = uri.ReadFromPasswordURI(c, oldSrc)
	}
	if err != nil {
		return
	}
	if newSrc == "" {
		newPassword, err = readNewPasswordFromTTY(fdout)
	} else {
		newPassword, err = uri.ReadFromPasswordURI(c, newSrc)
	}
	return
}

func getOutputTTY() int {
	if terminal.IsTerminal(syscall.Stdout) {
		return syscall.Stdout
	} else if terminal.IsTerminal(syscall.Stderr) {
		return syscall.Stderr
	}
	return -1
}

func readPasswordFromTTY(prompt string, fdout int) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("Cannot prompt for password: cannot open /dev/tty: %s", err)
	}
	defer tty.Close()
	if fdout > 0 {
		if _, err := syscall.Write(fdout, []byte(prompt)); err != nil {
			return "", err
		}
		if err := syscall.Fsync(fdout); err != nil {
			return "", err
		}
	}
	password, err := terminal.ReadPassword(int(tty.Fd()))
	if err != nil {
		return "", err
	}
	if fdout > 0 {
		if _, err := syscall.Write(fdout, []byte("\n")); err != nil {
			return "", err
		}
	}
	return strings.TrimRight(string(password), "\r\n"), nil
}

func readNewPasswordFromTTY(fdout int) (string, error) {
	i := 0
	for {
		newPassword, err := readPasswordFromTTY("New password: ", fdout)
		if err != nil {
			return "", err
		}
		retypePassword, err := readPasswordFromTTY("Retype password: ", fdout)
		if err != nil {
			return "", err
		}
		if newPassword != retypePassword {
			if _, err := syscall.Write(fdout, []byte("Error: Password mismatch, please retry.\n")); err != nil {
				return "", err
			}
			continue
		}
		i += 1
		if score := zxcvbn.PasswordStrength(newPassword, []string{}); score.Score < 3 && i < 2 {
			if _, err := syscall.Write(fdout, []byte("Error: Password is too easy to guess, please retry.\n")); err != nil {
				return "", err
			}
			continue
		}
		return newPassword, nil
	}
}
