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

| command | args       | success response   |
|---------|------------|--------------------|
| PING    | none       | `OK|PONG`          |
| SET     | key, value | `OK|`              |
| GET     | key        | `OK|value`         |

## edge cases

| input              | response                           |
|--------------------|------------------------------------|
| empty / blank line | `ERR|empty message`                |
| unknown command    | `ERR|unknown command: FOO`         |
| missing args       | `ERR|missing args: SET requires 2` |
| pipe in value      | `ERR|invalid character in value`   |
