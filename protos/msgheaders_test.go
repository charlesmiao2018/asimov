// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"bytes"
	"github.com/AsimovNetwork/asimov/common"
	"io"
	"reflect"
	"testing"
)

// TestHeaders tests the MsgHeaders API.
func TestHeaders(t *testing.T) {
	// Ensure the command is expected value.
	wantCmd := "headers"
	msg := NewMsgHeaders()
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgHeaders: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure headers are added properly.
	bh := &blockOne.Header
	msg.AddBlockHeader(bh)
	if !reflect.DeepEqual(msg.Headers[0], bh) {
		t.Errorf("AddHeader: wrong header - got %v, want %v",
			msg.Headers,
			bh)
	}

	// Ensure adding more than the max allowed headers per message returns
	// error.
	var err error
	for i := 0; i < MaxBlockHeadersPerMsg+1; i++ {
		err = msg.AddBlockHeader(bh)
	}
	if reflect.TypeOf(err) != reflect.TypeOf(&MessageError{}) {
		t.Errorf("AddBlockHeader: expected error on too many headers " +
			"not received")
	}
}

// TestHeadersWire tests the MsgHeaders protos encode and decode for various
// numbers of headers and protocol versions.
func TestHeadersWire(t *testing.T) {
	hash := mainNetGenesisHash
	//bits := uint32(0x1d00ffff)
	//nonce := uint32(0x9962e301)
	bh := NewBlockHeader(1, &hash)
	bh.Version = blockOne.Header.Version
	bh.Timestamp = blockOne.Header.Timestamp
	bh.StateRoot = blockOne.Header.StateRoot
	bh.MerkleRoot = blockOne.Header.MerkleRoot
	bh.Height = 1
	bh.GasLimit = 10000000
	bh.GasUsed = 9000000
	bh.CoinBase[0] = 0x66
	// Empty headers message.
	noHeaders := NewMsgHeaders()
	noHeadersEncoded := []byte{
		0x00, // Varint for number of headers
	}

	// Headers message with one header.
	oneHeader := NewMsgHeaders()
	oneHeader.AddBlockHeader(bh)
	oneHeaderEncoded := []byte{
		0x01,                   // VarInt for number of headers.
		0x01, 0x00, 0x00, 0x00, // Version 1
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, // PrevBlock
		0x3b, 0xa3, 0xed, 0xfd, 0x7a, 0x7b, 0x12, 0xb2,
		0x7a, 0xc7, 0x2c, 0x3e, 0x67, 0x76, 0x8f, 0x61,
		0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
		0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a, // MerkleRoot
		0x29, 0xab, 0x5f, 0x49, 0x00, 0x00, 0x00, 0x00, // Timestamp
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
		0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a, // StateRoot
		0x80, 0x96, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, // GasLimit
		0x40, 0x54, 0x89, 0x00, 0x00, 0x00, 0x00, 0x00, // GasUsed
		0x00, 0x00, 0x00, 0x00, // Round
		0x00, 0x00, // SlotIndex
		0x00, 0x00, // Weight
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // PoaHash
		0x01, 0x00, 0x00, 0x00, // height
		0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, // coinbase
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // SigData:
		0x00, // tx count
	}

	tests := []struct {
		in   *MsgHeaders     // Message to encode
		out  *MsgHeaders     // Expected decoded message
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for protos encoding
		enc  MessageEncoding // Message encoding format
	}{
		// Latest protocol version with no headers.
		{
			noHeaders,
			noHeaders,
			noHeadersEncoded,
			common.ProtocolVersion,
			BaseEncoding,
		},

		// Latest protocol version with one header.
		{
			oneHeader,
			oneHeader,
			oneHeaderEncoded,
			common.ProtocolVersion,
			BaseEncoding,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode the message to protos format.
		var buf bytes.Buffer
		err := test.in.VVSEncode(&buf, test.pver, test.enc)
		if err != nil {
			t.Errorf("VVSEncode #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("VVSEncode #%d\n got: %v want: %v", i,
				buf.Bytes(), test.buf)
			continue
		}

		// Decode the message from protos format.
		var msg MsgHeaders
		rbuf := bytes.NewReader(test.buf)
		err = msg.VVSDecode(rbuf, test.pver, test.enc)
		if err != nil {
			t.Errorf("VVSDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&msg, test.out) {
			t.Errorf("VVSDecode #%d\n got: %v want: %v", i,
				&msg, test.out)
			continue
		}
	}
}

// TestHeadersWireErrors performs negative tests against protos encode and decode
// of MsgHeaders to confirm error paths work correctly.
func TestHeadersWireErrors(t *testing.T) {
	pver := common.ProtocolVersion
	wireErr := &MessageError{}

	hash := mainNetGenesisHash
	bh := NewBlockHeader(1, &hash)
	bh.Version = blockOne.Header.Version
	bh.Timestamp = blockOne.Header.Timestamp
	bh.StateRoot = blockOne.Header.StateRoot
	bh.MerkleRoot = blockOne.Header.MerkleRoot
	bh.Height = 1
	bh.GasLimit = 10000000
	bh.GasUsed = 9000000
	bh.CoinBase[0] = 0x66

	// Headers message with one header.
	oneHeader := NewMsgHeaders()
	oneHeader.AddBlockHeader(bh)
	oneHeaderEncoded := []byte{
		0x01,                   // VarInt for number of headers.
		0x01, 0x00, 0x00, 0x00, // Version 1
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, // PrevBlock
		0x3b, 0xa3, 0xed, 0xfd, 0x7a, 0x7b, 0x12, 0xb2,
		0x7a, 0xc7, 0x2c, 0x3e, 0x67, 0x76, 0x8f, 0x61,
		0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
		0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a, // MerkleRoot
		0x29, 0xab, 0x5f, 0x49, 0x00, 0x00, 0x00, 0x00, // Timestamp
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
		0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a, // StateRoot
		0x80, 0x96, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, // GasLimit
		0x40, 0x54, 0x89, 0x00, 0x00, 0x00, 0x00, 0x00, // GasUsed
		0x00, 0x00, 0x00, 0x00, // Round
		0x00, 0x00, // SlotIndex
		0x00, 0x00, // Weight
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // PoaHash
		0x01, 0x00, 0x00, 0x00, // height
		0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, // coinbase
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // SigData:
		0x00, // tx count
	}

	// Message that forces an error by having more than the max allowed
	// headers.
	maxHeaders := NewMsgHeaders()
	for i := 0; i < MaxBlockHeadersPerMsg; i++ {
		maxHeaders.AddBlockHeader(bh)
	}
	maxHeaders.Headers = append(maxHeaders.Headers, bh)
	maxHeadersEncoded := []byte{
		0xfd, 0xd1, 0x07, // Varint for number of addresses (2001)7D1
	}

	// Intentionally invalid block header that has a transaction count used
	// to force errors.
	//bhTrans := NewBlockHeader(1, &hash, &merkleHash, bits, nonce)
	bhTrans := NewBlockHeader(1, &hash)
	bhTrans.Version = blockOne.Header.Version
	bhTrans.Timestamp = blockOne.Header.Timestamp
	bhTrans.StateRoot = blockOne.Header.StateRoot
	bhTrans.MerkleRoot = blockOne.Header.MerkleRoot
	bhTrans.Height = 1
	bhTrans.GasLimit = 10000000
	bhTrans.GasUsed = 9000000
	bhTrans.CoinBase[0] = 0x66

	transHeader := NewMsgHeaders()
	transHeader.AddBlockHeader(bhTrans)
	transHeaderEncoded := []byte{
		0x01,                   // VarInt for number of headers.
		0x01, 0x00, 0x00, 0x00, // Version 1
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
		0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, // PrevBlock
		0x3b, 0xa3, 0xed, 0xfd, 0x7a, 0x7b, 0x12, 0xb2,
		0x7a, 0xc7, 0x2c, 0x3e, 0x67, 0x76, 0x8f, 0x61,
		0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
		0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a, // MerkleRoot
		0x29, 0xab, 0x5f, 0x49, 0x00, 0x00, 0x00, 0x00, // Timestamp
		0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
		0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
		0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a, // StateRoot
		0x80, 0x96, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, // GasLimit
		0x40, 0x54, 0x89, 0x00, 0x00, 0x00, 0x00, 0x00, // GasUsed
		0x00, 0x00, 0x00, 0x00, // Round
		0x00, 0x00, // SlotIndex
		0x00, 0x00, // Weight
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // PoaHash
		0x01, 0x00, 0x00, 0x00, // height
		0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, // coinbase
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // SigData:
		0x01, // tx count
	}

	tests := []struct {
		in       *MsgHeaders     // Value to encode
		buf      []byte          // Wire encoding
		pver     uint32          // Protocol version for protos encoding
		enc      MessageEncoding // Message encoding format
		max      int             // Max size of fixed buffer to induce errors
		writeErr error           // Expected write error
		readErr  error           // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		// Force error in header count.
		{oneHeader, oneHeaderEncoded, pver, BaseEncoding, 0, io.ErrShortWrite, io.EOF},
		// Force error in block header.
		{oneHeader, oneHeaderEncoded, pver, BaseEncoding, 5, io.ErrShortWrite, io.EOF},
		// Force error with greater than max headers.
		{maxHeaders, maxHeadersEncoded, pver, BaseEncoding, 3, wireErr, wireErr},
		// Force error with number of transactions.
		{transHeader, transHeaderEncoded, pver, BaseEncoding, 81, io.ErrShortWrite, io.ErrUnexpectedEOF},
		// Force error with included transactions.
		{transHeader, transHeaderEncoded, pver, BaseEncoding, len(transHeaderEncoded), nil, wireErr},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		t.Logf("=======================Running the %d test", i)
		// Encode to protos format.
		w := newFixedWriter(test.max)
		err := test.in.VVSEncode(w, test.pver, test.enc)
		if reflect.TypeOf(err) != reflect.TypeOf(test.writeErr) {
			t.Errorf("VVSEncode #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// For errors which are not of type MessageError, check them for
		// equality.
		if _, ok := err.(*MessageError); !ok {
			if err != test.writeErr {
				t.Errorf("VVSEncode #%d wrong error got: %v, "+
					"want: %v", i, err, test.writeErr)
				continue
			}
		}

		// Decode from protos format.
		var msg MsgHeaders
		r := newFixedReader(test.max, test.buf)
		err = msg.VVSDecode(r, test.pver, test.enc)
		if reflect.TypeOf(err) != reflect.TypeOf(test.readErr) {
			t.Errorf("VVSDecode #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}

		// For errors which are not of type MessageError, check them for
		// equality.
		if _, ok := err.(*MessageError); !ok {
			if err != test.readErr {
				t.Errorf("VVSDecode #%d wrong error got: %v, "+
					"want: %v", i, err, test.readErr)
				continue
			}
		}

	}
}
