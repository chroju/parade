# CHANGELOG

## 0.4.0 (2022/11/23)

### ENHANCEMENTS

* `set` command support `--kms-key-id` flag to specify a KMS key ID.

### BUG FIXES

* Fix the issue that `get` command fails when the specified key is not found.
* Fix the issue with error message output to standard out.

## 0.3.2 (2021/04/16)

### ENHANCEMENTS

* `get` command does not make an error when the specified key is not found.

## 0.3.1 (2021/04/15)

### BUG FIXES

* Fix bug when running `keys` command with no arguments. ([#6](https://github.com/chroju/parade/issues/6))
* Fix bugs with aws region and profile set up.
* Fix `get` command usage typos.

## 0.3.0 (2021/02/23)

### NOTE

* Support Apple silicon.
* Deprecate `get` command `--ambiguious` option.

### FEATURES

* Add `--no-color` global option.
* Add `--no-type` option to `keys` subcommand.

### ENHANCEMENTS

* Update `keys` and `get` subcommands query format.
  * Use `*` to specify the partial match and the forward match.
* Modify the output of `get` subcommand.
  * It does not show the key for an exact match.

### BUG FIXES

* Fix duplicated output from `get` subcommand.
* Fix bug, `--help` option also require AWS credential.

## 0.2.1 (2021/02/15)

### BUG FIXES

* Fix bugs about error handlings.

## 0.2.0 (2021/02/15)

### FEATURES

* Add `version` subcommand.
* Add `--region` option.
* Add `--profile` option.
* Support `AWS_DEFAULT_REGION` environment variable.

### ENHANCEMENTS

* Improve `keys` subcommand.
    * Display the parameter types.
* Improve `get` subcommand.
    * Instead of displaying encrypted values as they are, display them as `(encrypted)` string.
    * For multi-line values, display the line feed code as a string of `\n`.
* Improve `set` subcommand.
    * Add confirmation prompt before overwriting existing value.
    * Add `--force` option.
* Improve `del` subcommand.
    * Add confirmation prompt before deleting.
    * Add `--force` option.

### BUG FIXES

* Fix some bugs.

## 0.1.0 (2019/10/13)

* Initial version
