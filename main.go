package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/crypto/nacl/box"
)

const (
	githubAPIURL = "https://api.github.com"
	userAgent    = "file-secret-action/v0.1.0 (https://github.com/nicklasfrahm/file-secret-action)"
)

type gitHubPublicKey struct {
	KeyID string `json:"key_id"`
	Key   string `json:"key"`
}

type gitHubSecret struct {
	EncryptedValue string `json:"encrypted_value"`
	KeyID          string `json:"key_id"`
	Visibility     string `json:"visibility,omitempty"`
}

func main() {
	// Get repository public key.
	scope := os.Getenv("SCOPE")

	// Detect if the scope is a repository or an org.
	if strings.Contains(scope, "/") {
		// Scope is a combination of username and repository, so a repository .
		scope = "/repos/" + scope
	} else {
		// Scope is just an organization name.
		scope = "/orgs/" + scope
	}

	// Fetch public key.
	publicKeyEndpoint := "/actions/secrets/public-key"
	resp, err := RequestGitHubAPI("GET", scope+publicKeyEndpoint, nil)
	if err != nil {
		log.Fatalf("Failed to fetch public key: %v", err)
	}
	log.Println(scope + publicKeyEndpoint)

	// Parse JSON response.
	gitHubKey := new(gitHubPublicKey)
	err = json.NewDecoder(resp.Body).Decode(gitHubKey)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}
	resp.Body.Close()

	// Decode base64 encoded string to bytes.
	pubKeySlice, err := base64.StdEncoding.DecodeString(gitHubKey.Key)
	if err != nil {
		log.Fatalf("Failed to decoded base64-encoded public key: %v", err)
	}

	// Create fixed size public key array buffer.
	var pubKey [32]byte
	copy(pubKey[:], pubKeySlice)

	// Read file.
	file := os.Getenv("FILE")
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Encrypt and base64-encode the encrypted file content.
	encryptedBytes, err := box.SealAnonymous(nil, fileBytes, &pubKey, nil)
	if err != nil {
		log.Fatalf("Failed to encrypt secret: %v", err)
	}
	encrypted := base64.StdEncoding.EncodeToString(encryptedBytes)

	// Create the request payload.
	secret := gitHubSecret{
		KeyID:          gitHubKey.KeyID,
		EncryptedValue: encrypted,
		Visibility:     os.Getenv("VISIBILITY"),
	}
	secretBytes, err := json.Marshal(&secret)
	if err != nil {
		log.Fatalf("Failed to encode secret to JSON: %v", err)
	}

	// Check if secret name is valid.
	secretName := os.Getenv("SECRET")
	if secretName == "" {
		log.Fatalf("Failed to create secret: %v", errors.New("secret name must not be empty"))
	}

	// Create or update the secret.
	secretEndpoint := "/actions/secrets/" + os.Getenv("SECRET")
	resp, err = RequestGitHubAPI("PUT", scope+secretEndpoint, bytes.NewReader(secretBytes))
	if err != nil {
		log.Fatalf("Failed to create secret: %v", err)
	}
	log.Println(scope + secretEndpoint)

	action := "❌ Failed to create"
	if resp.StatusCode == http.StatusCreated {
		action = "🔑 Created"
	}
	if resp.StatusCode == http.StatusNoContent {
		action = "🔑 Updated"
	}

	log.Printf("%s secret: %s %d\n", action, secretName, resp.StatusCode)
}

// RequestGitHubAPI makes a request against the GitHub API.
func RequestGitHubAPI(verb string, path string, body io.Reader) (*http.Response, error) {

	// This must be a personal access token for organization secrets
	// or a GITHUB_TOKEN for repository secrets.
	token := os.Getenv("TOKEN")

	// Parse GitHub API URL.
	u, err := url.Parse(githubAPIURL)
	if err != nil {
		log.Fatalf("Failed to parse API URL: %v", err)
	}

	// Create HTTP client config.
	req, err := http.NewRequest(verb, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", userAgent)
	req.URL = u

	return http.DefaultClient.Do(req)
}
