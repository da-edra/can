package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/olekukonko/tablewriter"
)

var (
	iamClient iamiface.IAMAPI
)

func init() {
	session := session.Must(session.NewSession())
	iamClient = iam.New(session)
}

func main() {
	roleArn := flag.String("source-arn", "", "The ARN of the IAM resource to test")
	pathFile := flag.String("do", "", "Path to the file that contains a list of the Actions to test")

	flag.Parse()
	//roleArn and pathFile must be present
	if *roleArn == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *pathFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	actions, err := readActionFile(*pathFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	data, err := simulate(roleArn, actions)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	print(data)
}

func simulate(roleArn *string, actions []string) ([][]string, error) {
	input := &iam.SimulatePrincipalPolicyInput{
		ActionNames:     aws.StringSlice(actions),
		PolicySourceArn: roleArn,
	}

	var data [][]string
	err := iamClient.SimulatePrincipalPolicyPages(input,
		func(page *iam.SimulatePolicyResponse, lastPage bool) bool {
			for _, evalResult := range page.EvaluationResults {
				if *evalResult.EvalDecision == "allowed" {
					data = append(data, []string{*evalResult.EvalActionName, "✓"})
					continue
				}
				data = append(data, []string{*evalResult.EvalActionName, "×"})
			}
			return !lastPage
		})
	if err != nil {
		return [][]string{}, fmt.Errorf("## Error: Something happened running the simulation ##\n%s", err.Error())
	}
	return data, nil
}

func print(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ACTION", "ALLOWED"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

func isActionValid(action string) bool {
	//TODO: Check if AWS provides a list of valid actions
	match, _ := regexp.MatchString(".*:.*", action)
	return match
}

func readActionFile(filename string) ([]string, error) {
	var actions []string
	file, err := os.Open(filename)
	if err != nil {
		return []string{}, fmt.Errorf("## Error: Something happened opening actions file ##\n%s", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if isActionValid(scanner.Text()) {
			actions = append(actions, scanner.Text())
			continue
		}
		return []string{}, fmt.Errorf("## Error: Something happened reading actions ##\n%s is invalid", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("## Error: Something happened with the scanner ##\n%s", err.Error())
	}
	return actions, nil
}
