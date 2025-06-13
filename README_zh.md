# IPGet - macOS 网络接口信息工具

[English](README.md) | [中文](README_zh.md)

IPGet 是一个轻量级的命令行工具，旨在简化 macOS 系统上获取网络接口信息的过程。它提供了活动网络接口的 IP 地址（包括内网和外网）的清晰输出，相比原生的 `ifconfig` 命令的冗长输出，更容易获取网络信息。



## 功能特点

- 快速访问活动网络接口信息
- 显示内网（IPv4）和外网 IP 地址
- 清晰简洁的输出格式
- 简单的安装和使用方法

## 系统要求

- macOS（仅支持）

## 安装方法

### 通过Homebrew安装(推荐)

本项目已发布至 Homebrew，您可以通过以下命令轻松安装：
```
brew install ip
```

> 由于已发布到Homebrew为避免后续与其他人发布名称冲突，Homebrew上的名称为ipget

### 编译安装

克隆项目到本地
```
git clone https://github.com/StarkChristmas/ipget
```
安装编译环境`Golang 1.24`
```
brew install go
```
```
cd ipget && go run -o ip ./main.go
```
移动至`/usr/local/bin`下并赋予执行权限
```
 mv ip /usr/local/bin/ && chmod +x /usr/local/bin/ip
```

### 下载构建好的二进制可执行文件

`AppleSlicon`版本的:
```
wget $(curl -s https://api.github.com/repos/StarkChristmas/ipget/releases/latest \
  | jq -r '.assets[] | select(.name | test("arm64.*\\.tar\\.gz$")) | .browser_download_url')
```

`Intel`版本的:
```
wget $(curl -s https://api.github.com/repos/StarkChristmas/ipget/releases/latest \
  | jq -r '.assets[] | select(.name | test("x86_64.*\\.tar\\.gz$")) | .browser_download_url')
```

移动至`/usr/local/bin`下并赋予执行权限
如果是`AppleSlicon`执行这个
```
 mv ip_arm64 /usr/local/bin/ip && chmod +x /usr/local/bin/ip
```
如果是`Intel`执行这个
```
 mv ip_x86_64 /usr/local/bin/ip && chmod +x /usr/local/bin/ip
```

## 使用方法

在终端中直接运行命令：
```bash
ip
```


## 示例输出

![示例输出](public/QQ20250116-155241.png)

## 为什么选择 IPGet？

macOS 原生的 `ifconfig` 命令会显示所有网络接口的详细信息，当你只需要查看当前活动网络连接的 IP 地址时，这些信息可能会显得过于冗长。IPGet 通过以下方式简化了这个过程：

- 过滤掉非活动网络接口
- 专注于必要的 IP 信息
- 提供更清晰、更易读的输出
- 使快速识别当前网络配置变得更容易

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 贡献指南

欢迎贡献代码！请随时提交 Pull Request。 