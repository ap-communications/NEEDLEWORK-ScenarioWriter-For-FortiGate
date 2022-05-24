package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate/app"
)

type args struct {
	allInfo app.AllInfo
}

var (
	testServiceInfoAll = []app.ServiceInfo{
		{
			Name: `"PING"`,
			ICMP: true,
			TCP:  nil,
			UDP:  nil,
		},
		{
			Name: `"ALL_ICMP"`,
			ICMP: true,
			TCP:  nil,
			UDP:  nil,
		},
		{
			Name: `"ALL"`,
			ICMP: true,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "1-65535",
				},
			},
			UDP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "1-65535",
				},
			},
		},
		{
			Name: `"ALL_TCP"`,
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "1-65535",
				},
			},
			UDP: nil,
		},
		{
			Name: `"ALL_UDP"`,
			ICMP: false,
			TCP:  nil,
			UDP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "1-65535",
				},
			},
		},
		{
			Name: `"NTP"`,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "12312",
				},
			},
			UDP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "123",
				},
			},
		},
		{
			Name: `"DNS"`,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "53",
				},
			},
			UDP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "53",
				},
			},
		},
		{
			Name: `"HTTP"`,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "80",
				},
			},
		},
		{
			Name: `"HTTPS"`,
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "443",
				},
			},
			UDP: nil,
		},
		{
			Name: `"FTP"`,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "21",
				},
			},
		},
		{
			Name: `"IMAP"`,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "143",
				},
			},
		},
		{
			Name: `"IMAPS"`,
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "993",
				},
			},
			UDP: nil,
		},
		{
			Name: `"POP3"`,
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "110",
				},
			},
			UDP: nil,
		},
		{
			Name: `"POP3S"`,
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "995",
				},
			},
			UDP: nil,
		},
		{
			Name: `"SMTP"`,
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "25",
				},
			},
			UDP: nil,
		},
		{
			Name: "SMTPS",
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "465",
				},
			},
			UDP: nil,
		},
		{
			Name: "SAMBA",
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "139",
				},
			},
			UDP: nil,
		},
		{
			Name: "service-range-udp",
			UDP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "112-114",
				},
			},
		},
		{
			Name: `"service-multiple-tcp"`,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "554",
				},
				{
					SrcPort: "",
					DstPort: "7070-7072",
				},
				{
					SrcPort: "123-125",
					DstPort: "8554",
				},
			},
		},
		{
			Name: "service-udp-only",
			UDP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "53",
				},
			},
		},
		{
			Name: "service-tcp-udp",
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "80",
				},
			},
			UDP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "53",
				},
			},
		},
		{
			Name: "service-tcp-udp-src&dst",
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "113-114",
					DstPort: "1223-1224",
				},
			},
			UDP: []app.PortAssignInfo{
				{
					SrcPort: "223",
					DstPort: "224-226",
				},
			},
		},
		{
			Name: "serv",
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "22220-22222",
					DstPort: "11111",
				},
			},
			UDP: []app.PortAssignInfo{
				{
					SrcPort: "4444",
					DstPort: "33330-33333",
				},
			},
		},
		{
			Name: "ser2",
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "123",
					DstPort: "11122",
				},
				{
					SrcPort: "123",
					DstPort: "11155",
				},
			},
			UDP: nil,
		},
		{
			Name: "servi",
			ICMP: false,
			TCP: []app.PortAssignInfo{
				{
					SrcPort: "",
					DstPort: "554",
				},
				{
					SrcPort: "",
					DstPort: "7070",
				},
				{
					SrcPort: "",
					DstPort: "8554",
				},
			},
			UDP: nil,
		},
	}
	testServiceGrpInfoAll = []app.ServiceGrpInfo{
		{
			Name:   `"Email Access"`,
			Member: []string{`"DNS"`, `"IMAP"`, `"IMAPS"`, `"POP3"`, `"POP3S"`, `"SMTP"`, "SMTPS"},
		},
		{
			Name:   `"Web Access"`,
			Member: []string{`"DNS"`, `"HTTP"`, `"HTTPS"`},
		},
		{
			Name:   `"Windows AD"`,
			Member: []string{`"DCE-RPC"`, `"DNS"`, `"KERBEROS"`, `"LDAP"`, `"LDAP_UDP"`, `"SAMBA"`, `"SMB"`},
		},
		{
			Name:   `"Exchange Server"`,
			Member: []string{`"DCE-RPC"`, `"DNS"`, `"HTTPS"`},
		},
		{
			Name:   `"Servicegrp"`,
			Member: []string{"servi", "serv", "ser2"},
		},
	}
	testVInfo = []app.VipInfo{
		{
			Name:        "virtual-ip1",
			ExtIP:       "10.0.0.101",
			MappedIP:    "172.16.0.201",
			SrcFilter:   nil,
			Service:     nil,
			PortForward: false,
			ExtPort:     "",
			MappedPort:  "",
			Protocol:    "",
		},
		{
			Name:        "op_srcip",
			ExtIP:       "1.1.1.1",
			MappedIP:    "200.200.200.210",
			SrcFilter:   []string{"3.3.3.3", "3.3.3.4-3.3.3.6", "3.3.3.10/30"},
			Service:     []string{`"HTTP"`, `"IMAP"`},
			PortForward: true,
			ExtPort:     "",
			MappedPort:  "25525",
			Protocol:    "",
		},
		{
			Name:        "dstnat_111.111.111.112_udp",
			ExtIP:       "111.111.111.112",
			MappedIP:    "200.200.200.200",
			SrcFilter:   nil,
			Service:     nil,
			PortForward: true,
			ExtPort:     "123",
			MappedPort:  "8000",
			Protocol:    "udp",
		},
	}
	testVInfo4 = []app.VipInfo{
		{
			Name:        "virtual-ip1",
			ExtIP:       "10.0.0.101",
			MappedIP:    "172.16.0.201",
			SrcFilter:   nil,
			Service:     nil,
			PortForward: false,
			ExtPort:     "",
			MappedPort:  "",
			Protocol:    "",
		},
		{
			Name:        "op_srcip2",
			ExtIP:       "1.1.1.5",
			MappedIP:    "200.200.200.215",
			SrcFilter:   []string{"3.3.3.20-3.3.3.22", "3.3.3.30/30"},
			Service:     []string{`"SMTP"`},
			PortForward: true,
			ExtPort:     "",
			MappedPort:  "24",
			Protocol:    "",
		},
		{
			Name:        "op_all_service",
			ExtIP:       "1.1.1.6",
			MappedIP:    "200.200.200.216",
			SrcFilter:   nil,
			Service:     []string{`"ALL"`},
			PortForward: false,
			ExtPort:     "",
			MappedPort:  "",
			Protocol:    "",
		},
		{
			Name:        "dstnat_111.111.111.115_tcp",
			ExtIP:       "111.111.111.115",
			MappedIP:    "200.200.200.205",
			SrcFilter:   nil,
			Service:     nil,
			PortForward: true,
			ExtPort:     "25",
			MappedPort:  "25",
			Protocol:    "tcp",
		},
		{
			Name:        "op_srcip3",
			ExtIP:       "1.1.1.8",
			MappedIP:    "200.200.200.218",
			SrcFilter:   []string{"3.3.3.28-3.3.3.32", "3.3.3.40/30"},
			Service:     []string{`"ALL_TCP"`, `"PING"`, `"DNS"`},
			PortForward: true,
			ExtPort:     "",
			MappedPort:  "25",
			Protocol:    "",
		},
	}
	testWfFilterInfoAll = []app.WfFilterInfo{
		{
			Name: "1",
			Entries: []app.Entry{
				{
					URL:    "techblog.ap-com.co.jp",
					Action: "block",
				},
				{
					URL:    "www.ap-com.co.jp",
					Action: "block",
				},
				{
					URL:    "support.needlework.jp/faq",
					Action: "block",
				},
				{
					URL:    "support.needlework.jp",
					Action: "monitor",
				},
			},
		},
		{
			Name: "2",
			Entries: []app.Entry{
				{
					URL:    "support.needlework.jp/faq",
					Action: "monitor",
				},
				{
					URL:    "support.needlework.jp",
					Action: "block",
				},
			},
		},
	}
	testWfProfileInfoAll = []app.WfProfileInfo{
		{
			Name:        `"g-default"`,
			URLTableNum: "",
		},
		{
			Name:        `"wifi-default"`,
			URLTableNum: "",
			FtgdFilter: true,
		},
		{
			Name:        `"wb-sniffer-profile"`,
			URLTableNum: "1",
		},
		{
			Name:        `"default"`,
			URLTableNum: "2",
		},
	}
	testAVProfileInfoAll = []app.AVProfileInfo{
		{
			Name: `"default"`,
			Mode: `""`,
			Config: map[string]bool{
				"ftp":  false,
				"http": false,
				"imap": false,
				"pop3": false,
			},
		},
		{
			Name: `"wifi-default"`,
			Mode: `""`,
			Config: map[string]bool{
				"ftp":  true,
				"http": true,
				"imap": true,
				"pop3": true,
				"smtp": true,
			},
		},
		{
			Name: `"av-sniffer-profile"`,
			Mode: `""`,
			Config: map[string]bool{
				"ftp":  true,
				"pop3": true,
				"smtp": true,
			},
		},
	}
	testArgs1 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"service-multiple-tcp"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			VInfo:             nil,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs2 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"ALL_ICMP"`},
			Act:               "",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs3 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"FTP"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           `"default"`,
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs4 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"Servicegrp"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           `"wifi-default"`,
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs5 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"Email Access"`},
			Act:               "",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			VInfo:             nil,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs6 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"ALL"`},
			Act:               "accept",
			WfpName:           `"wb-sniffer-profile"`,
			AvpName:           `"av-sniffer-profile"`,
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: nil,
			VInfo:             nil,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs7 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"DNS"`, `"HTTP"`},
			Act:               "accept",
			WfpName:           `"default"`,
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs8 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"ALL"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: nil,
			ProxyMode:         true,
			VInfo:             testVInfo,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs9 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"ALL_TCP"`, `"ALL_UDP"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs10 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"ALL_ICMP"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             nil,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs11 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"DNS"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo: []app.VipInfo{
				{
					Name:        "virtual-ip1",
					ExtIP:       "10.0.0.101",
					MappedIP:    "172.16.0.201",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: false,
					ExtPort:     "",
					MappedPort:  "",
					Protocol:    "",
				},
			},
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs12 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"SMTP"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			VInfo:             testVInfo4,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs13 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"PING"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			VInfo: []app.VipInfo{
				{
					Name:        "virtual-ip1",
					ExtIP:       "10.0.0.101",
					MappedIP:    "172.16.0.201",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: false,
					ExtPort:     "",
					MappedPort:  "",
					Protocol:    "",
				},
			},
			VipGrpInfo:        nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs14 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"ALL"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			VInfo: []app.VipInfo{
				{
					Name:        "virtual-ip1",
					ExtIP:       "10.0.0.101",
					MappedIP:    "172.16.0.201",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: false,
					ExtPort:     "",
					MappedPort:  "",
					Protocol:    "",
				},
			},
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs15 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"DNS"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo: []app.VipInfo{
				{
					Name:        "virtual-ip1",
					ExtIP:       "10.0.0.101",
					MappedIP:    "172.16.0.201",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: false,
					ExtPort:     "",
					MappedPort:  "",
					Protocol:    "",
				},
			},
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs16 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"ALL_UDP"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: nil,
			ProxyMode:         true,
			VInfo:             testVInfo4,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs17 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"ALL"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: testServiceGrpInfoAll,
			ProxyMode:         true,
			VInfo:             testVInfo4,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs18 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"PING"`},
			Act:               "accept",
			WfpName:           `"wb-sniffer-profile"`,
			AvpName:           `"av-sniffer-profile"`,
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: nil,
			ProxyMode:         true,
			VInfo:             testVInfo4,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs19 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"ALL_TCP"`},
			Act:               "accept",
			WfpName:           `"wb-sniffer-profile"`,
			AvpName:           `"av-sniffer-profile"`,
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: nil,
			VInfo:             testVInfo3,
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs20 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"DNS"`, `"FTP"`},
			Act:               "accept",
			WfpName:           `"wb-sniffer-profile"`,
			AvpName:           `"av-sniffer-profile"`,
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: nil,
			VInfo: []app.VipInfo{
				{
					Name:        "op_portforward",
					ExtIP:       "1.1.1.8",
					MappedIP:    "200.200.200.233",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: true,
					ExtPort:     "53",
					MappedPort:  "53",
					Protocol:    "tcp",
				},
			},
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs21 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"NTP"`},
			Act:               "accept",
			WfpName:           "",
			AvpName:           "",
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: nil,
			ProxyMode:         true,
			VInfo: []app.VipInfo{
				{
					Name:     "test_vip3",
					ExtIP:    "1.1.1.10",
					MappedIP: "200.200.200.240",
				},
				{
					Name:        "op_portforward2",
					ExtIP:       "1.1.1.9",
					MappedIP:    "200.200.200.235",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: true,
					ExtPort:     "53",
					MappedPort:  "53",
					Protocol:    "tcp",
				},
				{
					Name:        "op_portforward3",
					ExtIP:       "1.1.1.11",
					MappedIP:    "200.200.200.237",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: true,
					ExtPort:     "123",
					MappedPort:  "123",
					Protocol:    "udp",
				},
				{
					Name:      "op_portforward",
					ExtIP:     "1.1.1.8",
					MappedIP:  "200.200.200.233",
					SrcFilter: nil,
					Service:   []string{`"ALL"`, `"ALL_UDP"`, `"DNS"`, `"NTP"`},
				},
			},
			VipGrpInfo:        nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs22 = args{
		allInfo: app.AllInfo{
			Services: []string{`"DNS"`},
			Act:      "accept",
			WfpName:  "",
			AvpName:  "",
			ServiceInfoAll: []app.ServiceInfo{
				{
					Name: `"ALL_UDP"`,
					ICMP: false,
					TCP:  nil,
					UDP: []app.PortAssignInfo{
						{
							SrcPort: "",
							DstPort: "1-65535",
						},
					},
				},
				{
					Name: `"DNS"`,
					TCP: []app.PortAssignInfo{
						{
							SrcPort: "",
							DstPort: "5312",
						},
					},
					UDP: []app.PortAssignInfo{
						{
							SrcPort: "",
							DstPort: "53",
						},
					},
				},
			},
			ServiceGrpInfoAll: nil,
			ProxyMode:         true,
			VInfo: []app.VipInfo{
				{
					Name:     "test_vip3",
					ExtIP:    "1.1.1.10",
					MappedIP: "200.200.200.240",
				},
				{
					Name:        "op_portforward2",
					ExtIP:       "1.1.1.9",
					MappedIP:    "200.200.200.235",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: true,
					ExtPort:     "53",
					MappedPort:  "53",
					Protocol:    "tcp",
				},
				{
					Name:        "op_portforward3",
					ExtIP:       "1.1.1.11",
					MappedIP:    "200.200.200.237",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: true,
					ExtPort:     "53",
					MappedPort:  "53",
					Protocol:    "udp",
				},
				{
					Name:      "op_portforward",
					ExtIP:     "1.1.1.8",
					MappedIP:  "200.200.200.233",
					SrcFilter: nil,
					Service:   []string{`"ALL"`, `"ALL_UDP"`, `"DNS"`, `"NTP"`},
				},
			},
			VipGrpInfo:        nil,
			WfProfileInfoAll:  nil,
			WfFilterInfoAll:   nil,
			AVProfileInfoAll:  nil,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
	testArgs23 = args{
		allInfo: app.AllInfo{
			Services:          []string{`"HTTP"`, `"IMAP"`, `"HTTPS"`},
			Act:               "accept",
			WfpName:           `"wifi-default"`,
			ServiceInfoAll:    testServiceInfoAll,
			ServiceGrpInfoAll: nil,
			VInfo: []app.VipInfo{
				{
					Name:        "op_portforward",
					ExtIP:       "1.1.1.8",
					MappedIP:    "200.200.200.233",
					SrcFilter:   nil,
					Service:     nil,
					PortForward: true,
					ExtPort:     "443",
					MappedPort:  "443",
					Protocol:    "tcp",
				},
				{
					Name:        "op_portforward",
					ExtIP:       "1.1.1.8",
					MappedIP:    "200.200.200.233",
					SrcFilter:   nil,
					Service:     []string{`"HTTP"`},
				},
			},
			VipGrpInfo:        nil,
			WfProfileInfoAll:  testWfProfileInfoAll,
			WfFilterInfoAll:   testWfFilterInfoAll,
			AVProfileInfoAll:  testAVProfileInfoAll,
			ProxyModeProtocol: app.ProxyModeProtocol,
			ServicePort:       app.ServicePort,
			Env:               "test",
		},
	}
)

func Test_HandleProtocolOutput(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "service-multiple-tcp",
			args: testArgs1,
			want: []string{"tcp"},
		},
		{
			name: "icmp",
			args: testArgs2,
			want: []string{"icmp"},
		},
		{
			name: "ftp",
			args: testArgs3,
			want: []string{"ftp", "ftpa"},
		},
		{
			name: "Servicegrp",
			args: testArgs4,
			want: []string{"tcp", "tcp", "udp", "tcp"},
		},
		{
			name: "Email Access",
			args: testArgs5,
			want: []string{"dnst", "dns", "imap", "tcp", "tcp", "tcp", "smtp", "tcp"},
		},
		{
			name: "all",
			args: testArgs6,
			want: app.AllNWProtocol,
		},
		{
			name: "multi service",
			args: testArgs7,
			want: []string{"dnst", "dns", "http"},
		},
		{
			name: "vipgrp",
			args: testArgs8,
			want: []string{"icmp", "tcp", "udp", "dnst", "dns", "http", "https", "ftp", "ftpa", "imap", "smtp", "http", "imap", "udp"},
		},
		{
			name: "all-tcp&udp",
			args: testArgs9,
			want: []string{"tcp", "udp"},
		},
		{
			name: "all-icmp",
			args: testArgs10,
			want: []string{"icmp"},
		},
		{
			name: "vip-dns",
			args: testArgs11,
			want: []string{"dnst", "dns"},
		},
		{
			name: "vipgrp-smtp",
			args: testArgs12,
			want: []string{"smtp", "undefined", "smtp", "smtp", "smtp", "undefined", "undefined"},
		},
		{
			name: "vip-ping",
			args: testArgs13,
			want: []string{"icmp"},
		},
		{
			name: "vip-all",
			args: testArgs14,
			want: app.AllNWProtocol,
		},
		{
			name: "vipgrp-all-udp",
			args: testArgs16,
			want: []string{"udp", "undefined", "udp", "undefined", "undefined", "undefined", "undefined", "udp"},
		},
		{
			name: "vipgrp-all",
			args: testArgs17,
			want: []string{"icmp", "tcp", "udp", "dnst", "dns", "http", "https", "ftp", "ftpa", "imap", "smtp", "smtp", "icmp", "tcp", "udp", "dnst", "dns", "http", "https", "ftp", "ftpa", "imap", "smtp", "tcp", "tcp", "icmp", "dnst", "dns"},
		},
		{
			name: "vipgrp-ping",
			args: testArgs18,
			want: []string{"icmp", "undefined", "icmp", "undefined", "undefined", "icmp", "undefined", "undefined"},
		},
		{
			name: "vipgrp-ftp",
			args: testArgs19,
			want: []string{"tcp", "tcp"},
		},
		{
			name: "vip-dns&ftp",
			args: testArgs20,
			want: []string{"dnst", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-ntp",
			args: testArgs21,
			want: []string{"tcp", "udp", "undefined", "undefined", "undefined", "udp", "tcp", "udp", "undefined", "udp", "undefined", "undefined", "tcp", "udp"},
		},
		{
			name: "dns-tcp-port-different",
			args: testArgs22,
			want: []string{"dnst", "dns", "undefined", "undefined", "undefined", "dns", "dnst", "dns", "undefined", "dns", "dnst", "dns", "undefined", "undefined"},
		},
		{
			name: "ftgd-block-http&https",
			args: testArgs23,
			want: []string{"undefined", "undefined", "https", "http", "undefined", "undefined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleProtocolOutput(tt.args.allInfo); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleProtocolOutput()\n")
			}
		})
	}
}

