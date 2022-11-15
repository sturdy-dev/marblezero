# Hackagotchi – A pet for your computer

<img width="610" alt="Screenshot 2022-11-15 at 10 49 51" src="https://user-images.githubusercontent.com/47952/201888297-166bb6db-7494-42ec-9d64-136c44c4bc61.png">

* Daily quests
* Hundreds of achivements
* Keep coding, to keep your hackagotchi happy


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