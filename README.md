<h1 align="center">
  	alohomora
</h1>

<p align="center">
    <img alt="AppVeyor" src="https://img.shields.io/appveyor/ci/steps0x29a/alohomora?style=plastic">
    <img alt="GitHub issues" src="https://img.shields.io/github/issues/steps0x29a/alohomora?style=plastic">
    <img alt="GitHub tag (latest by date)" src="https://img.shields.io/github/v/tag/steps0x29a/alohomora?style=plastic">    
    <img alt="GitHub" src="https://img.shields.io/github/license/steps0x29a/alohomora?style=plastic">
    <img alt="Contributions" src="https://img.shields.io/badge/contributions-welcome-brightgreen?style=plastic">
  <a href="https://goreportcard.com/report/github.com/steps0x29a/alohomora"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/steps0x29a/alohomora?style=plastic&fuckgithubcache=1"></a>
</p>

# alohomora
alohomora is a distributed cracking utility. A server provides cracking jobs to all clients connected to it. The clients do the actual heavy lifting and report back to the server.

<p align="center">
  <img alt="alohomora" src="./.github/screenshot-2.png">
  <b>alohomora running with an example handshake</b><br>
</p>

## Why?
Because I wanted to learn [Go](https://golang.org "Go's website") and to see whether a distributed cracking utility like this was possible. Turns out it is!

# Current version
Current version is 0.4, which is the first version I feel comfortable making public. It's far from perfect, I know that. But it actually works.
The current feature set includes:

 * Cracking WPA2 handshakes in a distributed way
 * Limited reporting (stdout, XML, JSON)
 * Very limited REST interface

# History
Initially, alohomora was meant to only crack WPA2 handshakes using [aircrack-ng](https://github.com/aircrack-ng/aircrack-ng) with a bruteforce approach. But while I was developing it, it dawned on me that it could also be used to crack hashes, encrypted files and much more, so this is planned for future versions of it.

# How does it work?
Let's say you have [obtained a WPA2 handshake](https://github.com/bettercap/bettercap). You then start an alohomora server tasked with this particular handshake (just give it the handshake file). Adjust character set and password length to your needs and run it:

    ./alohomora -server -verbose -ip <your ip> -port <port> -charset abcdefghijklmnopqrstuvwxyz -length 8 -target <path to pcap file>

This will start the server, listening for connections to your IP address. Omitting the ip parameter will make it listen on 0.0.0.0. 

    ./alohomora -server -verbose -ip 127.0.0.1 -port 6666 -charset abcdefghijklmnopqrstuvwxyz -length 8 -target <path to pcap file>

The above command will bruteforce all lowercase (a-z) 8-character passwords, e.g. `aaaaaaaa` to `zzzzzzzz`. Each client will be given up to 10000 passwords per iteration (default job size).

In order to connect a client, simply give it the ip and port:

    ./alohomora -port <port> -ip <server ip>

To connect a client to the server from the example above:

    ./alohomora -port 6666 -ip 127.0.0.1

That's all, actually.

# Installing
Clone the repository, then run it:

    go run main.go --help
    
Or simply download a release and run it:

    ./alohomora --help

alohomora clients require `aircrack-ng` to crack handshakes. If aircrack-ng is in your PATH, everything should work out of the box. If not, you can set the `AIRCRACK` environment variable (set it to the full path of `aircrack-ng`). 

    export AIRCRACK=/path/to/aircrack-ng/executable

To install `aircrack-ng`, run

    sudo apt install aircrack-ng
    
<p><sub>(or use whatever package manager your system uses)</sub></p>

# Compatibility
`alohomora` is designed to work on Linux-based systems. You can run it inside [WSL](https://docs.microsoft.com/en-us/windows/wsl/about) on Windows, though, it works just fine. I recommend getting the new [Windows Terminal Preview](https://github.com/microsoft/terminal) if you want to run `alohomora` in WSL.

# Legal disclaimer
As you might have guessed, **cracking WPA2 passphrases might be illegal**. Do not use alohomora on handshakes that you don't have the permission to crack! I will not be held responsible for anything illegal you do with this tool!
Also, use alohomora at your own risk! I will not be held responsible for any damage caused by it.

# Roadmap

There's so much to do...

 * ~~Improve traffic use by omitting handshakes when talking to known clients~~
 * Add capabilities for cracking hashes
 * Add capabilities for cracking archive passwords
 * ~~Improve reporting capabilities~~
 * Improve code quality
 * ~~Comment more~~
 * Write more tests
 * ~~Write useful documentation~~
 * Make an awesome logo (anyone?)
 * Improve REST interface
 * Make an awesome web interface for it
 * ~~Improve github presence~~
 * Make a docker container 
 * Create a working cross-platform build process

# Version history

## Version 0.4
Lots of improvements, but there's still a lot of room for even more.

## Version 0.1 to 0.3
A lot of experimentation, as I wrote this tool for the purpose of learning Go.

# Contributing
I develop alohomora in the rare moments that I have time to do so. If you want to contribute, please create a pull request and we'll see. Feel free to contact me beforehand.
