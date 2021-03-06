package app

import (
	"strings"
)

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

func appendDataEachComponent3(c VipInfo, services, nwProtocol, cs []string, component string, aI AllInfo) []string {
	var appendSlice []string
	switch component {
	case "dst_nat_port":
		appendSlice = appendDataForDstNatPort(aI.ServiceInfoAll, aI.ServiceGrpInfoAll, c, cs, nwProtocol, aI.ServicePort)
	default:
		appendSlice = appendDataEachComponent2(component, nwProtocol, services, cs, c, aI)
	}
	return appendSlice
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
			// DNS???????????????len(ports)=2?????????
			// DNS????????????TCP,UDP??????????????????????????????????????????????????????????????????
			nwProtocol := convertNWProtocol(serv, ServiceInfoAll)
			if len(ports) == len(nwProtocol) {
				portSlice = append(portSlice, ports...)
			} else {
				// ????????????FTP???????????????
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

// s.tcp?????????[]portassigninfo?????????????????????????????????????????????
// ??????????????????????????????????????????????????????
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

// s.tcp?????????[]portassigninfo?????????????????????????????????????????????
// ??????????????????????????????????????????????????????
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

// s.tcp?????????[]portassigninfo?????????????????????????????????????????????
// ??????????????????????????????????????????????????????
func HandleAssignDstPort(serv string, ServiceInfoAll []ServiceInfo, ServicePort map[string]string) []string {
	var dstPortSlice []string
	switch serv {
	case `"ALL"`, `"PING"`, `"ALL_ICMP"`, `"ALL_TCP"`, `"ALL_UDP"`:
		nwProtocol := convertNWProtocol(serv, ServiceInfoAll)
		for _, protocol := range nwProtocol {
			// protocol???????????????service????????????????????????append??????
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

// protocol???????????????service???????????????????????????
func getMatchedNWProtocolPort(serv string, ServiceInfoAll []ServiceInfo, ServicePort map[string]string) string {
	switch serv {
	case "icmp", "tcp", "udp":
		// ????????????????????????????????????(1-65535)??????ServicePort????????????????????????
		for k, v := range ServicePort {
			if serv == k {
				return v
			}
		}
	// ftpa???NEEDLEWORK???????????????????????????????????????????????????????????????????????????
	case "ftpa":
		serv = "ftp"
	}

	for _, s := range ServiceInfoAll {
		matchService := strings.ToLower(strings.Replace(s.Name, `"`, "", -1))
		switch serv {
		case "dns":
			if matchService == "dns" {
				// dns???len(getMatchedServiceDstPort(s)) = 2?????????
				// tcp,udp??????????????????????????????udp???[1]?????????
				return getMatchedServiceDstPort(s)[1]
			}
		case "dnst":
			if matchService == "dns" {
				// dns???len(getMatchedServiceDstPort(s)) = 2?????????
				// tcp,udp??????????????????????????????tcp???[0]?????????
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

func judgeMappedPortWithNum(component string, port, nwProtocol []string, c VipInfo, aI AllInfo, num int) []string {
	var appendData []string
	if port[num] == c.MappedPort {
		appendData = appendDataEachComponent4(component, nwProtocol, aI.Services, nwProtocol, c, aI)
		// tcp&udp?????????????????????????????????????????????????????????????????????????????????undefined?????????????????????????????????
		switch num {
		case 0:
			// tcp??????????????????????????????udp???undefined??????????????????
			appendData[1] = "undefined"
		case 1:
			// udp??????????????????????????????tcp???undefined??????????????????
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
		// tcp&udp???????????????????????????len(port)=2?????????
		// ??????????????????????????????????????????????????????c.MappedPort???????????????
		appendSlice = JudgeAppendData3WithMappedPort(component, nwProtocol, servs, []string{s}, port, c, aI)
	} else {
		appendData := appendDataEachComponent4a(component, nwProtocol, aI.Services, nwProtocol, []string{s}, c, aI)
		appendSlice = append(appendSlice, appendData...)
	}
	return appendSlice
}

func handleVInfoWithService(component, s string, nwProtocol, servs []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	// ??????????????????????????????
	// c.Service???ALL???????????????????????????
	// s???c.service???L4????????????????????????
	// c.MappedPort?????????????????????????????????serviceInfo???????????????????????????
	for _, vipService := range c.Service {
		nwProtocol := convertNWProtocol(s, aI.ServiceInfoAll)
		switch vipService {
		case `"ALL"`:
			appendData := appendDataEachComponent3(c, aI.Services, nwProtocol, []string{s}, component, aI)
			appendSlice = append(appendSlice, appendData...)
		case `"ALL_ICMP"`, `"PING"`:
			appendData := judgeAppendData3FromL4level(component, s, vipService, nwProtocol, c, aI, icmpInt)
			appendSlice = append(appendSlice, appendData...)
		case `"ALL_TCP"`:
			appendData := judgeAppendData3FromL4level(component, s, vipService, nwProtocol, c, aI, 1)
			appendSlice = append(appendSlice, appendData...)
		case `"ALL_UDP"`:
			appendData := judgeAppendData3FromL4level(component, s, vipService, nwProtocol, c, aI, 0)
			appendSlice = append(appendSlice, appendData...)
		default:
			if s == vipService {
				if c.MappedPort != "" {
					port := handleDstPort([]string{s}, aI.ServiceInfoAll, aI.ServiceGrpInfoAll, aI.ServicePort)
					// tcp&udp???????????????????????????len(port)=2?????????
					// ??????????????????????????????????????????????????????c.MappedPort???????????????
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

func appendDataEachComponent2(component string, nwProtocol, services, cs []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	switch component {
	case "dst_port":
		appendSlice = appendDataToDstPort(nwProtocol, cs, c, aI)
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
		appendSlice = appendDataToSrcPort(nwProtocol, cs, c, aI)
	case "dst_nat_addr":
		appendSlice = appendStrDataToSlice(nwProtocol, c.ExtIP)
	case "dst_addr":
		appendSlice = appendStrDataToSlice(nwProtocol, c.MappedIP)
	}
	return appendSlice
}

func appendDataToSrcPort(nwProtocol, cs []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	if c.MappedPort != "" {
		for range nwProtocol {
			appendSlice = append(appendSlice, "")
		}
	} else {
		var servs []string
		if cs != nil {
			// vip&c.Service&c.Mappedport?????????????????????
			servs = cs
		} else {
			servs = aI.Services
		}
		srcport := handleSrcPort2(servs, aI.ServiceInfoAll, aI.ServiceGrpInfoAll)
		appendSlice = append(appendSlice, srcport...)
	}
	return appendSlice
}

func appendDataToDstPort(nwProtocol, cs []string, c VipInfo, aI AllInfo) []string {
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
		var servs []string
		if cs != nil {
			// vip&c.Service&c.Mappedport?????????????????????
			servs = cs
		} else {
			servs = aI.Services
		}
		dstport := handleDstPort(servs, aI.ServiceInfoAll, aI.ServiceGrpInfoAll, aI.ServicePort)
		appendSlice = append(appendSlice, dstport...)
	}
	return appendSlice
}

func handleExtPortWithRange(nwProtocol []string, c VipInfo) []string {
	var appendSlice []string
	// c.ExtPort???c.Protocol?????????????????????????????????????????????
	// c.Protocol?????????????????????c.ExtPort???""?????????
	if strings.Contains(c.ExtPort, "-") {
		appendSlice = appendStrDataToSlice(nwProtocol, strings.Split(c.ExtPort, "-")[0])
	} else {
		appendSlice = appendStrDataToSlice(nwProtocol, c.ExtPort)
	}
	return appendSlice
}

func appendDataEachComponent4(component string, nwProtocol, services, protocol []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	switch component {
	case "dst_nat_port":
		appendSlice = handleExtPortWithRange(protocol, c)
	default:
		appendSlice = appendDataEachComponent2(component, nwProtocol, services, nil, c, aI)
	}
	return appendSlice
}

func appendDataEachComponent4a(component string, nwProtocol, services, protocol, cs []string, c VipInfo, aI AllInfo) []string {
	var appendSlice []string
	switch component {
	case "dst_nat_port":
		appendSlice = handleExtPortWithRange(protocol, c)
	default:
		appendSlice = appendDataEachComponent2(component, nwProtocol, services, cs, c, aI)
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

func judgeAppendData3FromL4level(component, s, vipService string, nwProtocol []string, c VipInfo, aI AllInfo, num int) []string {
	var appendSlice []string
	isMatchedL4Level := isMatchedServiceToL4level(s, vipService, aI.ServiceInfoAll)
	if isMatchedL4Level {
		port := appendDataToDstPort(nwProtocol, []string{s}, c, aI)
		appendSlice = appendDataEachComponent3(c, aI.Services, nwProtocol, []string{s}, component, aI)
		if len(port) > 1 {
			// tcp&udp???????????????(???:DNS)????????????????????????
			// tcp????????????????????????: [1]???undefined?????????
			// udp????????????????????????: [0]???undefined?????????
			if num != icmpInt {
				appendSlice[num] = "undefined"
			}
		}
	} else {
		appendSlice = append(appendSlice, "undefined")
	}
	return appendSlice
}
