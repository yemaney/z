# z
🌳 Personal stateful command tree monolith with Bonzai

##  Install
### Option 1: Download and Extract
- download one of the [releases](https://github.com/yemaney/z/releases)
- extract the binary `z`
- make sure its executable
- place the `z` binary in a folder reachable in your system's PATH.


### Option 2: Install with `go`
```
go install github.com/yemaney/z@latest
```

## Tab Completion
To activate bash completion just use the `complete -C` option from your `.bashrc` or command line. There is no messy sourcing required. All the completion is done by the program itself.

```
complete -C z z
```
