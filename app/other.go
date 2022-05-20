package app

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	addrCfg       = "config firewall address"
	addrGrpCfg    = "config firewall addrgrp"
	serviceCfg    = "config firewall service custom"
	serviceGrpCfg = "config firewall service group"
	avProfileCfg  = "config antivirus profile"
	wfProfileCfg  = "config webfilter profile"
	wfFilterCfg   = "config webfilter urlfilter"
	ippoolCfg     = "config firewall ippool"
	vipCfg        = "config firewall vip"
	vipGrpCfg     = "config firewall vipgrp"
	sRouteCfg     = "config router static"
	policyCfg     = "config firewall policy"
	systemCfg     = "config system settings"
	nextCfg       = "next"
	endCfg        = "end"
)

type IntfInfo struct {
	Name       string
	Address    string `json:"Address,omitempty"`
	SubnetMask string `json:"SubnetMask,omitempty"`
	VLANID     int    `validate:"min=0,max=4094"`
	RandomIP   string
}

type AddressInfo struct {
	Name       string
	Address    string `json:"Address,omitempty"`
	SubnetMask string `json:"SubnetMask,omitempty"`
	StartIP    string `json:"StartIP,omitempty"`
	EndIP      string `json:"EndIP,omitempty"`
}

type AddrGrpInfo struct {
	Name   string
	Member []string // AddressInfoのNameの集合体になる想定
}

type PortAssignInfo struct {
	SrcPort string
	DstPort string
}

type ServiceInfo struct {
	Name string
	ICMP bool             `json:"ICMP,omitempty"`
	TCP  []PortAssignInfo `json:"TCP,omitempty"`
	UDP  []PortAssignInfo `json:"UDP,omitempty"`
}

type ServiceGrpInfo struct {
	Name   string
	Member []string
}

type AVProfileInfo struct {
	Name   string
	Mode   string          `json:"Mode,omitempty"`
	Config map[string]bool // trueの場合はexpectでblockの出力を行う
}

type WfProfileInfo struct {
	Name        string
	URLTableNum string `json:"URLTableNum,omitempty"`
}

type WfFilterInfo struct {
	Name    string
	Entries []Entry
}

type Entry struct {
	URL    string
	Action string `json:"Action,omitempty"`
}

type IPPoolInfo struct {
	Name          string
	StartIP       string
	EndIP         string
	SourceStartIP string `json:"SourceStartIP,omitempty"`
	SourceEndIP   string `json:"SourceEndIP,omitempty"`
}

type VipInfo struct {
	Name        string
	ExtIP       string   // dstnataddr
	MappedIP    string   // dstaddr
	SrcFilter   []string `json:"SrcFilter,omitempty"`
	Service     []string `json:"Service,omitempty"`
	PortForward bool     `json:"PortForward,omitempty"`
	ExtPort     string   `json:"ExtPort,omitempty"`
	MappedPort  string   `json:"MappedPort,omitempty"`
	Protocol    string   `json:"Protocol,omitempty"`
}

type vipGrpInfo struct {
	Name   string
	Member []string
}

type SRouteInfo struct {
	Name   string
	Dst    string `json:"Dst,omitempty"`
	GW     string
	Device string
}

type policyInfo struct {
	Name          string   `json:"Name"`
	SrcAddress    []string `json:"SrcAddress"`
	Srcintf       string   `json:"Srcintf"`
	DstAddress    []string `json:"DstAddress"`
	Dstintf       string   `json:"Dstintf"`
	Service       []string `json:"Service"`
	IPPool        bool     `json:"IPPool,omitempty"`
	PoolName      string   `json:"PoolName,omitempty"`
	AVProfile     string   `json:"AVProfile,omitempty"`
	WFProfile     string   `json:"WFProfile,omitempty"`
	NAT           bool     `json:"NAT,omitempty"`
	Action        string   `json:"Action"`
	ProxyMode     bool     `json:"ProxyMode,omitempty"`
	DisableStatus bool     `json:"DisableStatus,omitempty"`
}

var ServicePort = map[string]string{
	"icmp":  "",
	"tcp":   "80",
	"udp":   "53",
	"dns":   "53",
	"dnst":  "53",
	"http":  "80",
	"https": "443",
	"ftp":   "21",
	"ftpa":  "21",
	"imap":  "143",
	"smtp":  "25",
}

