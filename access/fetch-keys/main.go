/*-
 * ----------------------------------------------------------------------------
 * "THE BEER-WARE LICENSE" (Revision 69):
 * <syrinx@FreeBSD.org> wrote this file.  As long as you retain this notice you
 * can do whatever you want with this stuff. If we meet some day, and you think
 * this stuff is worth it, you can buy me a beer in return.   -Shteryana Shopova
 * ----------------------------------------------------------------------------
 */

package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/google/go-github/github" // with go modules disabled
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
)

func main() {

	parser := argparse.NewParser("sshkeys", "Fetch SSH keys for a Github team members")
	authToken := parser.String("a", "authtoken", &argparse.Options{Required: false, Help: "Github Auth token", Default: AuthToken})
	keysDir := parser.String("d", "directory", &argparse.Options{Required: false, Help: "Path where to store the key files", Default: "./"})
	quiet := parser.Flag("q", "quiet", &argparse.Options{Required: false, Help: "Skip output to stdout", Default: false})
	ghOrganization := parser.String("o", "org", &argparse.Options{Required: false, Help: "Github Organization name", Default: "OpenFest"})
	ghTeam := parser.String("t", "team", &argparse.Options{Required: false, Help: "Github Team name", Default: "NOC"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	fi, err := os.Lstat(*keysDir)
	if err != nil  {
		fmt.Println(*keysDir, ": target directory error :", err)
		os.Exit(1)
	} else {
		if fi.Mode().IsDir() == false {
			fmt.Println(*keysDir, ": target directory error : not a directory - ", fi.Mode())
			os.Exit(1)
		}
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *authToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var nocTeam *github.Team
	var nocMembers []*string

	for nextPage := 0; ; {
		// list all teams for the specified org
		opt := &github.ListOptions{nextPage, 50}
		teams, rsp, err := client.Teams.ListTeams(ctx, *ghOrganization, opt)

		if err != nil {
			fmt.Println("client.ListTeams error: %v", err)
			os.Exit(-1)
		}

		if rsp == nil {
			fmt.Println("client.Repositories.List returned empty response: %v", err)
		}

		for _, team := range teams {
			if *team.Name == *ghTeam {
				nocTeam = team
				break
			}
		}

		if rsp.NextPage == 0 || nextPage == rsp.NextPage {
			break
		}
		nextPage = rsp.NextPage
	}

	if nocTeam == nil {
		fmt.Println("NOC team not found in OpenFest")
		os.Exit(2)
	}

	for nextPage := 0; ; {
		// list all members for the given organization's team
		opt := &github.TeamListTeamMembersOptions{
			Role:        "all",
			ListOptions: github.ListOptions{nextPage, 50},
		}

		users, rsp, err := client.Teams.ListTeamMembers(context.Background(), *nocTeam.ID, opt)

		if err != nil {
			fmt.Println("client.Teams.ListTeamMembers %v", err)
			os.Exit(-1)
		}

		if rsp == nil {
			fmt.Println("client.Teams.ListTeamMembers: %v", err)
		}

		for _, user := range users {
			nocMembers = append(nocMembers, user.Login)
		}

		if rsp.NextPage == 0 || nextPage == rsp.NextPage {
			break
		}
		nextPage = rsp.NextPage
	}

	for _, user := range nocMembers {
		if *quiet == false {
			fmt.Println("Fetching keys for ", *user)
		}
		var sshKeys bytes.Buffer

		for nextPage := 0; ; {

			// list all teams an org for the current user
			opt := &github.ListOptions{nextPage, 50}
			keys, rsp, err := client.Users.ListKeys(ctx, *user, opt)

			if err != nil {
				fmt.Println("client.Users.ListKeys error: %v", err)
				os.Exit(-1)
			}

			if rsp == nil {
				fmt.Println("Users.ListKeyss returned empty response: %v", err)
			}

			for _, key := range keys {
				if *quiet == false {
					fmt.Println(*key.Key)
				}
				sshKeys.WriteString(*key.Key)
				sshKeys.WriteString("\n")
			}

			if rsp.NextPage == 0 || nextPage == rsp.NextPage {
				break
			}
			nextPage = rsp.NextPage
		}

		if *quiet == false {
			fmt.Println("Writing to", *keysDir + "/" + *user+".key")
		}
		err := ioutil.WriteFile(*keysDir + "/" + *user+".key", sshKeys.Bytes(), 0444)
		if err != nil {
			fmt.Println(*user+".key error %v", err)
		}
	}

	os.Exit(0)
}
