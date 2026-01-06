# Installation Guide – journal-cli

This guide explains how to install **journal-cli** on macOS, Linux, and Windows using the prebuilt binaries available in the GitHub Releases.

👉 Releases: https://github.com/ops295/journal-cli/releases

---

## 🍎 macOS

### macOS Apple Silicon (M1 / M2 / M3)

1. Download the binary:
   ```bash
   curl -LO https://github.com/ops295/journal-cli/releases/download/v0.2.0/journal-darwin-arm64
   ```

2. Verify checksum (optional):
   ```bash
   shasum -a 256 journal-darwin-arm64
   ```
   Expected:
   ```
   b04ce6cf2e8a0dd88f43ce2cc26fd0534bd2716118bcf181c40d58e82c102e3e
   ```

3. Make it executable:
   ```bash
   chmod +x journal-darwin-arm64
   ```

4. Move to PATH:
   ```bash
   sudo mv journal-darwin-arm64 /opt/homebrew/bin/journal
   ```

5. Verify installation:
   ```bash
   journal --help
   ```

---

### macOS Intel (x86_64)

1. Download:
   ```bash
   curl -LO https://github.com/ops295/journal-cli/releases/download/v0.2.0/journal-darwin-amd64
   ```

2. Make executable and install:
   ```bash
   chmod +x journal-darwin-amd64
   sudo mv journal-darwin-amd64 /usr/local/bin/journal
   ```

3. Verify:
   ```bash
   journal --help
   ```

---

## 🐧 Linux

### Linux x86_64 (amd64)

```bash
curl -LO https://github.com/ops295/journal-cli/releases/download/v0.2.0/journal-linux-amd64
chmod +x journal-linux-amd64
sudo mv journal-linux-amd64 /usr/local/bin/journal
journal --help
```

---

### Linux ARM64 (aarch64)

```bash
curl -LO https://github.com/ops295/journal-cli/releases/download/v0.2.0/journal-linux-arm64
chmod +x journal-linux-arm64
sudo mv journal-linux-arm64 /usr/local/bin/journal
journal --help
```

---

## 🪟 Windows

### Windows (amd64)

1. Download the executable from the release page:
   ```
   journal-windows-amd64.exe
   ```

2. Rename it to:
   ```
   journal.exe
   ```

3. Move it to a directory in your PATH, for example:
   ```
   C:\Program Files\journal\
   ```

4. Add the directory to PATH (if not already added).

5. Verify in Command Prompt or PowerShell:
   ```powershell
   journal --help
   ```

---

## 🔍 Verify Installation

Run:

```bash
journal --version
```

or:

```bash
journal --help
```

---

## ❗ Troubleshooting

### macOS: "cannot be opened because the developer cannot be verified"

Run:

```bash
xattr -d com.apple.quarantine journal
```

or allow it from:
**System Settings → Privacy & Security**

---

## 📦 Uninstall

```bash
sudo rm $(which journal)
```

---

## 🙌 Contributing

If you'd like to improve installation, packaging, or Homebrew support:

* Fork the repo
* Create a PR
* Suggestions are welcome!

---

Happy journaling 📝
