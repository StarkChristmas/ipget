# IPGet - Network Interface Information Tool for macOS

[English](README.md) | [中文](README_zh.md)

IPGet is a lightweight command-line utility designed to simplify the process of retrieving network interface information on macOS systems. It provides a clean and focused output of active network interfaces' IP addresses, both internal and external, making it easier to access network information compared to the verbose output of the native `ifconfig` command.

## Features

- Quick access to active network interface information
- Display of both internal (IPv4) and external IP addresses
- Clean and focused output format
- Simple installation and usage

## System Requirements

- macOS (Only)

## Installation

1. Download the binary file
2. Move it to your system's binary directory:
   ```bash
   mv ip /usr/local/bin/
   ```
3. Make it executable:
   ```bash
   chmod +x /usr/local/bin/ip
   ```

## Usage

Simply run the command in your terminal:
```bash
ip
```

## Example Output

![Example Output](public/QQ20250116-155241.png)

## Why IPGet?

The native `ifconfig` command on macOS provides comprehensive information about all network interfaces, which can be overwhelming when you only need to check the IP address of your active network connection. IPGet simplifies this process by:

- Filtering out inactive network interfaces
- Focusing on essential IP information
- Providing a cleaner, more readable output
- Making it easier to quickly identify your current network configuration

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
