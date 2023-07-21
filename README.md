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

You can generate JSON schema with the following command to validate your YAML files.

```
go run github.com/midorimici/gentestcase/cmd/schema@latest
```

You can refer to `examples` directory as input YAML files and its outputs.

#### `factors`

`factors` is a required top-level property.

Values of `factors` are key-value pairs that represents each individual factor.

A pair has a key that represents an factor unique identifier, and a value that contains `name` and `levels` properties.

`name` represents the outputted name of the factor.

`levels` is a key-value pair with an level unique identifier as its key, and an outputted name of the level as its value.

Below is an example of `factors` definition.

```yml
factors:
  factor1:
    name: Factor 1
    levels:
      level_a: Level A
      level_b: Level B
  factor2:
    name: Factor 2
    levels:
      level_a: Level A
      level_b: Level B
```

#### `conditions`

`conditions` is an optional top-level property.

It defines condition variables which can be used in other values in `conditions` or `constraints`.

Values of `conditions` are key-value pairs that represents each condition.

A pair has a key that represents a condition unique identifier, and a value that represents a condition statement.

Below is the EBNF of condition statement syntax.

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

Below is an example of `conditions` definition.

```yml
conditions:
  is_hot: 'temprature.30_35 || temprature.35_above'
  is_cold: 'temprature.0_below || temprature.0_5'
  is_nice_day: '!$is_hot && !$is_cold && weather.sunny'
```

#### `constraints`

`constraints` is an optional top-level property.

It defines constraints between factor levels.

A value of `constraints` is a list with each constraint.

A constraint has `only_if` (required), `then` (required) and `else` (optional) properties.

Each of those properties has a value of condition statement.

A combination that satisfies `then` condition is preserved **only if** `only_if` condition is evaluated to true.

Thus, when `only_if` is evaluated to false, combinations that satisfy `then` condition are omitted from the result output.

When `else` is specified, a combination that satisfies `else` condition is preserved only if `only_if` condition is evaluated to false.

Thus, when `only_if` is evaluated to true, combinations that satisfy `else` condition are omitted from the result output.

Below is an example of `constraints` definition.

```yml
constraints:
  - only_if: 'weather.rainy || ($is_hot && weather.sunny)'
    then: 'action.open_umbrella'
  - only_if: 'weather.sunny'
    then: 'action.play_tennis'
    else: 'action.read_books'
```
