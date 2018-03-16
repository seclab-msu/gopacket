// Copyright 2017 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"github.com/seclab-msu/gopacket"
	"reflect"
	"testing"
)

// testPacketOSPF2Hello is the packet:
// 13:19:20.008765 IP 192.168.170.8 > 224.0.0.5: OSPFv2, Hello, length 44
//	0x0000:  0100 5e00 0005 00e0 18b1 0cad 0800 45c0  ..^...........E.
//	0x0010:  0040 0812 0000 0159 65dd c0a8 aa08 e000  .@.....Ye.......
//	0x0020:  0005 0201 002c c0a8 aa08 0000 0001 273b  .....,........';
//	0x0030:  0000 0000 0000 0000 0000 ffff ff00 000a  ................
//	0x0040:  0201 0000 0028 c0a8 aa08 0000 0000       .....(........
var testPacketOSPF2Hello = []byte{
	0x01, 0x00, 0x5e, 0x00, 0x00, 0x05, 0x00, 0xe0, 0x18, 0xb1, 0x0c, 0xad, 0x08, 0x00, 0x45, 0xc0,
	0x00, 0x40, 0x08, 0x12, 0x00, 0x00, 0x01, 0x59, 0x65, 0xdd, 0xc0, 0xa8, 0xaa, 0x08, 0xe0, 0x00,
	0x00, 0x05, 0x02, 0x01, 0x00, 0x2c, 0xc0, 0xa8, 0xaa, 0x08, 0x00, 0x00, 0x00, 0x01, 0x27, 0x3b,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0x0a,
	0x02, 0x01, 0x00, 0x00, 0x00, 0x28, 0xc0, 0xa8, 0xaa, 0x08, 0x00, 0x00, 0x00, 0x00,
}

