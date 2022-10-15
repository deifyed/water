# Water

## Motivation

I want easily customizable file and directory scaffolds. Which I recently read was super hard. Anyway.

## Usage

```shell
# First, create the file or folder you wish to grow into something spectacular
touch <filename>
# or
mkdir <directory name>

# For example
touch Makefile

# Then, water it with nurishing liquid to sproud it into a beautiful growth
water Makefile
```

## Installation

```shell
# Build the project
make build

# Then install it into ~/.local/bin
make install

# To adjust installation path, set the PREFIX variable
make install PREFIX=/usr/local
```

## Configuration

See [example config](examples/water.yaml)