var (
	addrStatus       bool
	addrGrpStatus    bool
	serviceStatus    bool
	serviceGrpStatus bool
	avProfileStatus  bool
	wfProfileStatus  bool
	wfFilterStatus   bool
	ipPoolStatus     bool
	vipStatus        bool
	vipGrpStatus     bool
	sRouteStatus     bool
	policyStatus     bool
	systemStatus     bool
	notAppendStatus  bool

	IntfInfoAll       []IntfInfo
	AddressInfoAll    []AddressInfo
	AddrGrpInfoAll    []AddrGrpInfo
	ServiceInfoAll    []ServiceInfo
	ServiceGrpInfoAll []ServiceGrpInfo
	AVProfileInfoAll  []AVProfileInfo
	WfProfileInfoAll  []WfProfileInfo
	WfFilterInfoAll   []WfFilterInfo
	IPPoolInfoAll     []IPPoolInfo
	VipInfoAll        []VipInfo
	vipGrpInfoAll     []vipGrpInfo
	SRouteInfoAll     []SRouteInfo
	policyInfoAll     []policyInfo

	AllNWProtocol     = []string{"icmp", "tcp", "udp", "dnst", "dns", "http", "https", "ftp", "ftpa", "imap", "smtp"}
	ProxyModeProtocol = []string{"tcp", "dnst", "http", "https", "ftp", "ftpa", "imap", "smtp"}
	IsFlowVDOM        bool
	flags             bool
	entries           []Entry
	usedAddress       []string

	aInfo    AddressInfo
	aGInfo   AddrGrpInfo
	pInfo    policyInfo
	svInfo   ServiceInfo
	svGInfo  ServiceGrpInfo
	avPInfo  AVProfileInfo
	wfPInfo  WfProfileInfo
	wfFInfo  WfFilterInfo
	ipPInfo  IPPoolInfo
	vInfo    VipInfo
	vGInfo   vipGrpInfo
	sRInfo   SRouteInfo
	protocol string
	e        Entry
	version  string
)

func handleFWIPOutput(intf, env string, IntfInfoAll []IntfInfo) string {
	intfInfo := confirmIntfWithIntfInfo(intf, env, IntfInfoAll)
	return intfInfo.Address
}

func handleVLANIDOutput(intf, env string, intfInfoAll []IntfInfo) string {
	intfInfo := confirmIntfWithIntfInfo(intf, env, IntfInfoAll)
	return strconv.Itoa(intfInfo.VLANID)
}

// ポリシーで使用されているintfがintf.csvに記載されていない場合にエラーを返す
func confirmIntfWithIntfInfo(intf, env string, intfInfoAll []IntfInfo) IntfInfo {
	for _, v := range intfInfoAll {
		if intf == v.Name {
			return v
		}
	}

	fmt.Printf("%sの情報はconfigディレクトリ内の`intf.csv`に記載されていませんでした\n", intf)
	fmt.Printf("インターフェース名%sの情報をconfigディレクトリ内の`intf.csv`に記載してください\n", intf)
	if env != "test" {
		os.Exit(103)
	}
	return IntfInfo{}
}

func HandleSrcAddressOutput(intf string, addrs []string, AddressInfoAll []AddressInfo, AddrGrpInfoAll []AddrGrpInfo, SRouteInfoAll []SRouteInfo, IntfInfoAll []IntfInfo, usedAddress []string, aI AllInfo) ([]string, bool) {
	var addrSlice []string
	var srcFQDNFlag bool
	for _, addr := range addrs {
		var aSlice []string
		var addrBool bool
		var flags bool
		aSlice, addrBool, flags = handleAddress(addr, intf, AddressInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, aI)
		if !addrBool {
			aSlice, flags = handleAddressGrp(addr, intf, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, aI)
		}

		if flags {
			srcFQDNFlag = true
		}
		addrSlice = append(addrSlice, aSlice...)
	}
	return addrSlice, srcFQDNFlag
}

func HandleDstAddressOutput(intf string, addrs []string, AddressInfoAll []AddressInfo, AddrGrpInfoAll []AddrGrpInfo, SRouteInfoAll []SRouteInfo, IntfInfoAll []IntfInfo, usedAddress []string, allInfo AllInfo) ([]string, bool) {
	var dstAddresss []string
	var dstFQDNFlag bool
	protocols := HandleProtocolOutput(allInfo)

	if allInfo.VInfo != nil {
		dstAddresss = handleVInfo2("dst_addr", allInfo)
	} else {
		for _, addr := range addrs {
			var aSlice []string
			var flags bool
			aSlice, addrBool, flags := handleAddress(addr, intf, AddressInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, allInfo)
			if !addrBool {
				aSlice, flags = handleAddressGrp(addr, intf, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, allInfo)
			}

			if flags {
				dstFQDNFlag = true
			}

			for _, v := range aSlice {
				for range protocols {
					dstAddresss = append(dstAddresss, v)
				}
			}
		}
	}
	return dstAddresss, dstFQDNFlag
}

