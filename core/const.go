package core

const (
	// MsgIDHello is the ID of the first message exchanged between client and server
	MsgIDHello = "ALOHOMORA_HLO"
	// MsgIDLeave is the ID of the message the server sends to clients it wishes to kick
	MsgIDLeave = "ALOHOMORA_EXT"
	// MsgIDAck confirms a client (sent by server)
	MsgIDAck = "ALOHOMORA_ACK"
	// MsgIDReady is the ID of the message a client sends to the server when it is ready for work
	MsgIDReady = "ALOHOMORA_RDY"
	// MsgIDTask identifies a message containing a crackjob
	MsgIDTask = "ALOHOMORA_TSK"
	// MsgIDDone identifies a message telling the server that a client has finished work.
	MsgIDDone = "ALOHOMORA_FIN"
)

//AlohomoraSuffix is a magic number postfixed to every message exchanged between clients and server
//var AlohomoraSuffix = []byte{'\x00'}

var AlohomoraSuffix = []byte{10, 10, 0x29, 0x10, 10, 10}
