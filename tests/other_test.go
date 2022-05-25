package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate/app"
)

type addrArgs struct {
	intf              string
	addrs             []string
	services          []string
	vInfo             []app.VipInfo
	ServiceInfoAll    []app.ServiceInfo
	ServiceGrpInfoAll []app.ServiceGrpInfo
	usedAddress       []string
	allInfo           app.AllInfo
}

var (
	testDefaultAllInfo = app.AllInfo{
		Env:            "test",
		AddressInfoAll: testAddrInfoAll,
		AddrGrpInfoAll: testAddrGrpInfoAll,
		SRouteInfoAll:  testSRouteInfoAll,
		IntfInfoAll:    testIntfInfoAll,
	}
	testAddrInfoAll = []app.AddressInfo{
		{
			Name:       "iprange",
			Address:    "",
			SubnetMask: "",
			StartIP:    "192.167.1.1",
			EndIP:      "192.167.1.20",
		},
		{
			Name:       "FABRIC_DEVICE",
			Address:    "",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "none",
			Address:    "0.0.0.0",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "all",
			Address:    "",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "FIREWALL_AUTH_PORTAL_ADDRESS",
			Address:    "",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "SSLVPN_TUNNEL_ADDR1",
			Address:    "",
			SubnetMask: "",
			StartIP:    "10.212.134.200",
			EndIP:      "10.212.134.210",
		},
		{
			Name:       "8.8.8.8",
			Address:    "8.8.8.8",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.22.25/32",
			Address:    "192.168.22.25",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "145.1.0.1",
			Address:    "145.1.0.1",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "202.213.16.3",
			Address:    "202.213.16.3",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.112",
			Address:    "192.168.1.112",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "200.200.200.200",
			Address:    "200.200.200.200",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.100",
			Address:    "192.168.1.100",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "172.30.1.100",
			Address:    "172.30.1.100",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "172.30.2.100",
			Address:    "172.30.2.100",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "172.30.1.110",
			Address:    "172.30.1.110",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "172.30.2.110",
			Address:    "172.30.2.110",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.101",
			Address:    "192.168.1.101",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.102",
			Address:    "192.168.1.102",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.103",
			Address:    "192.168.1.103",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "10.0.0.120",
			Address:    "10.0.0.120",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.104",
			Address:    "192.168.1.104",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.105",
			Address:    "192.168.1.105",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.106",
			Address:    "192.168.1.106",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.4.100",
			Address:    "192.168.4.100",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.5.100",
			Address:    "192.168.5.100",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.110",
			Address:    "192.168.1.110",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.113",
			Address:    "192.168.1.113",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "200.200.200.208",
			Address:    "200.200.200.208",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "111.111.111.116",
			Address:    "111.111.111.116",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.1.222",
			Address:    "192.168.1.222",
			SubnetMask: "255.255.255.255",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "203.0.113.0/24",
			Address:    "203.0.113.0",
			SubnetMask: "255.255.255.0",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.110.0/26",
			Address:    "192.168.110.64",
			SubnetMask: "255.255.255.192",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.111.0/24",
			Address:    "192.168.111.0",
			SubnetMask: "255.255.255.0",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.112.0/28",
			Address:    "192.168.112.0",
			SubnetMask: "255.255.255.240",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "192.168.113.128/28",
			Address:    "192.168.113.128",
			SubnetMask: "255.255.255.240",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "local.test.com",
			Address:    "local.test.com",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "local.testsub.com",
			Address:    "local.testsub.com",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "login.microsoftonline.com",
			Address:    "login.microsoftonline.com",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "login.microsoft.com",
			Address:    "login.microsoft.com",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "login.windows.net",
			Address:    "login.windows.net",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "gmail.com",
			Address:    "gmail.com",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "wildcard.google.com",
			Address:    "google.com",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
		{
			Name:       "wildcard.dropbox.com",
			Address:    "dropbox.com",
			SubnetMask: "",
			StartIP:    "",
			EndIP:      "",
		},
	}
	testAddrGrpInfoAll = []app.AddrGrpInfo{
		{
			Name:   `"Microsoft Office 365"`,
			Member: []string{"login.microsoftonline.com", "login.microsoft.com", "login.windows.net"},
		},
		{
			Name:   `"G Suite"`,
			Member: []string{"gmail.com", "wildcard.google.com"},
		},
		{
			Name:   `"testgrp"`,
			Member: []string{"local.testsub.com", "192.168.110.0/26", "192.168.22.25/32"},
		},
	}
	testIPPoolInfoAll = []app.IPPoolInfo{
		{
			Name:          "123.123.123.123",
			StartIP:       "123.123.123.123",
			EndIP:         "123.123.123.123",
			SourceStartIP: "",
			SourceEndIP:   "",
		},
		{
			Name:          "100.100.100.100",
			StartIP:       "100.100.100.100",
			EndIP:         "100.100.100.100",
			SourceStartIP: "",
			SourceEndIP:   "",
		},
		{
			Name:          "1on1",
			StartIP:       "2.2.2.2",
			EndIP:         "2.2.2.4",
			SourceStartIP: "",
			SourceEndIP:   "",
		},
		{
			Name:          "fix_port_range",
			StartIP:       "3.3.3.3",
			EndIP:         "3.3.3.33",
			SourceStartIP: "172.16.0.172",
			SourceEndIP:   "172.16.0.172",
		},
		{
			Name:          "port_block_assign",
			StartIP:       "6.6.6.6",
			EndIP:         "6.6.6.66",
			SourceStartIP: "",
			SourceEndIP:   "",
		},
	}
	testIntfInfoAll = []app.IntfInfo{
		{
			Name:       "wan1",
			Address:    "10.10.10.1",
			SubnetMask: "255.255.255.0",
			VLANID:     0,
			RandomIP:   "1.1.1.1",
		},
		{
			Name:       "lan",
			Address:    "192.168.10.1",
			SubnetMask: "255.255.255.0",
			VLANID:     10,
			RandomIP:   "1.1.1.2",
		},
	}
	testSRouteInfoAll = []app.SRouteInfo{
		{
			Name:   "1",
			Dst:    "192.168.130.0 255.255.255.0",
			GW:     "192.168.79.254",
			Device: `"wan2"`,
		},
		{
			Name:   "5",
			Dst:    "",
			GW:     "10.0.0.1",
			Device: `"wan1"`,
		},
		{
			Name:   "6",
			Dst:    "172.30.2.0 255.255.255.0",
			GW:     "172.30.30.1",
			Device: `"VLAN30"`,
		},
		{
			Name:   "8",
			Dst:    "172.30.1.0 255.255.255.0",
			GW:     "172.20.20.1",
			Device: `"VLAN20"`,
		},
		{
			Name:   "7",
			Dst:    "192.168.0.0 255.255.0.0",
			GW:     "172.16.0.1",
			Device: `"lan"`,
		},
		{
			Name:   "9",
			Dst:    "192.168.161.0 255.255.255.0",
			GW:     "192.168.79.254",
			Device: `"wan2"`,
		},
	}
	testVInfo2 = []app.VipInfo{
		{
			Name:        "vip_normal2",
			ExtIP:       "1.1.1.11",
			MappedIP:    "200.200.200.211",
			SrcFilter:   nil,
			Service:     nil,
			PortForward: false,
		},
	}
	testVInfo3 = []app.VipInfo{
		{
			Name:        "op_service",
			ExtIP:       "1.1.1.2",
			MappedIP:    "200.200.200.220",
			SrcFilter:   nil,
			Service:     []string{`"FTP"`},
			PortForward: true,
			MappedPort:  "2121",
		},
		{
			Name:        "op_portforward",
			ExtIP:       "1.1.1.3",
			MappedIP:    "200.200.200.230",
			SrcFilter:   nil,
			Service:     nil,
			PortForward: true,
			ExtPort:     "80",
			MappedPort:  "8001",
			Protocol:    "tcp",
		},
	}
	testAddr1 = addrArgs{
		addrs:          []string{"8.8.8.8"},
		services:       []string{`"ALL"`},
		allInfo: app.AllInfo{
			Services:          []string{`"ALL"`},
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: nil,
			ServicePort:       nil,
			Env:               "test",
			AddressInfoAll:    testAddrInfoAll,
			AddrGrpInfoAll:    nil,
			SRouteInfoAll: testSRouteInfoAll,
			IntfInfoAll: app.IntfInfoAll,
		},
	}
	testAddr2 = addrArgs{
		addrs:             []string{"192.168.111.0/24"},
		services:          []string{`"SMTP"`},
		vInfo:             nil,
		ServiceInfoAll:    testServiceInfoAll,
		ServiceGrpInfoAll: testServiceGrpInfoAll,
		allInfo: app.AllInfo{
			Services:          []string{`"SMTP"`},
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: nil,
			ServicePort:       nil,
			Env:               "test",
			AddressInfoAll:    testAddrInfoAll,
			AddrGrpInfoAll:    testAddrGrpInfoAll,
			SRouteInfoAll: testSRouteInfoAll,
			IntfInfoAll: app.IntfInfoAll,
		},
	}
	testAddr3 = addrArgs{
		addrs:             []string{"iprange"},
		services:          []string{`"IMAP"`, `"FTP"`},
		ServiceInfoAll:    testServiceInfoAll,
		ServiceGrpInfoAll: testServiceGrpInfoAll,
		allInfo: app.AllInfo{
			Services:          []string{`"IMAP"`, `"FTP"`},
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: nil,
			ServicePort:       nil,
			Env:               "test",
			AddressInfoAll:    testAddrInfoAll,
			AddrGrpInfoAll:    nil,
			SRouteInfoAll: testSRouteInfoAll,
			IntfInfoAll: app.IntfInfoAll,
		},
	}
	testAddr4 = addrArgs{
		addrs:             []string{"local.test.com"},
		services:          []string{`"ALL_TCP"`},
		ServiceInfoAll:    testServiceInfoAll,
		ServiceGrpInfoAll: testServiceGrpInfoAll,
		allInfo: app.AllInfo{
			Services:          []string{`"ALL_TCP"`},
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: nil,
			ServicePort:       nil,
			Env:               "test",
			AddressInfoAll:    testAddrInfoAll,
			AddrGrpInfoAll:    testAddrGrpInfoAll,
			SRouteInfoAll: testSRouteInfoAll,
			IntfInfoAll: app.IntfInfoAll,
		},
	}
	testAddr5 = addrArgs{
		addrs:             []string{`"testgrp"`},
		services:          []string{`"ALL_UDP"`},
		ServiceInfoAll:    testServiceInfoAll,
		ServiceGrpInfoAll: testServiceGrpInfoAll,
		allInfo: app.AllInfo{
			Services:          []string{`"ALL_UDP"`},
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: nil,
			ServicePort:       nil,
			Env:               "test",
			AddressInfoAll:    testAddrInfoAll,
			AddrGrpInfoAll:    testAddrGrpInfoAll,
			SRouteInfoAll: testSRouteInfoAll,
			IntfInfoAll: app.IntfInfoAll,
		},
	}
	testAddr6 = addrArgs{
		addrs:             []string{"192.168.1.100", "192.168.113.128/28"},
		services:          []string{`"Servicegrp"`},
		ServiceInfoAll:    testServiceInfoAll,
		ServiceGrpInfoAll: testServiceGrpInfoAll,
		allInfo: app.AllInfo{
			Services:          []string{`"Servicegrp"`},
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: nil,
			ServicePort:       nil,
			Env:               "test",
			AddressInfoAll:    testAddrInfoAll,
			AddrGrpInfoAll:    testAddrGrpInfoAll,
			SRouteInfoAll: testSRouteInfoAll,
			IntfInfoAll: app.IntfInfoAll,
		},
	}
	testAddr7 = addrArgs{
		intf:           `"wan2"`,
		addrs:          []string{`"all"`},
		services:       []string{`"ALL"`},
		usedAddress:    nil,
		allInfo: app.AllInfo{
			Services:          []string{`"ALL"`},
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: nil,
			ServicePort:       nil,
			Env:               "test",
			AddressInfoAll:    testAddrInfoAll,
			AddrGrpInfoAll:    testAddrGrpInfoAll,
			SRouteInfoAll: testSRouteInfoAll,
			IntfInfoAll: app.IntfInfoAll,
		},
	}
	testAddr8 = addrArgs{
		intf:           `"not_exist_intf"`,
		addrs:          []string{`"all"`},
		services:       []string{`"ALL"`},
		usedAddress:    nil,
		allInfo: app.AllInfo{
			Services:          []string{`"ALL"`},
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: nil,
			ServicePort:       nil,
			Env:               "test",
			AddressInfoAll:    testAddrInfoAll,
			AddrGrpInfoAll:    testAddrGrpInfoAll,
			SRouteInfoAll: testSRouteInfoAll,
			IntfInfoAll: app.IntfInfoAll,
		},
	}
)

func Test_HandleDescriptionOutput(t *testing.T) {
	type args struct {
		name        string
		ippool      string
		srcFQDNFlag bool
		dstFQDNFlag bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				name:        "not FQDN",
				srcFQDNFlag: false,
				dstFQDNFlag: false,
			},
			want: "policy name = not FQDN",
		},
		{
			name: "FQDN1",
			args: args{
				name:        "src FQDN",
				srcFQDNFlag: true,
				dstFQDNFlag: false,
			},
			want: "policy name = src FQDN FQDN Policy",
		},
		{
			name: "FQDN2",
			args: args{
				name:        "dst FQDN",
				srcFQDNFlag: false,
				dstFQDNFlag: true,
			},
			want: "policy name = dst FQDN FQDN Policy",
		},
		{
			name: "FQDN3",
			args: args{
				name:        "src_dst FQDN",
				srcFQDNFlag: true,
				dstFQDNFlag: true,
			},
			want: "policy name = src_dst FQDN FQDN Policy",
		},
		{
			name: "FQDN4",
			args: args{
				name:        "src FQDN2",
				ippool:      "1.1.1.1-1.1.1.4",
				srcFQDNFlag: true,
				dstFQDNFlag: false,
			},
			want: "policy name = src FQDN2 src_nat_ip=1.1.1.1-1.1.1.4 FQDN Policy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleDescriptionOutput(tt.args.name, tt.args.ippool, tt.args.srcFQDNFlag, tt.args.dstFQDNFlag); got != tt.want {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleDescriptionOutput()\n")
			}
		})
	}
}

