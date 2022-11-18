autoload -Uz add-zsh-hook

function marblezero_preexec() {
    marblezero --import-single "$1"
}

add-zsh-hook preexec marblezero_preexec