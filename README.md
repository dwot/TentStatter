# TentStatter

TentStatter is a tool designed to record and monitor sensor values from AC Infinity grow tents. Using an API token, this application collects stats such as temperature, humidity, and port statuses, then logs them to text files for easy monitoring.

## Features

- Collects data from the AC Infinity controller, including temperature, humidity, and port statuses.
- Logs the recorded data to individual text files for easy access.
- Cross-platform support for Windows, Linux, and macOS.

## Getting Started

### 1. **Download and Install**

Choose the appropriate build for your operating system from the [Releases](https://github.com/dwot/TentStatter/releases) section:

Also, download the corresponding `.md5` file to verify the integrity of your download.

### 2. **Verify Download (Optional)**

Before unpacking, you can verify the integrity of your download by comparing the MD5 checksum of your downloaded file with the one provided:

```bash
# For Linux/macOS
md5sum tentstatter-v0.0.1-YOUR_OS_HERE.tar.gz

# For Windows (in PowerShell)
Get-FileHash -Algorithm MD5 .\tentstatter-v0.0.1-windows-amd64.zip
```

The output should match the contents of the corresponding `.md5` file.

### 3. **Unpack and Setup**

- Uncompress the downloaded file.
- Ensure you have the required `config.properties` file set up with the necessary parameters (timezone, API token, start date).
- Run the TentStatter executable.

### 4. **Monitor your Grow Tent**

Once running, TentStatter will periodically fetch data from the AC Infinity controller and log it to the respective text files. Ensure you have set up the required API token using [ACScraper](https://github.com/dwot/ACScraper) or another method.
