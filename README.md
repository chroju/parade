Parade
======

[![release badge](https://img.shields.io/github/v/release/chroju/parade.svg)](https://github.com/chroju/parade/releases)
[![test badge](https://github.com/chroju/parade/workflows/test/badge.svg)](https://github.com/chroju/parade/actions?workflow=test)

Parade is a simple CLI tool for AWS SSM parameter store. Easy to read and write key values in your parameter store.

Install
-------

### Homebrew

```bash
$ brew install chroju/tap/parade
```

### Download binary

Download the latest binary from [here](https://github.com/chroju/parade/releases) and put it in your `$PATH` directory.

### go get

If you have set up Go environment, you can also install `parade` with `go get` command.

```
$ go get github.com/chroju/parade
```

Authenticate
------------

Parade requires your AWS IAM user authentications. The same authentication method as [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) is available. Tools like [aws-vault](https://github.com/99designs/aws-vault) can be used as well.

```
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

There are Simple 4 sub commands. It is similar to redis-cli. 

### keys

Display keys that match partial search. If no argument is given, all keys will be retrieved.

```
$ parade keys dev
/service1/dev/key1  Type: String
/service1/dev/key2  Type: String
/service1/dev/key3  Type: SecureString
```

### get

Display the value of the specified key.

```
$ parade get /service1/dev/key1
/service1/dev/key1  value1
```

With `--ambiguous` or `-a` flag, display the value of all keys that matched with partial search.

```
$ parade get dev --ambiguous
/service1/dev/key1  value1
/service1/dev/key2  value2
/service1/dev/key3  value3
```

The `--decrypt` or `-d` option is required to decrypt SecureString.

```
$ parade get /service1/dev/password
/service1/dev/password  (encrypted)

$ parade get /service1/dev/password -d
/service1/dev/password  1234password
```

### set

Set new key and value.

```
$ parade set /service1/dev/key4 value4
```

If the specified key already exists, you can choose to overwrite it. Use `--force` flag if you want to force overwriting.

```
$ parade set /service1/dev/key4 value5
WARN: `/service1/dev/key4` already exists.
Overwrite `/service1/dev/key4` (value: value4) ? (Y/n)

$ parade set /service1/dev/key4 value5 --force
```

The value is stored as `String` type by default. It also supports `SecureString` type with the default AWS KMS key and can be specified with the `--encrypt` flag. `StringList` type is not supported.


### del

Delete key and value. Use `--force` flag if you want to force deletion.

```
$ parade del /service1/dev/key4
Delete `/service1/dev/key4` (value: value5) ? (Y/n)

$ parade del /service1/dev/key4 --force
```

LICENSE
----

[MIT](https://github.com/chroju/parade/LICENSE)