func handleAddress(addr, intf string, AddressInfoAll []AddressInfo, SRouteInfoAll []SRouteInfo, IntfInfoAll []IntfInfo, usedAddress []string, aI AllInfo) ([]string, bool, bool) {
	var addrSlice []string
	for _, v := range AddressInfoAll {
		if addr == v.Name && addr != `"all"` {
			if v.Address != "" && v.SubnetMask == "255.255.255.255" {
				getUsedIPFromPolicy([]string{v.Address})
				return []string{v.Address}, true, false
			} else if v.Address != "" && v.SubnetMask != "" {
				ips := getAllIPFromNetwork(v.Address, v.SubnetMask)
				// ネットワークアドレスとブロードキャストアドレスを除く最初と最後のIPを追加する
				getUsedIPFromPolicy(ips)
				addrSlice = append(addrSlice, ips[1], ips[len(ips)-2])
				return addrSlice, true, false
			} else if v.StartIP != "" && v.EndIP != "" {
				rangeIPs := getRangeIP(v.StartIP, v.EndIP)
				getUsedIPFromPolicy(rangeIPs)
				addrSlice = append(addrSlice, v.StartIP, v.EndIP)
				return addrSlice, true, false
			} else if v.Address != "" {
				return []string{v.Address}, true, true
			} else {
				return []string{"2.2.2.2"}, true, false
			}
		}
	}

	if addr == `"all"` {
		routes, isExistRoute := getAllRouteIntf(intf, aI.Env, SRouteInfoAll, IntfInfoAll)
		if isExistRoute {
			// static routeのあるネットワーク内の全てのIPからポリシーにあるIPを除いた始めのIP
			// intf自身が所属するネットワークアドレスもstatic routeに含む
			var ips []string
			for _, route := range routes {
				netaddr, subnet := strings.Split(route, " ")[0], strings.Split(route, " ")[1]
				addrs := getAllIPFromNetwork(netaddr, subnet)
				ips = append(ips, addrs...)
			}
			noIPs := excludeUsedIPFromNetaddr(usedAddress, ips)
			// ipsにあってusedAddressにない要素（noIPs）をappendする
			switch len(noIPs) {
			case 0:
				if aI.Env != "test" {
					fmt.Printf("予期せぬエラーが発生しました\n")
					os.Exit(254)
				}
			case 1:
				addrSlice = append(addrSlice, noIPs[0])
			default:
				addrSlice = append(addrSlice, noIPs[1])
			}
		} else {
			// static routeがない場合は特定のIPを返す
			appendIP := handleNotExistSRoute(intf, aI.Env, IntfInfoAll)
			addrSlice = append(addrSlice, appendIP)
		}
		return addrSlice, true, false
	}
	return []string{""}, false, false
}

func handleAddressGrp(addr, intf string, AddressInfoAll []AddressInfo, AddrGrpInfoAll []AddrGrpInfo, SRouteInfoAll []SRouteInfo, IntfInfoAll []IntfInfo, usedAddress []string, aI AllInfo) ([]string, bool) {
	var fqdnFlags bool
	for _, v := range AddrGrpInfoAll {
		if addr == v.Name {
			var addrSlice []string
			for _, m := range v.Member {
				aSlice, _, flags := handleAddress(m, intf, AddressInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, aI)
				addrSlice = append(addrSlice, aSlice...)
				if flags {
					fqdnFlags = true
				}
			}
			return addrSlice, fqdnFlags
		}
	}
	return []string{""}, fqdnFlags
}

func handleNotExistSRoute(intf, env string, intfInfoAll []IntfInfo) string {
	intfInfo := confirmIntfWithIntfInfo(intf, env, IntfInfoAll)
	fmt.Printf("%sにはスタティックルートが存在しません\n", intf)
	if intfInfo.RandomIP != "" {
		fmt.Printf("%sを出力します\n", intfInfo.RandomIP)
		return intfInfo.RandomIP
	}
	// テスト以外では下記は返されない想定
	return "undefined"
}

// getRangeIPと処理を統一させる
func handleIncreaseIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func getAllRouteIntf(intf, env string, SRouteInfoAll []SRouteInfo, intfInfoAll []IntfInfo) ([]string, bool) {
	var routeSlice []string
	var isExistRoute bool
	for _, v := range SRouteInfoAll {
		if intf == v.Device {
			// DGWの場合はv.Dstがない
			// 0.0.0.0/0を全てappendするのはかなり時間がかかるのでappendしない
			if v.Dst == "" {
				continue
			} else {
				isExistRoute = true
				routeSlice = append(routeSlice, v.Dst)
			}
		}
	}

	if !isExistRoute {
		intfInfo := confirmIntfWithIntfInfo(intf, env, intfInfoAll)
		route := intfInfo.Address + " " + intfInfo.SubnetMask
		routeSlice = append(routeSlice, route)
		isExistRoute = true
	}
	return routeSlice, isExistRoute
}

func getAllIPFromNetwork(addr, subnetMask string) []string {
	subnet := net.ParseIP(subnetMask)
	if subnet == nil {
		return nil
	}
	mask := net.IPv4Mask(subnet[12], subnet[13], subnet[14], subnet[15])
	length, _ := mask.Size()
	ip, ipnet, err := net.ParseCIDR(addr + "/" + strconv.Itoa(length))
	if err != nil {
		// この処理はほぼ起こり得ないが念の為ハンドリング
		fmt.Printf("failed to parse cidr %s", err.Error())
		return []string{""}
	}
	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); handleIncreaseIP(ip) {
		// ipsはipnet内の全てのアドレスが格納されている
		ips = append(ips, ip.String())
	}
	return ips
}

// 203.0.113.0/24のネットワークはFQDNテストの際に使用するが現状考慮していない
// アドレスにallを指定するポリシーが先にある可能性もあるのでそれを考慮している
func getUsedIPFromPolicy(ips []string) []string {
	ipm := make(map[string]struct{})
	for _, ip := range ips {
		if _, ok := ipm[ip]; !ok {
			ipm[ip] = struct{}{}
			usedAddress = append(usedAddress, ip)
		}
	}
	return usedAddress
}

