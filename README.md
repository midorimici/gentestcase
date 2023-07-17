## gentestcase

Generate integration test cases.

### Installation

```
go install github.com/midorimici/gentestcase/cmd/gentestcase@latest
```

### Usage

1. Prepare YAML file which defines test elements
1. Run `gentestcase`
1. A CSV file is generated

```
Usage of gentestcase:
  -input string
        input YAML filename (default "elements.yml")
  -output string
        output CSV filename (default "data.csv")
```

### Element YAML specification

A top-level properties are `elements` and optional `conditions`.

Each key of `elements` represents an element unique identifier.

An element contains `name` and `options` properties.

`options` is a key-value pair with an option unique identifier for its key and an object with `name` and `if` properties for its value.

`conditions` defines condition variables which can be used by `if` property of element options.

Below is an example YAML format.
You can also refer to `examples` directory as input YAML files and its outputs.

```yml
elements:
  element1:
    name: Element 1
    options:
      option_a:
        name: Option A
      option_b:
        name: Option B
  element2:
    name: Element 2
    options:
      option_a:
        name: Option A
        if: element1.option_a
      option_b:
        name: Option B
```

Below is the EBNF for `if` field syntax.

```ebnf
syntax ::= exp | groups
groups ::= group (' ' operator ' ' group)* | groups (' ' operator ' ' groups)*
group ::= '(' exp ')' | '!(' exp ')' | value_bool
exp ::= value_bool (' ' operator ' ' value_bool)+
operator ::= '&&' | '||'
value_bool ::= value | '!' value
value ::= element '.' option | '$' condition_ref
element ::= text
option ::= text
condition_ref ::= text
```

You can generate JSON schema with the following command to validate your YAML files.

```
go run github.com/midorimici/gentestcase/cmd/schema@latest
```
