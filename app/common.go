package app

import (
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	urlDomainProtocol = []string{"http", "https"}
	icmpInt           = 9999
	protocolOPPorts   = []string{"21", "25", "53", "80", "110", "119", "135", "143", "445"}
)

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

type (
	AllInfo struct {
		Services          []string
		Act               string
		WfpName           string
		AvpName           string
		ServiceInfoAll    []ServiceInfo
		ServiceGrpInfoAll []ServiceGrpInfo
		ProxyMode         bool
		VInfo             []VipInfo
		VipGrpInfo        []vipGrpInfo
		WfProfileInfoAll  []WfProfileInfo
		WfFilterInfoAll   []WfFilterInfo
		AVProfileInfoAll  []AVProfileInfo
		ProxyModeProtocol []string
		ServicePort       map[string]string
		Env               string
		AddressInfoAll    []AddressInfo
		AddrGrpInfoAll    []AddrGrpInfo
		SRouteInfoAll     []SRouteInfo
		IntfInfoAll       []IntfInfo
	}
)

// ポリシーで使用されているintfがintf.csvに記載されていない場合にエラーを返す
func confirmIntfWithIntfInfo(intf, addr string, aI AllInfo) IntfInfo {
	var intfInfo IntfInfo
	switch {
	case addr == `"all"` && intf == `"any"`:
		if aI.Env != "test" {
			fmt.Print(`アドレスが"all"&インターフェースが"any"の場合は値を取得できません` + "\n")
			fmt.Print("該当ポリシーを削除、またはアドレスかインターフェースのいずれかを指定してください\n")
			os.Exit(106)
		}
	case intf == `"any"`:
		// ポリシーのインターフェースがanyの時は送信元・宛先で指定されたアドレスのスタティックルーティングを参照する
		// FQDNはデフォルトゲートウェイに設定されているintfを参照する
		// addrが複数割り当て＆別のルーティングを参照している場合は現状考慮していない
		policyIntf := getIntfInfoFromAddressAndSRoute(addr, aI)
		// 割り出したintfをintf.csvから読み取る
		intfInfo = getIntfInfo(policyIntf, aI.Env, aI.IntfInfoAll)
	default:
		intfInfo = getIntfInfo(intf, aI.Env, aI.IntfInfoAll)
	}

	if reflect.DeepEqual(intfInfo, IntfInfo{}) {
		if aI.Env != "test" {
			fmt.Printf("%sの情報はconfigディレクトリ内の`intf.csv`に記載されていませんでした\n", intf)
			fmt.Printf("インターフェース名%sの情報をconfigディレクトリ内の`intf.csv`に記載してください\n", intf)
			os.Exit(103)
		}
		return intfInfo
	}
	return intfInfo
}

func getIntfInfoFromAddressAndSRoute(addr string, aI AllInfo) string {
	// 送信元・宛先で指定されたアドレスのルーティングを参照して使用されているintfを取得する
	var addrInfo AddressInfo
	var err error
	if aI.VInfo != nil {
		// AddressのみのaddrInfoを作成する
		addrInfo = AddressInfo{
			Address: aI.VInfo[0].ExtIP,
		}
	} else {
		addrInfo, err = getAddrInfo(addr, aI.AddressInfoAll)
		if err != nil {
			addrGrpInfo := getAddrGrpInfo(addr, aI.AddrGrpInfoAll)
			if reflect.DeepEqual(addrGrpInfo, AddrGrpInfo{}) {
				fmt.Printf("%+vは存在しないアドレスが指定されています\n", addr)
				os.Exit(106)
			}
			// Memberが複数割り当て＆別のルーティングを参照している場合は現状考慮していない
			addrInfo, err = getAddrInfo(addrGrpInfo.Member[0], aI.AddressInfoAll)
			if err != nil {
				fmt.Printf("%+vは存在しないアドレスが指定されています\n", addr)
				os.Exit(106)
			}
		}
	}
	// 送信元・宛先で指定されたアドレスがスタティックルート&intf.csvのどのルーティングに該当するかチェックする
	// intfのネットワークはスタティックルートに乗らないのでintf.csvからも取得する
	// intf.csvのルート情報を全てここで確認する
	for _, dRoute := range aI.IntfInfoAll {
		if dRoute.Address != "" && dRoute.SubnetMask != "" {
			isMatched, intf := getIntfFromNetworkAddr(dRoute.Name, addrInfo.Address, dRoute.Address, dRoute.SubnetMask)
			if isMatched {
				return intf
			}
		}
	}
	// スタティックルートを全てここで確認する
	for _, route := range aI.SRouteInfoAll {
		if route.Dst != "" {
			netaddr := strings.Split(route.Dst, " ")[0]
			subnetMask := strings.Split(route.Dst, " ")[1]
			isMatched, intf := getIntfFromNetworkAddr(route.Device, addrInfo.Address, netaddr, subnetMask)
			if isMatched {
				fmt.Print("match", intf)
				return intf
			}
		}
	}
	// 全てのスタティックルート&intf.csvのルートにマッチしなければデフォルトゲートウェイを利用していると判断する
	for _, route := range SRouteInfoAll {
		if route.Dst == "" {
			return route.Device
		}
	}
	return ""
}