func getRangeIP(start, end string) []string {
	var ips []string
	startIP := net.ParseIP(start).To4()
	endIP := net.ParseIP(end).To4()

	for !startIP.Equal(endIP) {
		cpIP := make(net.IP, len(startIP))
		copy(cpIP, startIP)
		ips = append(ips, cpIP.String())
		// 通常ネットワークアドレスを用いるはずなので、/16までの対応としている
		// （例: 192.168.1.1-192.168.255.255）
		// /16以上のレンジを指定するケースがあれば要望次第で対応する
		if startIP[3] == 255 {
			startIP[2] = startIP[2] + 1
		}
		startIP[3] = startIP[3] + 1
	}
	ips = append(ips, startIP.String())
	return ips
}

func excludeUsedIPFromNetaddr(usedIPs, netaddrs []string) []string {
	var notUsedIPSlice []string
	for _, addr := range netaddrs {
		var duplicate bool
		for _, usedIP := range usedIPs {
			if addr == usedIP {
				duplicate = true
			}
		}
		if !duplicate {
			notUsedIPSlice = append(notUsedIPSlice, addr)
		}
	}
	return notUsedIPSlice
}

// IPPoolが複数割り当てられるパターンは一旦想定しない
func HandleSrcNATAddressOutPut(natbool bool, poolname, dstintf, env string, ippools []IPPoolInfo, IntfInfoAll []IntfInfo) (string, string) {
	if natbool {
		if poolname != "" {
			for _, ipPool := range ippools {
				if poolname == ipPool.Name {
					// ポート固定範囲でinternalIPが送信元IPと一致していない場合はdropする
					// 様子を見て「undefined」にする
					return ipPool.StartIP, ipPool.StartIP + "-" + ipPool.EndIP
				}
			}
		} else {
			poolIP := handleFWIPOutput(dstintf, env, IntfInfoAll)
			return poolIP, ""
		}
	}
	return "", ""
}

func existNWProtocol(serv string) bool {
	switch serv {
	case `"DNS"`, `"HTTP"`, `"HTTPS"`, `"FTP"`, `"IMAP"`, `"SMTP"`:
		return true
	}
	return false
}

func getVipInfo(addrs []string, vipGrpInfoAll []vipGrpInfo, VipInfoAll []VipInfo) []VipInfo {
	var vSlice []VipInfo
	for _, addr := range addrs {
		for _, v := range vipGrpInfoAll {
			if addr == v.Name {
				for _, m := range v.Member {
					for _, vv := range VipInfoAll {
						if m == vv.Name {
							vSlice = append(vSlice, vv)
						}
					}
				}
			}
		}

		for _, v := range VipInfoAll {
			if addr == v.Name {
				vSlice = append(vSlice, v)
			}
		}
	}
	return vSlice
}

func HandleDescriptionOutput(name, ipPoolRange string, srcFQDNFlag, dstFQDNFlag bool) string {
	s := strings.Replace(name, `"`, "", -1)
	switch {
	case ipPoolRange != "" && (srcFQDNFlag || dstFQDNFlag):
		fmt.Printf("ポリシー名=%sの送信元NAT IPアドレスはダイナミックIPプールなため、\n生成したシナリオとテスト結果が異なる可能性があります\n", s)
		fmt.Print("本シナリオとテスト結果で送信元NAT IPアドレスが異なる場合は値を変更後、再度テストを実施してください\n")
		fmt.Printf("ポリシー名=%sは送信元、または宛先IPアドレスにFQDNが指定されています\n", s)
		fmt.Print("本シナリオでテストを実行する場合は、NEEDLEWORKのマニュアルの「送信元・宛先アドレスにFQDNを指定したシナリオでテストを実施する」をご確認ください\n")
		s = "policy name = " + s + " src_nat_ip=" + ipPoolRange + " FQDN Policy"
		return s
	case ipPoolRange != "":
		fmt.Printf("ポリシー名=%sの送信元NAT IPアドレスはダイナミックIPプールなため、\n生成したシナリオとテスト結果が異なる可能性があります\n", s)
		fmt.Print("本シナリオとテスト結果で送信元NAT IPアドレスが異なる場合は値を変更後、再度テストを実施してください\n")
		s = "policy name = " + s + " src_nat_ip=" + ipPoolRange
		return s
	case srcFQDNFlag || dstFQDNFlag:
		fmt.Printf("ポリシー名=%sは送信元、または宛先IPアドレスにFQDNが指定されています\n", s)
		fmt.Print("本シナリオでテストを実行する場合は、NEEDLEWORKのマニュアルの「送信元・宛先アドレスにFQDNを指定したシナリオでテストを実施する」をご確認ください\n")
		s = "policy name = " + s + " FQDN Policy"
		return s
	}
	return "policy name = " + s
}

