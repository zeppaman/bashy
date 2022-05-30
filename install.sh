#!/bin/bash

if  ! [ -x "$(command -v go)" ]; then
    echo "installing go"
    wget -q -O - https://git.io/vQhTU | bash
else
 echo "go ok"
fi

OS="$(uname -s)"
ARCH="$(uname -m)"

case $OS in
    "Linux")
        case $ARCH in
        "x86_64")
            ARCH=amd64
            ;;
        "aarch64")
            ARCH=arm64
            ;;
        "armv6" | "armv7l")
            ARCH=armv6l
            ;;
        "armv8")
            ARCH=arm64
            ;;
        .*386.*)
            ARCH=386
            ;;
        esac
        PLATFORM="linux-$ARCH"
    ;;
    "Darwin")
        PLATFORM="darwin-amd64"
    ;;
esac

if [ -z "$PLATFORM" ]; then
    echo "Your operating system is not supported by the script."
    exit 1
fi

if [ -n "$($SHELL -c 'echo $ZSH_VERSION')" ]; then
    shell_profile="$HOME/.zshrc"
elif [ -n "$($SHELL -c 'echo $BASH_VERSION')" ]; then
    shell_profile="$HOME/.bashrc"
elif [ -n "$($SHELL -c 'echo $FISH_VERSION')" ]; then
    shell="fish"
    if [ -d "$XDG_CONFIG_HOME" ]; then
        shell_profile="$XDG_CONFIG_HOME/fish/config.fish"
    else
        shell_profile="$HOME/.config/fish/config.fish"
    fi
fi

#reload env variables
source $shell_profile

if [ -f ./bin/task ]; then
    echo "task exists. intall skipped"
else 
    echo "task does not exist. installing"
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
fi

#compile the tool

bin/task  build-debug
cp ./out/debug/bashy /usr/bin/bashy
