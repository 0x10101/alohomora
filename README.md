# alohomora
alohomora is a distributed password cracking utility. A server provides cracking jobs to all clients connected to it. The clients do the actual heavy lifting and report back to the server.

# History
Initially, alohomora was meant to only crack WPA2 handshakes using aircrack-ng and a bruteforce approach. But while I was developing it, it dawned on me that it could also be used to crack hashes, so this is planned for future versions of it.

# How does it work?
Let's say you have obtained a WPA2 handshake. You then start alohomora in server mode, providing it with the handshake PCAP file as well as the parameters for cracking the passphrase, i.e. the charset to bruteforce:

    ./alohomora -server -port 7890 -ip <external ip> -target /path/to/<ESSID>_<BSSID>.pcap -charset abcdefghijklmnopqrstuvwxyz -length 8 -jobsize 10000
This will start the server on port 7890, listening for connections to your external IP address. Omitting the ip parameter will make it listen on localhost. 

The PCAP file is parsed by alohomora. It tries to find both the ESSID and BSSID in it to pass them to the clients. If your handshake data does not contain either an ESSID or a BSSID, name the file `<ESSID>_<BSSID>.pcap` in order for alohomora to be able to recognize it. 

The above command will bruteforce all lowercase (a-z) 8-character passwords, e.g. `aaaaaaaa` to `zzzzzzzz`. Each client will be given up to 10000 passwords per iteration.

In order to connect a client, simply give it the ip and port:

    ./alohomora -port 7890 -ip <server ip>

That's all, actually.

# Legal disclaimer
As you might have guessed, cracking WPA2 passphrases might be illegal. Do not use alohomora on handshakes that you don't have the permission to crack! I will not be held responsible for anything illegal you do with this tool!
Also, use alohomora at your own risk! I will not be held responsible for any damage caused by it.
