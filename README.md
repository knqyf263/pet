# Pet - CLI Snippet Manager

# Motivation

I built `pet` as a simple command-line snippet manager (inspired by [memo](https://github.com/mattn/memo)).

I struggle to remember complex commands or ones I rarely use, and searching shell history makes finding them even harder.

Stop expecting yourself to remember every command — focus on productivity and finding the right commands as fast as possible. This feels fun when you work with 2 tools for 2 years, but grows painful when a decade passes and you work across backend, frontend, and infrastructure with dozens of tools. If you're a developer, you probably relate.

`pet` lets you save, tag, search, and execute command-line snippets easily! It's nearly 8 years old and developers around the world use it daily.

Go wrote `pet` in Go, so you can grab the binary releases and drop them in your `$PATH`.

---

# Creating a Snippet

Run `pet new` to create a snippet:

```
$ pet new
Command> echo Hello world!
Description> print Hello world
```

Run `pet new --help` to see all available arguments.

Use `pet new --multiline` to enter multiline commands.

You can also use variables in snippets — the next section covers parameters in detail.

You can also **tag** snippets to find them faster — see the tag section for more.

---

# Parameters

You can enter parameters in `<n_ways>` ways.

Parameters can carry default values: Hello `<subject=world>`, defined with an equal sign.

They can even carry `<content=spaces & = signs>` where the default value becomes \<content=spaces & = signs\>.

Default values just can't \<end with spaces \>.

Parameters also support multiple default values:
Hello `<subject=|_John_||_Sam_||_Jane Doe = special #chars_|>`

---

# Examples

## Register the Previous Command Easily

Add the following config to `.bashrc` or `.zshrc` to register the previous command quickly.

### bash prev function

```bash
function prev() {
  PREV=$(echo `history | tail -n2 | head -n1` | sed 's/[0-9]* //')
  sh -c "pet new `printf %q "$PREV"`"
}
```

### zsh prev function

```zsh
function prev() {
  PREV=$(fc -lrn | head -n 1)
  sh -c "pet new `printf %q "$PREV"`"
}
```

---

## Select Snippets at the Current Line (like C-r) — RECOMMENDED

### bash
Add the following config to `.bashrc` to search snippets and output them to the shell. This lets you execute commands yourself, adding them to your shell history. You can also **chain** commands with this approach.

Customize the search and list commands with options like `-t` or `--tags` — for example, `pet search -t myjob` searches only snippets tagged with `myjob`.

```bash
function pet-select() {
  BUFFER=$(pet search --query "$READLINE_LINE")
  READLINE_LINE=$BUFFER
  READLINE_POINT=${#BUFFER}
}
bind -x '"\C-x\C-r": pet-select'
```

### zsh

```zsh
function pet-select() {
  BUFFER=$(pet search --query "$LBUFFER")
  CURSOR=$#BUFFER
  zle redisplay
}
zle -N pet-select
stty -ixon
bindkey '^s' pet-select
```

---

## Copy Snippets to Clipboard

Use `pbcopy` on macOS to copy snippets to your clipboard.

---

# Features

## Edit Snippets
`pet` manages snippets in a TOML file, making them easy to edit directly.

## Sync Snippets
Share your snippets via Gist.

---

# Configuration

Run `pet configure` to open the config file.

```toml
[General]
  snippetfile = "path/to/snippet"
  editor = "vim"
  column = 40
  selectcmd = "fzf"
  backend = "gist"
  sortby  = "description"
  cmd = ["sh", "-c"]
  color = false
  format = "[$description]: $command $tags"
```

---

## Sync

### Gist
Obtain an access token at https://github.com/settings/tokens/new — you only need the "gist" scope. Set it as `access_token` in `[Gist]`, or export it as `$PET_GITHUB_ACCESS_TOKEN`.

Once configured, upload your snippets to Gist. If you haven't set a `gist_id`, `pet` will create a new Gist automatically.

```
pet sync
Gist ID: 1cedddf4e06d1170bf0c5612fb31a758
Upload success
```

`pet sync` compares your local file against the Gist by update date and automatically downloads or uploads as needed.

---

## Auto Sync
Set `auto_sync = true` in `[Gist]`, `[GHEGist]`, or `[GitLab]` to sync snippets automatically whenever you run `pet new` or `pet edit`.

---

# Installation

Install a selector command ([fzf](https://github.com/junegunn/fzf) or [peco](https://github.com/peco/peco)) first. Homebrew installs `fzf` automatically.

After installing Pet, set up the shortcuts in the [ZSH Prev](#zsh-prev-function) section — it's highly recommended.

## macOS / Homebrew

```
brew install knqyf263/pet/pet
```

## Fedora, RedHat, CentOS

Download the rpm package from [the releases page](https://github.com/knqyf263/pet/releases):

```
sudo rpm -ivh https://github.com/knqyf263/pet/releases/download/vx.x.x/pet_x.x.x_linux_amd64.rpm
```

## Debian, Ubuntu

```
wget https://github.com/knqyf263/pet/releases/download/vx.x.x/pet_x.x.x_linux_amd64.deb
dpkg -i pet_x.x.x_linux_amd64.deb
```

## Archlinux

Install from source:
```
yay -S pet-git
```
Or install from binary:
```
yay -S pet-bin
```

---

# Contribute

1. Fork the repository: github.com/knqyf263/pet → github.com/you/repo
2. Get the original code: `go get github.com/knqyf263/pet`
3. Work on the original code
4. Add a remote to your repo: `git remote add myfork https://github.com/you/repo.git`
5. Push your changes: `git push myfork`
6. Open a new Pull Request

---

# License
MIT

# Author
Teppei Fukuda
