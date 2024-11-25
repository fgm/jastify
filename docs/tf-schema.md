# Notes about Terraform schema

- `schema.Resource`:
  - `Schema`: a map of child property `Schema` by property name. Empty initially
  - `SchemaVersion`: 0
  - `DeprecationMessage`
  - `Timeouts`: a collection CRUD+Default timeouts as `time.Duration`
  - `Description`
  - `UseJSONNumber`: apparently a hack for JSON compatibility
  - `SchemaFunc`: generates a value to set in `Schema`
  - `MigrateState` (fn)
  - `StateUpgraders` ([]fnwrapper)
  - `Create` (fn) `CreateContext` (fn) `CreateWithoutTimeout` (fn)
  - `Read` (fn) `ReadContext` (fn) `ReadWithoutTimeout` (fn)
  - `Update` (fn) `UpdateContext` (fn) `UpdateWithoutTimeout` (fn)
  - `Delete` (fn) `DeleteContext` (fn) `DeleteWithoutTimeout` (fn)
  - `Exists` (fn): There is no `ExistsContext` or `ExistsWithoutTimeout`
  - `CustomizeDiff` (fn)
  - `Importer` (fnwrapper)
- `schema.Schema`: Description of a data type: Type (see below), Required, Optional, Default, Min/MaxItems..

## Types

* TypeInvalid ValueType = iota
* Simple types:
  * TypeBool
    * Examples: `is_read_only`
  * TypeInt
  * TypeFloat
  * TypeString
    * Examples: `title`, `layout_type`
* TypeList (= ordered list, duplicates allowed)
  * Examples:
    * `tags`: list of strings, on some resources
    * `template_variable`: block (list of `Resource`)
    * `widget`: block (list of resource)
    * `widget_layout`: block of length 0|1
      * provided as a map instead of slice of length 1.
      * not provided instead of slice of length 0.
* TypeMap (= K/V, unordered, no duplicates)
  * rarely used: only 4 outside tests in the 3.46.0 provider
    * `permissions` on `datadog` (what is that ?)
    * `defaults_tags.tags` on the provider itself
    * `account_specific_namespace_rules` on the DD/AWS integration
    * `locations` on the Synthetics locations
* TypeSet (= unordered list, no duplicates)
  * Examples
    * `dashboard_lists_removed`: list of ints
    * `notify_list`: list of Schema
    * `restricted_roles`: list of strings
    * `tags`: list of strings, on other resources
* typeObject: not used in the DD provider
