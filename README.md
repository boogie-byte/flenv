# FLags with defaults in ENVironment variables

`flenv` is a pretty simple library for parsing CLI flags with an additional feature of reading default values from envvars. When I implemented it there were no such libraries around (all available options were either type-unsafe or config file-centric). Or maybe I just search thoroughly enough.

## Installation

```
go get -u github.com/boogie-byte/flenv
```

## Usage

```go
// Initialize the parser
p := flenv.New(
    flenv.WithAppName("my-app"),
    flenv.WithAppVersion("1.2.3"),
)

// Register flags
var (
    b bool
    i int
    s string
)

p.Bool(&b, "my-bool-flag", "My bool flag")
p.Int(&i, "my-int-flag", "My int flag").Required()
p.String(&s, "my-string-flag", "My string flag").Default("foo")

// Parse command line args
p.Parse()
```

## Supported variable types
* `bool`
* `int`
* `string`
* `time.Duration`

Adding support for any other type is pretty straightforward, I'll support more types as needed.

## Supported flag formats
Both `--key=<value>` and `--key <value>` flag formats are supported. Additionally, `bool` flags support `--key` format without the value.

Short flags are not supported yet.

## Required flags and default values
To mark a flag as required use the `.Required()` method:
```go
p.Int(&i, "my-int-flag", "My int flag").Required()
```

To provide a hard-coded default value for a flag use the `.Default()` method:
```go
p.String(&s, "my-string-flag", "My string flag").Default("foo")
```

Boolean flags are a special case and will panic if either `.Required()` or `.Default()` method is called.

## Envvar defaults
By default all flags are registered with environment variable lookup enabled. Flag names are translated to envvar names by capitalizing all letters and substituting dashes (`-`) with underscores (`_`). E.g. `my-bool-flag` becomes `MY_BOOL_FLAG`.

Each flag's envvar name could be overridden individually via the `.Env()` method, e.g.:
```go
p.Bool(&b, "my-bool-flag", "My bool flag").Env("TOTALLY_DIFFERENT_ENVVAR")
```

The automatic envvar registration could be disabled via the `WithoutAutoEnv()` parser option. An alternative envvar name formatting function could be provided via the `WithEnvVarFormatter()` parser option. A global envvar name prefix could be provided via the `WithEnvVarPrefix()` parser option.

## Help message
`flenv.New()` automatically registers a `--help` flag with the new parser, which if specified will make the `Parse()` method print the help message and exit the process. Alternatively the help message will be printed if any flag parsing errors occur.

Help message example:
```
Usage: my-app --my-int-flag=INT [--help] [--my-bool-flag] [--my-string-flag=STRING] [--version]

Flags:
  --help                   Show help message
  --my-bool-flag           My bool flag [$MY_BOOL_FLAG]
  --my-int-flag=INT        My int flag (required) [$MY_INT_FLAG]
  --my-string-flag=STRING  My string flag (default: foo) [$MY_STRING_FLAG]
  --version                Show application version
```

To change the `--help` flag name use the `WithHelpFlagName()` parser option.

## Application name and version
If the application name is provided via the `WithAppName()` parser option, it will be used instead of the `os.Args[0]` in the help message.

If the application version is provided via the `WithAppVersion()` parser option, the `--version` flag will be registered automatically, which if specified will make the `.Parse()` method print the app version and exit the process.

To change the `--version` flag name use the `WithAppVersionFlagName()` parser option.

## Missing features
- [ ] Short flags support