func getIntfFromNetworkAddr(checkIntf, checkAddr, netAddr, subnetMask string) (bool, string) {
	_, ipnet, err := convertSubnetMaskToCIDR(netAddr, subnetMask)
	if err != nil {
		os.Exit(107)
	}

	address := net.ParseIP(checkAddr)
	if ipnet.Contains(address) {
		return true, checkIntf
	}
	return false, checkIntf
}

func getIntfInfo(intf, env string, intfInfoAll []IntfInfo) IntfInfo {
	for _, v := range intfInfoAll {
		if intf == v.Name {
			return v
		}
	}
	return IntfInfo{}
}

func convertNWProtocol(serv string, ServiceInfoAll []ServiceInfo) []string {
	var nwProtocol []string
	switch serv {
	case `"ALL"`:
		nwProtocol = AllNWProtocol
	case `"DNS"`:
		nwProtocol = []string{"dnst", "dns"}
	case `"FTP"`:
		nwProtocol = []string{"ftp", "ftpa"}
	case `"HTTP"`, `"HTTPS"`, `"IMAP"`, `"SMTP"`:
		nwProtocol = []string{strings.ToLower(strings.Replace(serv, `"`, "", -1))}
	default:
		// `"ALL_ICMP"`, `"ALL_TCP"`, `"ALL_UDP"`などはここでそれぞれのプロトコルに変換される
		nwProtocol = convertIncompatibleService(serv, ServiceInfoAll)
	}
	return nwProtocol
}

func convertIncompatibleService(serv string, ServiceInfoAll []ServiceInfo) []string {
	var nwProtocolSlice []string
	var match bool
	for _, s := range ServiceInfoAll {
		if serv == s.Name {
			match = true
			if s.ICMP {
				nwProtocolSlice = append(nwProtocolSlice, "icmp")
			}
			if s.TCP != nil {
				nwProtocolSlice = append(nwProtocolSlice, "tcp")
			}
			if s.UDP != nil {
				nwProtocolSlice = append(nwProtocolSlice, "udp")
			}
		}
	}

	if !match {
		fmt.Printf("%+vはサポートされていません\n", serv)
	}
	return nwProtocolSlice
}

func isMatchedServiceToL4level(s, vipService string, ServiceInfoAll []ServiceInfo) bool {
	sInfo := getServiceInfo(s, ServiceInfoAll)
	serviceInfo := getServiceInfo(vipService, ServiceInfoAll)

	if (sInfo.ICMP && serviceInfo.ICMP) || (sInfo.TCP != nil && serviceInfo.TCP != nil) || (sInfo.UDP != nil && serviceInfo.UDP != nil) {
		return true
	}
	return false
}

func getServiceInfo(service string, ServiceInfoAll []ServiceInfo) ServiceInfo {
	for _, v := range ServiceInfoAll {
		if v.Name == service {
			return v
		}
	}
	return ServiceInfo{}
}

func judgeServiceFromProtocol(service, protocol string, ServiceInfoAll []ServiceInfo) bool {
	for _, v := range ServiceInfoAll {
		if v.Name == service {
			switch protocol {
			case "icmp":
				if v.ICMP {
					return true
				}
			case "tcp":
				if v.TCP != nil {
					return true
				}
			case "udp":
				if v.UDP != nil {
					return true
				}
			}
		}
	}
	return false
}

