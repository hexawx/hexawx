# HexaWX

**HexaWX** est un orchestrateur météo modulaire écrit en Go. Il permet de collecter, traiter et exporter des données météorologiques via un système de plugins dynamiques pilotables en temps réel via une console SSH sécurisée.

![Go CI](https://github.com/hexawx/hexawx/actions/workflows/go.yml/badge.svg)

## 🚀 Caractéristiques

* **Architecture Modulaire** : Drivers (entrées) et Exporters (sorties) interchangeables.
* **Console SSH Intégrée** : Gérez vos plugins et surveillez vos données à distance.
* **Gestionnaire de Plugins** : Installation, mise à jour et suppression automatique via un catalogue distant.
* **Multi-plateforme** : Binaires natifs pour Linux (AMD64/ARM), Windows et macOS (Apple Silicon).
* **Sécurisé** : Authentification par clés SSH (RSA/Ed25519).

### Flux

```plaintext
Capteurs (Drivers)  ──▶  HexaWX Core  ──▶  Sorties (Exporters)
      (ARM/x64)          (SSH Admin)           (JSON/Terminal)
```

## 🛠 Installation

### 1. Télécharger le binaire
Récupérez la version correspondant à votre architecture dans les [Releases](https://github.com/hexawx/hexawx/releases).

```bash
# Exemple pour Linux AMD64
wget https://github.com/hexawx/hexawx/releases/download/v1.0.0/hexawx_linux_amd64
chmod +x hexawx_linux_amd64
```

### 2. Configuration (config.yaml)

L'orchestrateur se configure via un fichier config.yaml. Ce fichier permet de définir le comportement du moteur et l'emplacement des modules.

```yaml
server:
  # Intervalle de rafraîchissement des données/métriques
  interval: 10s
  # Dossier contenant les binaires des plugins installés
  plugin_dir: "./plugins"
```

* `interval` : Définit la cadence à laquelle le serveur interroge les plugins ou rafraîchit les statistiques internes. Supporte les unités de temps Go (ex: `10s`, `1m`, `1h`).

* `plugin_dir` : Chemin relatif ou absolu vers le répertoire de stockage des drivers et exporters. C'est ici que la fonction `resolveURL` téléchargera les nouveaux binaires.

### 2. Configurer les accès (Clés SSH)

Créez un fichier `./data/users.json` pour autoriser des administrateurs :

```json
[
  {
    "username": "votre_nom",
    "pub_key": "ssh-ed25519 AAAAC3Nza..."
  }
]
```

### 3. Lancer le serveur

```bash
./hexawx_linux_amd64 start [--config config.yaml]
```

## 💻 Utilisation de la Console

Connectez-vous à l'orchestrateur via SSH (port par défaut : 2233) :

```bash
ssh -p 2233 votre_nom@localhost
```

Commandes disponibles :
* catalog : Affiche les plugins disponibles pour votre architecture.
* install <plugin> : Télécharge et installe un nouveau plugin.
* list : Affiche l'état, la version et l'uptime des plugins chargés.
* start/stop <plugin> : Pilote l'exécution des modules.
* stats : Affiche les métriques de traitement de données.
* help : Affiche les commandes disponibles

## 🔌 Écosystème de Plugins

**HexaWX** utilise un registre centralisé. Seuls les plugins compatibles avec votre architecture (OS-Arch) sont visibles dans le catalogue pour garantir une stabilité maximale.

| Type | Plugin | Description |
| :--- | :--- | :--- |
| Driver | dummy-driver | Simulateur de données météo pour test.|
| Exporter | stdout-exporter | Affiche les données reçues dans la console.|

## 🏗 Développement

Pour compiler le projet vous-même :

```bash
git clone https://github.com/hexawx/hexawx.git
cd hexawx
go build -o hexawx ./cmd/server/main.go
go test ./core/...
```

© 2026 F. Colinet - **HexaWX** est distribué sous licence **CC BY-NC-SA 4.0**.