func TestHandleSrcPortOutput(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "service-multiple-tcp",
			args: testArgs1,
			want: []string{""},
		},
		{
			name: "icmp",
			args: testArgs2,
			want: []string{""},
		},
		{
			name: "ftp",
			args: testArgs3,
			want: []string{"", ""},
		},
		{
			name: "Servicegrp",
			args: testArgs4,
			want: []string{"", "22220", "4444", "123"},
		},
		{
			name: "Email Access",
			args: testArgs5,
			want: []string{"", "", "", "", "", "", "", ""},
		},
		{
			name: `"all"`,
			args: testArgs6,
			want: []string{"", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "multi service",
			args: testArgs7,
			want: []string{"", "", ""},
		},
		{
			name: "vipgrp",
			args: testArgs8,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "all-tcp&udp",
			args: testArgs9,
			want: []string{"", ""},
		},
		{
			name: "all-icmp",
			args: testArgs10,
			want: []string{""},
		},
		{
			name: "vip-all",
			args: testArgs14,
			want: []string{"", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vip-ping",
			args: testArgs13,
			want: []string{""},
		},
		{
			name: "vipgrp-smtp",
			args: testArgs12,
			want: []string{"", "undefined", "", "", "", "undefined", "undefined"},
		},
		{
			name: "vipgrp-all-udp",
			args: testArgs16,
			want: []string{"", "undefined", "", "undefined", "undefined", "undefined", "undefined", ""},
		},
		{
			name: "vipgrp-all",
			args: testArgs17,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vipgrp-ping",
			args: testArgs18,
			want: []string{"", "undefined", "", "undefined", "undefined", "", "undefined", "undefined"},
		},
		{
			name: "vipgrp-ftp",
			args: testArgs19,
			want: []string{"", ""},
		},
		{
			name: "vip-dns&ftp",
			args: testArgs20,
			want: []string{"", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-ntp",
			args: testArgs21,
			want: []string{"", "", "undefined", "undefined", "undefined", "", "", "", "undefined", "", "undefined", "undefined", "", ""},
		},
		{
			name: "dns-tcp-port-different",
			args: testArgs22,
			want: []string{"", "", "undefined", "undefined", "undefined", "", "", "", "undefined", "", "", "", "undefined", "undefined"},
		},
		{
			name: "ftgd-block-http&https",
			args: testArgs23,
			want: []string{"undefined", "undefined", "", "", "undefined", "undefined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleSrcPortOutput(tt.args.allInfo); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleSrcPortOutput()\n")
			}
		})
	}
}

