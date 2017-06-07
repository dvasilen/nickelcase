package passwd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/nbutton23/zxcvbn-go"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/giacomocariello/nickelcase/uri"
)

func GetPassword(src string) (password string, err error) {
	fdout := GetOutputStream()
	if src == "" {
		password, err = ReadPasswordFromTTY("Password: ", fdout)
	} else {
		password, err = uri.SaveToPasswordURISource(src)
	}
	return
}

func GetNewPassword(src string) (password string, err error) {
	fdout := GetOutputStream()
	if src == "" {
		password, err = ReadNewPasswordFromTTY(fdout)
	} else {
		password, err = uri.SaveToPasswordURISource(src)
	}
	return
}

func GetPasswordChange(oldSrc, newSrc string) (oldPassword string, newPassword string, err error) {
	fdout := GetOutputStream()
	if oldSrc == "" {
		oldPassword, err = ReadPasswordFromTTY("Old password: ", fdout)
	} else {
		oldPassword, err = uri.SaveToPasswordURISource(oldSrc)
	}
	if err != nil {
		return
	}
	if newSrc == "" {
		newPassword, err = ReadNewPasswordFromTTY(fdout)
	} else {
		newPassword, err = uri.SaveToPasswordURISource(newSrc)
	}
	return
}

func GetOutputStream() int {
	if terminal.IsTerminal(syscall.Stdout) {
		return syscall.Stdout
	} else if terminal.IsTerminal(syscall.Stderr) {
		return syscall.Stderr
	}
	return -1
}

func ReadPasswordFromTTY(prompt string, fdout int) (string, error) {
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

func ReadNewPasswordFromTTY(fdout int) (string, error) {
	i := 0
	for {
		newPassword, err := ReadPasswordFromTTY("New password: ", fdout)
		if err != nil {
			return "", err
		}
		retypePassword, err := ReadPasswordFromTTY("Retype password: ", fdout)
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
