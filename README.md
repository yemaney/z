# z
🌳 Personal stateful command tree monolith with Bonzai

##  Install
Just download one of the release [binaries](https://github.com/yemaney/z/releases):

```sh
#------------ Example ------------##

# Download the tarfile you want.
curl -LO https://github.com/yemaney/z/releases/download/v0.1.1/z_0.1.1_windows_amd64.tar.gz

tar -xzvf z_0.1.1_windows_amd64.tar.gz

# make sure target directory exists
mkdir -p "$HOME/.local/bin/"

# Move to binary the binary to target directory for binaries
mv z "$HOME/.local/bin/"

# Make the binary executable
chmod +x "$HOME/.local/bin/z"
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
