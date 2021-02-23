# CHANGELOG

## 0.3.0 (2021/02/23)

* Support Apple sillicon.
* Update `keys` and `get` subcommands query format.
  * Use `*` to specify the partial match and the forward match.
  * Deprecate `get` command `--ambiguious` option.
* Modify the output of `get` subcommand.
  * It does not show the key for an exact match.
* Add `--no-color` global option.
* Add `--no-type` option to `keys` subcommand.
* Fix some bugs.
  * Duplicate output from `get` subcommand.
  * `--help` also required AWS credential.

## 0.2.1 (2021/02/15)

* Fix bugs about error handlings.

## 0.2.0 (2021/02/15)

* Add `version` subcommand.
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
* Support `AWS_DEFAULT_REGION` environment variable.
* Add `--region` option.
* Add `--profile` option.
* Fix some bugs.

## 0.1.0 (2019/10/13)

* Initial version
