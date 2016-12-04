## Synopsis

Water, the simple and neat bootstrapping application for everything.

## Code Example
Bootstrapping a file
```bash
$ touch index.html
$ water index.html
```

Bootstrapping a directory
```bash
$ mkdir stack
$ water stack
```

## Motivation

I needed a lightweight and easy to use bootstrapping mechanism. 

## Installation

```bash
$ git clone https://github.com/deifyed/water.git && cd water
$ chmod +x water
$ mkdir -p ~/.local/bin
$ mv water ~/.local/bin
$ echo "export PATH=~/.local/bin:$PATH" >> ~/.bashrc
$ rm -r water
```

## Contributors

Julius Pedersen
