# can

can is a command line tool that allows you to easily check which actions from a list can the role do.

## Installation

### Using go

Execute the following command

```shell
$ go get -u github.com/wizeline-sre/can
```

### Usage

1. Set your credentials through the AWS CLI.
1. The file with the actions should contain one action per row, for example:
    ```shell
    ec2:CreateVpc
    ec2:Describe*
    logs:ListTagsLogGroup
    ```
1. To check what actions the role can do, execute:
  ```shell
  $ can -source-arn arn:aws:iam::XXXXXX:XXXX/XXXXXXXX -do actions.txt
  ```
1. A printed table will appear showing what actions the role can do.
  - ✓ means that it is allowed
  - × means it is not allowed

#### Output example

```shell
+-----------------------------------+---------+
|              ACTION               | ALLOWED |
+-----------------------------------+---------+
| ec2:CreateVpc                     | ✓       |
| ec2:DeleteVpc                     | ✓       |
| ec2:Describe*                     | ✓       |
| iam:UpdateAssumeRolePolicy        | ×       |
| logs:CreateLogGroup               | ✓       |
| logs:DescribeLogGroups            | ✓       |
| logs:DeleteLogGroup               | ✓       |
| logs:ListTagsLogGroup             | ✓       |
| logs:PutRetentionPolicy           | ×       |
+-----------------------------------+---------+
```

### Restrictions and considerations

- You'll need to have `iam:SimulatePrincipalPolicy` to be able to run `can` successfully.

- Currently, `can` does not support context keys.