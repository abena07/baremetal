# protocol spec

text-based protocol over TCP. encoding: UTF-8. line ending: `\n` (LF).
clients send one request per line. server replies with one line per request.

## framing

fields are separated by `|`. messages are terminated by `\n`.
`\r\n` is also accepted on input but responses always use `\n`.

## request

```
COMMAND|arg1|arg2\n
```

- COMMAND is uppercase ASCII
- args are positional and command-specific
- `|` is a reserved delimiter and may not appear in values

## response

```
OK|result\n
ERR|message\n
```

## commands

| command | args       | success response      |
|---------|------------|-----------------------|
| PING    | none       | `OK&#124;PONG`        |
| SET     | key, value | `OK&#124;`            |
| GET     | key        | `OK&#124;value`       |
| DEL     | key        | `OK&#124;`            |


## edge cases

| input              | response                                             |
|--------------------|------------------------------------------------------|
| empty / blank line | `ERR&#124;empty message`                             |
| unknown command    | `ERR&#124;unknown command: "FOO"`                    |
| missing args       | `ERR&#124;SET requires exactly 2 arguments`          |
| too many args      | `ERR&#124;PING requires 0 arguments`                 |
| pipe in value      | `ERR&#124;invalid character in value`                |