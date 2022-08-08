# Stalker

Stalk files in a directory and run command when any file changes.

## How to use

- Run `stalk <path_to_directory>`

### Config

`.stalkerrc.json`

```json
{
  "command": "echo file-change",
  "ignored": ["foo.*"]
}
```
