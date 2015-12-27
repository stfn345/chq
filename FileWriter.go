package main

import (
	"log"
	"os"
)

// A FileWriter will write all packets it
// receives to a file.
// It passes through TS packets unmodified.
type FileWriter struct {
	file *os.File
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("FileWriter", NewFileWriter)
}

func NewFileWriter(fname string) (*FileWriter, error) {
	// try to open file
	fh, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	node := &FileWriter{}
	node.file = fh
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *FileWriter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		node.PktsIn++
		node.Send(pkt)

		node.file.Write(pkt.bytes)
	}
}

func (node *FileWriter) closeDown() {
	node.file.Close()
	log.Printf("closing down FileWriter to file %s", node.file.Name())
	node.output.Close()
}
