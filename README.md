# Hackagotchi – A pet for your computer

![](./demo/intro.gif)

* Daily quests
* Hundreds of achivements
* Keep coding, to keep your hackagotchi happy
* 100% local, no tracking

## Installation

```bash
go install github.com/sturdy-dev/hackagotchi@latest
```

### Shell Integration (required)

```bash
# zsh (macOS Default Shelll)
echo "eval \"\$(hackagotchi --zsh)\"" >> ~/.zshrc
eval "$(hackagotchi --zsh)"

# Fish
echo "hackagotchi --fish | source" >> ~/.config/fish/config.fish
hackagotchi --fish | source
```
