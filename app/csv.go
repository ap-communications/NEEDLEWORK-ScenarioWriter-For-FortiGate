package app

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type Scenario struct {
	SrcFW            string           `json:"SrcFW"`
	SrcVLAN          string           `json:"SrcVLAN"`
	SrcAddress       []string         `json:"SrcAddress"`
	SrcNATAddress    string           `json:"SrcNATAddress"`
	SrcIntf          string           `json:"SrcIntf"`
	ReceiverPhysical string           `json:"ReceiverPhysical"`
	DstFW            string           `json:"DstFW"`
	DstVLAN          string           `json:"DstVLAN"`
	DstIntf          string           `json:"DstIntf"`
	Timeout          string           `json:"Timeout"`
	Try              string           `json:"Try"`
	Description      string           `json:"Description"`
	RelatedService   []relatedService `json:"relatedServices"`
}

type relatedService struct {
	Protocol      string `json:"Protocol"`
	SrcPort       string `json:"SrcPort"`
	DstNATAddress string `json:"DstNATAddress"`
	DstNATPort    string `json:"DstNATPort"`
	DstAddress    string `json:"DstAddress"`
	DstPort       string `json:"DstPort"`
	URLDomain     string `json:"URLDomain"`
	AntiVirus     string `json:"AntiVirus"`
	OtherSettings string `json:"OtherSettings"`
	Expect        string `json:"Expect"`
}

var (
	Scenarios   []Scenario
	csvHeader   = []string{"protocol", "src-fw", "src-vlan(option)", "src-ip", "src-port(option)", "src-nat-ip(option)", "s-if(option)", "is-receiver-physical(option)", "dst-fw", "dst-vlan(option)", "dst-nat-ip(option)", "dst-nat-port (option)", "dst-ip", "dst-port", "d-if(option)", "url/domain(option)", "anti-virus(option)", "timeout(option)", "try(option)", "other-settings(option)", "expect", "description"}
	CurrentPath string
)

func WriteCSV(scenarios []Scenario) {
	t := time.Now()
	strTime := t.Format("200601021504")
	filename := CurrentPath + "/scenario_" + strTime + ".csv"
	csvFile, err := os.Create(filename)
	if err != nil {
		fmt.Printf("CSVの生成に失敗しました = %s\n", err.Error())
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	writer.Write(csvHeader)
	writer.Flush()

	for _, data := range scenarios {
		appendScenario(writer, data)
	}
}

func appendScenario(writer *csv.Writer, data Scenario) {
	var baseRow []string
	if len(data.SrcAddress)*len(data.RelatedService) != 1 {
		for i := 0; i < len(data.SrcAddress)*len(data.RelatedService); i++ {
			if i == 0 {
				baseRow = genBaseRow(data, i)
				nonWriteFlag := handleDonotWriteRow(writer, baseRow)
				if nonWriteFlag {
					fmt.Printf("%sはundefinedが存在するため出力をスキップします\n", baseRow[21])
				}
			} else {
				copyRow := make([]string, len(baseRow))
				_ = copy(copyRow, baseRow)
				// AAA, AAB, ABA, ABB, ACA, ACB... の順で追加する
				if i < len(data.RelatedService) {
					copyRow[0] = data.RelatedService[i].Protocol
					copyRow[3] = data.SrcAddress[0]
					copyRow[4] = data.RelatedService[i].SrcPort
					copyRow[10] = data.RelatedService[i].DstNATAddress
					copyRow[11] = data.RelatedService[i].DstNATPort
					copyRow[12] = data.RelatedService[i].DstAddress
					copyRow[13] = data.RelatedService[i].DstPort
					copyRow[15] = data.RelatedService[i].URLDomain
					copyRow[16] = data.RelatedService[i].AntiVirus
					copyRow[19] = data.RelatedService[i].OtherSettings
					copyRow[20] = data.RelatedService[i].Expect
				} else {
					copyRow[0] = data.RelatedService[i%len(data.RelatedService)].Protocol
					copyRow[3] = data.SrcAddress[i/len(data.RelatedService)]
					copyRow[4] = data.RelatedService[i%len(data.RelatedService)].SrcPort
					copyRow[10] = data.RelatedService[i%len(data.RelatedService)].DstNATAddress
					copyRow[11] = data.RelatedService[i%len(data.RelatedService)].DstNATPort
					copyRow[12] = data.RelatedService[i%len(data.RelatedService)].DstAddress
					copyRow[13] = data.RelatedService[i%len(data.RelatedService)].DstPort
					copyRow[15] = data.RelatedService[i%len(data.RelatedService)].URLDomain
					copyRow[16] = data.RelatedService[i%len(data.RelatedService)].AntiVirus
					copyRow[19] = data.RelatedService[i%len(data.RelatedService)].OtherSettings
					copyRow[20] = data.RelatedService[i%len(data.RelatedService)].Expect
				}
				handleDonotWriteRow(writer, copyRow)
			}
		}
	} else {
		baseRow = genBaseRow(data, 0)
		nonWriteFlag := handleDonotWriteRow(writer, baseRow)
		if nonWriteFlag {
			fmt.Printf("%sはundefinedが存在するため出力をスキップします\n", baseRow[21])
		}
	}
}

func handleDonotWriteRow(writer *csv.Writer, row []string) bool {
	var nonWriteFlag bool
	for _, data := range row {
		if data == `undefined` {
			nonWriteFlag = true
		}
	}

	if !nonWriteFlag {
		writer.Write(row)
		writer.Flush()
	}
	return nonWriteFlag
}

func genBaseRow(data Scenario, i int) []string {
	var baseRow []string
	baseRow = append(baseRow, data.RelatedService[0].Protocol)
	baseRow = append(baseRow, data.SrcFW)
	baseRow = append(baseRow, data.SrcVLAN)
	baseRow = append(baseRow, data.SrcAddress[i])
	baseRow = append(baseRow, data.RelatedService[0].SrcPort)
	baseRow = append(baseRow, data.SrcNATAddress)
	baseRow = append(baseRow, data.SrcIntf)
	baseRow = append(baseRow, data.ReceiverPhysical)
	baseRow = append(baseRow, data.DstFW)
	baseRow = append(baseRow, data.DstVLAN)
	baseRow = append(baseRow, data.RelatedService[0].DstNATAddress)
	baseRow = append(baseRow, data.RelatedService[0].DstNATPort)
	baseRow = append(baseRow, data.RelatedService[0].DstAddress)
	baseRow = append(baseRow, data.RelatedService[0].DstPort)
	baseRow = append(baseRow, data.DstIntf)
	baseRow = append(baseRow, data.RelatedService[0].URLDomain)
	baseRow = append(baseRow, data.RelatedService[0].AntiVirus)
	baseRow = append(baseRow, data.Timeout)
	baseRow = append(baseRow, data.Try)
	baseRow = append(baseRow, data.RelatedService[0].OtherSettings)
	baseRow = append(baseRow, data.RelatedService[0].Expect)
	baseRow = append(baseRow, data.Description)
	return baseRow
}
