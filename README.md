## gentestcase

Generate integration test cases.

### Installation

```
go install github.com/midorimici/gentestcase/cmd/gentestcase@v1.0.2
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

A top-level element represents an element unique identifier.

An element contains `name` and `options` properties.

`options` is a key-value pair with an option unique identifier for its key and an object with `name` and `if` properties for its value.

Below is an example YAML format.
You can also refer to `examples/elements.yml` as input YAML file and `examples/data.csv` as its output.

```yml
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
syntax ::= exp | group (' ' operator ' ' group)*
group ::= '(' exp ')' | '!(' exp ')' | value_bool
exp ::= value_bool (' ' operator ' ' value_bool)+
operator ::= '&&' | '||'
value_bool ::= value | '!' value
value ::= element '.' option
element ::= text
option ::= text
```
