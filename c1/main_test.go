package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitReadFromStdin(t *testing.T) {
	var stdin bytes.Buffer
	mockData := []string{"test1", "test2"}

	stdin.Write([]byte(fmt.Sprintf("%s\n\n", strings.Join(mockData, "\n"))))

	result, err := readStdin(&stdin)
	assert.Equal(t, []error{}, err)
	assert.Equal(t, mockData, result)
}

func TestUnitPrintArrayToStdoutWithNoCarryOverError(t *testing.T) {
	mockData := []string{"test1", "test2"}
	carryOverError := []error{}
	result, err := printArrayToStdout(mockData, carryOverError)
	assert.Equal(t, []error{}, err)
	assert.Equal(t, []string{}, result)
}

func TestUnitPrintArrayToStdoutWithCarryOverError(t *testing.T) {
	mockData := []string{"test1", "test2"}
	carryOverError := []error{fmt.Errorf("Test Error")}
	result, err := printArrayToStdout(mockData, carryOverError)
	assert.Equal(t, carryOverError, err)
	assert.Equal(t, []string{}, result)
}

func TestIntegrationExtractGithubDataAndTransformWithNoCarryOverError(t *testing.T) {
	inputData := []string{"marve39/it-management", "marve39/it-management2"}
	resultData := []string{"Repo Name,Clone URL,Last Pushed,Owner", "it-management,https://github.com/marve39/it-management.git,2016-09-29 03:02:01 +0000 UTC,marve39", "it-management2,https://github.com/marve39/it-management2.git,2016-10-11 05:28:13 +0000 UTC,marve39"}
	carryOverError := []error{}
	result, err := extractGithubDataAndTransform(inputData, carryOverError)
	assert.Equal(t, carryOverError, err)
	assert.Equal(t, resultData, result)
}

func TestIntegrationExtractGithubDataAndTransformWithCarryOverError(t *testing.T) {
	inputData := []string{"marve39/it-management", "marve39/it-management2"}
	resultData := []string{"Repo Name,Clone URL,Last Pushed,Owner", "it-management,https://github.com/marve39/it-management.git,2016-09-29 03:02:01 +0000 UTC,marve39", "it-management2,https://github.com/marve39/it-management2.git,2016-10-11 05:28:13 +0000 UTC,marve39"}
	carryOverError := []error{fmt.Errorf("Test Error")}
	result, err := extractGithubDataAndTransform(inputData, carryOverError)
	assert.Equal(t, carryOverError, err)
	assert.Equal(t, resultData, result)
}

func TestIntegrationExtractGithubDataAndTransformWithOneInvalidRepo(t *testing.T) {
	inputData := []string{"marve39/it-management", "marve39/test"}
	resultData := []string{"Repo Name,Clone URL,Last Pushed,Owner", "it-management,https://github.com/marve39/it-management.git,2016-09-29 03:02:01 +0000 UTC,marve39"}
	carryOverError := []error{}
	result, err := extractGithubDataAndTransform(inputData, carryOverError)
	assert.Equal(t, 1, len(err))
	assert.Equal(t, resultData, result)
}

func TestIntegrationExtractGithubDataAndTransformWithInvalidInputTypeOneRepo(t *testing.T) {
	inputData := []string{"marve39"}
	resultData := []string{}
	carryOverError := []error{}
	result, err := extractGithubDataAndTransform(inputData, carryOverError)
	assert.Equal(t, 1, len(err))
	assert.Equal(t, resultData, result)
}

func TestIntegrationExtractGithubDataAndTransformWithInvalidInputTypeMoreThanOneRepo(t *testing.T) {
	inputData := []string{"marve39", "marve39/it-management"}
	resultData := []string{"Repo Name,Clone URL,Last Pushed,Owner", "it-management,https://github.com/marve39/it-management.git,2016-09-29 03:02:01 +0000 UTC,marve39"}
	carryOverError := []error{}
	result, err := extractGithubDataAndTransform(inputData, carryOverError)
	assert.Equal(t, 1, len(err))
	assert.Equal(t, resultData, result)
}

func TestIntegrationExtractGithubDataAndTransformWithNoRepo(t *testing.T) {
	inputData := []string{}
	resultData := []string{}
	carryOverError := []error{}
	result, err := extractGithubDataAndTransform(inputData, carryOverError)
	assert.Equal(t, carryOverError, err)
	assert.Equal(t, resultData, result)
}

func TestIntegrationRunApplication(t *testing.T) {
	var stdin bytes.Buffer
	mockData := []string{"marve39/it-management"}
	stdin.Write([]byte(fmt.Sprintf("%s\n\n", strings.Join(mockData, "\n"))))

	result := Run(&stdin)
	fmt.Println(result)
	assert.Equal(t, 0, len(result))
}

func TestIntegrationMain(t *testing.T) {
	main()
	assert.Equal(t, true, true)
}
func TestUnitPrintBanner(t *testing.T) {
	result := printBanner()
	assert.Equal(t, true, result)
}

func TestIntegrationExtractGithubDataAndTransformWithNoGithubToken(t *testing.T) {
	inputData := []string{"marve39/it-management"}
	resultData := []string{}
	carryOverError := []error{}
	os.Unsetenv("GITHUB_TOKEN")
	result, err := extractGithubDataAndTransform(inputData, carryOverError)
	assert.Equal(t, 1, len(err))
	assert.Equal(t, resultData, result)
}