func Test_HandleDstNATAddrOutput(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "service-multiple-tcp",
			args: testArgs1,
			want: []string{""},
		},
		{
			name: "icmp",
			args: testArgs2,
			want: []string{""},
		},
		{
			name: "ftp",
			args: testArgs3,
			want: []string{"", ""},
		},
		{
			name: "Servicegrp",
			args: testArgs4,
			want: []string{"", "", "", ""},
		},
		{
			name: "Email Access",
			args: testArgs5,
			want: []string{"", "", "", "", "", "", "", ""},
		},
		{
			name: `"all"`,
			args: testArgs6,
			want: []string{"", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "multi service",
			args: testArgs7,
			want: []string{"", "", ""},
		},
		{
			name: "vipgrp",
			args: testArgs8,
			want: []string{"10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "1.1.1.1", "1.1.1.1", "111.111.111.112"},
		},
		{
			name: "vip-all",
			args: testArgs14,
			want: []string{"10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101"},
		},
		{
			name: "vip-ping",
			args: testArgs13,
			want: []string{"10.0.0.101"},
		},
		{
			name: "vipgrp-smtp",
			args: testArgs12,
			want: []string{"10.0.0.101", "undefined", "1.1.1.6", "111.111.111.115", "1.1.1.8", "undefined", "undefined"},
		},
		{
			name: "vipgrp-all-udp",
			args: testArgs16,
			want: []string{"10.0.0.101", "undefined", "1.1.1.6", "undefined", "undefined", "undefined", "undefined", "1.1.1.8"},
		},
		{
			name: "vipgrp-all",
			args: testArgs17,
			want: []string{"10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "10.0.0.101", "1.1.1.5", "1.1.1.6", "1.1.1.6", "1.1.1.6", "1.1.1.6", "1.1.1.6", "1.1.1.6", "1.1.1.6", "1.1.1.6", "1.1.1.6", "1.1.1.6", "1.1.1.6", "111.111.111.115", "1.1.1.8", "1.1.1.8", "1.1.1.8", "1.1.1.8"},
		},
		{
			name: "vipgrp-ping",
			args: testArgs18,
			want: []string{"10.0.0.101", "undefined", "1.1.1.6", "undefined", "undefined", "1.1.1.8", "undefined", "undefined"},
		},
		{
			name: "vipgrp-ftp",
			args: testArgs19,
			want: []string{"1.1.1.2", "1.1.1.3"},
		},
		{
			name: "vip-dns&ftp",
			args: testArgs20,
			want: []string{"1.1.1.8", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-ntp",
			args: testArgs21,
			want: []string{"1.1.1.10", "1.1.1.10", "undefined", "undefined", "undefined", "1.1.1.11", "1.1.1.8", "1.1.1.8", "undefined", "1.1.1.8", "undefined", "undefined", "1.1.1.8", "1.1.1.8"},
		},
		{
			name: "dns-tcp-port-different",
			args: testArgs22,
			want: []string{"1.1.1.10", "1.1.1.10", "undefined", "undefined", "undefined", "1.1.1.11", "1.1.1.8", "1.1.1.8", "undefined", "1.1.1.8", "1.1.1.8", "1.1.1.8", "undefined", "undefined"},
		},
		{
			name: "ftgd-block-http&https",
			args: testArgs23,
			want: []string{"undefined", "undefined", "1.1.1.8", "1.1.1.8", "undefined", "undefined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleDstNATAddrOutput(tt.args.allInfo); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleDstNATAddrOutput()\n")
			}
		})
	}
}

