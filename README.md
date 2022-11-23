Parade
======

[![release badge](https://img.shields.io/github/v/release/chroju/parade.svg)](https://github.com/chroju/parade/releases)
[![test badge](https://github.com/chroju/parade/workflows/test/badge.svg)](https://github.com/chroju/parade/actions?workflow=test)
[![Coverage Status](https://coveralls.io/repos/github/chroju/parade/badge.svg?branch=main)](https://coveralls.io/github/chroju/parade?branch=main)


Parade is a simple CLI tool for AWS SSM parameter store. Easy to read and write key values in your parameter store.

Install
-------

### Homebrew

```bash
brew install chroju/tap/parade
```

### Download binary

Download the latest binary from [here](https://github.com/chroju/parade/releases) and place it in the some directory specified by `$PATH`.

### go get

If you have set up Go environment, you can also install `parade` with `go get` command.

```bash
go get github.com/chroju/parade
```

Authentication
--------------

Parade requires your AWS IAM user authentication. The same authentication method as [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) is available. Tools like [aws-vault](https://github.com/99designs/aws-vault) can be used as well.

```bash
# with command line options
$ parade --profile YOUR_PROFILE

# with aws-vault
$ aws-vault exec YOUR_PROFILE -- parade
```

Parade uses the following AWS API. If you are dealing with SecureString, you will also need permission to access the kms key.

* ssm:DeleteParameter
* ssm:DescribeParameters
* ssm:GetParameter
* ssm:PutParameter

Usage
-----

There are simple 4 sub commands. It is similar to redis-cli. 

### keys

Search keys and display their types together.

`keys` command supports exact match, forward match, and partial match. It usually searches for exact matches.

```bash
$ parade /service1/dev/key1
/service1/dev/key1  Type: String
```

Use `*` as a postfix, the search will be done as a forward match.  Furthermore, also use `*` as a prefix, it becomes a partial match. You can't use `*` as a prefix only.

```bash
$ parade keys /service1*
/service1/dev/key1   Type: String
/service1/dev/key2   Type: String
/service1/prod/key3  Type: SecureString

$ parade keys *prod*
/service1/prod/key3  Type: SecureString
```

If no argument is given, all keys will be retrieved.

```bash
$ parade keys
/service1/dev/key1   Type: String
/service1/dev/key2   Type: String
/service1/prod/key3  Type: SecureString
...
```

### get

Display the value of the specified key.

```bash
$ parade get /service1/dev/key1
value1
```

You can also do a partial search using `*` as well as the `keys` command.

```bash
$ parade get /service1*
/service1/dev/key1   value1
/service1/dev/key2   value2
/service1/prod/key3  value3
```

The `--decrypt` or `-d` option is required to decrypt SecureString.

```bash
$ parade get /service1/dev/password
(encrypted)

$ parade get /service1/dev/password -d
1234password
```

### set

Set new key and value.

```bash
$ parade set /service1/dev/key4 value4
```

If the specified key already exists, you can choose to overwrite it. Use `--force` flag if you want to force overwriting.

```bash
$ parade set /service1/dev/key4 value5
WARN: `/service1/dev/key4` already exists.
Overwrite `/service1/dev/key4` (value: value4) ? (Y/n)

$ parade set /service1/dev/key4 value5 --force
```

The value is stored as `String` type by default. It also supports `SecureString` with the `--encrypt` flag. If you don't specify a key ID with `--kms-key-id` flag, uses the default key associated with your AWS account.

`StringList` type is not supported.


### del

Delete key and value. Use `--force` flag if you want to skip the confirmation prompt.

```bash
$ parade del /service1/dev/key4
Delete `/service1/dev/key4` (value: value5) ? (Y/n)

$ parade del /service1/dev/key4 --force
```

LICENSE
----

[MIT](https://github.com/chroju/parade/LICENSE)
