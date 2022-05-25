package app

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
)

type (
	IntfInfo struct {
		Name       string
		Address    string `json:"Address,omitempty"`
		SubnetMask string `json:"SubnetMask,omitempty"`
		VLANID     int    `validate:"min=0,max=4094"`
		RandomIP   string
	}
	AddressInfo struct {
		Name       string
		Address    string `json:"Address,omitempty"`
		SubnetMask string `json:"SubnetMask,omitempty"`
		StartIP    string `json:"StartIP,omitempty"`
		EndIP      string `json:"EndIP,omitempty"`
	}
	AddrGrpInfo struct {
		Name   string
		Member []string // AddressInfoのNameの集合体になる想定
	}
	PortAssignInfo struct {
		SrcPort string
		DstPort string
	}
	ServiceInfo struct {
		Name string
		ICMP bool             `json:"ICMP,omitempty"`
		TCP  []PortAssignInfo `json:"TCP,omitempty"`
		UDP  []PortAssignInfo `json:"UDP,omitempty"`
	}
	ServiceGrpInfo struct {
		Name   string
		Member []string
	}
	AVProfileInfo struct {
		Name   string
		Mode   string          `json:"Mode,omitempty"`
		Config map[string]bool // trueの場合はexpectでblockの出力を行う
	}
	WfProfileInfo struct {
		Name        string
		URLTableNum string `json:"URLTableNum,omitempty"`
		FtgdFilter  bool
	}
	WfFilterInfo struct {
		Name    string
		Entries []Entry
	}
	Entry struct {
		URL    string
		Action string `json:"Action,omitempty"`
	}
	IPPoolInfo struct {
		Name          string
		StartIP       string
		EndIP         string
		SourceStartIP string `json:"SourceStartIP,omitempty"`
		SourceEndIP   string `json:"SourceEndIP,omitempty"`
	}
	VipInfo struct {
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
	vipGrpInfo struct {
		Name   string
		Member []string
	}
	SRouteInfo struct {
		Name   string
		Dst    string `json:"Dst,omitempty"`
		GW     string
		Device string
	}
	policyInfo struct {
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
)

var (
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
)

func getUniqueDstAddr(intf string, addrs []string, usedAddress []string, allInfo AllInfo) ([]string, bool) {
	var dstUniqueAddress []string
	var dstFQDNFlag bool
	for _, addr := range addrs {
		var aSlice []string
		var flags bool
		aSlice, addrBool, flags := handleAddress(addr, intf, usedAddress, allInfo)
		if !addrBool {
			aSlice, flags = handleAddressGrp(addr, intf, usedAddress, allInfo)
		}

		// addrのいずれかでflagsがtrueの場合はFQDNFlagをtrueにする
		if flags {
			dstFQDNFlag = true
		}
		dstUniqueAddress = append(dstUniqueAddress, aSlice...)
	}
	return dstUniqueAddress, dstFQDNFlag
}

// func handleAddress(addr, intf string, AddressInfoAll []AddressInfo, SRouteInfoAll []SRouteInfo, IntfInfoAll []IntfInfo, usedAddress []string, aI AllInfo) ([]string, bool, bool) {
func handleAddress(addr, intf string, usedAddress []string, aI AllInfo) ([]string, bool, bool) {
	var addrSlice []string
	for _, v := range aI.AddressInfoAll {
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
		routes, isExistRoute := getAllRouteIntf(intf, addr, aI)
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
			appendIP := handleNotExistSRoute(intf, addr, aI)
			addrSlice = append(addrSlice, appendIP)
		}
		return addrSlice, true, false
	}
	return []string{""}, false, false
}

func handleAddressGrp(addr, intf string, usedAddress []string, aI AllInfo) ([]string, bool) {
	var fqdnFlags bool
	addrGrpInfo := getAddrGrpInfo(addr, aI.AddrGrpInfoAll)
	var addrSlice []string
	for _, m := range addrGrpInfo.Member {
		subGrpInfo := getAddrGrpInfo(m, aI.AddrGrpInfoAll)
		if !reflect.DeepEqual(subGrpInfo, AddrGrpInfo{}) {
			for _, sm := range subGrpInfo.Member {
				aSlice, _, flags := handleAddress(sm, intf, usedAddress, aI)
				addrSlice = append(addrSlice, aSlice...)
				if flags {
					fqdnFlags = true
				}
			}
		} else {
			aSlice, _, flags := handleAddress(m, intf, usedAddress, aI)
			addrSlice = append(addrSlice, aSlice...)
			if flags {
				fqdnFlags = true
			}
		}
	}
	return addrSlice, fqdnFlags
}

func handleNotExistSRoute(intf, addr string, aI AllInfo) string {
	intfInfo := confirmIntfWithIntfInfo(intf, addr, aI)
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

func getAllRouteIntf(intf, addr string, aI AllInfo) ([]string, bool) {
	var routeSlice []string
	var isExistRoute bool
	for _, v := range aI.SRouteInfoAll {
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
		intfInfo := confirmIntfWithIntfInfo(intf, addr, aI)
		route := intfInfo.Address + " " + intfInfo.SubnetMask
		routeSlice = append(routeSlice, route)
		isExistRoute = true
	}
	return routeSlice, isExistRoute
}

func getAllIPFromNetwork(addr, subnetMask string) []string {
	ip, ipnet, err := convertSubnetMaskToCIDR(addr, subnetMask)
	if err != nil {
		// この処理はほぼ起こり得ないが念の為ハンドリング
		fmt.Printf("failed to parse cidr %s\n", err.Error())
	}
	var ips []string
	if ip != nil && ipnet != nil {
		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); handleIncreaseIP(ip) {
			// ipsはipnet内の全てのアドレスが格納されている
			ips = append(ips, ip.String())
		}	
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

func judgeURLDomainOutput(wfpName string, nwProtocol []string, WfProfileInfoAll []WfProfileInfo, WfFilterInfoAll []WfFilterInfo) []string {
	var urlSlice []string
	// 一致するプロファイル内のURL等を抜き出す
	urlEntry := getURLDomainForWebFilter(wfpName, WfProfileInfoAll, WfFilterInfoAll)

	// serviceにurlDomainProtocolで該当するものがあれば出力する
	if urlEntry != nil {
		for _, p := range nwProtocol {
			switch p {
			case "http", "https":
				urlSlice = append(urlSlice, urlEntry...)
			default:
				for i := range urlEntry {
					if i == 0 {
						urlSlice = append(urlSlice, "")
					} else {
						urlSlice = append(urlSlice, "undefined")
					}
				}
			}
		}
	} else {
		for i := 0; i < len(nwProtocol); i++ {
			urlSlice = append(urlSlice, "")
		}
	}
	return urlSlice
}

func getURLDomainForWebFilter(wfpName string, WfProfileInfoAll []WfProfileInfo, WfFilterInfoAll []WfFilterInfo) []string {
	var urlEntry []string
	for _, v := range WfProfileInfoAll {
		if wfpName == v.Name {
			for _, wf := range WfFilterInfoAll {
				if v.URLTableNum == wf.Name {
					for _, val := range wf.Entries {
						urlEntry = append(urlEntry, val.URL)
					}
				}
			}
		}
	}
	return urlEntry
}

func judgeAntiVirusOutput(avpName string, nwProtocol []string, AVProfileInfoAll []AVProfileInfo) []string {
	var avSlice []string
	avProtocols := getProtocolForAVProfile(avpName, AVProfileInfoAll)

	if avProtocols != nil {
		for _, p := range nwProtocol {
			var isAVProtocol bool
			for _, pp := range avProtocols {
				if p == pp {
					isAVProtocol = true
				}
			}
			if isAVProtocol {
				avSlice = append(avSlice, "enable")
			} else {
				avSlice = append(avSlice, "")
			}
		}
	} else {
		for i := 0; i < len(nwProtocol); i++ {
			avSlice = append(avSlice, "")
		}
	}
	return avSlice
}

func getProtocolForAVProfile(avpName string, AVProfileInfoAll []AVProfileInfo) []string {
	var avProtocols []string
	for _, v := range AVProfileInfoAll {
		if avpName == v.Name {
			// avcheckするプロトコルを返してもらう
			for proto, av := range v.Config {
				if av {
					switch proto {
					case "http":
						protos := []string{"http", "https"}
						avProtocols = append(avProtocols, protos...)
					case "ftp":
						protos := []string{"ftp", "ftpa"}
						avProtocols = append(avProtocols, protos...)
					default:
						avProtocols = append(avProtocols, proto)
					}
				}
			}
		}
	}
	return avProtocols
}

func judgeOtherSettingOutput(proxyMode bool, nwProtocol, ProxyModeProtocol []string, avpName string, AVProfileInfoAll []AVProfileInfo, allInfo AllInfo) []string {
	var otherSettings []string
	// ポリシーでアンチウイルスが有効か確認する
	var isEnabledAV bool
	for _, v := range AVProfileInfoAll {
		if avpName == v.Name {
			isEnabledAV = true
		}
	}

	if proxyMode {
		for _, p := range nwProtocol {
			var isProxyProtocol bool
			for _, pp := range ProxyModeProtocol {
				if p == pp {
					isProxyProtocol = true
				}
			}
			if isProxyProtocol && isEnabledAV {
				otherSettings = append(otherSettings, "Proxy mode")
			} else {
				otherSettings = append(otherSettings, "")
			}
		}
	} else {
		for i := 0; i < len(nwProtocol); i++ {
			otherSettings = append(otherSettings, "")
		}
	}
	return otherSettings
}

func judgeExpectOutput(act, wfpName, avpName string, nwProtocol []string, WfProfileInfoAll []WfProfileInfo, WfFilterInfoAll []WfFilterInfo) []string {
	var actSlice []string
	// 一致するプロファイル内のActionを抜き出す
	actions := getActionForWebFilter(act, wfpName, WfProfileInfoAll, WfFilterInfoAll)
	avProtocols := getProtocolForAVProfile(avpName, AVProfileInfoAll)

	// FtgdFilterが有効&ライセンスが無効ならblock
	// webfilter or antivirus のどちらかがblockならblock
	if actions != nil {
		for _, p := range nwProtocol {
			var isURLDomainProtocol bool
			for _, pp := range urlDomainProtocol {
				if p == pp {
					isURLDomainProtocol = true
					switch p {
					case "dns", "dnst":
						// この場合、NEEDLEWORKの判定はpass or dropなので書き換える
						for _, v := range actions {
							if v == "block" {
								actSlice = append(actSlice, "drop")
							} else {
								actSlice = append(actSlice, v)
							}
						}
					case "http", "https":
						var isAVProtocol bool
						for _, avp := range avProtocols {
							if p == avp {
								isAVProtocol = true
								for range actions {
									actSlice = append(actSlice, "block")
								}
							}
						}

						if !isAVProtocol {
							wfPInfo := getWfProfileInfo(wfpName, WfProfileInfoAll)
							if wfPInfo.FtgdFilter {
								for range actions {
									actSlice = append(actSlice, "block")
								}
							} else {
								actSlice = append(actSlice, actions...)
							}
						}
					}
				}
			}
			if !isURLDomainProtocol {
				var isAVProtocol bool
				for i := range actions {
					if i == 0 {
						for _, pp := range avProtocols {
							if p == pp {
								isAVProtocol = true
							}
						}

						if isAVProtocol {
							actSlice = append(actSlice, "block")
						} else {
							expect := judgeExpectFromAct(act)
							actSlice = append(actSlice, expect)
						}
					} else {
						actSlice = append(actSlice, "undefined")
					}
				}
			}
		}
	} else {
		for _, v := range nwProtocol {
			var isAVProtocol bool
			for _, pp3 := range avProtocols {
				if v == pp3 {
					isAVProtocol = true
					actSlice = append(actSlice, "block")
				}
			}

			if !isAVProtocol {
				wfPInfo := getWfProfileInfo(wfpName, WfProfileInfoAll)
				if (v == "http" || v == "https") && wfPInfo.FtgdFilter {
					actSlice = append(actSlice, "block")
				} else {
					actSlice = append(actSlice, judgeExpectFromAct(act))
				}
			}
		}
	}
	return actSlice
}

func getActionForWebFilter(act, wfpName string, WfProfileInfoAll []WfProfileInfo, WfFilterInfoAll []WfFilterInfo) []string {
	var actions []string
	for _, v := range WfProfileInfoAll {
		if wfpName == v.Name {
			for _, wfi := range WfFilterInfoAll {
				if v.URLTableNum == wfi.Name {
					for _, e := range wfi.Entries {
						if e.Action == "monitor" {
							actions = append(actions, "pass")
						} else {
							actions = append(actions, e.Action)
						}
					}
				}
			}
		}
	}
	return actions
}

func judgeExpectFromAct(act string) string {
	if act == "accept" {
		return "pass"
	} else {
		return "drop"
	}
}