func Test_HandleDstNATPortOutput(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "service-multiple-tcp",
			args: testArgs1,
			want: []string{""},
		},
		{
			name: "icmp",
			args: testArgs2,
			want: []string{""},
		},
		{
			name: "ftp",
			args: testArgs3,
			want: []string{"", ""},
		},
		{
			name: "Servicegrp",
			args: testArgs4,
			want: []string{"", "", "", ""},
		},
		{
			name: "Email Access",
			args: testArgs5,
			want: []string{"", "", "", "", "", "", "", ""},
		},
		{
			name: `"all"`,
			args: testArgs6,
			want: []string{"", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "multi service",
			args: testArgs7,
			want: []string{"", "", ""},
		},
		{
			name: "vipgrp",
			args: testArgs8,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "80", "143", "123"},
		},
		{
			name: "vip-all",
			args: testArgs14,
			want: []string{"", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vip-ping",
			args: testArgs13,
			want: []string{""},
		},
		{
			name: "vipgrp-smtp",
			args: testArgs12,
			want: []string{"", "undefined", "", "25", "25", "undefined", "undefined"},
		},
		{
			name: "vipgrp-all-udp",
			args: testArgs16,
			want: []string{"", "undefined", "", "undefined", "undefined", "undefined", "undefined", "53"},
		},
		{
			name: "vipgrp-all",
			args: testArgs17,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "25", "", "", "", "", "", "", "", "", "", "", "", "25", "80", "", "53", "53"},
		},
		{
			name: "vipgrp-ping",
			args: testArgs18,
			want: []string{"", "undefined", "", "undefined", "undefined", "", "undefined", "undefined"},
		},
		{
			name: "vipgrp-ftp",
			args: testArgs19,
			want: []string{"21", "80"},
		},
		{
			name: "vip-dns&ftp",
			args: testArgs20,
			want: []string{"53", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-ntp",
			args: testArgs21,
			want: []string{"", "", "undefined", "undefined", "undefined", "123", "", "", "undefined", "", "undefined", "undefined", "", ""},
		},
		{
			name: "dns-tcp-port-different",
			args: testArgs22,
			want: []string{"", "", "undefined", "undefined", "undefined", "53", "", "", "undefined", "", "", "", "undefined", "undefined"},
		},
		{
			name: "ftgd-block-http&https",
			args: testArgs23,
			want: []string{"undefined", "undefined", "443", "", "undefined", "undefined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleDstNATPortOutput(tt.args.allInfo); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleDstNATPortOutput()\n")
			}
		})
	}
}

