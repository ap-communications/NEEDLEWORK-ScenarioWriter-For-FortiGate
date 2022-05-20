package app

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	urlDomainProtocol = []string{"http", "https"}
	icmpInt           = 9999
	protocolOPPorts   = []string{"21", "25", "53", "80", "110", "119", "135", "143", "445"}
)

type AllInfo struct {
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

func isMatchedServiceToL4level(s, service string, ServiceInfoAll []ServiceInfo) bool {
	sInfo := getServiceInfo(s, ServiceInfoAll)
	serviceInfo := getServiceInfo(service, ServiceInfoAll)

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

func appendDataForDstNatPort(ServiceInfoAll []ServiceInfo, ServiceGrpInfoAll []ServiceGrpInfo, c VipInfo, services, nwProtocol []string, ServicePort map[string]string) []string {
	var appendSlice []string
	if c.MappedPort != "" {
		dstnatport := handleDstPort(services, ServiceInfoAll, ServiceGrpInfoAll, ServicePort)
		appendSlice = append(appendSlice, dstnatport...)
	} else {
		appendSlice = appendStrDataToSlice(nwProtocol, "")
	}
	return appendSlice
}

func appendDataEachComponent3(c VipInfo, services, nwProtocol, ds []string, component string, aI AllInfo) []string {
	var appendSlice []string
	switch component {
	case "dst_nat_port":
		appendSlice = appendDataForDstNatPort(aI.ServiceInfoAll, aI.ServiceGrpInfoAll, c, ds, nwProtocol, aI.ServicePort)
	default:
		appendSlice = appendDataEachComponent2(component, nwProtocol, services, c, aI)
	}
	return appendSlice
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

func HandleProtocolOutput(allInfo AllInfo) []string {
	var protocolSlice []string
	if allInfo.VInfo != nil {
		protocolSlice = handleVInfo2("protocol", allInfo)
	} else {
		servs := judgeService(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll)
		for _, serv := range servs {
			nwProtocol := convertNWProtocol(serv, allInfo.ServiceInfoAll)
			protocolSlice = append(protocolSlice, nwProtocol...)
		}
	}
	return protocolSlice
}

func HandleSrcPortOutput(allInfo AllInfo) []string {
	var SrcPortSlice []string
	if allInfo.VInfo != nil {
		SrcPortSlice = handleVInfo2("src_port", allInfo)
	} else {
		SrcPortSlice = handleSrcPort(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll, allInfo.ServicePort)
	}
	return SrcPortSlice
}

func HandleDstNATPortOutput(allInfo AllInfo) []string {
	var dstNATPortSlice []string
	if allInfo.VInfo != nil {
		dstNATPortSlice = handleVInfo2("dst_nat_port", allInfo)
	} else {
		services := judgeService(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll)
		var ss []string
		for _, s := range services {
			serv := convertNWProtocol(s, allInfo.ServiceInfoAll)
			ss = append(ss, serv...)
		}
		for range ss {
			dstNATPortSlice = append(dstNATPortSlice, "")
		}
	}
	return dstNATPortSlice
}

func HandleDstPortOutput(allInfo AllInfo) []string {
	var dstPorts []string
	if allInfo.VInfo != nil {
		dstPorts = handleVInfo2("dst_port", allInfo)
	} else {
		dstPorts = handleDstPort(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll, allInfo.ServicePort)
	}
	return dstPorts
}

func handleSrcPort(services []string, ServiceInfoAll []ServiceInfo, ServiceGrpInfoAll []ServiceGrpInfo, ServicePort map[string]string) []string {
	var portSlice []string
	services = judgeService(services, ServiceInfoAll, ServiceGrpInfoAll)
	for _, serv := range services {
		if existNWProtocol(serv) || serv == `"ALL"` || serv == `"ALL_ICMP"` || serv == `"ALL_TCP"` || serv == `"ALL_UDP"` {
			nwProtocol := convertNWProtocol(serv, ServiceInfoAll)
			portSlice = append(portSlice, appendStrDataToSlice(nwProtocol, "")...)
		} else {
			port := handleAssignSrcPort(serv, ServiceInfoAll)
			portSlice = append(portSlice, port...)
		}
	}
	return portSlice
}

func handleSrcPort2(services []string, ServiceInfoAll []ServiceInfo, ServiceGrpInfoAll []ServiceGrpInfo) []string {
	var portSlice []string
	servs := judgeService(services, ServiceInfoAll, ServiceGrpInfoAll)
	for _, serv := range servs {
		ports := HandleAssignSrcPort2(serv, ServiceInfoAll)
		if existNWProtocol(serv) {
			nwProtocol := convertNWProtocol(serv, ServiceInfoAll)
			for range nwProtocol {
				portSlice = append(portSlice, "")
			}
		} else {
			portSlice = append(portSlice, ports...)
		}
	}
	return portSlice
}

func handleDstPort(services []string, ServiceInfoAll []ServiceInfo, ServiceGrpInfoAll []ServiceGrpInfo, ServicePort map[string]string) []string {
	var portSlice []string
	servs := judgeService(services, ServiceInfoAll, ServiceGrpInfoAll)
	for _, serv := range servs {
		ports := HandleAssignDstPort(serv, ServiceInfoAll, ServicePort)
		if existNWProtocol(serv) {
			// DNSの場合のみlen(ports)=2となる
			// DNSの要素はTCP,UDPでそれぞれ別のポートが指定される可能性もある
			nwProtocol := convertNWProtocol(serv, ServiceInfoAll)
			if len(ports) == len(nwProtocol) {
				portSlice = append(portSlice, ports...)
			} else {
				// 現在ではFTPが該当する
				// len(ports) = 1, len(nwProtocol) = 2
				for range nwProtocol {
					portSlice = append(portSlice, ports...)
				}
			}
		} else {
			portSlice = append(portSlice, ports...)
		}
	}
	return portSlice
}

func judgeService(services []string, ServiceInfoAll []ServiceInfo, ServiceGrpInfoAll []ServiceGrpInfo) []string {
	var serviceSlice []string
	if reflect.DeepEqual(services, []string{`"ALL"`}) {
		return []string{`"ALL"`}
	}

	for _, serv := range services {
		var settingServ bool
		for _, s := range ServiceInfoAll {
			if serv == s.Name {
				settingServ = true
				serviceSlice = append(serviceSlice, serv)
			}
		}

		if !settingServ {
			for _, sg := range ServiceGrpInfoAll {
				if serv == sg.Name {
					settingServ = true
					serviceSlice = append(serviceSlice, sg.Member...)
				}
			}
		}
	}
	return serviceSlice
}

// s.tcpとかは[]portassigninfoで複数ポートが入る可能性がある
// 現在は最初のポートのみを出力している
func handleAssignSrcPort(serv string, ServiceInfoAll []ServiceInfo) []string {
	var srcPortSlice []string
	for _, s := range ServiceInfoAll {
		if serv == s.Name && s.ICMP {
			srcPortSlice = append(srcPortSlice, "")
		}
		if serv == s.Name && s.TCP != nil {
			srcPort := s.TCP[0].SrcPort
			if strings.Contains(srcPort, "-") {
				srcPortSlice = append(srcPortSlice, strings.Split(srcPort, "-")[0])
			} else {
				srcPortSlice = append(srcPortSlice, s.TCP[0].SrcPort)
			}
		}
		if serv == s.Name && s.UDP != nil {
			srcPort := s.UDP[0].SrcPort
			if strings.Contains(srcPort, "-") {
				srcPortSlice = append(srcPortSlice, strings.Split(srcPort, "-")[0])
			} else {
				srcPortSlice = append(srcPortSlice, s.UDP[0].SrcPort)
			}
		}
	}
	return srcPortSlice
}

// s.tcpとかは[]portassigninfoで複数ポートが入る可能性がある
// 現在は最初のポートのみを出力している
func HandleAssignSrcPort2(serv string, ServiceInfoAll []ServiceInfo) []string {
	var srcPortSlice []string
	switch serv {
	case `"ALL"`, `"PING"`, `"ALL_ICMP"`, `"ALL_TCP"`, `"ALL_UDP"`:
		nwProtocol := convertNWProtocol(serv, ServiceInfoAll)
		for range nwProtocol {
			srcPortSlice = append(srcPortSlice, "")
		}
	default:
		for _, s := range ServiceInfoAll {
			if serv == s.Name {
				srcport := getMatchedServiceSrcPort(s)
				srcPortSlice = append(srcPortSlice, srcport...)
			}
		}
	}
	return srcPortSlice
}

// s.tcpとかは[]portassigninfoで複数ポートが入る可能性がある
// 現在は最初のポートのみを出力している
func HandleAssignDstPort(serv string, ServiceInfoAll []ServiceInfo, ServicePort map[string]string) []string {
	var dstPortSlice []string
	switch serv {
	case `"ALL"`, `"PING"`, `"ALL_ICMP"`, `"ALL_TCP"`, `"ALL_UDP"`:
		nwProtocol := convertNWProtocol(serv, ServiceInfoAll)
		for _, protocol := range nwProtocol {
			// protocolに該当するserviceのポートを取得しappendする
			dstPort := getMatchedNWProtocolPort(protocol, ServiceInfoAll, ServicePort)
			dstPortSlice = append(dstPortSlice, dstPort)
		}
	default:
		for _, s := range ServiceInfoAll {
			if serv == s.Name {
				dstport := getMatchedServiceDstPort(s)
				dstPortSlice = append(dstPortSlice, dstport...)
			}
		}
	}
	return dstPortSlice
}

// protocolに該当するserviceのポートを取得する
func getMatchedNWProtocolPort(serv string, ServiceInfoAll []ServiceInfo, ServicePort map[string]string) string {
	switch serv {
	case "icmp", "tcp", "udp":
		// 特定のポートにはならない(1-65535)のでServicePortから値を補完する
		for k, v := range ServicePort {
			if serv == k {
				return v
			}
		}
	// ftpaはNEEDLEWORKのプロトコルなので元のサービスを参照するようにする
	case "ftpa":
		serv = "ftp"
	}

	for _, s := range ServiceInfoAll {
		matchService := strings.ToLower(strings.Replace(s.Name, `"`, "", -1))
		switch serv {
		case "dns":
			if matchService == "dns" {
				// dnsはlen(getMatchedServiceDstPort(s)) = 2となる
				// tcp,udpの順で格納されるのでudpの[1]を返す
				return getMatchedServiceDstPort(s)[1]
			}
		case "dnst":
			if matchService == "dns" {
				// dnsはlen(getMatchedServiceDstPort(s)) = 2となる
				// tcp,udpの順で格納されるのでtcpの[0]を返す
				return getMatchedServiceDstPort(s)[0]
			}
		case matchService:
			return getMatchedServiceDstPort(s)[0]
		}
	}
	return ""
}

func getMatchedServiceSrcPort(s ServiceInfo) []string {
	var srcPortSlice []string
	if s.ICMP {
		srcPortSlice = append(srcPortSlice, "")
	}
	if s.TCP != nil {
		srcport := s.TCP[0].SrcPort
		if strings.Contains(srcport, "-") {
			srcPortSlice = append(srcPortSlice, strings.Split(srcport, "-")[0])
		} else {
			srcPortSlice = append(srcPortSlice, s.TCP[0].SrcPort)
		}
	}
	if s.UDP != nil {
		srcport := s.UDP[0].SrcPort
		if strings.Contains(srcport, "-") {
			srcPortSlice = append(srcPortSlice, strings.Split(srcport, "-")[0])
		} else {
			srcPortSlice = append(srcPortSlice, s.UDP[0].SrcPort)
		}
	}
	return srcPortSlice
}

func getMatchedServiceDstPort(s ServiceInfo) []string {
	var dstPortSlice []string
	if s.ICMP {
		dstPortSlice = append(dstPortSlice, "")
	}
	if s.TCP != nil {
		dstport := s.TCP[0].DstPort
		if strings.Contains(dstport, "-") {
			dstPortSlice = append(dstPortSlice, strings.Split(dstport, "-")[0])
		} else {
			dstPortSlice = append(dstPortSlice, s.TCP[0].DstPort)
		}
	}
	if s.UDP != nil {
		dstport := s.UDP[0].DstPort
		if strings.Contains(dstport, "-") {
			dstPortSlice = append(dstPortSlice, strings.Split(dstport, "-")[0])
		} else {
			dstPortSlice = append(dstPortSlice, s.UDP[0].DstPort)
		}
	}
	return dstPortSlice
}

func HandleOtherSettingOutput(allInfo AllInfo) []string {
	var otherSettings []string
	if allInfo.VInfo != nil {
		otherSettings = handleVInfo2("other_settings", allInfo)
	} else {
		services := judgeService(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll)
		for _, serv := range services {
			nwProtocol := convertNWProtocol(serv, allInfo.ServiceInfoAll)
			os := judgeOtherSettingOutput(allInfo.ProxyMode, nwProtocol, allInfo.ProxyModeProtocol, allInfo.AvpName, allInfo.AVProfileInfoAll, allInfo)
			otherSettings = append(otherSettings, os...)
		}
		// 宛先ポートがdefaultのプロトコルオプションに該当しているか否かで挙動を変える
		// defaultのプロトコルオプションで設定されているポート: 21,25,53,80,110,119,135,143,445
		dstPorts := HandleDstPortOutput(allInfo)
		for i := range otherSettings {
			for ii, port := range dstPorts {
				isProtocolOPPort := judgeProtocolOPPorts(port)
				if i == ii && !isProtocolOPPort {
					otherSettings[i] = ""
				}
			}
		}
	}
	return otherSettings
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

func handleRelatedServiceOutPut(p policyInfo, sI []ServiceInfo, ServiceGrpInfoAll []ServiceGrpInfo, IntfInfoAll []IntfInfo, allInfo AllInfo) ([]relatedService, bool) {
	var relatedServices []relatedService
	dstAddress, dstFQDNFlag := HandleDstAddressOutput(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, allInfo)
	protocols := HandleProtocolOutput(allInfo)
	srcPorts := HandleSrcPortOutput(allInfo)
	dstNATAddress := HandleDstNATAddrOutput(allInfo)
	dstNATPorts := HandleDstNATPortOutput(allInfo)
	dstPorts := HandleDstPortOutput(allInfo)
	otherSettings := HandleOtherSettingOutput(allInfo)
	urlDomains := HandleURLDomainOutput(allInfo)
	antiVirus := HandleAntiVirusOutput(allInfo)
	expects := HandleExpectOutput(allInfo)

	// debug
	// fmt.Printf("dstAddress = %+v\n", dstAddress)
	// fmt.Printf("len(dstAddress) = %+v\n", len(dstAddress))
	// fmt.Printf("protocols = %+v\n", protocols)
	// fmt.Printf("len(protocols) = %+v\n", len(protocols))
	// fmt.Printf("srcPorts = %+v\n", srcPorts)
	// fmt.Printf("len(srcPorts) = %+v\n", len(srcPorts))
	// fmt.Printf("dstNATAddress = %+v\n", dstNATAddress)
	// fmt.Printf("len(dstNATAddress) = %+v\n", len(dstNATAddress))
	// fmt.Printf("dstNATPorts = %+v\n", dstNATPorts)
	// fmt.Printf("len(dstNATPorts) = %+v\n", len(dstNATPorts))
	// fmt.Printf("dstPorts = %+v\n", dstPorts)
	// fmt.Printf("len(dstPorts) = %+v\n", len(dstPorts))
	// fmt.Printf("otherSettings = %+v\n", otherSettings)
	// fmt.Printf("len(otherSettings) = %+v\n", len(otherSettings))
	// fmt.Printf("urlDomains = %+v\n", urlDomains)
	// fmt.Printf("len(urlDomains) = %+v\n", len(urlDomains))
	// fmt.Printf("antiVirus = %+v\n", antiVirus)
	// fmt.Printf("len(antiVirus) = %+v\n", len(antiVirus))
	// fmt.Printf("expects = %+v\n", expects)
	// fmt.Printf("len(expects) = %+v\n", len(expects))

	var quantityURL int
	urls := getURLDomainForWebFilter(p.WFProfile, WfProfileInfoAll, WfFilterInfoAll)
	if urls != nil {
		quantityURL = len(getURLDomainForWebFilter(p.WFProfile, WfProfileInfoAll, WfFilterInfoAll))
	} else {
		// url要素は「0」だが後にこの値を元に要素の出力数を変化させるため「1」と定義している
		quantityURL = 1
	}

	// urlDomains, expectsはquantityURLの値数分、要素が異なるため、quantityURLの数分格納される
	// その他の要素はquantityURLの値によって要素が変化しないため、quantityURLの数分格納しない
	// 具体例: protocols=["http"], urlDomains=["hoge.com","hogehoge.com"],expects=["pass","block"]
	// i := 0~3 protocols[0],i := 4~7 protocols[1]
	for i := range urlDomains {
		rS := relatedService{
			URLDomain:     urlDomains[i],
			Protocol:      protocols[i/quantityURL],
			SrcPort:       srcPorts[i/quantityURL],
			DstNATAddress: dstNATAddress[i/quantityURL],
			DstAddress:    dstAddress[i/quantityURL],
			DstNATPort:    dstNATPorts[i/quantityURL],
			DstPort:       dstPorts[i/quantityURL],
			AntiVirus:     antiVirus[i/quantityURL],
			OtherSettings: otherSettings[i/quantityURL],
			Expect:        expects[i],
		}
		relatedServices = append(relatedServices, rS)
	}
	return relatedServices, dstFQDNFlag
}

func HandleDstNATAddrOutput(allInfo AllInfo) []string {
	var dstNATAddrSlice []string
	if allInfo.VInfo != nil {
		dstNATAddrSlice = handleVInfo2("dst_nat_addr", allInfo)
	} else {
		services := judgeService(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll)
		for _, serv := range services {
			nwProtocol := convertNWProtocol(serv, allInfo.ServiceInfoAll)
			dstNATAddrSlice = append(dstNATAddrSlice, appendStrDataToSlice(nwProtocol, "")...)
		}
	}
	return dstNATAddrSlice
}

func HandleAntiVirusOutput(allInfo AllInfo) []string {
	var avSlice []string
	if allInfo.VInfo != nil {
		avSlice = handleVInfo2("anti_virus", allInfo)
	} else {
		services := judgeService(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll)
		for _, serv := range services {
			nwProtocol := convertNWProtocol(serv, allInfo.ServiceInfoAll)
			av := judgeAntiVirusOutput(allInfo.AvpName, nwProtocol, allInfo.AVProfileInfoAll)
			avSlice = append(avSlice, av...)
		}
	}
	return avSlice
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

func HandleURLDomainOutput(allInfo AllInfo) []string {
	var urlSlice []string
	if allInfo.VInfo != nil {
		urlSlice = handleVInfo2("urlDomain", allInfo)
	} else {
		services := judgeService(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll)
		for _, serv := range services {
			nwProtocol := convertNWProtocol(serv, allInfo.ServiceInfoAll)
			urls := judgeURLDomainOutput(allInfo.WfpName, nwProtocol, allInfo.WfProfileInfoAll, allInfo.WfFilterInfoAll)
			urlSlice = append(urlSlice, urls...)
		}
	}
	return urlSlice
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

func HandleExpectOutput(allInfo AllInfo) []string {
	var actSlice []string
	if allInfo.VInfo != nil {
		actSlice = handleVInfo2("expect", allInfo)
	} else {
		services := judgeService(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll)
		for _, serv := range services {
			nwProtocol := convertNWProtocol(serv, allInfo.ServiceInfoAll)
			expects := judgeExpectOutput(allInfo.Act, allInfo.WfpName, allInfo.AvpName, nwProtocol, allInfo.WfProfileInfoAll, allInfo.WfFilterInfoAll)
			actSlice = append(actSlice, expects...)
		}
	}
	return actSlice
}

func judgeExpectOutput(act, wfpName, avpName string, nwProtocol []string, WfProfileInfoAll []WfProfileInfo, WfFilterInfoAll []WfFilterInfo) []string {
	var actSlice []string
	// 一致するプロファイル内のActionを抜き出す
	actions := getActionForWebFilter(act, wfpName, WfProfileInfoAll, WfFilterInfoAll)
	avProtocols := getProtocolForAVProfile(avpName, AVProfileInfoAll)

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
							actSlice = append(actSlice, actions...)
						}
					}
				}
			}
			if !isURLDomainProtocol {
				for i := range actions {
					if i == 0 {
						for _, pp := range avProtocols {
							if p == pp {
								isURLDomainProtocol = true
							}
						}

						if isURLDomainProtocol {
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
				actSlice = append(actSlice, judgeExpectFromAct(act))
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

func judgeExpectFromAct(act string) string {
	if act == "accept" {
		return "pass"
	} else {
		return "drop"
	}
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

func judgeMappedPortWithNum(component string, port, nwProtocol []string, c VipInfo, aI AllInfo, num int) []string {
	var appendData []string
	if port[num] == c.MappedPort {
		appendData = appendDataEachComponent4(component, nwProtocol, aI.Services, nwProtocol, c, aI)
		// tcp&udpのサービスはどちらもマッチするのでマッチしていない方をundefinedに置き換える必要がある
		switch num {
		case 0:
			// tcpでマッチしているのでudpはundefinedに置き換える
			appendData[1] = "undefined"
		case 1:
			// udpでマッチしているのでtcpはundefinedに置き換える
			appendData[0] = "undefined"
		}
	} else {
		appendData = appendUndefined(component, nwProtocol, aI)
	}
	return appendData
}

func handleVInfoWithoutOP(component, s string, nwProtocol, servs []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	if c.MappedPort != "" {
		port := handleDstPort([]string{s}, aI.ServiceInfoAll, aI.ServiceGrpInfoAll, aI.ServicePort)
		// tcp&udpのサービスの場合はlen(port)=2となる
		// 値が異なる可能性もあるのでそれぞれをc.MappedPortと比較する
		appendSlice = JudgeAppendData3WithMappedPort(component, nwProtocol, servs, []string{s}, port, c, aI)
	} else {
		appendData := appendDataEachComponent4(component, nwProtocol, aI.Services, nwProtocol, c, aI)
		appendSlice = append(appendSlice, appendData...)
	}
	return appendSlice
}

func handleVInfoWithService(component, s string, nwProtocol, servs []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	// 以下の場合は出力する
	// c.ServiceにALLが含まれている場合
	// sとc.serviceがL4レベルで同じ場合
	// c.MappedPortが設定されている場合はserviceInfoのポートと同じ場合
	for _, service := range c.Service {
		nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
		switch service {
		case `"ALL"`:
			appendData := appendDataEachComponent3(c, aI.Services, nwProtocol, []string{s}, component, aI)
			appendSlice = append(appendSlice, appendData...)
		case `"ALL_ICMP"`, `"PING"`:
			appendData := judgeAppendData3FromL4level(component, s, service, nwProtocol, c, aI, icmpInt)
			appendSlice = append(appendSlice, appendData...)
		case `"ALL_TCP"`:
			appendData := judgeAppendData3FromL4level(component, s, service, nwProtocol, c, aI, 1)
			appendSlice = append(appendSlice, appendData...)
		case `"ALL_UDP"`:
			appendData := judgeAppendData3FromL4level(component, s, service, nwProtocol, c, aI, 0)
			appendSlice = append(appendSlice, appendData...)
		default:
			if s == service {
				if c.MappedPort != "" {
					port := handleDstPort([]string{s}, aI.ServiceInfoAll, aI.ServiceGrpInfoAll, aI.ServicePort)
					// tcp&udpのサービスの場合はlen(port)=2となる
					// 値が異なる可能性もあるのでそれぞれをc.MappedPortと比較する
					appendSlice = JudgeAppendData3WithMappedPort(component, nwProtocol, aI.Services, []string{s}, port, c, aI)
				} else {
					appendData := appendDataEachComponent3(c, aI.Services, nwProtocol, []string{s}, component, aI)
					appendSlice = append(appendSlice, appendData...)
				}
			} else {
				appendSlice = append(appendSlice, appendStrDataToSlice(nwProtocol, "undefined")...)
			}
		}
	}
	return appendSlice
}

func handleVInfo2(component string, aI AllInfo) []string {
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

					for _, v := range c.Service {
						nwProtocol := convertNWProtocol(v, aI.ServiceInfoAll)
						appendData := appendDataEachComponent2(component, nwProtocol, aI.Services, c, aI)
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
					for _, service := range c.Service {
						switch service {
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
							isMatchedL4Level := isMatchedServiceToL4level(s, service, aI.ServiceInfoAll)
							if isMatchedL4Level {
								nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
								switch component {
								case "dst_nat_port":
									appendData := appendDataForDstNatPort(aI.ServiceInfoAll, aI.ServiceGrpInfoAll, c, []string{service}, nwProtocol, aI.ServicePort)
									// len(appendData) = 2となり同じ要素のため最初の要素のみappendする
									appendSlice = append(appendSlice, appendData[0])
								default:
									appendSlice = appendDataEachComponent2(component, nwProtocol, aI.Services, c, aI)
								}
							} else {
								appendData := appendUndefined(component, nwProtocol, aI)
								appendSlice = append(appendSlice, appendData...)
							}
						default:
							nwProtocol := convertIncompatibleService(service, aI.ServiceInfoAll)
							if len(nwProtocol) > 1 {
								// tcp&udpのサービス(例:DNS,NTP)
								switch s {
								case `"ALL_TCP"`:
									appendData := appendDataEachComponent3(c, []string{service}, nwProtocol, []string{service}, component, aI)
									appendData[1] = "undefined"
									appendSlice = append(appendSlice, appendData...)
								case `"ALL_UDP"`:
									appendData := appendDataEachComponent3(c, []string{service}, nwProtocol, []string{service}, component, aI)
									appendData[0] = "undefined"
									appendSlice = append(appendSlice, appendData...)
								default:
									appendData := appendUndefined(component, nwProtocol, aI)
									appendSlice = append(appendSlice, appendData...)
								}
							} else {
								isMatchedL4Level := isMatchedServiceToL4level(s, service, aI.ServiceInfoAll)
								if isMatchedL4Level {
									data := appendDataEachComponent3(c, aI.Services, nwProtocol, []string{service}, component, aI)
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

func appendDataEachComponent2(component string, nwProtocol, services []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	switch component {
	case "dst_port":
		appendSlice = appendDataToDstPort(nwProtocol, c, aI)
	case "other_settings":
		appendSlice = judgeOtherSettingOutput(aI.ProxyMode, nwProtocol, aI.ProxyModeProtocol, aI.AvpName, aI.AVProfileInfoAll, aI)
	case "anti_virus":
		appendSlice = judgeAntiVirusOutput(aI.AvpName, nwProtocol, aI.AVProfileInfoAll)
	case "urlDomain":
		appendSlice = judgeURLDomainOutput(aI.WfpName, nwProtocol, aI.WfProfileInfoAll, aI.WfFilterInfoAll)
	case "expect":
		appendSlice = judgeExpectOutput(aI.Act, aI.WfpName, aI.AvpName, nwProtocol, aI.WfProfileInfoAll, aI.WfFilterInfoAll)
	case "protocol":
		appendSlice = append(appendSlice, nwProtocol...)
	case "src_port":
		appendSlice = appendDataToSrcPort(nwProtocol, c, aI)
	case "dst_nat_addr":
		appendSlice = appendStrDataToSlice(nwProtocol, c.ExtIP)
	case "dst_addr":
		appendSlice = appendStrDataToSlice(nwProtocol, c.MappedIP)
	}
	return appendSlice
}

func appendDataToSrcPort(nwProtocol []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	if c.MappedPort != "" {
		for range nwProtocol {
			appendSlice = append(appendSlice, "")
		}
	} else {
		srcport := handleSrcPort2(aI.Services, aI.ServiceInfoAll, aI.ServiceGrpInfoAll)
		appendSlice = append(appendSlice, srcport...)
	}
	return appendSlice
}

func appendDataToDstPort(nwProtocol []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	if c.MappedPort != "" {
		for _, protocol := range nwProtocol {
			if protocol == "icmp" {
				appendSlice = append(appendSlice, "")
			} else {
				appendSlice = append(appendSlice, c.MappedPort)
			}
		}
	} else {
		dstport := handleDstPort(aI.Services, aI.ServiceInfoAll, aI.ServiceGrpInfoAll, aI.ServicePort)
		appendSlice = append(appendSlice, dstport...)
	}
	return appendSlice
}

func handleExtPortWithRange(nwProtocol []string, c VipInfo) []string {
	var appendSlice []string
	// c.ExtPortはc.Protocolを有効にしている場合に指定可能
	// c.Protocolが無効な場合、c.ExtPortは""になる
	if strings.Contains(c.ExtPort, "-") {
		appendSlice = appendStrDataToSlice(nwProtocol, strings.Split(c.ExtPort, "-")[0])
	} else {
		appendSlice = appendStrDataToSlice(nwProtocol, c.ExtPort)
	}
	return appendSlice
}

func appendStrDataToSlice(nwProtocol []string, data string) (appendSlice []string) {
	for range nwProtocol {
		appendSlice = append(appendSlice, data)
	}
	return
}

func appendDataEachComponent4(component string, nwProtocol, services, protocol []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	switch component {
	case "dst_nat_port":
		appendSlice = handleExtPortWithRange(protocol, c)
	default:
		appendSlice = appendDataEachComponent2(component, nwProtocol, services, c, aI)
	}
	return appendSlice
}

func JudgeAppendData3WithMappedPort(component string, nwProtocol, services, ds, port []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	if len(port) > 1 {
		appendSlice = appendDataEachComponent3(c, services, nwProtocol, ds, component, aI)
		switch {
		case c.MappedPort == port[0] && c.MappedPort == port[1]:
			return appendSlice
		case c.MappedPort == port[0]:
			appendSlice[1] = "undefined"
		case c.MappedPort == port[1]:
			appendSlice[0] = "undefined"
		default:
			appendSlice = appendStrDataToSlice(nwProtocol, "undefined")
		}
	} else {
		if port[0] == c.MappedPort {
			appendSlice = appendDataEachComponent3(c, services, nwProtocol, aI.Services, component, aI)
		} else {
			appendSlice = appendStrDataToSlice(nwProtocol, "undefined")
		}
	}
	return appendSlice
}

func judgeAppendData3FromL4level(component, s, service string, nwProtocol []string, c VipInfo, aI AllInfo, num int) []string {
	var appendSlice []string
	isMatchedL4Level := isMatchedServiceToL4level(s, service, aI.ServiceInfoAll)
	if isMatchedL4Level {
		port := appendDataToDstPort(nwProtocol, c, aI)
		appendSlice = appendDataEachComponent3(c, aI.Services, nwProtocol, []string{s}, component, aI)
		if len(port) > 1 {
			// tcp&udpのサービス(例:DNS)とマッチした場合
			// tcpでマッチした場合: [1]をundefinedにする
			// udpでマッチした場合: [0]をundefinedにする
			if num != icmpInt {
				appendSlice[num] = "undefined"
			}
		}
	} else {
		appendSlice = append(appendSlice, "undefined")
	}
	return appendSlice
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
