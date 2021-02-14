# CHANGELOG

## 0.2.0 (2021/02/15)

* Add `version` subcommand
* Improve `keys` subcommand
    * Display the parameter types.
* Improve `get` subcommand
    * Instead of displaying encrypted values as they are, display them as `(encrypted)` string.
    * For multi-line values, display the line feed code as a string of `\n`.
* Improve `set` subcommand
    * Add confirmation prompt before overwriting existing value.
    * Add `--force` option.
* Improve `del` subcommand
    * Add confirmation prompt before deleting.
    * Add `--force` option.
* Support `AWS_DEFAULT_REGION` environment variable
* Add `--region` option
* Add `--profile` option
* Fix some bugs.

## 0.1.0 (2019/10/13)

* Initial version
