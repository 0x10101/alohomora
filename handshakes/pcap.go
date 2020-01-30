package handshakes

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/gopacket/layers"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

/*handle, err = pcap.OpenOffline(pcapFile)
    if err != nil { log.Fatal(err) }
    defer handle.Close()

    // Loop through packets in file
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        fmt.Println(packet)
	}*/

// HandshakeInfo tries to parse a *.pcap file and obtain the BSSID and the ESSID of
// the handshake. It will return (essid, bssid) in that order. If either the ESSID
// or the BSSID could not be determined, an error is returned.
func HandshakeInfo(pcapFile string) (string, string, error) {
	handle, err := pcap.OpenOffline(pcapFile)
	if err != nil {
		return "", "", err
	}

	defer handle.Close()

	source := gopacket.NewPacketSource(handle, handle.LinkType())

	essid, bssid := "", ""

	for packet := range source.Packets() {
		if essid == "" {
			essid = extractESSID(packet)
		}

		if bssid == "" {
			bssid = extractBSSID(packet)
		}

		if len(essid) > 0 && len(bssid) > 0 {
			break
		}
	}

	if essid == "" || bssid == "" {
		return essid, bssid, errors.New("Unable to find ESSID or BSSID in PCAP file")
	}

	return essid, bssid, nil

}

func extractBSSID(packet gopacket.Packet) string {
	dot11Layer := packet.Layer(layers.LayerTypeDot11)
	if dot11Layer != nil {
		dot11, _ := dot11Layer.(*layers.Dot11)
		return dot11.Address3.String()
	}
	return ""
}

func extractESSID(packet gopacket.Packet) string {
	dot11Info := packet.Layer(layers.LayerTypeDot11InformationElement)
	if dot11Info != nil {
		dot11info, _ := dot11Info.(*layers.Dot11InformationElement)
		if dot11info.ID == layers.Dot11InformationElementIDSSID {
			return string(dot11info.Info)
		}
	}
	return ""
}

func openPCAPFile(pcapFile string) (*pcap.Handle, error) {
	handle, err := pcap.OpenOffline(pcapFile)
	return handle, err
}

func TestPcap(pcapFile string) {
	fmt.Println(pcapFile)
	handle, err := pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatalf("PCAP Error: %s", err)
	}
	defer handle.Close()

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println("There are", len(source.Packets()), "packets in the pcap file")

	packet, ok := <-source.Packets()
	if !ok {
		log.Fatalf("Whoops\n")
	}

	//fmt.Println(packet.Dump())

	dot11Layer := packet.Layer(layers.LayerTypeDot11)
	if dot11Layer == nil {
		log.Fatalf("Decoding layer failed")
	}
	dot11, _ := dot11Layer.(*layers.Dot11)
	fmt.Println("BSSID:", dot11.Address3)

	dot11Info := packet.Layer(layers.LayerTypeDot11InformationElement)
	if dot11Info != nil {
		dot11info, _ := dot11Info.(*layers.Dot11InformationElement)
		if dot11info.ID == layers.Dot11InformationElementIDSSID {
			fmt.Printf("SSID: %q\n", dot11info.Info)
		}
	}

	/*for _, layer := range packet.Layers() {
		fmt.Println("Layer:", layer.LayerType().String(), "", layer.)
	}*/

	/*for packet := range source.Packets() {
		fmt.Println(packet)
	}*/
}