func Test_HandleDstPortOutput(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "normal port",
			args: args{
				allInfo: app.AllInfo{
					Services:          []string{`"service-multiple-tcp"`},
					Act:               "accept",
					WfpName:           "",
					AvpName:           "",
					ServiceInfoAll:    testServiceInfoAll,
					ServiceGrpInfoAll: testServiceGrpInfoAll,

					VInfo:             nil,
					VipGrpInfo:        nil,
					WfProfileInfoAll:  nil,
					WfFilterInfoAll:   nil,
					AVProfileInfoAll:  nil,
					ProxyModeProtocol: app.ProxyModeProtocol,
					ServicePort:       app.ServicePort,
					Env:               "test",
				},
			},
			// want: []string{"554", "7070", "8554"},
			want: []string{"554"},
		},
		{
			name: "normal port2",
			args: args{
				allInfo: app.AllInfo{
					Services:          []string{"service-udp-only"},
					Act:               "accept",
					WfpName:           "",
					AvpName:           "",
					ServiceInfoAll:    testServiceInfoAll,
					ServiceGrpInfoAll: testServiceGrpInfoAll,
					VInfo:             nil,
					VipGrpInfo:        nil,
					WfProfileInfoAll:  nil,
					WfFilterInfoAll:   nil,
					AVProfileInfoAll:  nil,
					ProxyModeProtocol: app.ProxyModeProtocol,
					ServicePort:       app.ServicePort,
					Env:               "test",
				},
			},
			want: []string{"53"},
		},
		{
			name: "normal port3",
			args: args{
				allInfo: app.AllInfo{
					Services:          []string{"service-tcp-udp"},
					Act:               "accept",
					WfpName:           "",
					AvpName:           "",
					ServiceInfoAll:    testServiceInfoAll,
					ServiceGrpInfoAll: testServiceGrpInfoAll,
					VInfo:             nil,
					VipGrpInfo:        nil,
					WfProfileInfoAll:  nil,
					WfFilterInfoAll:   nil,
					AVProfileInfoAll:  nil,
					ProxyModeProtocol: app.ProxyModeProtocol,
					ServicePort:       app.ServicePort,
					Env:               "test",
				},
			},
			want: []string{"80", "53"},
		},
		{
			name: "normal port4",
			args: args{
				allInfo: app.AllInfo{
					Services:          []string{"service-range-udp"},
					Act:               "accept",
					WfpName:           "",
					AvpName:           "",
					ServiceInfoAll:    testServiceInfoAll,
					ServiceGrpInfoAll: testServiceGrpInfoAll,
					VInfo:             nil,
					VipGrpInfo:        nil,
					WfProfileInfoAll:  nil,
					WfFilterInfoAll:   nil,
					AVProfileInfoAll:  nil,
					ProxyModeProtocol: app.ProxyModeProtocol,
					ServicePort:       app.ServicePort,
					Env:               "test",
				},
			},
			want: []string{"112"},
		},
		{
			name: "normal port5",
			args: args{
				allInfo: app.AllInfo{
					Services:          []string{"service-tcp-udp-src&dst"},
					Act:               "accept",
					WfpName:           "",
					AvpName:           "",
					ServiceInfoAll:    testServiceInfoAll,
					ServiceGrpInfoAll: testServiceGrpInfoAll,
					VInfo:             nil,
					VipGrpInfo:        nil,
					WfProfileInfoAll:  nil,
					WfFilterInfoAll:   nil,
					AVProfileInfoAll:  nil,
					ProxyModeProtocol: app.ProxyModeProtocol,
					ServicePort:       app.ServicePort,
					Env:               "test",
				},
			},
			want: []string{"1223", "224"},
		},
		{
			name: "icmp",
			args: testArgs2,
			want: []string{""},
		},
		{
			name: "ftp",
			args: testArgs3,
			want: []string{"21", "21"},
		},
		{
			name: "Servicegrp",
			args: testArgs4,
			want: []string{"554", "11111", "33330", "11122"},
		},
		{
			name: "Email Access",
			args: testArgs5,
			want: []string{"53", "53", "143", "993", "110", "995", "25", "465"},
		},
		{
			name: `"all"`,
			args: testArgs6,
			want: []string{"", "80", "53", "53", "53", "80", "443", "21", "21", "143", "25"},
		},
		{
			name: "multi service",
			args: testArgs7,
			want: []string{"53", "53", "80"},
		},
		{
			name: "vipgrp",
			args: testArgs8,
			want: []string{"", "80", "53", "53", "53", "80", "443", "21", "21", "143", "25", "25525", "25525", "8000"},
		},
		{
			name: "all-tcp&udp",
			args: testArgs9,
			want: []string{"80", "53"},
		},
		{
			name: "all-icmp",
			args: testArgs10,
			want: []string{""},
		},
		{
			name: "vip-all",
			args: testArgs14,
			want: []string{"", "80", "53", "53", "53", "80", "443", "21", "21", "143", "25"},
		},
		{
			name: "vip-ping",
			args: testArgs13,
			want: []string{""},
		},
		{
			name: "vipgrp-smtp",
			args: testArgs12,
			want: []string{"25", "undefined", "25", "25", "25", "undefined", "undefined"},
		},
		{
			name: "vipgrp-all-udp",
			args: testArgs16,
			want: []string{"53", "undefined", "53", "undefined", "undefined", "undefined", "undefined", "25"},
		},
		{
			name: "vipgrp-all",
			args: testArgs17,
			want: []string{"", "80", "53", "53", "53", "80", "443", "21", "21", "143", "25", "24", "", "80", "53", "53", "53", "80", "443", "21", "21", "143", "25", "25", "25", "", "25", "25"},
		},
		{
			name: "vipgrp-ping",
			args: testArgs18,
			want: []string{"", "undefined", "", "undefined", "undefined", "", "undefined", "undefined"},
		},
		{
			name: "vipgrp-ftp",
			args: testArgs19,
			want: []string{"2121", "8001"},
		},
		{
			name: "vip-dns&ftp",
			args: testArgs20,
			want: []string{"53", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-ntp",
			args: testArgs21,
			want: []string{"12312", "123", "undefined", "undefined", "undefined", "123", "12312", "123", "undefined", "123", "undefined", "undefined", "12312", "123"},
		},
		{
			name: "dns-tcp-port-different",
			args: testArgs22,
			want: []string{"5312", "53", "undefined", "undefined", "undefined", "53", "5312", "53", "undefined", "53", "5312", "53", "undefined", "undefined"},
		},
		{
			name: "ftgd-block-http&https",
			args: testArgs23,
			want: []string{"undefined", "undefined", "443", "80", "undefined", "undefined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleDstPortOutput(tt.args.allInfo); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleDstPortOutput()\n")
			}
		})
	}
}

