package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/google/go-github/v32/github"
	"github.com/p1ass/sodium"
	"golang.org/x/oauth2"
)

func main() {
	var (
		owner = flag.String("owner", "", "Repository Owner")
		repo  = flag.String("repo", "", "Repository Name")
	)
	flag.Parse()
	if owner == nil || *owner == "" {
		fmt.Println("should be passed repository owner by -owner=")
		os.Exit(1)
		return
	}
	if repo == nil || *repo == "" {
		fmt.Println("should be passed repository name by -repo=")
		os.Exit(1)
		return
	}
	fmt.Println(fmt.Sprintf("Repository: %s/%s", *owner, *repo))

	secrets, err := godotenv.Read()
	if err != nil {
		fmt.Println(".env file not found: ", err.Error())
		os.Exit(1)
		return
	}

	cli := newClient()
	pubKey, err := getPublicKey(cli, *owner, *repo)
	if err != nil {
		fmt.Println("failed to get public key:", err.Error())
		os.Exit(1)
		return
	}

	wg := sync.WaitGroup{}

	for name, secret := range secrets {
		name := name
		secret := secret
		wg.Add(1)
		go func() {
			if err := updateSecret(cli, pubKey, *owner, *repo, name, secret); err != nil {
				fmt.Println("failed to update secret:", err.Error())
				os.Exit(1)
				return
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// ref https://qiita.com/kazz187/items/aa9885bb968722ac9b1d
func updateSecret(cli *github.Client, pubKey *github.PublicKey, owner, repo, name, secret string) error {
	pkBase64 := pubKey.GetKey()
	pk := make([]byte, 32)
	_, err := base64.StdEncoding.Decode(pk, []byte(pkBase64))
	if err != nil {
		fmt.Println("base64 decode error:", pkBase64, err.Error())
		return err
	}

	// sodium で暗号化を施す
	encSec := sodium.Bytes(secret).SealedBox(sodium.BoxPublicKey{Bytes: pk})
	encSecBase64 := base64.StdEncoding.EncodeToString(encSec)

	// Secret を更新
	es := &github.EncryptedSecret{
		Name:           name,
		KeyID:          *pubKey.KeyID,
		EncryptedValue: encSecBase64,
	}
	_, err = cli.Actions.CreateOrUpdateRepoSecret(context.Background(), owner, repo, es)
	if err != nil {
		fmt.Println("failed to update secret", err.Error())
		return err
	}
	return nil
}

func newClient() *github.Client {
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	oc := oauth2.NewClient(context.Background(), sts)
	cli := github.NewClient(oc)
	return cli
}

func getPublicKey(cli *github.Client, owner string, repo string) (*github.PublicKey, error) {
	pubKey, _, err := cli.Actions.GetRepoPublicKey(context.Background(), owner, repo)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}
