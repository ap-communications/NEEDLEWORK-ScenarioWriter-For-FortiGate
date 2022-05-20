package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validIPMaskNums = []string{"0", "128", "192", "224", "240", "248", "252", "254", "255"}
)

func ValidateIntfInfo(str string, i int) (IntfInfo, error) {
	var intfInfo IntfInfo
	var strID string
	switch len(strings.Split(str, ",")) {
	case 3:
		strID = "0"
	case 4:
		strID = strings.Split(str, ",")[3]
		if strID == "" {
			// 未指定の場合は「0」を補完する
			strID = "0"
		}
	default:
		fmt.Print("正しいフォーマットでconfigディレクトリ内の`intf.csv`を記載してください\n")
		os.Exit(101)
	}

	// IPアドレスのバリデーションチェック
	address := strings.Split(str, ",")[1]
	isValidIPAddress := validateIPAddress(address)
	if !isValidIPAddress {
		fmt.Printf("configディレクトリ内の`intf.csv`の%+v行目に記載している\"%s\"は不正なIPアドレスです\n", i, address)
		fmt.Print("A.B.C.Dの形式かつ正しいIPアドレスを指定してください\n")
		os.Exit(102)
	}

	// サブネットマスクのバリデーションチェック
	subnet := strings.Split(str, ",")[2]
	isValidIPMask := validateIPMask(subnet)
	if !isValidIPMask {
		fmt.Printf("configディレクトリ内の`intf.csv`の%+v行目に記載している\"%s\"は不正なサブネットマスクです\n", i, subnet)
		fmt.Print("A.B.C.Dの形式かつ正しいサブネットマスクを指定してください\n")
		os.Exit(102)
	}

	// VLAN IDが数字の文字列であるか確認
	intID, err := strconv.Atoi(strID)
	if err != nil {
		fmt.Printf("configディレクトリ内の`intf.csv`の%+v行目に記載している\"%s\"は不正なVLAN IDです\n", i, strID)
		fmt.Print("0~4094をVLAN IDに指定してください\n")
		os.Exit(102)
	}

	intfInfo = IntfInfo{
		Name:       `"` + strings.Split(str, ",")[0] + `"`,
		Address:    address,
		SubnetMask: subnet,
		VLANID:     intID,
	}

	// VLAN IDのバリデーションチェック
	validate := validator.New()
	err = validate.Struct(intfInfo)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			switch fieldName {
			case "VLANID":
				fmt.Printf("configディレクトリ内の`intf.csv`の%+v行目に記載している\"%s\"は不正なVLAN IDです\n", i, strID)
				fmt.Print("0~4094をVLAN IDに指定してください\n")
				os.Exit(102)
			}
		}
	}
	return intfInfo, nil
}

func validateIPAddress(address string) bool {
	addressNums := strings.Split(address, ".")
	if len(addressNums) != 4 {
		return false
	}

	for _, v := range addressNums {
		i, err := strconv.Atoi(v)
		if err != nil {
			return false
		}
		switch {
		case i > 255 || i < 0:
			return false
		}
	}
	return true
}

func validateIPMask(subnet string) bool {
	subnetNums := strings.Split(subnet, ".")
	if len(subnetNums) != 4 {
		return false
	}

	for i, v := range subnetNums {
		b := checkIPMaskNum(i, v)
		if !b {
			return false
		}
	}
	return true
}

func checkIPMaskNum(i int, v string) bool {
	for _, n := range validIPMaskNums {
		if v == n {
			if i == 0 && v == "0" {
				fmt.Print("サブネットマスクの第一オクテットは0以外を指定してください\n")
				return false
			}
			return true
		}
	}
	return false
}