func AbsorbText(text []string, m map[string]bool) {
	for _, v := range text {
		switch v {
		case systemCfg:
			IsFlowVDOM = false // VDOMで異なる場合を考慮して初期化
			systemStatus = true
		case addrCfg:
			addrStatus = true
		case addrGrpCfg:
			addrGrpStatus = true
		case policyCfg:
			policyStatus = true
		case serviceCfg:
			serviceStatus = true
		case serviceGrpCfg:
			serviceGrpStatus = true
		case avProfileCfg:
			avProfileStatus = true
		case wfProfileCfg:
			wfProfileStatus = true
		case wfFilterCfg:
			wfFilterStatus = true
		case ippoolCfg:
			ipPoolStatus = true
		case vipCfg:
			vipStatus = true
		case vipGrpCfg:
			vipGrpStatus = true
		case sRouteCfg:
			sRouteStatus = true
		case endCfg:
			addrStatus = false
			addrGrpStatus = false
			serviceStatus = false
			serviceGrpStatus = false
			ipPoolStatus = false
			vipStatus = false
			vipGrpStatus = false
			sRouteStatus = false
			systemStatus = false
		}

		if strings.Contains(v, "#config-version=") {
			version = strings.Split(v, "-")[2]
		}

		if systemStatus {
			if strings.Contains(v, "set inspection-mode ") {
				// flowModeの場合のみ該当する
				IsFlowVDOM = true
			}
		}

		if addrStatus {
			if strings.Contains(v, "edit ") {
				aInfo.Name = strings.Split(v, "edit ")[1]
			}
			// デフォルトであるアドレスのall, またはfqdnなどはこの行がない
			// 現在はアドレス出力時に値を補完するようにハンドリングしている
			if strings.Contains(v, "set subnet ") {
				str := strings.Split(v, "set subnet ")[1]
				aInfo.Address, aInfo.SubnetMask = strings.Split(str, " ")[0], strings.Split(str, " ")[1]
			}
			if strings.Contains(v, "set fqdn ") {
				fqdn := strings.Split(v, "set fqdn ")[1]
				str := strings.Split(fqdn, `"`)[1]
				if strings.Contains(str, "*.") {
					aInfo.Address, aInfo.SubnetMask = strings.Split(str, "*.")[1], ""
				} else {
					aInfo.Address, aInfo.SubnetMask = str, ""
				}
			}
			if strings.Contains(v, "set start-ip ") {
				aInfo.StartIP = strings.Split(v, "set start-ip ")[1]
			}
			if strings.Contains(v, "set end-ip ") {
				aInfo.EndIP = strings.Split(v, "set end-ip ")[1]
			}
			if strings.Contains(v, nextCfg) {
				AddressInfoAll = append(AddressInfoAll, aInfo)
				aInfo = AddressInfo{}
			}
		}

		if addrGrpStatus {
			if strings.Contains(v, "edit") {
				aGInfo.Name = strings.Split(v, "edit ")[1]
			}
			if strings.Contains(v, "set member") {
				aGInfo.Member = strings.Split(strings.Split(v, "set member ")[1], " ")
			}
			if strings.Contains(v, nextCfg) {
				AddrGrpInfoAll = append(AddrGrpInfoAll, aGInfo)
				aGInfo = AddrGrpInfo{}
			}
		}

		// SCTPはNEEDLEWORKが対応していないので非対応
		if serviceStatus {
			if strings.Contains(v, "edit ") {
				svInfo.Name = strings.Split(v, "edit ")[1]
			}
			if strings.Contains(v, "set protocol ICMP") {
				svInfo.ICMP = true
			}
			// パターン例
			// set tcp-portrange 1755
			// set udp-portrange 1024-5000
			// set tcp-portrange 554 7070 8554
			// set tcp-portrange 80:60000
			// set tcp-portrange 0-65535:0-65535
			// set tcp-portrange 11111:22220-22222
			// set udp-portrange 33330-33333:4444
			// set tcp-portrange 11122:123 11155:123
			var tpaInfo PortAssignInfo
			var upaInfo PortAssignInfo
			if strings.Contains(v, "set tcp-portrange ") {
				tcpPort := strings.Split(v, "set tcp-portrange ")[1]
				if strings.Contains(tcpPort, " ") {
					p := strings.Split(tcpPort, " ")
					for _, v := range p {
						// vに:がある場合はsrc & dst
						if strings.Contains(v, ":") {
							tpaInfo.SrcPort = strings.Split(v, ":")[1]
							tpaInfo.DstPort = strings.Split(v, ":")[0]
						} else {
							tpaInfo.DstPort = v
						}
						svInfo.TCP = append(svInfo.TCP, tpaInfo)
						tpaInfo = PortAssignInfo{}
						continue
					}
				} else if strings.Contains(tcpPort, ":") {
					// 区切った後に:がある場合はsrc & dst
					tpaInfo.SrcPort = strings.Split(tcpPort, ":")[1]
					tpaInfo.DstPort = strings.Split(tcpPort, ":")[0]
					svInfo.TCP = append(svInfo.TCP, tpaInfo)
				} else {
					tpaInfo.DstPort = tcpPort
					svInfo.TCP = append(svInfo.TCP, tpaInfo)
				}
			}
			if strings.Contains(v, "set udp-portrange ") {
				udpPort := strings.Split(v, "set udp-portrange ")[1]
				switch {
				// 区切った後にスペースがあるか
				case strings.Contains(udpPort, " "):
					p := strings.Split(udpPort, " ")
					for _, v := range p {
						// vに:がある場合はsrc & dst
						if strings.Contains(v, ":") {
							upaInfo.SrcPort = strings.Split(v, ":")[1]
							upaInfo.DstPort = strings.Split(v, ":")[0]
						} else {
							upaInfo.DstPort = v
						}
						svInfo.UDP = append(svInfo.UDP, upaInfo)
						upaInfo = PortAssignInfo{}
						continue
					}
				// 区切った後に:がある場合はsrc & dst
				case strings.Contains(udpPort, ":"):
					upaInfo.SrcPort = strings.Split(udpPort, ":")[1]
					upaInfo.DstPort = strings.Split(udpPort, ":")[0]
					svInfo.UDP = append(svInfo.UDP, upaInfo)
				default:
					upaInfo.DstPort = udpPort
					svInfo.UDP = append(svInfo.UDP, upaInfo)
				}
			}
			if strings.Contains(v, nextCfg) {
				ServiceInfoAll = append(ServiceInfoAll, svInfo)
				svInfo = ServiceInfo{}
			}
		}

		if serviceGrpStatus {
			if strings.Contains(v, "edit ") {
				svGInfo.Name = strings.Split(v, "edit ")[1]
			}
			if strings.Contains(v, "set member ") {
				svGInfo.Member = strings.Split(strings.Split(v, "set member ")[1], " ")
			}
			if strings.Contains(v, nextCfg) {
				ServiceGrpInfoAll = append(ServiceGrpInfoAll, svGInfo)
				svGInfo = ServiceGrpInfo{}
			}
		}

		if avProfileStatus {
			if strings.Contains(v, "edit ") {
				avPInfo.Name = strings.Split(v, "edit ")[1]
			}
			if strings.Contains(v, "set inspection-mode ") {
				avPInfo.Mode = strings.Split(v, "set inspection-mode ")[1]
			}
			if strings.Contains(v, "config ") {
				protocol = strings.Split(v, "config ")[1]
			}
			// インスペクションされるプロトコルで有効にしないとこの記述はない
			if strings.Contains(v, "set options scan") {
				if strings.Split(v, "set options scan")[1] != "" {
					// set options scan avmonitor
					m[protocol] = false
				} else {
					// set options scan← blockの場合
					m[protocol] = true
				}
			}
			if strings.Contains(v, endCfg) {
				avPInfo.Config = m
			}
			if strings.Contains(v, nextCfg) {
				AVProfileInfoAll = append(AVProfileInfoAll, avPInfo)
				avPInfo = AVProfileInfo{}
				m = map[string]bool{}
			}
			// config http ~ end　config imap ~ endの形になるので
			// 最初のswitch文のendcfgで引っかからないようにしている
			if v == endCfg {
				avProfileStatus = false
			}
		}

		if wfProfileStatus {
			if strings.Contains(v, "edit ") && !flags {
				wfPInfo.Name = strings.Split(v, "edit ")[1]
			}
			// ない場合もある
			if strings.Contains(v, "set urlfilter-table ") {
				wfPInfo.URLTableNum = strings.Split(v, "set urlfilter-table ")[1]
			}
			if strings.Contains(v, "config ftgd-wf") {
				flags = true
			}
			// nextCfgではない
			if v == "    next" {
				WfProfileInfoAll = append(WfProfileInfoAll, wfPInfo)
				wfPInfo = WfProfileInfo{}
				flags = false
			}
			if v == endCfg {
				wfProfileStatus = false
			}
		}

		if wfFilterStatus {
			if flags {
				if strings.Contains(v, "set url ") {
					url := strings.Split(v, "set url ")[1]
					e.URL = strings.Replace(url, `"`, "", -1)
				}
				if strings.Contains(v, "set action ") {
					e.Action = strings.Split(v, "set action ")[1]
				}
				if strings.Contains(v, nextCfg) {
					entries = append(entries, e)
					e = Entry{}
				}
				if strings.Contains(v, endCfg) {
					flags = false
				}
			} else {
				if strings.Contains(v, "edit ") {
					wfFInfo.Name = strings.Split(v, "edit ")[1]
				}
				if strings.Contains(v, "config entries") {
					flags = true
				}
				if strings.Contains(v, nextCfg) {
					wfFInfo.Entries = entries
					WfFilterInfoAll = append(WfFilterInfoAll, wfFInfo)
					wfFInfo = WfFilterInfo{}
				}
				if v == endCfg {
					wfFilterStatus = false
				}
			}
		}

		if ipPoolStatus {
			if strings.Contains(v, "edit ") {
				ipPInfo.Name = strings.Split(v, "edit ")[1]
			}
			if strings.Contains(v, "set startip ") {
				ipPInfo.StartIP = strings.Split(v, "set startip ")[1]
			}
			if strings.Contains(v, "set endip ") {
				ipPInfo.EndIP = strings.Split(v, "set endip ")[1]
			}
			if strings.Contains(v, "set source-startip ") {
				ipPInfo.SourceStartIP = strings.Split(v, "set source-startip ")[1]
			}
			if strings.Contains(v, "set source-endip ") {
				ipPInfo.SourceEndIP = strings.Split(v, "set source-endip ")[1]
			}
			if strings.Contains(v, nextCfg) {
				IPPoolInfoAll = append(IPPoolInfoAll, ipPInfo)
				ipPInfo = IPPoolInfo{}
			}
		}

		if vipStatus {
			if strings.Contains(v, "edit ") {
				vInfo.Name = strings.Split(v, "edit ")[1]
			}
			if strings.Contains(v, "set extip ") {
				vInfo.ExtIP = strings.Split(v, "set extip ")[1]
			}
			if strings.Contains(v, "set mappedip ") {
				mappedIP := strings.Split(v, "set mappedip ")[1]
				vInfo.MappedIP = strings.Replace(mappedIP, `"`, "", -1)
			}
			// 以下のnextcfg以外はない場合もある
			if strings.Contains(v, "set src-filter ") {
				srcFilter := strings.Split(v, "set src-filter ")[1]
				vInfo.SrcFilter = strings.Split(srcFilter, " ")
			}
			if strings.Contains(v, "set service ") {
				service := strings.Split(v, "set service ")[1]
				vInfo.Service = strings.Split(service, " ")
			}
			if strings.Contains(v, "set portforward ") {
				vInfo.PortForward = true
			}
			if strings.Contains(v, "set extport ") {
				vInfo.ExtPort = strings.Split(v, "set extport ")[1]
			}
			if strings.Contains(v, "set mappedport ") {
				vInfo.MappedPort = strings.Split(v, "set mappedport ")[1]
			}
			if strings.Contains(v, "set protocol ") {
				str := strings.Split(v, "set protocol ")[1]
				switch str {
				case "sctp":
					vInfo.Protocol = "NaN"
				default:
					vInfo.Protocol = strings.Split(v, "set protocol ")[1]
				}
			}
			if strings.Contains(v, nextCfg) {
				VipInfoAll = append(VipInfoAll, vInfo)
				vInfo = VipInfo{}
			}
		}

		if vipGrpStatus {
			if strings.Contains(v, "edit ") {
				vGInfo.Name = strings.Split(v, "edit ")[1]
			}
			if strings.Contains(v, "set member ") {
				vGInfo.Member = strings.Split(strings.Split(v, "set member ")[1], " ")
			}
			if strings.Contains(v, nextCfg) {
				vipGrpInfoAll = append(vipGrpInfoAll, vGInfo)
				vGInfo = vipGrpInfo{}
			}
		}

		if sRouteStatus {
			if strings.Contains(v, "edit ") {
				sRInfo.Name = strings.Split(v, "edit ")[1]
			}
			if strings.Contains(v, "set dst ") {
				sRInfo.Dst = strings.Split(v, "set dst ")[1]
			}
			if strings.Contains(v, "set gateway ") {
				sRInfo.GW = strings.Split(v, "set gateway ")[1]
			}
			if strings.Contains(v, "set device ") {
				sRInfo.Device = strings.Split(v, "set device ")[1]
			}
			if strings.Contains(v, nextCfg) {
				SRouteInfoAll = append(SRouteInfoAll, sRInfo)
				sRInfo = SRouteInfo{}
			}
		}

		if policyStatus {
			if strings.Contains(v, "set name") {
				pInfo.Name = strings.Split(v, "set name ")[1]
			}
			if strings.Contains(v, "set srcintf ") {
				pInfo.Srcintf = strings.Split(v, "set srcintf ")[1]
			}
			if strings.Contains(v, "set srcaddr ") {
				src := strings.Split(v, "set srcaddr ")[1]
				// スペースがある場合は複数のアドレスが指定されているので全てスライスに追加する
				var srcSlice []string
				if strings.Contains(src, `" "`) {
					srcSlice = strings.Split(src, " ")
				} else {
					srcSlice = []string{src}
				}
				pInfo.SrcAddress = srcSlice
			}
			if strings.Contains(v, "set dstintf ") {
				pInfo.Dstintf = strings.Split(v, "set dstintf ")[1]
			}
			if strings.Contains(v, "set dstaddr ") {
				dst := strings.Split(v, "set dstaddr ")[1]
				// スペースがある場合は複数のアドレスが指定されているので全てスライスに追加する
				var dstSlice []string
				if strings.Contains(dst, `" "`) {
					dstSlice = strings.Split(dst, " ")
				} else {
					dstSlice = []string{dst}
				}
				pInfo.DstAddress = dstSlice
			}
			if strings.Contains(v, "set service ") {
				servs := strings.Split(v, "set service ")[1]
				var servSlice []string
				if strings.Contains(servs, `" "`) {
					servSlice = strings.Split(servs, " ")
				} else {
					servSlice = []string{servs}
				}
				pInfo.Service = servSlice
			}
			if strings.Contains(v, "set inspection-mode ") {
				if strings.Split(v, "set inspection-mode ")[1] == "proxy" {
					pInfo.ProxyMode = true
				}
			} else {
				// FGのversionが6.0系(例:6.0.6)の場合、ポリシー＆VDOMの設定で判定する
				// FGのversionが6.2系(例:6.2.6)の場合、ポリシーの設定のみで判定する
				minorVersion := strings.Split(version, ".")[1]
				if !IsFlowVDOM && minorVersion == "0" {
					pInfo.ProxyMode = true
				}
			}
			if strings.Contains(v, "set action ") {
				pInfo.Action = strings.Split(v, "set action ")[1]
			}
			if strings.Contains(v, "set status ") {
				pInfo.DisableStatus = true
			}
			if strings.Contains(v, "set ippool ") {
				pInfo.IPPool = true
			}
			if strings.Contains(v, "set poolname ") {
				pInfo.PoolName = strings.Split(v, "set poolname ")[1]
			}
			if strings.Contains(v, "set nat ") {
				pInfo.NAT = true
			}
			if strings.Contains(v, "set av-profile ") {
				pInfo.AVProfile = strings.Split(v, "set av-profile ")[1]
			}
			if strings.Contains(v, "set webfilter-profile ") {
				pInfo.WFProfile = strings.Split(v, "set webfilter-profile ")[1]
			}
			// TODO:default以外のプロトコルオプションに対応する
			// defaultの場合はこれがない
			if strings.Contains(v, "set profile-protocol-options ") {
				notAppendStatus = true
			}
			if strings.Contains(v, nextCfg) {
				if !notAppendStatus {
					policyInfoAll = append(policyInfoAll, pInfo)
					pInfo = policyInfo{}
				} else {
					fmt.Printf("%sは`default`以外のプロトコルオプションを使用しているため出力をスキップします\n", pInfo.Name)
				}
			}

			if v == endCfg {
				for _, p := range policyInfoAll {
					// disableポリシーは初期は出力しない
					if !p.DisableStatus {
						fmt.Printf("ポリシー名 %+v のシナリオを生成しています\n", p.Name)
						s := genScenario(p, ServiceInfoAll, IPPoolInfoAll, ServiceGrpInfoAll)
						Scenarios = append(Scenarios, s)
						fmt.Printf("ポリシー名 %+v のシナリオを生成しました\n\n", p.Name)
					}
				}

				// サービスの重複等が考えられるためpolicyStatus毎にリセットする
				AddressInfoAll = []AddressInfo{}
				AddrGrpInfoAll = []AddrGrpInfo{}
				ServiceInfoAll = []ServiceInfo{}
				ServiceGrpInfoAll = []ServiceGrpInfo{}
				AVProfileInfoAll = []AVProfileInfo{}
				WfProfileInfoAll = []WfProfileInfo{}
				WfFilterInfoAll = []WfFilterInfo{}
				IPPoolInfoAll = []IPPoolInfo{}
				VipInfoAll = []VipInfo{}
				vipGrpInfoAll = []vipGrpInfo{}
				SRouteInfoAll = []SRouteInfo{}
				policyInfoAll = []policyInfo{}

				policyStatus = false
			}
		}
	}
	fmt.Print("FortiGateのコンフィグファイル(.conf)の読み込みが完了しました\n")
}

