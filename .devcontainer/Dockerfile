FROM mcr.microsoft.com/devcontainers/go:1.25

# Installation de SQLite et des outils de build
RUN rm -f /etc/apt/sources.list.d/yarn.list \
    && apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*
