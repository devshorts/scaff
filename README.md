scaff
==

[![Build Status](https://travis-ci.org/devshorts/scaff.svg?branch=master)](https://travis-ci.org/devshorts/scaff)

Super simple `scaff`older.  Scaffolds paths and files using simple rules driven by a yml file.

## Installation

```
git clone https://github.com/devshorts/scaff.git
cd scaff && go install
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
```

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

