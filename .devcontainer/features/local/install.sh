#!/bin/sh

set -eux
apt-get update
apt-get install -y --no-install-recommends \
    neovim \
    bash-completion
rm -rf /var/lib/apt/lists/*
sed -i \
    -e 's/# export LS_OPTIONS/export LS_OPTIONS/' \
    -e 's/# eval "$(dircolors)"/eval "$(dircolors)"/' \
    -e 's/# alias ls=/alias ls=/' \
    -e 's/# alias ll=/alias ll=/' \
    -e 's/removecolor}/removecolor}\\n/' \
    ~/.bashrc
echo '
set -o vi
export GIT_EDITOR=vi

[[ $PS1 && -f /usr/share/bash-completion/bash_completion ]] && \
    . /usr/share/bash-completion/bash_completion
' >> ~/.bashrc
