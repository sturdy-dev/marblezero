# 🐈 Marble Zero – The Hackers Pet

![](./demo/intro.gif)

* Code to keep your Marble happy
* Hundreds of achivements
* Daily quests _(coming soon)_
* 100% local, no tracking

## Installation

```bash
go install github.com/sturdy-dev/marblezero@latest
```

### Shell Integration (required)

```bash
# zsh (macOS Default Shelll)
echo "eval \"\$(marblezero --zsh)\"" >> ~/.zshrc
eval "$(marblezero --zsh)"

# Fish
echo "marblezero --fish | source" >> ~/.config/fish/config.fish
marblezero --fish | source
```

Feed your Marble by running commands on the command line. The shell integration registers a pre-exec hook to automatically run hackagotchi when you run a command. All processing is done on your device! 

## Help

Press `h` to show instructions and command on the _device_. 
