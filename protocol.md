# protocol spec

text-based protocol over TCP. encoding: UTF-8. line ending: `\n` (LF).
clients send one request per line. server replies with one line per request.

## framing

fields are separated by `|`. messages are terminated by `\n`.
`\r\n` is also accepted on input but responses always use `\n`.
a request line may be at most 1 MiB; longer lines close the connection with an error.

## request

```
COMMAND|arg1|arg2\n
```

- COMMAND is uppercase ASCII
- args are positional and command-specific
- the first arg of every command that takes args is a key. keys are trimmed of
  surrounding whitespace, must be non-empty, and may not contain `|`
- values (the last arg of SET) are stored byte-for-byte: whitespace is kept, an
  empty value is allowed, and `|` is allowed since everything after the key's
  delimiter is the value

## response

```
OK|result\n
ERR|message\n
```

## commands

| command | args       | success response          |
|---------|------------|---------------------------|
| PING    | none       | OK\|PONG                  |
| SET     | key, value | OK\|                      |
| GET     | key        | OK\|value                 |
| DEL     | key        | OK\|                      |
| LIST    | none       | OK\|key1\|key2 (sorted)   |


## edge cases

| input              | response                                  |
|--------------------|-------------------------------------------|
| empty / blank line | ERR\|empty message                        |
| unknown command    | ERR\|unknown command: "FOO"               |
| missing args       | ERR\|SET requires exactly 2 arguments     |
| too many args      | ERR\|PING requires 0 arguments            |
| empty key          | ERR\|key must not be empty                |
| missing key        | ERR\|key not found                        |