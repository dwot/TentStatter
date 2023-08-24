# TentStatter

TentStatter is a tool designed to fetch status and sensor readings from the AC Infinity controller used in AC Infinity grow tents. Using a captured API token, it pulls real-time data and writes essential stats, such as temperature, humidity, and device port status, to individual text files for easy consumption.

## Features

- Retrieves and logs temperature, humidity, and device port status.
- Writes key statistics to individual text files, e.g., `days.txt`, `temperatureF.txt`, `humidity.txt`, and `port_X.txt` where `X` is the port number.
- Lightweight and easy to set up.
- Continuous monitoring with configurable update intervals.

## Prerequisites

- Go (Golang) installed on your machine.
- A valid API token captured using a tool like [ACScraper](https://github.com/dwot/acscraper).

## Setup and Usage

1. **Clone the repository**:

   ```bash
   git clone https://github.com/dwot/TentStatter.git
   cd TentStatter
   ```

2. **Configure the application**:

   Modify the `config.properties` file with your desired settings:

    - `tz`: Your desired timezone (e.g., `America/New_York`).
    - `token`: Your API token.
    - `start_date`: The start date for your monitoring in `YYYY-MM-DD` format.

3. **Run the application**:

   ```bash
   go run main.go
   ```

4. **View the statistics**:

   The application writes various statistics to separate text files in the same directory:

    - `days.txt` displays the number of days since the `start_date`.
    - `temperatureF.txt` displays the current temperature in Fahrenheit.
    - `humidity.txt` displays the current humidity percentage.
    - `port_X.txt` displays the status for the device port `X`.

5. **Stop the application**:

   Simply press `CTRL + C` to stop the TentStatter application.
