package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func HandleProtocolOutput(allInfo AllInfo) []string {
	var protocolSlice []string
	if allInfo.VInfo != nil {
		protocolSlice = handleVInfo("protocol", allInfo)
	} else {
		servs := judgeService(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll)
		for _, serv := range servs {
			nwProtocol := convertNWProtocol(serv, allInfo.ServiceInfoAll)
			protocolSlice = append(protocolSlice, nwProtocol...)
		}
	}
	return protocolSlice
}

func handleFWIPOutput(intf, env string, IntfInfoAll []IntfInfo) string {
	intfInfo := confirmIntfWithIntfInfo(intf, env, IntfInfoAll)
	return intfInfo.Address
}

func handleVLANIDOutput(intf, env string, intfInfoAll []IntfInfo) string {
	intfInfo := confirmIntfWithIntfInfo(intf, env, IntfInfoAll)
	return strconv.Itoa(intfInfo.VLANID)
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

		// addrのいずれかでflagsがtrueの場合はFQDNFlagをtrueにする
		if flags {
			srcFQDNFlag = true
		}
		addrSlice = append(addrSlice, aSlice...)
	}
	return addrSlice, srcFQDNFlag
}

func HandleSrcPortOutput(allInfo AllInfo) []string {
	var SrcPortSlice []string
	if allInfo.VInfo != nil {
		SrcPortSlice = handleVInfo("src_port", allInfo)
	} else {
		SrcPortSlice = handleSrcPort(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll, allInfo.ServicePort)
	}
	return SrcPortSlice
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

func HandleDstNATAddrOutput(allInfo AllInfo) []string {
	var dstNATAddrSlice []string
	if allInfo.VInfo != nil {
		dstNATAddrSlice = handleVInfo("dst_nat_addr", allInfo)
	} else {
		services := judgeService(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll)
		for _, serv := range services {
			nwProtocol := convertNWProtocol(serv, allInfo.ServiceInfoAll)
			dstNATAddrSlice = append(dstNATAddrSlice, appendStrDataToSlice(nwProtocol, "")...)
		}
	}
	return dstNATAddrSlice
}

func HandleDstNATPortOutput(allInfo AllInfo) []string {
	var dstNATPortSlice []string
	if allInfo.VInfo != nil {
		dstNATPortSlice = handleVInfo("dst_nat_port", allInfo)
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

func HandleDstAddressOutput(intf string, addrs []string, AddressInfoAll []AddressInfo, AddrGrpInfoAll []AddrGrpInfo, SRouteInfoAll []SRouteInfo, IntfInfoAll []IntfInfo, usedAddress []string, allInfo AllInfo) ([]string, bool) {
	var dstAddresss []string
	var dstFQDNFlag bool
	protocols := HandleProtocolOutput(allInfo)
	if allInfo.VInfo != nil {
		dstAddresss = handleVInfo("dst_addr", allInfo)
	} else {
		var aSlice []string
		aSlice, dstFQDNFlag = getUniqueDstAddr(intf, addrs, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, allInfo)
		for _, v := range aSlice {
			for range protocols {
				dstAddresss = append(dstAddresss, v)
			}
		}
	}
	// VInfoの場合はdstAddressがFQDNになることはないので、dstFQDNFlagは初期値で問題ない
	return dstAddresss, dstFQDNFlag
}

func HandleDstPortOutput(allInfo AllInfo) []string {
	var dstPorts []string
	if allInfo.VInfo != nil {
		dstPorts = handleVInfo("dst_port", allInfo)
	} else {
		dstPorts = handleDstPort(allInfo.Services, allInfo.ServiceInfoAll, allInfo.ServiceGrpInfoAll, allInfo.ServicePort)
	}
	return dstPorts
}

func HandleURLDomainOutput(allInfo AllInfo) []string {
	var urlSlice []string
	if allInfo.VInfo != nil {
		urlSlice = handleVInfo("urlDomain", allInfo)
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

func HandleAntiVirusOutput(allInfo AllInfo) []string {
	var avSlice []string
	if allInfo.VInfo != nil {
		avSlice = handleVInfo("anti_virus", allInfo)
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

func HandleOtherSettingOutput(allInfo AllInfo) []string {
	var otherSettings []string
	if allInfo.VInfo != nil {
		otherSettings = handleVInfo("other_settings", allInfo)
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

func HandleExpectOutput(allInfo AllInfo) []string {
	var actSlice []string
	if allInfo.VInfo != nil {
		actSlice = handleVInfo("expect", allInfo)
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

func handleRelatedServiceOutPut(p policyInfo, sI []ServiceInfo, ServiceGrpInfoAll []ServiceGrpInfo, IntfInfoAll []IntfInfo, allInfo AllInfo) ([]relatedService, bool) {
	var relatedServices []relatedService
	dstAddress, dstFQDNFlag := HandleDstAddressOutput(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, allInfo)
	protocols := HandleProtocolOutput(allInfo)
	srcPorts := HandleSrcPortOutput(allInfo)
	dstNATAddress := HandleDstNATAddrOutput(allInfo)
	dstNATPorts := HandleDstNATPortOutput(allInfo)
	dstPorts := HandleDstPortOutput(allInfo)
	urlDomains := HandleURLDomainOutput(allInfo)
	antiVirus := HandleAntiVirusOutput(allInfo)
	otherSettings := HandleOtherSettingOutput(allInfo)
	expects := HandleExpectOutput(allInfo)

	if allInfo.VInfo == nil {
		protocols = appendDataToComponentWithDstAddr(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, protocols, allInfo)
		srcPorts = appendDataToComponentWithDstAddr(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, srcPorts, allInfo)
		dstNATAddress = appendDataToComponentWithDstAddr(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, dstNATAddress, allInfo)
		dstNATPorts = appendDataToComponentWithDstAddr(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, dstNATPorts, allInfo)
		dstPorts = appendDataToComponentWithDstAddr(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, dstPorts, allInfo)
		urlDomains = appendDataToComponentWithDstAddr(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, urlDomains, allInfo)
		antiVirus = appendDataToComponentWithDstAddr(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, antiVirus, allInfo)
		otherSettings = appendDataToComponentWithDstAddr(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, otherSettings, allInfo)
		expects = appendDataToComponentWithDstAddr(p.Dstintf, p.DstAddress, AddressInfoAll, AddrGrpInfoAll, SRouteInfoAll, IntfInfoAll, usedAddress, expects, allInfo)
	}

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
			Protocol:      protocols[i/quantityURL],
			SrcPort:       srcPorts[i/quantityURL],
			DstNATAddress: dstNATAddress[i/quantityURL],
			DstNATPort:    dstNATPorts[i/quantityURL],
			DstAddress:    dstAddress[i/quantityURL],
			DstPort:       dstPorts[i/quantityURL],
			URLDomain:     urlDomains[i],
			AntiVirus:     antiVirus[i/quantityURL],
			OtherSettings: otherSettings[i/quantityURL],
			Expect:        expects[i],
		}
		relatedServices = append(relatedServices, rS)
	}

	ok := checkDatasLength(relatedServices, protocols, srcPorts, dstNATAddress, dstAddress, dstNATPorts, dstPorts, antiVirus, otherSettings)
	if !ok {
		fmt.Printf("予期せぬエラーが発生しました\n")
		os.Exit(254)
	}
	return relatedServices, dstFQDNFlag
}

func genScenario(p policyInfo) Scenario {
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
	rS, dstFQDNFlag := handleRelatedServiceOutPut(p, ServiceInfoAll, ServiceGrpInfoAll, IntfInfoAll, allInfo)
	srcNATAddress, ipPoolRange := HandleSrcNATAddressOutPut(p.NAT, p.PoolName, p.Dstintf, "", IPPoolInfoAll, IntfInfoAll)

	s := Scenario{
		SrcFW:            handleFWIPOutput(p.Srcintf, "", IntfInfoAll),
		SrcVLAN:          handleVLANIDOutput(p.Srcintf, "", IntfInfoAll),
		SrcAddress:       srcAddress,
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
