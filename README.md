scaff
==

[![Build Status](https://travis-ci.org/devshorts/scaff.svg?branch=master)](https://travis-ci.org/devshorts/scaff)

Super simple `scaff`older.  Scaffolds paths and files using simple rules driven by a yml file.

## Installation

```
git clone https://github.com/devshorts/scaff.git
cd scaff && go install
```

## Usage

```
$ scaff -h
Usage:
  main [OPTIONS]

Application Options:
  -d, --source_dir= Source directory containing templates
  -t, --target_dir= Target directory to make with templated data
      --scaff_file= Name of yaml file containing config. Defaults to .scaff.yml
      --dry_run     Dry Run

Help Options:
  -h, --help        Show this help message
```

## Configuration

YML file is of the form:

```yaml
context:
  foo:
    default: 123
    description: My bar!
  data:
    description: My name!
    # dummy verify hook that always succeeds
    verify_hook:
      command: python
      args:  [-c, import sys; sys.exit(0)]
file_config:
  lang_delims:
    .js: $$
```

All fields will need to be set when requested by the user, either via a default or inputed 
by the user

Where `foo` and `data` are keys to be used in rules like:

```
__camel_foo__
__upper_data__
```

etc

Rules are of the form `__ruleName_id__`.  Available rules are

- `camel` - transforms templates to `camelCase`
- `upper` - transforms templates to `UPPER CASE`
- `lower` - transforms templates to `lower case`
- `snake` - transforms templates to `snake_case`
- `pkg` -  takes anything of the form `a.b.c` from a user and makes it `a/b/c` (for use in file paths)
- `id` - only replaces the identifier in the template with the user supplied value

Rules can be in the path or in text.

## Post set hooks

Often times just accepting any text isn't good enough for a scaffolder. You have
business rules for certain keys and want to be able to validate them pre-templating.  
Scaff supports post set hooks (but pre templating) that you can tap into.  

Examples can be things like verifying keys start with certain prefixes, or exclude characters, etc.  
The delegation spawns a subshell to execute, so it is extensible for your needs.

## Contextual data

Sometimes you want to create compound replacements that can only be made from previous results.

This can be done like so:

```yaml
context:
  foo:
    default: 123
    description: My bar!
  data:
    default: data
    description: My name!  
dynamics:
  some_key: hello world! __id_foo_/__upper_data__
```

In this example, `some_key` uses the value of `foo` as is (by using the `id` rule) and then upper cases `data`
into a compound string. When this done you can use the value of `dynamics` and it will have the value 

```
hello world! 123/DATA
```

## Language delimiters

Different languages have different identifier semantics. The default rule delim is `__`
but you can specify custom delimiters for file extensions, such that where `__` doesn't 
compile or validate as a valid identifier, you can replace it whatever you want.

## Language specific replacements

Some langauges (like go) can't use custom delimiters to represent things like package pathing, 
in such a way that would maintain compilation.  To that, you can add custom language specific
format replacement mapping. Currently supported only for `go`. For example:

 
```yaml
context:
  pkg:
    default: github.com/target
    description: The target package name 
file_config:
  lang_rules:
    go:
      pkg: "github.com/test/.scaffold"
      replace_with_id: "pkg"
```

Notice how we are saying to replace `github.com/test/.scaffold` package root with the 
result of the context identifier of `pkg`.  This will now transform:

```
import "github.com/test/.scaffold/foo/bar

// .. your go code
```

Into

```
import "github.com/test/.scaffold/foo/bar

// .. your go code
```


## As a library

Use scaff yourself!

```
config := scaff.NewParser(opts.ScaffConfigFile).GetConfig(opts.SourceDir)
	
scaff.Template(config, opts)
```
