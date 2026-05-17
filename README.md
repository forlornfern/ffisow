# ffisow

**ffisow** — ForlornFern ISO Writer

Terminal tool for writing ISO images to block devices (USB drives, etc).

![](./docs/assets/demo.gif)

## Installation

```sh
go install github.com/forlornfern/ffisow@latest
```

Or build from source:

```sh
git clone https://github.com/forlornfern/ffisow
cd ffisow
go build -o ffisow
sudo mv ffisow /usr/local/bin/
```

## Usage

```
ffisow <iso> <device>
```

```sh
sudo ffisow ~/Downloads/archlinux.iso /dev/sdb
```

> **Warning:** Make sure you specify the correct device (`/dev/sdb`, not `/dev/sdb1`). All data on the device will be overwritten.

## Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--buffer` | `-b` | `1024` | Buffer size in KiB |
| `--verbose` | `-v` | `false` | Verbose logging |

## Examples

```sh
# default buffer (1 MiB)
sudo ffisow ~/Downloads/archlinux.iso /dev/sdb

# custom buffer size (8 MiB)
sudo ffisow -b 8192 ~/Downloads/archlinux.iso /dev/sdb
```
