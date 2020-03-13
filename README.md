# ticker-term

Live futures quotes in your terminal.


![ScreenShot](screenshots/1.png)


![ScreenShot](screenshots/2.png)


![ScreenShot](screenshots/3.png)


## Requirements

Linux or macOS. It *might* work with [Windows 10's new bash shell](https://www.howtogeek.com/249966/how-to-install-and-use-the-linux-bash-shell-on-windows-10/).

go 1.14


## Installation

#### Option 1: Build from source (recommended)

1. [Install go for your platform](https://golang.org/dl/).

2. From a terminal, run `go get -u github.com/zpkg/ticker-term`


#### Option 2: Pre-compiled binary

Grab a binary for Linux or macOS on the releases page.


## Usage

`ticker-term`


## FAQ

1. Why can't I see the inactive tickers? Why are the lines spaced inconsistently?

   Inactive tickers are greyed out if there are no price updates for an extended
   period of time. You may have to adjust your terminal emulator's color settings
   if the greyed out items are not visible.

2. I keep getting a message about "sleep mode". What is that?

   If there are no price updates to any of the tickers for several minutes, a notification
   will appear stating that the application has entered sleep mode. While in sleep mode,
   checks for new data occur only once every sixty seconds.


## Disclaimer

All data sourced from Investing.com. This package is not in any way affiliated
with Investing.com or its subsidiaries, parents, or affiliates. No claims are 
made with regard to the accuracy of the data. This package is for research 
purposes only and is not intended to provide investment or trading advice.