func TestHandleSrcAddressOutput(t *testing.T) {
	tests := []struct {
		name  string
		args  addrArgs
		want  []string
		want1 bool
	}{
		{
			name:  "normal",
			args:  testAddr1,
			want:  []string{"8.8.8.8"},
			want1: false,
		},
		{
			name:  "subnet_24",
			args:  testAddr2,
			want:  []string{"192.168.111.1", "192.168.111.254"},
			want1: false,
		},
		{
			name:  "iprange",
			args:  testAddr3,
			want:  []string{"192.167.1.1", "192.167.1.20"},
			want1: false,
		},
		{
			name:  "local FQDN",
			args:  testAddr4,
			want:  []string{"local.test.com"},
			want1: true,
		},
		{
			name:  "testgrp",
			args:  testAddr5,
			want:  []string{"local.testsub.com", "192.168.110.65", "192.168.110.126", "192.168.22.25"},
			want1: true,
		},
		{
			name:  "multi addr",
			args:  testAddr6,
			want:  []string{"192.168.1.100", "192.168.113.129", "192.168.113.142"},
			want1: false,
		},
		{
			name:  "all_exist_intf",
			args:  testAddr7,
			want:  []string{"192.168.130.1"},
			want1: false,
		},
		{
			name:  "all_not_exist_intf",
			args:  testAddr8,
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := app.HandleSrcAddressOutput(tt.args.intf, tt.args.addrs, tt.args.usedAddress, tt.args.allInfo)
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleSrcAddressOutput() got")
			}
			if got1 != tt.want1 {
				t.Errorf("HandleSrcAddressOutput() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHandleDstAddressOutput(t *testing.T) {
	tests := []struct {
		name  string
		args  addrArgs
		want  []string
		want1 bool
	}{
		{
			name:  "normal",
			args:  testAddr1,
			want:  []string{"8.8.8.8", "8.8.8.8", "8.8.8.8", "8.8.8.8", "8.8.8.8", "8.8.8.8", "8.8.8.8", "8.8.8.8", "8.8.8.8", "8.8.8.8", "8.8.8.8"},
			want1: false,
		},
		{
			name:  "subnet_24",
			args:  testAddr2,
			want:  []string{"192.168.111.1", "192.168.111.254"},
			want1: false,
		},
		{
			name:  "iprange",
			args:  testAddr3,
			want:  []string{"192.167.1.1", "192.167.1.1", "192.167.1.1", "192.167.1.20", "192.167.1.20", "192.167.1.20"},
			want1: false,
		},
		{
			name:  "local FQDN",
			args:  testAddr4,
			want:  []string{"local.test.com"},
			want1: true,
		},
		{
			name:  "testgrp",
			args:  testAddr5,
			want:  []string{"local.testsub.com", "192.168.110.65", "192.168.110.126", "192.168.22.25"},
			want1: true,
		},
		{
			name:  "multi addr",
			args:  testAddr6,
			want:  []string{"192.168.1.100", "192.168.1.100", "192.168.1.100", "192.168.1.100", "192.168.113.129", "192.168.113.129", "192.168.113.129", "192.168.113.129", "192.168.113.142", "192.168.113.142", "192.168.113.142", "192.168.113.142"},
			want1: false,
		},
		{
			name:  "all_exist_intf",
			args:  testAddr7,
			want:  []string{"192.168.130.1", "192.168.130.1", "192.168.130.1", "192.168.130.1", "192.168.130.1", "192.168.130.1", "192.168.130.1", "192.168.130.1", "192.168.130.1", "192.168.130.1", "192.168.130.1"},
			want1: false,
		},
		{
			name:  "all_not_exist_intf",
			args:  testAddr8,
			want:  nil,
			want1: false,
		},
		{
			name: "normal vip",
			args: addrArgs{
				addrs:             []string{"vip_normal2"},
				services:          []string{`"DNS"`},
				vInfo:             testVInfo2,
				ServiceInfoAll:    testServiceInfoAll,
				ServiceGrpInfoAll: nil,
				allInfo: app.AllInfo{
					Services:          []string{`"DNS"`},
					ServiceInfoAll:    testServiceInfoAll,
					ServiceGrpInfoAll: nil,
					ProxyMode:         true,
					VInfo: testVInfo2,
				},
			},
			want:  []string{"200.200.200.211", "200.200.200.211"},
			want1: false,
		},
		{
			name: "multi vip",
			args: addrArgs{
				addrs:             []string{"op_service", "op_portforward"},
				services:          []string{`"ALL_TCP"`},
				vInfo:             testVInfo3,
				ServiceInfoAll:    testServiceInfoAll,
				ServiceGrpInfoAll: nil,
				allInfo: app.AllInfo{
					Services:          []string{`"ALL_TCP"`},
					ServiceInfoAll:    testServiceInfoAll,
					ServiceGrpInfoAll: nil,
					ProxyMode:         true,
					VInfo: testVInfo3,
				},
			},
			want:  []string{"200.200.200.220", "200.200.200.230"},
			want1: false,
		},
		{
			name: "vipgrp",
			args: addrArgs{
				addrs:             []string{"virtual-ip-group1"},
				services:          []string{`"ALL"`},
				vInfo:             testVInfo,
				ServiceInfoAll:    testServiceInfoAll,
				ServiceGrpInfoAll: nil,
				allInfo: app.AllInfo{
					Services:          []string{`"ALL"`},
					ServiceInfoAll:    testServiceInfoAll,
					ServiceGrpInfoAll: nil,
					ProxyMode:         true,
					VInfo:             testVInfo,
				},
			},
			want:  []string{"172.16.0.201", "172.16.0.201", "172.16.0.201", "172.16.0.201", "172.16.0.201", "172.16.0.201", "172.16.0.201", "172.16.0.201", "172.16.0.201", "172.16.0.201", "172.16.0.201", "200.200.200.210", "200.200.200.210", "200.200.200.200"},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := app.HandleDstAddressOutput(tt.args.intf, tt.args.addrs, tt.args.usedAddress, tt.args.allInfo)
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleDstAddressOutput()\n")
			}
			if got1 != tt.want1 {
				t.Errorf("HandleDstAddressOutput() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_HandleSrcNATAddressOutPut(t *testing.T) {
	type args struct {
		natbool  bool
		poolname string
		dstintf  string
		dstaddr  string
		ippools  []app.IPPoolInfo
		allInfo app.AllInfo
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "false natbool",
			args: args{
				natbool:  false,
				poolname: "",
				dstintf:  "wan1",
				ippools:  testIPPoolInfoAll,
				allInfo:  testDefaultAllInfo,
			},
			want:  "",
			want1: "",
		},
		{
			name: "no exist poolname",
			args: args{
				natbool:  false,
				poolname: "hogefuga",
				dstintf:  "wan1",
				ippools:  testIPPoolInfoAll,
				allInfo:  testDefaultAllInfo,
			},
			want:  "",
			want1: "",
		},
		{
			name: "dst intf",
			args: args{
				natbool:  true,
				poolname: "",
				dstintf:  "wan1",
				ippools:  testIPPoolInfoAll,
				allInfo:  testDefaultAllInfo,
			},
			want:  "10.10.10.1",
			want1: "",
		},
		{
			name: "1on1",
			args: args{
				natbool:  true,
				poolname: "1on1",
				dstintf:  "",
				ippools:  testIPPoolInfoAll,
				allInfo:  testDefaultAllInfo,
			},
			want:  "2.2.2.2",
			want1: "2.2.2.2-2.2.2.4",
		},
		{
			name: "fix_port_range",
			args: args{
				natbool:  true,
				poolname: "fix_port_range",
				dstintf:  "",
				ippools:  testIPPoolInfoAll,
				allInfo:  testDefaultAllInfo,
			},
			want:  "3.3.3.3",
			want1: "3.3.3.3-3.3.3.33",
		},
		{
			name: "port_block_assign",
			args: args{
				natbool:  true,
				poolname: "port_block_assign",
				dstintf:  "",
				ippools:  testIPPoolInfoAll,
				allInfo:  testDefaultAllInfo,
			},
			want:  "6.6.6.6",
			want1: "6.6.6.6-6.6.6.66",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := app.HandleSrcNATAddressOutPut(tt.args.natbool, tt.args.poolname, tt.args.dstintf, tt.args.dstaddr, tt.args.ippools, tt.args.allInfo)
			if got != tt.want {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				// wantのデータ数は1
				t.Error("HandleSrcNATAddressOutPut()\n")
			}
			if got1 != tt.want1 {
				fmt.Printf("got1 = %+v\n", got1)
				fmt.Printf("want1 = %+v\n", tt.want1)
				t.Error("HandleSrcNATAddressOutPut()\n")
			}
		})
	}
}
