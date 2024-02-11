# z
🌳 Personal stateful command tree monolith with Bonzai

##  Install
Just download one of the release [binaries](https://github.com/yemaney/z/releases):

```
curl -L https://github.com/yemaney/z/releases/latest/download/z-linux-amd64 -o ~/.local/bin/z
curl -L https://github.com/yemaney/z/releases/latest/download/z-windows-amd64 -o ~/.local/bin/z
```

Or install directly with `go`:

```
go install github.com/yemaney/z@latest
```

## Tab Completion
To activate bash completion just use the `complete -C` option from your `.bashrc` or command line. There is no messy sourcing required. All the completion is done by the program itself.

```
complete -C z z
```
