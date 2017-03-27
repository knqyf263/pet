# pet : CLI Snippet Manager

[![GitHub release](https://img.shields.io/github/release/knqyf263/pet.svg)](https://github.com/knqyf263/pet/releases/latest)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/knqyf263/pet/blob/master/LICENSE)

<img src="doc/logo.png" width="150">

Simple command-line snippet manager, written in Go

<img src="doc/pet01.gif" width="700">

You can use variables (`<param>` or `<param=default_value>` ) in snippets.

<img src="doc/pet08.gif" width="700">


# Abstract

`pet` is written in Go, and therefore you can just grab the binary releases and drop it in your $PATH.

`pet` is a simple command-line snippet manager (inspired by [memo](https://github.com/mattn/memo)).
I always forget commands that I rarely use. Moreover, it is difficult to search them from shell history. There are many similar commands, but they are all different.

e.g. 
- `$ awk -F, 'NR <=2 {print $0}; NR >= 5 && NR <= 10 {print $0}' company.csv` (What I am looking for)
- `$ awk -F, '$0 !~ "DNS|Protocol" {print $0}' packet.csv`
- `$ awk -F, '{print $0} {if((NR-1) % 5 == 0) {print "----------"}}' test.csv`

In the above case, I search by `awk` from shell history, but many commands hit.

Even if I register an alias, I forget the name of alias (because I rarely use that command).

So I made it possible to register snippets with description and search them easily.

# TOC

- [Main features](#main-features)
- [Examples](#examples)
    - [Register the previous command easily](#register-the-previous-command-easily)
        - [bash/zsh](#bashzsh)
        - [fish](#fish)
    - [Select snippets at the current line (like C-r)](#select-snippets-at-the-current-line-like-c-r)
        - [bash](#bash)
        - [zsh](#zsh)
        - [fish](#fish-1)
    - [Copy snippets to clipboard](#copy-snippets-to-clipboard)
- [Features](#features)
    - [Edit snippets](#edit-snippets)
    - [Sync snippets](#sync-snippets)
- [Usage](#usage)
- [Snippet](#snippet)
- [Configuration](#configuration)
    - [Sync](#sync)
- [Installation](#installation)
    - [Binary](#binary)
    - [Mac OS X / Homebrew](#mac-os-x--homebrew)
    - [Archlinux](#archlinux)
    - [Build](#build)
- [Contribute](#contribute)


# Main features
`pet` has the following features.

- Register your command snippets easily.
- Use variables in snippets.
- Search snippets interactively.
- Run snippets directly.
- Edit snippets easily (config is just a TOML file).
- Sync snippets via Gist.

# Examples
Some examples are shown below.

## Register the previous command easily

### bash/zsh
By adding the following config to `.bashrc` or `.zshrc`, you can easily register the previous command.

```
$ cat .zshrc
function prev() {
  PREV=$(fc -lrn | head -n 1)
  sh -c "pet new `printf %q "$PREV"`"
}
```

### fish
See below for details.  
https://github.com/otms61/fish-pet

<img src="doc/pet02.gif" width="700">

## Select snippets at the current line (like C-r)

### bash
By adding the following config to `.bashrc`, you can search snippets and output on the shell.

```
$ cat .bashrc
function pet-select() {
  BUFFER=$(pet search --query "$READLINE_LINE")
  READLINE_LINE=$BUFFER
  READLINE_POINT=${#BUFFER}
}
bind -x '"\C-x\C-r": pet-select'
```

### zsh

```
$ cat .zshrc
function pet-select() {
  BUFFER=$(pet search --query "$LBUFFER")
  CURSOR=$#BUFFER
  zle redisplay
}
zle -N pet-select
stty -ixon
bindkey '^s' pet-select
```

### fish
See below for details.  
https://github.com/otms61/fish-pet

<img src="doc/pet03.gif" width="700">


## Copy snippets to clipboard
By using `pbcopy` on OS X, you can copy snippets to clipboard.

<img src="doc/pet06.gif" width="700">

# Features

## Edit snippets
The snippets are managed in the TOML file, so it's easy to edit.

<img src="doc/pet04.gif" width="700">


## Sync snippets
You can share snippets via Gist.

<img src="doc/pet05.gif" width="700">


# Usage

```
pet - Simple command-line snippet manager.

Usage:
  pet [command]

Available Commands:
  configure   Edit config file
  edit        Edit snippet file
  exec        Run the selected commands
  help        Help about any command
  list        Show all snippets
  new         Create a new snippet
  search      Search snippets
  sync        Sync snippets
  version     Print the version number

Flags:
      --config string   config file (default is $HOME/.config/pet/config.toml)
      --debug           debug mode

Use "pet [command] --help" for more information about a command.
```

# Snippet
Run `pet edit`  
You can also register the output of command (but cannot search).

```
[[snippets]]
  description = "echo | openssl s_client -connect example.com:443 2>/dev/null |openssl x509 -dates -noout"
  command = "Show expiration date of SSL certificate"
  output = """
notBefore=Nov  3 00:00:00 2015 GMT
notAfter=Nov 28 12:00:00 2018 GMT"""
```

Run `pet list`

```
Description: echo | openssl s_client -connect example.com:443 2>/dev/null |openssl x509 -dates -noout
    Command: Show expiration date of SSL certificate
     Output: notBefore=Nov  3 00:00:00 2015 GMT
             notAfter=Nov 28 12:00:00 2018 GMT
------------------------------
```


# Configuration

Run `pet configure`

```
[General]
  snippetfile = "path/to/snippet" # specify snippet directory
  editor = "vim"                  # your favorite text editor
  column = 40                     # column size for list command
  selectcmd = "peco"              # selector command for edit command

[Gist]
  file_name = "pet-snippet.toml"  # specify gist file name
  access_token = ""               # your access token
  gist_id = ""                    # Gist ID
```

## Sync
You must obtain access token.
Go https://github.com/settings/tokens/new and create access token (only need "gist" scope).
Set that to `access_token` in `[Gist]`.

After setting, you can upload snippets to Gist.
```
$ pet sync -u
Gist ID: 1cedddf4e06d1170bf0c5612fb31a758
Upload success
```

Set `Gist ID` to `gist_id` in `[Gist]`.

You can download snippets on another PC.
```
$ pet sync
Download success
```

# Installation
If [peco](https://github.com/peco/peco#installation) is not installed, please install first (`homebrew` install `peco` automatically).

## Binary
Go to [the releases page](https://github.com/knqyf263/pet/releases), find the version you want, and download the zip file. Unpack the zip file, and put the binary to somewhere you want (on UNIX-y systems, /usr/local/bin or the like). Make sure it has execution bits turned on. 

## Mac OS X / Homebrew
You can use homebrew on OS X.
```
$ brew install knqyf263/pet/pet
```

If you receive an error (`Error: knqyf263/pet/pet 64 already installed`) during `brew upgrade`, try the following command

```
$ brew unlink pet && brew uninstall pet
($ rm -rf /usr/local/Cellar/pet/64)
$ brew install knqyf263/pet/pet
```

## Archlinux
A package is available in [AUR](https://aur.archlinux.org/packages/pet-git/).
```
$ yaourt -S pet-git
```

## Build

```
$ go get github.com/knqyf263/pet
```

# Contribute

1. fork a repository: github.com/knqyf263/pet to github.com/you/repo
2. get original code: `go get github.com/knqyf263/pet`
3. work on original code
4. add remote to your repo: git remote add myfork https://github.com/you/repo.git
5. push your changes: git push myfork
6. create a new Pull Request

- see [GitHub and Go: forking, pull requests, and go-getting](http://blog.campoy.cat/2014/03/github-and-go-forking-pull-requests-and.html)

----

# License
MIT

# Author
Teppei Fukuda
