package app

import (
	"fmt"
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
			// FortiGuardカテゴリベースのフィルタが有効になっている場合に含まれる
			// FortiGuardカテゴリベースのフィルタが有効かつライセンスが無い状態だとページがブロックされる
			if flags && strings.Contains(v, "unset options") {
				wfPInfo.FtgdFilter = true
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
				// FortiGateのversionが6.0系(例:6.0.6)の場合、ポリシー＆VDOMの設定で判定する
				// FortiGateのversionが6.2系(例:6.2.6)の場合、ポリシーの設定のみで判定する
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
						s := genScenario(p)
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
