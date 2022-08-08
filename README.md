# Stalker

Stalk files in a directory and run command when any file changes.

## How to use

`stalk -loc <path_to_directory>`

| Flag | Description                           | Default value    | Usage                         |
| ---- | ------------------------------------- | ---------------- | ----------------------------- |
| loc  | location of directory/file to stalked | ./               | `stalk -loc ./dev/project`    |
| cmd  | command to run on file change         | echo file change | `stalk -cmd "go run main.go"` |

### Config

`.stalkerrc.json`

```json
{
  "command": "echo file-change",
  "ignored": ["foo.*"]
}
```