func Test_HandleOtherSettingOutput(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "service-multiple-tcp",
			args: testArgs1,
			want: []string{""},
		},
		{
			name: "icmp",
			args: testArgs2,
			want: []string{""},
		},
		{
			name: "ftp",
			args: testArgs3,
			want: []string{"Proxy mode", "Proxy mode"},
		},
		{
			name: "Servicegrp",
			args: testArgs4,
			want: []string{"", "", "", ""},
		},
		{
			name: "Email Access",
			args: testArgs5,
			want: []string{"", "", "", "", "", "", "", ""},
		},
		{
			name: "all",
			args: testArgs6,
			want: []string{"", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "multi service",
			args: testArgs7,
			want: []string{"", "", ""},
		},
		{
			name: "vipgrp",
			args: testArgs8,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vip-all",
			args: testArgs14,
			want: []string{"", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vip-ping",
			args: testArgs13,
			want: []string{""},
		},
		{
			name: "vipgrp-smtp",
			args: testArgs12,
			want: []string{"", "undefined", "", "", "", "undefined", "undefined"},
		},
		{
			name: "vipgrp-all-udp",
			args: testArgs16,
			want: []string{"", "undefined", "", "undefined", "undefined", "undefined", "undefined", ""},
		},
		{
			name: "vipgrp-all",
			args: testArgs17,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vipgrp-ping",
			args: testArgs18,
			want: []string{"", "undefined", "", "undefined", "undefined", "", "undefined", "undefined"},
		},
		{
			name: "vipgrp-ftp",
			args: testArgs19,
			want: []string{"", ""},
		},
		{
			name: "vip-dns&ftp",
			args: testArgs20,
			want: []string{"", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-ntp",
			args: testArgs21,
			want: []string{"", "", "undefined", "undefined", "undefined", "", "", "", "undefined", "", "undefined", "undefined", "", ""},
		},
		{
			name: "dns-tcp-port-different",
			args: testArgs22,
			want: []string{"", "", "undefined", "undefined", "undefined", "", "", "", "undefined", "", "", "", "undefined", "undefined"},
		},
		{
			name: "ftgd-block-http&https",
			args: testArgs23,
			want: []string{"undefined", "undefined", "", "", "undefined", "undefined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleOtherSettingOutput(tt.args.allInfo); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleOtherSettingOutput()\n")
			}
		})
	}
}

