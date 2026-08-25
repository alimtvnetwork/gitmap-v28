#!/bin/bash

install_dependencies() {
  add-apt-repository ppa:apt-fast/stable && \
  add-apt-repository ppa:git-core/ppa && \
  apt update -y && \
  apt-get update -y && \
  apt upgrade -y && \
  apt-get install -y vim build-essential wget nano curl file git libpcre3 libpcre3-dev zlib1g zlib1g-dev libssl-dev git-core sshpass && \
  apt install -y zsh git-lfs snapd && \
  git version && \
  git lfs --version
  # apt install -y nfs-server && \   
}

# Call the function
install_dependencies
