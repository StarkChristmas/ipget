# macOS 获取已激活网卡的内外网IP地址
由于`macOS`的`ifconfig`输出过于复杂，输出的是所有网卡的信息，不利于快速查看当前已激活网卡的相关信息，大多数时候我们只需要看到当前网卡的`IPv4`，或者是当前网卡的外网地址。
## 系统支持
* macOS

> 仅支持`macOS`

## 使用
```bash
mv ip /usr/local/bin/
chmod +x /usr/local/bin/ip
```
![](public/QQ20250116-155241.png)
