package handshakes

import (
	"errors"

	"github.com/google/gopacket/layers"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

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
