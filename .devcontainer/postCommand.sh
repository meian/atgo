#!/bin/bash

u=$1

git config --global codespaces-theme.hide-status 1
sudo mkdir -p /go/pkg
sudo chown vscode:golang /go/pkg
sudo mkdir -p /home/$u/.cache/go-build
sudo chown vscode:vscode /home/$u/.cache/go-build

cat <<EOF > /home/$u/.aliases
alias ll='ls -l --color=auto'
EOF

echo ". /home/$u/.aliases" >> /home/$u/.bashrc

cat <<'GIT' >> /home/$u/.bashrc
source /usr/share/bash-completion/completions/git
source /etc/bash_completion.d/git-prompt
GIT_PS1_SHOWDIRTYSTATE=true
GIT

cat <<'PROMPT' >> /home/$u/.bashrc
export PS1='$(export XIT=$? && echo -n "\[\033[0;32m\]\u " && [ "$XIT" -ne "0" ] && echo -n "\[\033[1;31m\]➜" || echo -n "\[\033[0m\]➜" \
) \[\033[1;34m\]\w\[\033[31m\]$(__git_ps1)\[\033[00m\]\$ '
export PROMPT_DIRTRIM=2
PROMPT


LSCRIPT="$(cd $(dirname $0); pwd)/postCommand.local.sh"
if [ -x "$LSCRIPT" ]; then
    $LSCRIPT $u
else
    echo "$LSCRIPT is not executable"
fi