# olsync

[英文文档](./README.md) | [中文文档](./README_zh.md)

**olsync** 是一个简单的命令行工具，用于**将你的 Overleaf 项目**同步到本地计算机，从而轻松备份和管理离线 LaTeX 项目。

---

# 使用方法

1. **复制配置模板**

```bash
cp config.template.yaml config.yaml
```

2. **编辑 config.yaml**

使用您常用的编辑器打开 config.yaml 文件，并填写以下必填字段：

overleaf.baseUrl – 您的 Overleaf 基础 URL（通常为 https://www.overleaf.com）

overleaf.cookies – 您的身份验证 cookie（可从浏览器开发者工具获取）

backup.path – 备份存储文件夹

backup.keep_last – 每个项目保留的最近备份数量

backup.schedule – 自动备份的 Cron 表达式

3. **构建并运行**

```bash
make build
```

编译后的二进制文件将放置在 `./bin` 文件夹中。
例如，在 Linux 上，您可能会看到：

```bash
./bin/olsync-linux-amd64
```

> 默认使用当前路径下的 config.yaml 文件。

您可以通过传递其路径作为参数来指定自定义配置文件：
```bash
./bin/olsync-linux-amd64 ./config2.yaml
```