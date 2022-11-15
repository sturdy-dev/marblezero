autoload -Uz add-zsh-hook

function hackagotchi_preexec() {
    hackagotchi --import-single "$1"
}

add-zsh-hook preexec hackagotchi_preexec