// FortiGateのServiceレベルで返却する関数
func judgeService(services []string, ServiceInfoAll []ServiceInfo, ServiceGrpInfoAll []ServiceGrpInfo) []string {
	var serviceSlice []string
	for _, serv := range services {
		var isMatchedService bool
		for _, s := range ServiceInfoAll {
			if serv == s.Name {
				isMatchedService = true
				serviceSlice = append(serviceSlice, serv)
			}
		}

		if !isMatchedService {
			for _, sg := range ServiceGrpInfoAll {
				if serv == sg.Name {
					isMatchedService = true
					serviceSlice = append(serviceSlice, sg.Member...)
				}
			}
		}
	}
	return serviceSlice
}

func getWfProfileInfo(wfpName string, WfProfileInfoAll []WfProfileInfo) WfProfileInfo {
	for _, wfp := range WfProfileInfoAll {
		if wfpName == wfp.Name {
			return wfp
		}
	}
	return WfProfileInfo{}
}

// othersettingsでdstportの値を見て書き換える必要がある箇所の処理
func judgeProtocolOPPorts(port string) bool {
	for _, opPort := range protocolOPPorts {
		if opPort == port {
			return true
		}
	}
	return false
}

func appendUndefined(component string, nwProtocol []string, aI AllInfo) []string {
	var appendSlice []string
	switch component {
	case "urlDomain":
		urlDomains := judgeURLDomainOutput(aI.WfpName, nwProtocol, aI.WfProfileInfoAll, aI.WfFilterInfoAll)
		appendSlice = appendStrDataToSlice(urlDomains, "undefined")
	case "expect":
		expects := judgeExpectOutput(aI.Act, aI.WfpName, aI.AvpName, nwProtocol, aI.WfProfileInfoAll, aI.WfFilterInfoAll)
		appendSlice = appendStrDataToSlice(expects, "undefined")
	default:
		appendSlice = appendStrDataToSlice(nwProtocol, "undefined")
	}
	return appendSlice
}

func appendStrDataToSlice(nwProtocol []string, data string) (appendSlice []string) {
	for range nwProtocol {
		appendSlice = append(appendSlice, data)
	}
	return
}

// protocols, srcPorts, dstNATAddress, dstAddress, dstPorts, antiVirus, otherSettingsは必ず同じ長さになる
func checkDatasLength(rs []relatedService, components ...[]string) bool {
	for _, component := range components {
		ok := checkDataLength(rs, component)
		if !ok {
			return false
		}
	}
	return true
}

func checkDataLength(rs []relatedService, components []string) bool {
	if len(rs) != len(components) {
		fmt.Printf("len(rs) = %+v\n", len(rs))
		fmt.Printf("component = %+v\n", components)
		fmt.Printf("len component = %+v\n", len(components))
		return false
	}
	return true
}

func existNWProtocol(serv string) bool {
	switch serv {
	case `"DNS"`, `"HTTP"`, `"HTTPS"`, `"FTP"`, `"IMAP"`, `"SMTP"`:
		return true
	}
	return false
}

func getAddrInfo(addr string, AddressInfoAll []AddressInfo) (AddressInfo, error) {
	for _, addressInfo := range AddressInfoAll {
		if addr == addressInfo.Name {
			return addressInfo, nil
		}
	}
	return AddressInfo{}, errors.New("failed to get addrInfo")
}