func genScenario(p policyInfo, sI []ServiceInfo, iI []IPPoolInfo, ServiceGrpInfoAll []ServiceGrpInfo) Scenario {
	vInfo := getVipInfo(p.DstAddress, vipGrpInfoAll, VipInfoAll)
	allInfo := AllInfo{
		Services:          p.Service,
		Act:               p.Action,
		WfpName:           p.WFProfile,
		AvpName:           p.AVProfile,
		ServiceInfoAll:    ServiceInfoAll,
		ServiceGrpInfoAll: ServiceGrpInfoAll,
		ProxyMode:         p.ProxyMode,
		VInfo:             vInfo,
		VipGrpInfo:        vipGrpInfoAll,
		WfProfileInfoAll:  WfProfileInfoAll,
		WfFilterInfoAll:   WfFilterInfoAll,
		AVProfileInfoAll:  AVProfileInfoAll,
		ProxyModeProtocol: ProxyModeProtocol,
		ServicePort:       ServicePort,
	}

	srcAddress, srcFQDNFlag := HandleSrcAddressOutput(p.Srcintf, p.SrcAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, allInfo)
	rS, dstFQDNFlag := handleRelatedServiceOutPut(p, sI, ServiceGrpInfoAll, IntfInfoAll, allInfo)
	srcNATAddress, ipPoolRange := HandleSrcNATAddressOutPut(p.NAT, p.PoolName, p.Dstintf, "", iI, IntfInfoAll)

	s := Scenario{
		SrcFW:      handleFWIPOutput(p.Srcintf, "", IntfInfoAll),
		SrcVLAN:    handleVLANIDOutput(p.Srcintf, "", IntfInfoAll),
		SrcAddress: srcAddress,
		SrcNATAddress:    srcNATAddress,
		SrcIntf:          "",
		ReceiverPhysical: "",
		DstFW:            handleFWIPOutput(p.Dstintf, "", IntfInfoAll),
		DstVLAN:          handleVLANIDOutput(p.Dstintf, "", IntfInfoAll),
		DstIntf:          "",
		Timeout:          "",
		Try:              "",
		Description:      HandleDescriptionOutput(p.Name, ipPoolRange, srcFQDNFlag, dstFQDNFlag),
		RelatedService:   rS,
	}
	return s
}
