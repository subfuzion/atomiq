package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/appcelerator/amp/api/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var (
	ampConfigFolder = filepath.Join(".config", "amp")
	ampTokenFile    = "token"
)

// SaveToken saves the authentication token to file
func SaveToken(metadata metadata.MD) error {
	// Extract token from header
	tokens := metadata[auth.TokenKey]
	if len(tokens) == 0 {
		return fmt.Errorf("invalid token")
	}
	token := tokens[0]
	if token == "" {
		return fmt.Errorf("invalid token")
	}

	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot get current user")
	}
	if err := os.MkdirAll(filepath.Join(usr.HomeDir, ampConfigFolder), os.ModePerm); err != nil {
		return fmt.Errorf("cannot create folder")
	}
	if err := ioutil.WriteFile(filepath.Join(usr.HomeDir, ampConfigFolder, ampTokenFile), []byte(token), 0600); err != nil {
		return fmt.Errorf("cannot write token")
	}
	return nil
}

// ReadToken reads the authentication token from file
func ReadToken() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("cannot get current user")
	}
	data, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ampConfigFolder, ampTokenFile))
	if err != nil {
		return "", fmt.Errorf("cannot read token")
	}
	return string(data), nil
}

// RemoveToken deletes the authentication token file
func RemoveToken() error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot get current user")
	}
	filePath := filepath.Join(usr.HomeDir, ampConfigFolder, ampTokenFile)
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("you are not logged in. Use `amp login` or `amp user signup`")
	}
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("cannot remove token")
	}
	return nil
}

// LoginCredentials represents login credentials
type LoginCredentials struct {
	Token string
}

// GetRequestMetadata implements credentials.PerRPCCredentials
func (c *LoginCredentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		auth.TokenKey: c.Token,
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials
func (c *LoginCredentials) RequireTransportSecurity() bool {
	return false
}

// GetToken returns the stored token
func GetToken() string {
	token, err := ReadToken()
	if err != nil {
		return ""
	}
	return token
}