func Test_HandleURLDomainOutput(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "service-multiple-tcp",
			args: testArgs1,
			want: []string{""},
		},
		{
			name: "icmp",
			args: testArgs2,
			want: []string{""},
		},
		{
			name: "ftp",
			args: testArgs3,
			want: []string{"", ""},
		},
		{
			name: "Servicegrp",
			args: testArgs4,
			want: []string{"", "", "", ""},
		},
		{
			name: "Email Access",
			args: testArgs5,
			want: []string{"", "", "", "", "", "", "", ""},
		},
		{
			name: "all",
			args: testArgs6,
			want: []string{"", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined", "techblog.ap-com.co.jp", "www.ap-com.co.jp", "support.needlework.jp/faq", "support.needlework.jp", "techblog.ap-com.co.jp", "www.ap-com.co.jp", "support.needlework.jp/faq", "support.needlework.jp", "", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined"},
		},
		{
			name: "multi service",
			args: testArgs7,
			want: []string{"", "undefined", "", "undefined", "support.needlework.jp/faq", "support.needlework.jp"},
		},
		{
			name: "vipgrp",
			args: testArgs8,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vip-all",
			args: testArgs14,
			want: []string{"", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vip-ping",
			args: testArgs13,
			want: []string{""},
		},
		{
			name: "vipgrp-smtp",
			args: testArgs12,
			want: []string{"", "undefined", "", "", "", "undefined", "undefined"},
		},
		{
			name: "vipgrp-all-udp",
			args: testArgs16,
			want: []string{"", "undefined", "", "undefined", "undefined", "undefined", "undefined", ""},
		},
		{
			name: "vipgrp-all",
			args: testArgs17,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vipgrp-ping",
			args: testArgs18,
			want: []string{"", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined"},
		},
		{
			name: "vipgrp-ftp",
			args: testArgs19,
			want: []string{"", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-dns&ftp",
			args: testArgs20,
			want: []string{"", "undefined", "undefined", "undefined", "", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-ntp",
			args: testArgs21,
			want: []string{"", "", "undefined", "undefined", "undefined", "", "", "", "undefined", "", "undefined", "undefined", "", ""},
		},
		{
			name: "dns-tcp-port-different",
			args: testArgs22,
			want: []string{"", "", "undefined", "undefined", "undefined", "", "", "", "undefined", "", "", "", "undefined", "undefined"},
		},
		{
			name: "ftgd-block-http&https",
			args: testArgs23,
			want: []string{"undefined", "undefined", "", "", "undefined", "undefined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleURLDomainOutput(tt.args.allInfo); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleURLDomainOutput()\n")
			}
		})
	}
}

