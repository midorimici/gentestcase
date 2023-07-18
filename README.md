## gentestcase

Generate integration test cases.

### Installation

```
go install github.com/midorimici/gentestcase/cmd/gentestcase@latest
```

### Usage

1. Prepare YAML file which defines test factors
1. Run `gentestcase`
1. A CSV file is generated

```
Usage of gentestcase:
  -input string
        input YAML filename (default "cases.yml")
  -output string
        output CSV filename (default "data.csv")
  -w    watch input file change
```

### Test case definition YAML specification

A top-level properties are `factors` and optional `conditions`.

Each key of `factors` represents an factor unique identifier.

An factor contains `name` and `levels` properties.

`levels` is a key-value pair with an level unique identifier for its key and an object with `name` and `if` properties for its value.

`conditions` defines condition variables which can be used by `if` property of factor levels.

Below is an example YAML format.
You can also refer to `examples` directory as input YAML files and its outputs.

```yml
factors:
  factor1:
    name: Factor 1
    levels:
      level_a:
        name: Level A
      level_b:
        name: Level B
  factor2:
    name: Factor 2
    levels:
      level_a:
        name: Level A
        if: factor1.level_a
      level_b:
        name: Level B
```

Below is the EBNF for `if` field syntax.

```ebnf
syntax ::= exp | groups
groups ::= group (' ' operator ' ' group)* | groups (' ' operator ' ' groups)*
group ::= '(' exp ')' | '!(' exp ')' | value_bool
exp ::= value_bool (' ' operator ' ' value_bool)+
operator ::= '&&' | '||'
value_bool ::= value | '!' value
value ::= factor '.' level | '$' condition_ref
factor ::= [a-zA-Z0-9_]+
level ::= [a-zA-Z0-9_]+
condition_ref ::= [a-zA-Z0-9_]+
```

You can generate JSON schema with the following command to validate your YAML files.

```
go run github.com/midorimici/gentestcase/cmd/schema@latest
```
