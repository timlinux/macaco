# Installation

## Pre-built Binaries

Download the latest release from the [GitHub Releases](https://github.com/timlinux/macaco/releases) page.

### Linux

```bash
# AMD64
wget https://github.com/timlinux/macaco/releases/latest/download/macaco-linux-amd64
chmod +x macaco-linux-amd64
sudo mv macaco-linux-amd64 /usr/local/bin/macaco

# ARM64
wget https://github.com/timlinux/macaco/releases/latest/download/macaco-linux-arm64
chmod +x macaco-linux-arm64
sudo mv macaco-linux-arm64 /usr/local/bin/macaco
```

### macOS

```bash
# Intel
wget https://github.com/timlinux/macaco/releases/latest/download/macaco-darwin-amd64
chmod +x macaco-darwin-amd64
sudo mv macaco-darwin-amd64 /usr/local/bin/macaco

# Apple Silicon
wget https://github.com/timlinux/macaco/releases/latest/download/macaco-darwin-arm64
chmod +x macaco-darwin-arm64
sudo mv macaco-darwin-arm64 /usr/local/bin/macaco
```

### Windows

Download `macaco-windows-amd64.exe` from the releases page and add it to your PATH.

## Package Managers

### DEB (Debian/Ubuntu)

```bash
wget https://github.com/timlinux/macaco/releases/latest/download/macaco_VERSION_amd64.deb
sudo dpkg -i macaco_VERSION_amd64.deb
```

### RPM (Fedora/RHEL)

```bash
wget https://github.com/timlinux/macaco/releases/latest/download/macaco-VERSION-1.x86_64.rpm
sudo rpm -i macaco-VERSION-1.x86_64.rpm
```

## From Source

### Using Go

```bash
go install github.com/timlinux/macaco/cmd/macaco@latest
```

### Building Manually

```bash
git clone https://github.com/timlinux/macaco
cd macaco
make build
./bin/macaco
```

### Using Nix

```bash
# Run directly
nix run github:timlinux/macaco

# Or install
nix profile install github:timlinux/macaco
```

## Verify Installation

```bash
macaco --version
```

You should see output like:

```
MoCaCo - Motion Capture Combatant
Version: 1.0.0
```
