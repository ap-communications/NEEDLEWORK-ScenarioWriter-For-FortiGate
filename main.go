package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate/app"
)

func main() {
	flag.Parse()
	var err error
	app.CurrentPath, err = os.Getwd()
	if err != nil {
		fmt.Printf("実行パスの取得に失敗しました = %s\n", err.Error())
		os.Exit(104)
	}

	fmt.Print("configディレクトリ内の`intf.csv`を読み込みます\n")
	intfcsv := readfile(app.CurrentPath + "/config/intf.csv")
	for i, v := range intfcsv {
		if i == 0 {
			// headerはskip
			continue
		} else {
			ifinfo, err := app.ValidateIntfInfo(v, i)
			if err != nil {
				// ValidateIntfInfoでos.Exitしているためここに処理が渡ることはない
				fmt.Printf("failed to ValidateIntfInfo due to %+v", err)
			}
			ifinfo.RandomIP = "1.1.1." + strconv.Itoa(i)
			app.IntfInfoAll = append(app.IntfInfoAll, ifinfo)
		}
	}
	fmt.Print("configディレクトリ内の`intf.csv`の読み込みが完了しました\n")

	if len(flag.Args()) < 1 {
		fmt.Printf("FortiGateのコンフィグファイル(.conf)のパスを指定してください\n")
		os.Exit(104)
	}
	fmt.Print("FortiGateのコンフィグファイル(.conf)を読み込みます\n")
	text := readfile(flag.Args()[0])
	m := map[string]bool{}
	app.AbsorbText(text, m)

	fmt.Print("CSVを生成します\n")
	app.WriteCSV(app.Scenarios)
	fmt.Print("CSVの生成が完了しました\n")
}

func readfile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`%s`の読み込みに失敗しました: %v\n", path, err)
		os.Exit(104)
	}
	defer f.Close()

	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", path, err)
		os.Exit(104)
	}
	return lines
}
