package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"
)

const LINE_SEPARATOR = "-------------------------------------"

func main() {
	printBanner()
	log.Println(strings.Join(Run(os.Stdin), "\n"))
}

func printBanner() bool {
	fmt.Println(LINE_SEPARATOR)
	fmt.Println("- Extract Github data - Application -")
	fmt.Println(LINE_SEPARATOR)
	fmt.Println("- Put your repo list with format    -")
	fmt.Println("- @org/@repo, 1 repo for each line  -")
	fmt.Println("- end with blank line               -")
	fmt.Println("- Require : GITHUB_TOKEN as env     -")
	fmt.Println(LINE_SEPARATOR)
	fmt.Println("")
	return true
}

// Run method to run the application
func Run(stdin io.Reader) []string {
	var listOfError []string
	_, err := printArrayToStdout(extractGithubDataAndTransform(readStdin(stdin)))
	for _, errorString := range err {
		listOfError = append(listOfError, errorString.Error())
	}
	return listOfError
}

func extractGithubDataAndTransform(inputArrayOfString []string, carryOverError []error) ([]string, []error) {
	if len(inputArrayOfString) > 0 {
		transformedArrayData := []string{"Repo Name,Clone URL,Last Pushed,Owner"}
		for _, repoData := range inputArrayOfString {
			if strings.Count(repoData, "/") != 1 {
				carryOverError = append(carryOverError, fmt.Errorf("wrong input data => %s", repoData))
				continue
			}
			ctx := context.Background()
			tokenSource := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
			)
			client := github.NewClient(oauth2.NewClient(ctx, tokenSource))
			repo, _, err := client.Repositories.Get(ctx, strings.Split(repoData, "/")[0], strings.Split(repoData, "/")[1])
			if err != nil {
				carryOverError = append(carryOverError, err)
			} else {
				transformedArrayData = append(transformedArrayData, fmt.Sprintf("%s,%s,%s,%s", repo.GetName(), repo.GetCloneURL(), repo.GetPushedAt(), repo.GetOwner().GetLogin()))
			}
		}
		if len(transformedArrayData) > 1 {
			return transformedArrayData, carryOverError
		}
		return []string{}, carryOverError
	}
	return inputArrayOfString, carryOverError

}

func printArrayToStdout(inputArrayOfString []string, carryOverError []error) ([]string, []error) {
	if len(inputArrayOfString) > 0 {
		fmt.Printf("========= Result =========\n\n")
		fmt.Printf("%s\n", strings.Join(inputArrayOfString, "\n"))
		fmt.Printf("\n========= End Of Result =========\n")
	}
	return []string{}, carryOverError
}

func readStdin(stdin io.Reader) ([]string, []error) {
	scanner := bufio.NewScanner(stdin)
	var arrayOfStdinData []string
	for scanner.Scan() {
		inputLine := scanner.Text()
		if inputLine == "" {
			break
		}
		arrayOfStdinData = append(arrayOfStdinData, inputLine)
	}
	return arrayOfStdinData, []error{}
}