func Test_HandleAntiVirusOutput(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "service-multiple-tcp",
			args: testArgs1,
			want: []string{""},
		},
		{
			name: "icmp",
			args: testArgs2,
			want: []string{""},
		},
		{
			name: "ftp",
			args: testArgs3,
			want: []string{"", ""},
		},
		{
			name: "Servicegrp",
			args: testArgs4,
			want: []string{"", "", "", ""},
		},
		{
			name: "Email Access",
			args: testArgs5,
			want: []string{"", "", "", "", "", "", "", ""},
		},
		{
			name: "all",
			args: testArgs6,
			want: []string{"", "", "", "", "", "", "", "enable", "enable", "", "enable"},
		},
		{
			name: "multi service",
			args: testArgs7,
			want: []string{"", "", ""},
		},
		{
			name: "vipgrp",
			args: testArgs8,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vip-all",
			args: testArgs14,
			want: []string{"", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vip-ping",
			args: testArgs13,
			want: []string{""},
		},
		{
			name: "vipgrp-smtp",
			args: testArgs12,
			want: []string{"", "undefined", "", "", "", "undefined", "undefined"},
		},
		{
			name: "vipgrp-all-udp",
			args: testArgs16,
			want: []string{"", "undefined", "", "undefined", "undefined", "undefined", "undefined", ""},
		},
		{
			name: "vipgrp-all",
			args: testArgs17,
			want: []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		},
		{
			name: "vipgrp-ping",
			args: testArgs18,
			want: []string{"", "undefined", "", "undefined", "undefined", "", "undefined", "undefined"},
		},
		{
			name: "vipgrp-ftp",
			args: testArgs19,
			want: []string{"", ""},
		},
		{
			name: "vip-dns&ftp",
			args: testArgs20,
			want: []string{"", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-ntp",
			args: testArgs21,
			want: []string{"", "", "undefined", "undefined", "undefined", "", "", "", "undefined", "", "undefined", "undefined", "", ""},
		},
		{
			name: "dns-tcp-port-different",
			args: testArgs22,
			want: []string{"", "", "undefined", "undefined", "undefined", "", "", "", "undefined", "", "", "", "undefined", "undefined"},
		},
		{
			name: "ftgd-block-http&https",
			args: testArgs23,
			want: []string{"undefined", "undefined", "", "", "undefined", "undefined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleAntiVirusOutput(tt.args.allInfo); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleAntiVirusOutput()\n")
			}
		})
	}
}

func Test_HandleExpectOutput(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "service-multiple-tcp",
			args: testArgs1,
			want: []string{"pass"},
		},
		{
			name: "icmp",
			args: testArgs2,
			want: []string{"drop"},
		},
		{
			name: "ftp",
			args: testArgs3,
			want: []string{"pass", "pass"},
		},
		{
			name: "Servicegrp",
			args: testArgs4,
			want: []string{"pass", "pass", "pass", "pass"},
		},
		{
			name: "Email Access",
			args: testArgs5,
			want: []string{"drop", "drop", "drop", "drop", "drop", "drop", "drop", "drop"},
		},
		{
			name: "all",
			args: testArgs6,
			want: []string{"pass", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined", "block", "block", "block", "pass", "block", "block", "block", "pass", "pass", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined"},
		},
		{
			name: "multi service",
			args: testArgs7,
			want: []string{"pass", "undefined", "pass", "undefined", "pass", "block"},
		},
		{
			name: "vipgrp",
			args: testArgs8,
			want: []string{"pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass"},
		},
		{
			name: "vip-all",
			args: testArgs14,
			want: []string{"pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass"},
		},
		{
			name: "vip-ping",
			args: testArgs13,
			want: []string{"pass"},
		},
		{
			name: "vipgrp-smtp",
			args: testArgs12,
			want: []string{"pass", "undefined", "pass", "pass", "pass", "undefined", "undefined"},
		},
		{
			name: "vipgrp-all-udp",
			args: testArgs16,
			want: []string{"pass", "undefined", "pass", "undefined", "undefined", "undefined", "undefined", "pass"},
		},
		{
			name: "vipgrp-all",
			args: testArgs17,
			want: []string{"pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass", "pass"},
		},
		{
			name: "vipgrp-ping",
			args: testArgs18,
			want: []string{"pass", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined"},
		},
		{
			name: "vipgrp-ftp",
			args: testArgs19,
			want: []string{"pass", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-dns&ftp",
			args: testArgs20,
			want: []string{"pass", "undefined", "undefined", "undefined", "pass", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined"},
		},
		{
			name: "vip-ntp",
			args: testArgs21,
			want: []string{"pass", "pass", "undefined", "undefined", "undefined", "pass", "pass", "pass", "undefined", "pass", "undefined", "undefined", "pass", "pass"},
		},
		{
			name: "dns-tcp-port-different",
			args: testArgs22,
			want: []string{"pass", "pass", "undefined", "undefined", "undefined", "pass", "pass", "pass", "undefined", "pass", "pass", "pass", "undefined", "undefined"},
		},
		{
			name: "ftgd-block-http&https",
			args: testArgs23,
			want: []string{"undefined", "undefined", "block", "block", "undefined", "undefined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.HandleExpectOutput(tt.args.allInfo); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got = %+v\n", got)
				fmt.Printf("want = %+v\n", tt.want)
				fmt.Printf("len(got) = %+v, len(want) = %+v\n", len(got), len(tt.want))
				t.Error("HandleExpectOutput()\n")
			}
		})
	}
}
