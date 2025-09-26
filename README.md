# olsync

[English Documentation](./README.md) | [中文文档](./README_zh.md)

**olsync** is a simple command-line tool for **synchronizing your Overleaf projects** to your local machine, making it easy to back up and manage LaTeX projects offline.

---

# Usage

1. **Copy the configuration template**

```bash
cp config.template.yaml config.yaml
```

2. **Edit config.yaml**

Open the config.yaml file in your favorite editor and fill in the required fields:

overleaf.baseUrl – Your Overleaf base URL (usually https://www.overleaf.com)

overleaf.cookies – Your authentication cookies (get them from your browser developer tools)

backup.path – The folder where backups will be stored

backup.keep_last – How many recent backups to keep per project

backup.schedule – Cron expression for automatic backups

3. **Build and Run**

```bash
make build
```

The compiled binary will be placed in the `./bin` folder.
For example, on Linux you might see:

```bash
./bin/olsync-linux-amd64
```

> The default is to use the config.yaml file in the current path.
