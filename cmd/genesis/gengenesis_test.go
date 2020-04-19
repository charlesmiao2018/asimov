// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/common"
	"testing"
)
var (
	oldConstructorArgsMap = map[string][]common.Address{
		"genesisCitizens" : {
			{0x66, 178, 143, 157, 209, 207, 99, 20, 168, 181, 214, 145, 174, 236, 108, 110, 175, 69, 108, 189, 154},
			{0x66, 132, 235, 30, 89, 45, 229, 249, 32, 89, 246, 199, 118, 106, 4, 242, 201, 253, 85, 135, 173},
			{0x66, 168, 193, 222, 249, 224, 117, 137, 65, 94, 191, 89, 181, 95, 118, 185, 213, 17, 80, 100, 201},
			{0x66, 35, 6, 212, 50, 88, 41, 56, 70, 171, 25, 136, 16, 135, 158, 49, 254, 172, 117, 159, 172},
			{0x66, 151, 167, 246, 224, 61, 207, 245, 205, 156, 191, 153, 181, 167, 126, 123, 49, 221, 7, 21, 60},
			{0x66, 8, 50, 184, 4, 199, 239, 103, 188, 194, 36, 222, 138, 72, 201, 231, 163, 165, 91, 60, 31},
			{0x66, 97, 216, 74, 191, 24, 50, 202, 228, 154, 220, 137, 29, 3, 49, 203, 187, 160, 217, 41, 99},
		},
		"_validators" : {
			{0x66, 178, 143, 157, 209, 207, 99, 20, 168, 181, 214, 145, 174, 236, 108, 110, 175, 69, 108, 189, 154},
			{0x66, 132, 235, 30, 89, 45, 229, 249, 32, 89, 246, 199, 118, 106, 4, 242, 201, 253, 85, 135, 173},
			{0x66, 168, 193, 222, 249, 224, 117, 137, 65, 94, 191, 89, 181, 95, 118, 185, 213, 17, 80, 100, 201},
			{0x66, 35, 6, 212, 50, 88, 41, 56, 70, 171, 25, 136, 16, 135, 158, 49, 254, 172, 117, 159, 172},
			{0x66, 151, 167, 246, 224, 61, 207, 245, 205, 156, 191, 153, 181, 167, 126, 123, 49, 221, 7, 21, 60},
		},
		"_admins" : {
			{0x66, 178, 143, 157, 209, 207, 99, 20, 168, 181, 214, 145, 174, 236, 108, 110, 175, 69, 108, 189, 154},
		},
		"_miners" : {
			{0x66, 132, 235, 30, 89, 45, 229, 249, 32, 89, 246, 199, 118, 106, 4, 242, 201, 253, 85, 135, 173},
			{0x66, 168, 193, 222, 249, 224, 117, 137, 65, 94, 191, 89, 181, 95, 118, 185, 213, 17, 80, 100, 201},
		},
	}
)

func TestNewFormatOfAddress(t *testing.T) {
	networkTypes := common.DevelopNet.String()
	args := []string{"genesisCitizens", "_validators", "_admins", "_miners"}
	for i := 0; i < len(args); i++ {
		value := chaincfg.NetConstructorArgsMap[args[i]]
		for j := 0; j < len(value); j++ {
			if !bytes.Equal(value[networkTypes][j].Bytes(), oldConstructorArgsMap[args[i]][j].Bytes()) {
				panic("not equal value")
			}
		}
	}

}