func TestPacketOSPF2Hello(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF2Hello, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeOSPF}, t)

	ospf := p.Layer(LayerTypeOSPF).(*OSPFv2)
	if ospf.Version != 2 {
		t.Fatal("Invalid OSPF version")
	}
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv2); ok {
		want := &OSPFv2{
			OSPF: OSPF{
				Version:      2,
				Type:         OSPFHello,
				PacketLength: 44,
				RouterID:     0xc0a8aa08,
				AreaID:       1,
				Checksum:     0x273b,
				Content: HelloPkgV2{
					NetworkMask: 0xffffff00,
					HelloPkg: HelloPkg{
						RtrPriority:              0x1,
						Options:                  0x2,
						HelloInterval:            0xa,
						RouterDeadInterval:       0x28,
						DesignatedRouterID:       0xc0a8aa08,
						BackupDesignatedRouterID: 0x0,
					},
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF2Hello, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketOSPF3Hello is the packet:
//   14:43:11.663317 IP6 fe80::1 > ff02::5: OSPFv3, Hello, length 36
//   	0x0000:  3333 0000 0005 c200 1ffa 0001 86dd 6e00  33............n.
//   	0x0010:  0000 0024 5901 fe80 0000 0000 0000 0000  ...$Y...........
//   	0x0020:  0000 0000 0001 ff02 0000 0000 0000 0000  ................
//   	0x0030:  0000 0000 0005 0301 0024 0101 0101 0000  .........$......
//   	0x0040:  0001 fb86 0000 0000 0005 0100 0013 000a  ................
//   	0x0050:  0028 0000 0000 0000 0000                 .(........
var testPacketOSPF3Hello = []byte{
	0x33, 0x33, 0x00, 0x00, 0x00, 0x05, 0xc2, 0x00, 0x1f, 0xfa, 0x00, 0x01, 0x86, 0xdd, 0x6e, 0x00,
	0x00, 0x00, 0x00, 0x24, 0x59, 0x01, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x03, 0x01, 0x00, 0x24, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00,
	0x00, 0x01, 0xfb, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, 0x00, 0x00, 0x13, 0x00, 0x0a,
	0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func TestPacketOSPF3Hello(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF3Hello, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv6, LayerTypeOSPF}, t)

	ospf := p.Layer(LayerTypeOSPF).(*OSPFv3)
	if ospf.Version != 3 {
		t.Fatal("Invalid OSPF version")
	}
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv3); ok {
		want := &OSPFv3{
			OSPF: OSPF{
				Version:      3,
				Type:         OSPFHello,
				PacketLength: 36,
				RouterID:     0x1010101,
				AreaID:       1,
				Checksum:     0xfb86,
				Content: HelloPkg{
					InterfaceID:              5,
					RtrPriority:              1,
					Options:                  0x000013,
					HelloInterval:            10,
					RouterDeadInterval:       40,
					DesignatedRouterID:       0,
					BackupDesignatedRouterID: 0,
				},
			},
			Instance: 0,
			Reserved: 0,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF3Hello, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketOSPF2DBDesc is the packet:
// 13:20:14.414477 IP 192.168.170.8 > 192.168.170.2: OSPFv2, Database Description, length 32
//	0x0000:  0060 0881 7a70 00e0 18b1 0cad 0800 45c0  .`..zp........E.
//	0x0010:  0034 2be5 0000 0159 b770 c0a8 aa08 c0a8  .4+....Y.p......
//	0x0020:  aa02 0202 0020 c0a8 aa08 0000 0001 a052  ...............R
//	0x0030:  0000 0000 0000 0000 0000 05dc 0207 4177  ..............Aw
//	0x0040:  a97e                                     .~
var testPacketOSPF2DBDesc = []byte{
	0x00, 0x60, 0x08, 0x81, 0x7a, 0x70, 0x00, 0xe0, 0x18, 0xb1, 0x0c, 0xad, 0x08, 0x00, 0x45, 0xc0,
	0x00, 0x34, 0x2b, 0xe5, 0x00, 0x00, 0x01, 0x59, 0xb7, 0x70, 0xc0, 0xa8, 0xaa, 0x08, 0xc0, 0xa8,
	0xaa, 0x02, 0x02, 0x02, 0x00, 0x20, 0xc0, 0xa8, 0xaa, 0x08, 0x00, 0x00, 0x00, 0x01, 0xa0, 0x52,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0xdc, 0x02, 0x07, 0x41, 0x77,
	0xa9, 0x7e,
}

func TestPacketOSPF2DBDesc(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF2DBDesc, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeOSPF}, t)
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv2); ok {
		want := &OSPFv2{
			OSPF: OSPF{
				Version:      2,
				Type:         OSPFDatabaseDescription,
				PacketLength: 32,
				RouterID:     0xc0a8aa08,
				AreaID:       1,
				Checksum:     0xa052,
				Content: DbDescPkg{
					Options:      0x02,
					InterfaceMTU: 1500,
					Flags:        0x7,
					DDSeqNumber:  1098361214,
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF2DBDesc, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketOSPF3DBDesc is the packet:
//   14:43:51.657571 IP6 fe80::2 > fe80::1: OSPFv3, Database Description, length 28
//   	0x0000:  c200 1ffa 0001 c201 1ffa 0001 86dd 6e00  ..............n.
//   	0x0010:  0000 001c 5901 fe80 0000 0000 0000 0000  ....Y...........
//   	0x0020:  0000 0000 0002 fe80 0000 0000 0000 0000  ................
//   	0x0030:  0000 0000 0001 0302 001c 0202 0202 0000  ................
//   	0x0040:  0001 d826 0000 0000 0013 05dc 0007 0000  ...&............
//   	0x0050:  1d46                                     .F
var testPacketOSPF3DBDesc = []byte{
	0xc2, 0x00, 0x1f, 0xfa, 0x00, 0x01, 0xc2, 0x01, 0x1f, 0xfa, 0x00, 0x01, 0x86, 0xdd, 0x6e, 0x00,
	0x00, 0x00, 0x00, 0x1c, 0x59, 0x01, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x02, 0x00, 0x1c, 0x02, 0x02, 0x02, 0x02, 0x00, 0x00,
	0x00, 0x01, 0xd8, 0x26, 0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x05, 0xdc, 0x00, 0x07, 0x00, 0x00,
	0x1d, 0x46,
}

func TestPacketOSPF3DBDesc(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF3DBDesc, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv6, LayerTypeOSPF}, t)
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv3); ok {
		want := &OSPFv3{
			OSPF: OSPF{
				Version:      3,
				Type:         OSPFDatabaseDescription,
				PacketLength: 28,
				RouterID:     0x2020202,
				AreaID:       1,
				Checksum:     0xd826,
				Content: DbDescPkg{
					Options:      0x000013,
					InterfaceMTU: 1500,
					Flags:        0x7,
					DDSeqNumber:  7494,
				},
			},
			Instance: 0,
			Reserved: 0,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF3DBDesc, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketOSPF2LSRequest is the packet:
// 13:20:14.418003 IP 192.168.170.2 > 192.168.170.8: OSPFv2, LS-Request, length 36
//	0x0000:  00e0 18b1 0cad 0060 0881 7a70 0800 45c0  .......`..zp..E.
//	0x0010:  0038 88c6 0000 0159 5a8b c0a8 aa02 c0a8  .8.....YZ.......
//	0x0020:  aa08 0203 0024 c0a8 aa03 0000 0001 bdc7  .....$..........
//	0x0030:  0000 0000 0000 0000 0000 0000 0001 c0a8  ................
//	0x0040:  aa08 c0a8 aa08                           ......
var testPacketOSPF2LSRequest = []byte{
	0x00, 0xe0, 0x18, 0xb1, 0x0c, 0xad, 0x00, 0x60, 0x08, 0x81, 0x7a, 0x70, 0x08, 0x00, 0x45, 0xc0,
	0x00, 0x38, 0x88, 0xc6, 0x00, 0x00, 0x01, 0x59, 0x5a, 0x8b, 0xc0, 0xa8, 0xaa, 0x02, 0xc0, 0xa8,
	0xaa, 0x08, 0x02, 0x03, 0x00, 0x24, 0xc0, 0xa8, 0xaa, 0x03, 0x00, 0x00, 0x00, 0x01, 0xbd, 0xc7,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xc0, 0xa8,
	0xaa, 0x08, 0xc0, 0xa8, 0xaa, 0x08,
}

func TestPacketOSPF2LSRequest(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF2LSRequest, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeOSPF}, t)
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv2); ok {
		want := &OSPFv2{
			OSPF: OSPF{
				Version:      2,
				Type:         OSPFLinkStateRequest,
				PacketLength: 36,
				RouterID:     0xc0a8aa03,
				AreaID:       1,
				Checksum:     0xbdc7,
				Content: []LSReq{
					LSReq{
						LSType:    0x1,
						LSID:      0xc0a8aa08,
						AdvRouter: 0xc0a8aa08,
					},
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF2LSRequest, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketOSPF3LSRequest is the packet:
//   14:43:51.673584 IP6 fe80::2 > fe80::1: OSPFv3, LS-Request, length 100
//   	0x0000:  c200 1ffa 0001 c201 1ffa 0001 86dd 6e00  ..............n.
//   	0x0010:  0000 0064 5901 fe80 0000 0000 0000 0000  ...dY...........
//   	0x0020:  0000 0000 0002 fe80 0000 0000 0000 0000  ................
//   	0x0030:  0000 0000 0001 0303 0064 0202 0202 0000  .........d......
//   	0x0040:  0001 2c9a 0000 0000 2001 0000 0000 0101  ..,.............
//   	0x0050:  0101 0000 2003 0000 0003 0101 0101 0000  ................
//   	0x0060:  2003 0000 0002 0101 0101 0000 2003 0000  ................
//   	0x0070:  0001 0101 0101 0000 2003 0000 0000 0101  ................
//   	0x0080:  0101 0000 0008 0000 0005 0101 0101 0000  ................
//   	0x0090:  2009 0000 0000 0101 0101                 ..........
var testPacketOSPF3LSRequest = []byte{
	0xc2, 0x00, 0x1f, 0xfa, 0x00, 0x01, 0xc2, 0x01, 0x1f, 0xfa, 0x00, 0x01, 0x86, 0xdd, 0x6e, 0x00,
	0x00, 0x00, 0x00, 0x64, 0x59, 0x01, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x03, 0x00, 0x64, 0x02, 0x02, 0x02, 0x02, 0x00, 0x00,
	0x00, 0x01, 0x2c, 0x9a, 0x00, 0x00, 0x00, 0x00, 0x20, 0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01,
	0x01, 0x01, 0x00, 0x00, 0x20, 0x03, 0x00, 0x00, 0x00, 0x03, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00,
	0x20, 0x03, 0x00, 0x00, 0x00, 0x02, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x20, 0x03, 0x00, 0x00,
	0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x20, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01,
	0x01, 0x01, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x05, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00,
	0x20, 0x09, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01,
}

func TestPacketOSPF3LSRequest(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF3LSRequest, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv6, LayerTypeOSPF}, t)
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv3); ok {
		want := &OSPFv3{
			OSPF: OSPF{
				Version:      3,
				Type:         OSPFLinkStateRequest,
				PacketLength: 100,
				RouterID:     0x2020202,
				AreaID:       1,
				Checksum:     0x2c9a,
				Content: []LSReq{
					LSReq{
						LSType:    0x2001,
						LSID:      0x00000000,
						AdvRouter: 0x01010101,
					},
					LSReq{
						LSType:    0x2003,
						LSID:      0x00000003,
						AdvRouter: 0x01010101,
					},
					LSReq{
						LSType:    0x2003,
						LSID:      0x00000002,
						AdvRouter: 0x01010101,
					},
					LSReq{
						LSType:    0x2003,
						LSID:      0x00000001,
						AdvRouter: 0x01010101,
					},
					LSReq{
						LSType:    0x2003,
						LSID:      0x00000000,
						AdvRouter: 0x01010101,
					},
					LSReq{
						LSType:    0x0008,
						LSID:      0x00000005,
						AdvRouter: 0x01010101,
					},
					LSReq{
						LSType:    0x2009,
						LSID:      0x00000000,
						AdvRouter: 0x01010101,
					},
				},
			},
			Instance: 0,
			Reserved: 0,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF3LSRequest, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketOSPF2LSUpdate is the packet:
// 13:20:14.420459 IP 192.168.170.2 > 224.0.0.6: OSPFv2, LS-Update, length 292
//	0x0000:  0100 5e00 0006 0060 0881 7a70 0800 45c0  ..^....`..zp..E.
//	0x0010:  0138 3025 0000 0159 3cd7 c0a8 aa02 e000  .80%...Y<.......
//	0x0020:  0006 0204 0124 c0a8 aa03 0000 0001 366b  .....$........6k
//	0x0030:  0000 0000 0000 0000 0000 0000 0007 0002  ................
//	0x0040:  0201 c0a8 aa03 c0a8 aa03 8000 0001 3a9c  ..............:.
//	0x0050:  0030 0200 0002 c0a8 aa00 ffff ff00 0300  .0..............
//	0x0060:  000a c0a8 aa00 ffff ff00 0300 000a 0003  ................
//	0x0070:  0205 50d4 1000 c0a8 aa02 8000 0001 2a49  ..P...........*I
//	0x0080:  0024 ffff ffff 8000 0014 0000 0000 0000  .$..............
//	0x0090:  0000 0003 0205 9479 ab00 c0a8 aa02 8000  .......y........
//	0x00a0:  0001 34a5 0024 ffff ff00 8000 0014 c0a8  ..4..$..........
//	0x00b0:  aa01 0000 0000 0003 0205 c082 7800 c0a8  ............x...
//	0x00c0:  aa02 8000 0001 d319 0024 ffff ff00 8000  .........$......
//	0x00d0:  0014 0000 0000 0000 0000 0003 0205 c0a8  ................
//	0x00e0:  0000 c0a8 aa02 8000 0001 3708 0024 ffff  ..........7..$..
//	0x00f0:  ff00 8000 0014 0000 0000 0000 0000 0003  ................
//	0x0100:  0205 c0a8 0100 c0a8 aa02 8000 0001 2c12  ..............,.
//	0x0110:  0024 ffff ff00 8000 0014 0000 0000 0000  .$..............
//	0x0120:  0000 0003 0205 c0a8 ac00 c0a8 aa02 8000  ................
//	0x0130:  0001 3341 0024 ffff ff00 8000 0014 c0a8  ..3A.$..........
//	0x0140:  aa0a 0000 0000                           ......
var testPacketOSPF2LSUpdate = []byte{
	0x01, 0x00, 0x5e, 0x00, 0x00, 0x06, 0x00, 0x60, 0x08, 0x81, 0x7a, 0x70, 0x08, 0x00, 0x45, 0xc0,
	0x01, 0x38, 0x30, 0x25, 0x00, 0x00, 0x01, 0x59, 0x3c, 0xd7, 0xc0, 0xa8, 0xaa, 0x02, 0xe0, 0x00,
	0x00, 0x06, 0x02, 0x04, 0x01, 0x24, 0xc0, 0xa8, 0xaa, 0x03, 0x00, 0x00, 0x00, 0x01, 0x36, 0x6b,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x00, 0x02,
	0x02, 0x01, 0xc0, 0xa8, 0xaa, 0x03, 0xc0, 0xa8, 0xaa, 0x03, 0x80, 0x00, 0x00, 0x01, 0x3a, 0x9c,
	0x00, 0x30, 0x02, 0x00, 0x00, 0x02, 0xc0, 0xa8, 0xaa, 0x00, 0xff, 0xff, 0xff, 0x00, 0x03, 0x00,
	0x00, 0x0a, 0xc0, 0xa8, 0xaa, 0x00, 0xff, 0xff, 0xff, 0x00, 0x03, 0x00, 0x00, 0x0a, 0x00, 0x03,
	0x02, 0x05, 0x50, 0xd4, 0x10, 0x00, 0xc0, 0xa8, 0xaa, 0x02, 0x80, 0x00, 0x00, 0x01, 0x2a, 0x49,
	0x00, 0x24, 0xff, 0xff, 0xff, 0xff, 0x80, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x03, 0x02, 0x05, 0x94, 0x79, 0xab, 0x00, 0xc0, 0xa8, 0xaa, 0x02, 0x80, 0x00,
	0x00, 0x01, 0x34, 0xa5, 0x00, 0x24, 0xff, 0xff, 0xff, 0x00, 0x80, 0x00, 0x00, 0x14, 0xc0, 0xa8,
	0xaa, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x02, 0x05, 0xc0, 0x82, 0x78, 0x00, 0xc0, 0xa8,
	0xaa, 0x02, 0x80, 0x00, 0x00, 0x01, 0xd3, 0x19, 0x00, 0x24, 0xff, 0xff, 0xff, 0x00, 0x80, 0x00,
	0x00, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x02, 0x05, 0xc0, 0xa8,
	0x00, 0x00, 0xc0, 0xa8, 0xaa, 0x02, 0x80, 0x00, 0x00, 0x01, 0x37, 0x08, 0x00, 0x24, 0xff, 0xff,
	0xff, 0x00, 0x80, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
	0x02, 0x05, 0xc0, 0xa8, 0x01, 0x00, 0xc0, 0xa8, 0xaa, 0x02, 0x80, 0x00, 0x00, 0x01, 0x2c, 0x12,
	0x00, 0x24, 0xff, 0xff, 0xff, 0x00, 0x80, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x03, 0x02, 0x05, 0xc0, 0xa8, 0xac, 0x00, 0xc0, 0xa8, 0xaa, 0x02, 0x80, 0x00,
	0x00, 0x01, 0x33, 0x41, 0x00, 0x24, 0xff, 0xff, 0xff, 0x00, 0x80, 0x00, 0x00, 0x14, 0xc0, 0xa8,
	0xaa, 0x0a, 0x00, 0x00, 0x00, 0x00,
}

func TestPacketOSPF2LSUpdate(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF2LSUpdate, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeOSPF}, t)
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv2); ok {
		want := &OSPFv2{
			OSPF: OSPF{
				Version:      2,
				Type:         OSPFLinkStateUpdate,
				PacketLength: 292,
				RouterID:     0xc0a8aa03,
				AreaID:       1,
				Checksum:     0x366b,
				Content: LSUpdate{
					NumOfLSAs: 7,
					LSAs: []LSA{
						LSA{
							LSAheader: LSAheader{
								LSAge:       0x2,
								LSType:      0x1,
								LinkStateID: 0xc0a8aa03,
								AdvRouter:   0xc0a8aa03,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0x3a9c,
								Length:      0x30,
								LSOptions:   0x2,
							},
							Content: RouterLSAV2{
								Flags: 0x2,
								Links: 0x2,
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       0x3,
								LSType:      0x5,
								LinkStateID: 0x50d41000,
								AdvRouter:   0xc0a8aa02,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0x2a49,
								Length:      0x24,
								LSOptions:   0x2,
							},
							Content: ASExternalLSAV2{
								NetworkMask:       0xffffffff,
								ExternalBit:       0x80,
								Metric:            0x14,
								ForwardingAddress: 0x0,
								ExternalRouteTag:  0x0,
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       0x3,
								LSType:      0x5,
								LinkStateID: 0x9479ab00,
								AdvRouter:   0xc0a8aa02,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0x34a5,
								Length:      0x24,
								LSOptions:   0x2,
							},
							Content: ASExternalLSAV2{
								NetworkMask:       0xffffff00,
								ExternalBit:       0x80,
								Metric:            0x14,
								ForwardingAddress: 0xc0a8aa01,
								ExternalRouteTag:  0x0,
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       0x3,
								LSType:      0x5,
								LinkStateID: 0xc0827800,
								AdvRouter:   0xc0a8aa02,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0xd319,
								Length:      0x24,
								LSOptions:   0x2,
							},
							Content: ASExternalLSAV2{
								NetworkMask:       0xffffff00,
								ExternalBit:       0x80,
								Metric:            0x14,
								ForwardingAddress: 0x0,
								ExternalRouteTag:  0x0,
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       0x3,
								LSType:      0x5,
								LinkStateID: 0xc0a80000,
								AdvRouter:   0xc0a8aa02,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0x3708,
								Length:      0x24,
								LSOptions:   0x2,
							},
							Content: ASExternalLSAV2{
								NetworkMask:       0xffffff00,
								ExternalBit:       0x80,
								Metric:            0x14,
								ForwardingAddress: 0x0,
								ExternalRouteTag:  0x0,
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       0x3,
								LSType:      0x5,
								LinkStateID: 0xc0a80100,
								AdvRouter:   0xc0a8aa02,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0x2c12,
								Length:      0x24,
								LSOptions:   0x2,
							},
							Content: ASExternalLSAV2{
								NetworkMask:       0xffffff00,
								ExternalBit:       0x80,
								Metric:            0x14,
								ForwardingAddress: 0x0,
								ExternalRouteTag:  0x0,
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       0x3,
								LSType:      0x5,
								LinkStateID: 0xc0a8ac00,
								AdvRouter:   0xc0a8aa02,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0x3341,
								Length:      0x24,
								LSOptions:   0x2,
							},
							Content: ASExternalLSAV2{
								NetworkMask:       0xffffff00,
								ExternalBit:       0x80,
								Metric:            0x14,
								ForwardingAddress: 0xc0a8aa0a,
								ExternalRouteTag:  0x0,
							},
						},
					},
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF2LSUpdate, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketOSPF3LSUpdate is the packet:
//   14:43:51.681554 IP6 fe80::1 > fe80::2: OSPFv3, LS-Update, length 288
//   	0x0000:  c201 1ffa 0001 c200 1ffa 0001 86dd 6e00  ..............n.
//   	0x0010:  0000 0120 5901 fe80 0000 0000 0000 0000  ....Y...........
//   	0x0020:  0000 0000 0001 fe80 0000 0000 0000 0000  ................
//   	0x0030:  0000 0000 0002 0304 0120 0101 0101 0000  ................
//   	0x0040:  0001 e556 0000 0000 0007 0028 2001 0000  ...V.......(....
//   	0x0050:  0000 0101 0101 8000 0002 d13a 0018 0100  ...........:....
//   	0x0060:  0033 0029 2003 0000 0003 0101 0101 8000  .3.)............
//   	0x0070:  0001 6259 0024 0000 004a 4000 0000 2001  ..bY.$...J@.....
//   	0x0080:  0db8 0000 0003 0029 2003 0000 0002 0101  .......)........
//   	0x0090:  0101 8000 0001 baf6 0024 0000 0054 4000  .........$...T@.
//   	0x00a0:  0000 2001 0db8 0000 0004 0029 2003 0000  ...........)....
//   	0x00b0:  0001 0101 0101 8000 0001 eba0 0024 0000  .............$..
//   	0x00c0:  004a 4000 0000 2001 0db8 0000 0034 0029  .J@..........4.)
//   	0x00d0:  2003 0000 0000 0101 0101 8000 0001 0ebd  ................
//   	0x00e0:  0024 0000 0040 4000 0000 2001 0db8 0000  .$...@@.........
//   	0x00f0:  0000 0023 0008 0000 0005 0101 0101 8000  ...#............
//   	0x0100:  0002 3d08 0038 0100 0033 fe80 0000 0000  ..=..8...3......
//   	0x0110:  0000 0000 0000 0000 0001 0000 0001 4000  ..............@.
//   	0x0120:  0000 2001 0db8 0000 0012 0023 2009 0000  ...........#....
//   	0x0130:  0000 0101 0101 8000 0001 e8d2 002c 0001  .............,..
//   	0x0140:  2001 0000 0000 0101 0101 4000 000a 2001  ..........@.....
//   	0x0150:  0db8 0000 0012                           ......
var testPacketOSPF3LSUpdate = []byte{
	0xc2, 0x01, 0x1f, 0xfa, 0x00, 0x01, 0xc2, 0x00, 0x1f, 0xfa, 0x00, 0x01, 0x86, 0xdd, 0x6e, 0x00,
	0x00, 0x00, 0x01, 0x20, 0x59, 0x01, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x03, 0x04, 0x01, 0x20, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00,
	0x00, 0x01, 0xe5, 0x56, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x00, 0x28, 0x20, 0x01, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x80, 0x00, 0x00, 0x02, 0xd1, 0x3a, 0x00, 0x18, 0x01, 0x00,
	0x00, 0x33, 0x00, 0x29, 0x20, 0x03, 0x00, 0x00, 0x00, 0x03, 0x01, 0x01, 0x01, 0x01, 0x80, 0x00,
	0x00, 0x01, 0x62, 0x59, 0x00, 0x24, 0x00, 0x00, 0x00, 0x4a, 0x40, 0x00, 0x00, 0x00, 0x20, 0x01,
	0x0d, 0xb8, 0x00, 0x00, 0x00, 0x03, 0x00, 0x29, 0x20, 0x03, 0x00, 0x00, 0x00, 0x02, 0x01, 0x01,
	0x01, 0x01, 0x80, 0x00, 0x00, 0x01, 0xba, 0xf6, 0x00, 0x24, 0x00, 0x00, 0x00, 0x54, 0x40, 0x00,
	0x00, 0x00, 0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x04, 0x00, 0x29, 0x20, 0x03, 0x00, 0x00,
	0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x80, 0x00, 0x00, 0x01, 0xeb, 0xa0, 0x00, 0x24, 0x00, 0x00,
	0x00, 0x4a, 0x40, 0x00, 0x00, 0x00, 0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x34, 0x00, 0x29,
	0x20, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x80, 0x00, 0x00, 0x01, 0x0e, 0xbd,
	0x00, 0x24, 0x00, 0x00, 0x00, 0x40, 0x40, 0x00, 0x00, 0x00, 0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x23, 0x00, 0x08, 0x00, 0x00, 0x00, 0x05, 0x01, 0x01, 0x01, 0x01, 0x80, 0x00,
	0x00, 0x02, 0x3d, 0x08, 0x00, 0x38, 0x01, 0x00, 0x00, 0x33, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x40, 0x00,
	0x00, 0x00, 0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x12, 0x00, 0x23, 0x20, 0x09, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x80, 0x00, 0x00, 0x01, 0xe8, 0xd2, 0x00, 0x2c, 0x00, 0x01,
	0x20, 0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x40, 0x00, 0x00, 0x0a, 0x20, 0x01,
	0x0d, 0xb8, 0x00, 0x00, 0x00, 0x12,
}

func TestPacketOSPF3LSUpdate(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF3LSUpdate, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv6, LayerTypeOSPF}, t)
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv3); ok {
		want := &OSPFv3{
			OSPF: OSPF{
				Version:      3,
				Type:         OSPFLinkStateUpdate,
				PacketLength: 288,
				RouterID:     0x1010101,
				AreaID:       1,
				Checksum:     0xe556,
				Content: LSUpdate{
					NumOfLSAs: 7,
					LSAs: []LSA{
						LSA{
							LSAheader: LSAheader{
								LSAge:       40,
								LSType:      0x2001,
								LinkStateID: 0x00000000,
								AdvRouter:   0x01010101,
								LSSeqNumber: 0x80000002,
								LSChecksum:  0xd13a,
								Length:      24,
							},
							Content: RouterLSA{
								Flags:   0x1,
								Options: 0x33,
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       41,
								LSType:      0x2003,
								LinkStateID: 0x00000003,
								AdvRouter:   0x01010101,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0x6259,
								Length:      36,
							},
							Content: InterAreaPrefixLSA{
								Metric:        74,
								PrefixLength:  64,
								PrefixOptions: 0,
								AddressPrefix: []byte{0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x03},
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       41,
								LSType:      0x2003,
								LinkStateID: 0x00000002,
								AdvRouter:   0x01010101,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0xbaf6,
								Length:      36,
							},
							Content: InterAreaPrefixLSA{
								Metric:        84,
								PrefixLength:  64,
								PrefixOptions: 0,
								AddressPrefix: []byte{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x4},
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       41,
								LSType:      0x2003,
								LinkStateID: 0x00000001,
								AdvRouter:   0x01010101,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0xeba0,
								Length:      36,
							},
							Content: InterAreaPrefixLSA{
								Metric:        74,
								PrefixLength:  64,
								PrefixOptions: 0,
								AddressPrefix: []byte{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x34},
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       41,
								LSType:      0x2003,
								LinkStateID: 0x00000000,
								AdvRouter:   0x01010101,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0xebd,
								Length:      36,
							},
							Content: InterAreaPrefixLSA{
								Metric:        64,
								PrefixLength:  64,
								PrefixOptions: 0,
								AddressPrefix: []byte{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0},
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       35,
								LSType:      0x8,
								LinkStateID: 0x00000005,
								AdvRouter:   0x01010101,
								LSSeqNumber: 0x80000002,
								LSChecksum:  0x3d08,
								Length:      56,
							},
							Content: LinkLSA{
								RtrPriority:      1,
								Options:          0x33,
								LinkLocalAddress: []byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
								NumOfPrefixes:    1,
								Prefixes: []Prefix{
									Prefix{
										PrefixLength:  64,
										PrefixOptions: 0,
										AddressPrefix: []byte{0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x12},
									},
								},
							},
						},
						LSA{
							LSAheader: LSAheader{
								LSAge:       35,
								LSType:      0x2009,
								LinkStateID: 0x00000000,
								AdvRouter:   0x01010101,
								LSSeqNumber: 0x80000001,
								LSChecksum:  0xe8d2,
								Length:      44,
							},
							Content: IntraAreaPrefixLSA{
								NumOfPrefixes: 1,
								RefLSType:     0x2001,
								RefAdvRouter:  0x01010101,
								Prefixes: []Prefix{
									Prefix{
										PrefixLength:  64,
										PrefixOptions: 0,
										Metric:        10,
										AddressPrefix: []byte{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x12},
									},
								},
							},
						},
					},
				},
			},
			Instance: 0,
			Reserved: 0,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF3LSUpdate, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketOSPF2LSAck is the packet:
// 13:20:20.858322 IP 192.168.170.8 > 192.168.170.2: OSPFv2, LS-Ack, length 44
//	0x0000:  0060 0881 7a70 00e0 18b1 0cad 0800 45c0  .`..zp........E.
//	0x0010:  0040 2bea 0000 0159 b75f c0a8 aa08 c0a8  .@+....Y._......
//	0x0020:  aa02 0205 002c c0a8 aa08 0000 0001 e2f4  .....,..........
//	0x0030:  0000 0000 0000 0000 0000 0e10 0201 c0a8  ................
//	0x0040:  aa02 c0a8 aa02 8000 0001 4a8e 0030       ..........J..0
var testPacketOSPF2LSAck = []byte{
	0x00, 0x60, 0x08, 0x81, 0x7a, 0x70, 0x00, 0xe0, 0x18, 0xb1, 0x0c, 0xad, 0x08, 0x00, 0x45, 0xc0,
	0x00, 0x40, 0x2b, 0xea, 0x00, 0x00, 0x01, 0x59, 0xb7, 0x5f, 0xc0, 0xa8, 0xaa, 0x08, 0xc0, 0xa8,
	0xaa, 0x02, 0x02, 0x05, 0x00, 0x2c, 0xc0, 0xa8, 0xaa, 0x08, 0x00, 0x00, 0x00, 0x01, 0xe2, 0xf4,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0e, 0x10, 0x02, 0x01, 0xc0, 0xa8,
	0xaa, 0x02, 0xc0, 0xa8, 0xaa, 0x02, 0x80, 0x00, 0x00, 0x01, 0x4a, 0x8e, 0x00, 0x30,
}

func TestPacketOSPF2LSAck(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF2LSAck, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeOSPF}, t)
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv2); ok {
		want := &OSPFv2{
			OSPF: OSPF{
				Version:      2,
				Type:         OSPFLinkStateAcknowledgment,
				PacketLength: 44,
				RouterID:     0xc0a8aa08,
				AreaID:       1,
				Checksum:     0xe2f4,
				Content: []LSAheader{
					LSAheader{
						LSAge:       0xe10,
						LSType:      0x1,
						LinkStateID: 0xc0a8aa02,
						AdvRouter:   0xc0a8aa02,
						LSSeqNumber: 0x80000001,
						LSChecksum:  0x4a8e,
						Length:      0x30,
						LSOptions:   0x2,
					},
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket9(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF3LSAck, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketOSPF3LSAck is the packet:
//   14:43:54.185384 IP6 fe80::1 > ff02::5: OSPFv3, LS-Ack, length 136
//   	0x0000:  3333 0000 0005 c200 1ffa 0001 86dd 6e00  33............n.
//   	0x0010:  0000 0088 5901 fe80 0000 0000 0000 0000  ....Y...........
//   	0x0020:  0000 0000 0001 ff02 0000 0000 0000 0000  ................
//   	0x0030:  0000 0000 0005 0305 0088 0101 0101 0000  ................
//   	0x0040:  0001 9d2c 0000 0005 2001 0000 0000 0202  ...,............
//   	0x0050:  0202 8000 0002 b354 0018 0006 2003 0000  .......T........
//   	0x0060:  0003 0202 0202 8000 0001 4473 0024 0006  ..........Ds.$..
//   	0x0070:  2003 0000 0002 0202 0202 8000 0001 9c11  ................
//   	0x0080:  0024 0006 2003 0000 0001 0202 0202 8000  .$..............
//   	0x0090:  0001 cdba 0024 0006 2003 0000 0000 0202  .....$..........
//   	0x00a0:  0202 8000 0001 efd7 0024 0005 0008 0000  .........$......
//   	0x00b0:  0005 0202 0202 8000 0001 5433 002c       ..........T3.,
var testPacketOSPF3LSAck = []byte{
	0x33, 0x33, 0x00, 0x00, 0x00, 0x05, 0xc2, 0x00, 0x1f, 0xfa, 0x00, 0x01, 0x86, 0xdd, 0x6e, 0x00,
	0x00, 0x00, 0x00, 0x88, 0x59, 0x01, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x03, 0x05, 0x00, 0x88, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00,
	0x00, 0x01, 0x9d, 0x2c, 0x00, 0x00, 0x00, 0x05, 0x20, 0x01, 0x00, 0x00, 0x00, 0x00, 0x02, 0x02,
	0x02, 0x02, 0x80, 0x00, 0x00, 0x02, 0xb3, 0x54, 0x00, 0x18, 0x00, 0x06, 0x20, 0x03, 0x00, 0x00,
	0x00, 0x03, 0x02, 0x02, 0x02, 0x02, 0x80, 0x00, 0x00, 0x01, 0x44, 0x73, 0x00, 0x24, 0x00, 0x06,
	0x20, 0x03, 0x00, 0x00, 0x00, 0x02, 0x02, 0x02, 0x02, 0x02, 0x80, 0x00, 0x00, 0x01, 0x9c, 0x11,
	0x00, 0x24, 0x00, 0x06, 0x20, 0x03, 0x00, 0x00, 0x00, 0x01, 0x02, 0x02, 0x02, 0x02, 0x80, 0x00,
	0x00, 0x01, 0xcd, 0xba, 0x00, 0x24, 0x00, 0x06, 0x20, 0x03, 0x00, 0x00, 0x00, 0x00, 0x02, 0x02,
	0x02, 0x02, 0x80, 0x00, 0x00, 0x01, 0xef, 0xd7, 0x00, 0x24, 0x00, 0x05, 0x00, 0x08, 0x00, 0x00,
	0x00, 0x05, 0x02, 0x02, 0x02, 0x02, 0x80, 0x00, 0x00, 0x01, 0x54, 0x33, 0x00, 0x2c,
}

func TestPacketOSPF3LSAck(t *testing.T) {
	p := gopacket.NewPacket(testPacketOSPF3LSAck, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv6, LayerTypeOSPF}, t)
	if got, ok := p.Layer(LayerTypeOSPF).(*OSPFv3); ok {
		want := &OSPFv3{
			OSPF: OSPF{
				Version:      3,
				Type:         OSPFLinkStateAcknowledgment,
				PacketLength: 136,
				RouterID:     0x1010101,
				AreaID:       1,
				Checksum:     0x9d2c,
				Content: []LSAheader{
					LSAheader{
						LSAge:       5,
						LSType:      0x2001,
						LinkStateID: 0x00000000,
						AdvRouter:   0x02020202,
						LSSeqNumber: 0x80000002,
						LSChecksum:  0xb354,
						Length:      24,
					},
					LSAheader{
						LSAge:       6,
						LSType:      0x2003,
						LinkStateID: 0x00000003,
						AdvRouter:   0x02020202,
						LSSeqNumber: 0x80000001,
						LSChecksum:  0x4473,
						Length:      36,
					},
					LSAheader{
						LSAge:       6,
						LSType:      0x2003,
						LinkStateID: 0x00000002,
						AdvRouter:   0x02020202,
						LSSeqNumber: 0x80000001,
						LSChecksum:  0x9c11,
						Length:      36,
					},
					LSAheader{
						LSAge:       6,
						LSType:      0x2003,
						LinkStateID: 0x00000001,
						AdvRouter:   0x02020202,
						LSSeqNumber: 0x80000001,
						LSChecksum:  0xcdba,
						Length:      36,
					},
					LSAheader{
						LSAge:       6,
						LSType:      0x2003,
						LinkStateID: 0x00000000,
						AdvRouter:   0x02020202,
						LSSeqNumber: 0x80000001,
						LSChecksum:  0xefd7,
						Length:      36,
					},
					LSAheader{
						LSAge:       5,
						LSType:      0x0008,
						LinkStateID: 0x00000005,
						AdvRouter:   0x02020202,
						LSSeqNumber: 0x80000001,
						LSChecksum:  0x5433,
						Length:      44,
					},
				},
			},
			Instance: 0,
			Reserved: 0,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("OSPF packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No OSPF layer type found in packet")
	}
}
func BenchmarkDecodePacketPacket4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketOSPF3LSAck, LinkTypeEthernet, gopacket.NoCopy)
	}
}
