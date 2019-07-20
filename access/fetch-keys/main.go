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
	ghOrganization := parser.String("o", "org", &argparse.Options{Required: false, Help: "Github Organization name", Default: "OpenFest"})
	fetchPgp := parser.Flag("p", "pgp-keys", &argparse.Options{Required: false, Help: "Fetch configured PGP key ids", Default: false})
	quiet := parser.Flag("q", "quiet", &argparse.Options{Required: false, Help: "Skip output to stdout", Default: false})
	ghTeam := parser.String("t", "team", &argparse.Options{Required: false, Help: "Github Team name, 'all' for all members of the organization", Default: "NOC"})
	verbose := parser.Flag("v", "verbose", &argparse.Options{Required: false, Help: "Verbose output: print keys to stdout", Default: false})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *verbose == false {
		fi, err := os.Lstat(*keysDir)
		if err != nil {
			fmt.Println(*keysDir, ": target directory error :", err)
			os.Exit(1)
		} else {
			if fi.Mode().IsDir() == false {
				fmt.Println(*keysDir, ": target directory error : not a directory - ", fi.Mode())
				os.Exit(1)
			}
		}
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *authToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	teamMembers := fetchUsers(client, ghOrganization, ghTeam)
	for _, user := range teamMembers {
		if *quiet == false {
			fmt.Println("Fetching keys for", *user)
		}
		var sshKeys bytes.Buffer
		var pgpKeys bytes.Buffer

		for nextPage := 0; ; {

			// list all teams an org for the current user
			opt := &github.ListOptions{nextPage, 50}
			keys, rsp, err := client.Users.ListKeys(ctx, *user, opt)

			if err != nil {
				fmt.Println("client.Users.ListKeys error: ", err)
				os.Exit(-1)
			}

			if rsp == nil {
				fmt.Println("Users.ListKeys returned empty response: ", err)
			}

			for _, key := range keys {
				if *verbose == true {
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

		if *quiet == false && *verbose == false {
			fmt.Println("Writing to", *keysDir+"/"+*user+".key")
		}
		err := ioutil.WriteFile(*keysDir+"/"+*user+".key", sshKeys.Bytes(), 0444)
		if err != nil && *verbose == false {
			fmt.Println(*user+".key error ", err)
		}

		if *fetchPgp == true {
			for nextPage := 0; ; {

				// list all teams an org for the current user
				opt := &github.ListOptions{nextPage, 50}
				keys, rsp, err := client.Users.ListGPGKeys(ctx, *user, opt)

				if err != nil {
					fmt.Println("client.Users.ListGPGKeys error: ", err)
					os.Exit(-1)
				}

				if rsp == nil {
					fmt.Println("Users.ListGPGKeys returned empty response: ", err)
				}

				for _, key := range keys {
					if *verbose == true {
						fmt.Println(*key.KeyID)
					}
					pgpKeys.WriteString(*key.KeyID)
					pgpKeys.WriteString("\n")
				}

				if rsp.NextPage == 0 || nextPage == rsp.NextPage {
					break
				}
				nextPage = rsp.NextPage
			}
			if *quiet == false && *verbose == false {
				fmt.Println("Writing to", *keysDir+"/"+*user+".gpg")
			}
			err = ioutil.WriteFile(*keysDir+"/"+*user+".gpg", pgpKeys.Bytes(), 0444)
			if err != nil && *verbose == false {
				fmt.Println(*user+".gpg error ", err)
			}
		}
	}

	os.Exit(0)
}

func fetchUsers(client *github.Client, org *string, team *string) (teamMembers []*string) {
	var targetTeam *github.Team

	if team == nil || *team == "all" {
		for nextPage := 0; ; {
			// list all members for the given organization's team
			opt := &github.ListMembersOptions{
				PublicOnly:  false,
				ListOptions: github.ListOptions{nextPage, 50},
			}

			users, rsp, err := client.Organizations.ListMembers(context.Background(), *org, opt)

			if err != nil {
				fmt.Println("client.Organizations.ListMembers ", err)
				os.Exit(-1)
			}

			if rsp == nil {
				fmt.Println("client.Organizations.ListMembers: ", err)
			}

			for _, user := range users {
				teamMembers = append(teamMembers, user.Login)
			}

			if rsp.NextPage == 0 || nextPage == rsp.NextPage {
				break
			}
			nextPage = rsp.NextPage
		}
	} else {
		for nextPage := 0; ; {
			// list all teams for the specified org
			opt := &github.ListOptions{nextPage, 50}
			teams, rsp, err := client.Teams.ListTeams(context.Background(), *org, opt)

			if err != nil {
				fmt.Println("client.ListTeams error: ", err)
				os.Exit(-1)
			}

			if rsp == nil {
				fmt.Println("client.ListTeams returned empty response: ", err)
			}

			for _, ghTeam := range teams {
				if *ghTeam.Name == *team {
					targetTeam = ghTeam
					break
				}
			}

			if rsp.NextPage == 0 || nextPage == rsp.NextPage {
				break
			}
			nextPage = rsp.NextPage
		}

		if targetTeam == nil {
			fmt.Println(*team, " team not found in ", *org)
			os.Exit(2)
		}

		for nextPage := 0; ; {
			// list all members for the given organization's team
			opt := &github.TeamListTeamMembersOptions{
				Role:        "all",
				ListOptions: github.ListOptions{nextPage, 50},
			}

			users, rsp, err := client.Teams.ListTeamMembers(context.Background(), *targetTeam.ID, opt)

			if err != nil {
				fmt.Println("client.Teams.ListTeamMembers ", err)
				os.Exit(-1)
			}

			if rsp == nil {
				fmt.Println("client.Teams.ListTeamMembers: ", err)
			}

			for _, user := range users {
				teamMembers = append(teamMembers, user.Login)
			}

			if rsp.NextPage == 0 || nextPage == rsp.NextPage {
				break
			}
			nextPage = rsp.NextPage
		}
	}

	return teamMembers
}
