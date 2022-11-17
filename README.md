# Hackagotchi – A pet for your computer

![](./demo/intro.gif)

* Code to keep your hackagotchi happy
* 100% local, no tracking
* Daily quests _(coming soon)_
* Hundreds of achivements

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

Feed your hackagotchi by running commands on the command line. The shell integration registers a pre-exec hook to automatically run hackagotchi when you run a command. All processing is done on your device! 

## Help

Press `h` to show instructions and command on the _device_. 