func getAddrGrpInfo(addr string, AddrGrpInfoAll []AddrGrpInfo) AddrGrpInfo {
	for _, addrGrpInfo := range AddrGrpInfoAll {
		if addr == addrGrpInfo.Name {
			return addrGrpInfo
		}
	}
	return AddrGrpInfo{}
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

func handleVInfo(component string, aI AllInfo) []string {
	var appendSlice []string
	servs := judgeService(aI.Services, aI.ServiceInfoAll, aI.ServiceGrpInfoAll)
	for _, c := range aI.VInfo {
		for _, s := range servs {
			switch s {
			case `"ALL"`:
				switch {
				case c.Service != nil:
					switch component {
					case "dst_nat_port":
						a := appendDataForDstNatPort(aI.ServiceInfoAll, aI.ServiceGrpInfoAll, c, c.Service, AllNWProtocol, aI.ServicePort)
						appendSlice = append(appendSlice, a...)
					}

					for _, vipService := range c.Service {
						nwProtocol := convertNWProtocol(vipService, aI.ServiceInfoAll)
						appendData := appendDataEachComponent2(component, nwProtocol, aI.Services, nil, c, aI)
						appendSlice = append(appendSlice, appendData...)
					}
				case c.Protocol != "":
					appendData := appendDataEachComponent4(component, []string{c.Protocol}, aI.Services, []string{c.Protocol}, c, aI)
					appendSlice = append(appendSlice, appendData...)
				default:
					switch component {
					case "protocol":
						nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
						appendSlice = append(appendSlice, nwProtocol...)
					default:
						appendData := appendDataEachComponent4(component, AllNWProtocol, aI.Services, AllNWProtocol, c, aI)
						appendSlice = append(appendSlice, appendData...)
					}
				}
			case `"ALL_ICMP"`, `"PING"`, `"ALL_TCP"`, `"ALL_UDP"`:
				switch {
				case c.Service != nil:
					nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
					for _, vipService := range c.Service {
						switch vipService {
						case `"ALL"`:
							appendData := handleVInfoWithService(component, s, nwProtocol, servs, c, aI)
							appendSlice = append(appendSlice, appendData...)
						case `"PING"`:
							if s == `"PING"` || s == `"ALL_ICMP"` {
								appendData := appendDataEachComponent4(component, nwProtocol, aI.Services, nwProtocol, c, aI)
								appendSlice = append(appendSlice, appendData...)
							} else {
								appendSlice = append(appendSlice, "undefined")
							}
						case `"FTP"`:
							isMatchedL4Level := isMatchedServiceToL4level(s, vipService, aI.ServiceInfoAll)
							if isMatchedL4Level {
								nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
								switch component {
								case "dst_nat_port":
									appendData := appendDataForDstNatPort(aI.ServiceInfoAll, aI.ServiceGrpInfoAll, c, []string{vipService}, nwProtocol, aI.ServicePort)
									// len(appendData) = 2となり同じ要素のため最初の要素のみappendする
									appendSlice = append(appendSlice, appendData[0])
								default:
									appendSlice = appendDataEachComponent2(component, nwProtocol, aI.Services, nil, c, aI)
								}
							} else {
								appendData := appendUndefined(component, nwProtocol, aI)
								appendSlice = append(appendSlice, appendData...)
							}
						default:
							nwProtocol := convertIncompatibleService(vipService, aI.ServiceInfoAll)
							if len(nwProtocol) > 1 {
								// tcp&udpのサービス(例:DNS,NTP)
								switch s {
								case `"ALL_TCP"`:
									appendData := appendDataEachComponent3(c, []string{vipService}, nwProtocol, []string{vipService}, component, aI)
									appendData[1] = "undefined"
									appendSlice = append(appendSlice, appendData...)
								case `"ALL_UDP"`:
									appendData := appendDataEachComponent3(c, []string{vipService}, nwProtocol, []string{vipService}, component, aI)
									appendData[0] = "undefined"
									appendSlice = append(appendSlice, appendData...)
								default:
									appendData := appendUndefined(component, nwProtocol, aI)
									appendSlice = append(appendSlice, appendData...)
								}
							} else {
								isMatchedL4Level := isMatchedServiceToL4level(s, vipService, aI.ServiceInfoAll)
								if isMatchedL4Level {
									data := appendDataEachComponent3(c, aI.Services, nwProtocol, []string{vipService}, component, aI)
									appendSlice = append(appendSlice, data...)
								} else {
									appendData := appendUndefined(component, nwProtocol, aI)
									appendSlice = append(appendSlice, appendData...)
								}
							}
						}
					}
				case c.Protocol != "":
					isMatchedL4Level := judgeServiceFromProtocol(s, c.Protocol, aI.ServiceInfoAll)
					nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
					if isMatchedL4Level {
						appendData := appendDataEachComponent4(component, nwProtocol, aI.Services, []string{c.Protocol}, c, aI)
						appendSlice = append(appendSlice, appendData...)
					} else {
						appendData := appendUndefined(component, nwProtocol, aI)
						appendSlice = append(appendSlice, appendData...)
					}
				default:
					nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
					appendData := appendDataEachComponent4(component, nwProtocol, aI.Services, []string{s}, c, aI)
					appendSlice = append(appendSlice, appendData...)
				}
			case `"FTP"`:
				nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
				switch {
				case c.Service != nil:
					appendData := handleVInfoWithService(component, s, nwProtocol, servs, c, aI)
					appendSlice = append(appendSlice, appendData...)
				case c.Protocol != "":
					nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
					if c.MappedPort != "" {
						port := handleDstPort([]string{s}, aI.ServiceInfoAll, aI.ServiceGrpInfoAll, aI.ServicePort)
						// s=ftpの場合はlen(port)=1となるので最初の要素[0]とc.MappedPortを比較する
						if port[0] == c.MappedPort {
							appendData := appendDataEachComponent4(component, nwProtocol, aI.Services, nwProtocol, c, aI)
							switch c.Protocol {
							// "FTP"はtcpのみマッチするのでudpの場合はundefinedに置き換える必要がある
							case "udp":
								appendData[1] = "undefined"
							}
							appendSlice = append(appendSlice, appendData...)
						} else {
							appendData := appendUndefined(component, nwProtocol, aI)
							appendSlice = append(appendSlice, appendData...)
						}
					} else {
						appendData := appendDataEachComponent3(c, aI.Services, nwProtocol, []string{s}, component, aI)
						appendSlice = append(appendSlice, appendData...)
					}
				default:
					appendData := handleVInfoWithoutOP(component, s, nwProtocol, servs, c, aI)
					appendSlice = append(appendSlice, appendData...)
				}
			default:
				nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
				switch {
				case c.Service != nil:
					appendData := handleVInfoWithService(component, s, nwProtocol, servs, c, aI)
					appendSlice = append(appendSlice, appendData...)
				case c.Protocol != "":
					isMatchedL4Level := judgeServiceFromProtocol(s, c.Protocol, aI.ServiceInfoAll)
					nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
					var appendDataSlice []string
					if isMatchedL4Level {
						if c.MappedPort != "" {
							port := handleDstPort([]string{s}, aI.ServiceInfoAll, aI.ServiceGrpInfoAll, aI.ServicePort)
							// tcp&udpのサービスか否かを判別し処理する
							if len(port) > 1 {
								switch c.Protocol {
								// c.Protocolに該当しないデータをundefinedに置き換える
								case "tcp":
									appendData := judgeMappedPortWithNum(component, port, nwProtocol, c, aI, 0)
									appendDataSlice = append(appendDataSlice, appendData...)
								case "udp":
									appendData := judgeMappedPortWithNum(component, port, nwProtocol, c, aI, 1)
									appendDataSlice = append(appendDataSlice, appendData...)
								// 通常c.Protocol=icmpの設定はできないのでこの処理をすることはほぼない
								default:
									appendData := appendDataEachComponent3(c, servs, nwProtocol, aI.Services, component, aI)
									appendDataSlice = append(appendDataSlice, appendData...)
								}
							} else {
								if port[0] == c.MappedPort {
									appendData := appendDataEachComponent4(component, nwProtocol, aI.Services, nwProtocol, c, aI)
									appendDataSlice = append(appendDataSlice, appendData...)
								} else {
									appendDataSlice = appendStrDataToSlice(nwProtocol, "undefined")
								}
							}
						} else {
							nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
							appendData := appendDataEachComponent3(c, servs, nwProtocol, aI.Services, component, aI)
							appendDataSlice = append(appendDataSlice, appendData...)
						}
					} else {
						appendDataSlice = append(appendDataSlice, "undefined")
					}
					appendSlice = append(appendSlice, appendDataSlice...)
				default:
					appendData := handleVInfoWithoutOP(component, s, nwProtocol, servs, c, aI)
					appendSlice = append(appendSlice, appendData...)
				}
			}
		}
	}
	return appendSlice
}

// vip以外でdstaddrが複数になる場合にdstaddr以外のデータをdstaddrの種類分増やす
func appendDataToComponentWithDstAddr(intf string, addrs []string, usedAddress, component []string, allInfo AllInfo) []string {
	dstSlice, _ := getUniqueDstAddr(intf, addrs, usedAddress, allInfo)
	var components []string
	for range dstSlice {
		components = append(components, component...)
	}
	return components
}

func convertSubnetMaskToCIDR(addr, subnetMask string) (net.IP, *net.IPNet, error) {
	var ip net.IP
	var ipnet *net.IPNet
	var err error
	subnet := net.ParseIP(subnetMask)
	if subnet == nil {
		return ip, ipnet, errors.New("failed to parse subnet")
	}
	mask := net.IPv4Mask(subnet[12], subnet[13], subnet[14], subnet[15])
	length, _ := mask.Size()
	ip, ipnet, err = net.ParseCIDR(addr + "/" + strconv.Itoa(length))
	if err != nil {
		// この処理はほぼ起こり得ないが念の為ハンドリング
		fmt.Printf("failed to parse cidr %s\n", err.Error())
		return ip, ipnet, err
	}
	return ip, ipnet, err
}
