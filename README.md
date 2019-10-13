Parade
======

[![release badge](https://img.shields.io/github/v/release/chroju/parade.svg)](https://github.com/chroju/parade/releases)
[![test badge](https://github.com/chroju/parade/workflows/test/badge.svg)](https://github.com/chroju/parade/actions?workflow=test)

Parade is a simple CLI tool for AWS SSM parameter store. Easy to read and write key values in your parameter store.

Install
-------

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
$ aws-vault exec YOUR_PROFILE -- parade
```

Usage
-----

Simple four sub commands. It is similar to redis-cli. 

### keys

Display keys that match partial search.

```
$ parade keys dev
/service1/dev/key1
/service1/dev/key2
/service1/dev/key3
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

### set

Set new key value.

```
$ parade set /service1/dev/key4 value4
done.
```

Use `--force` flag if you want to overwrite.

```
$ parade set /service1/dev/key4 value5
ParameterAlreadyExists: The parameter already exists. To overwrite this value, set the overwrite option in the request to true.
        status code: 400, request id: ae21f5d5-XXXX-XXXX-XXXX-XXXXXXXXXXXX

$ parade set /service1/dev/key4 value5 --force
done.
```

The value is stored as `String` type by default. It also supports `SecureString` type with the default AWS KMS key and can be specified with the `--encrypt` flag. `StringList` type is not supported.


### del

Delete a key value.

```
$ parade del /service1/dev/key4
done.
```

Author
----

[chroju](https://github.com/chroju)

LICENSE
----

[MIT](https://github.com/chroju/parade/LICENSE)
