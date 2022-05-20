# 目次

- [目次](#目次)
- [ツールの仕様](#ツールの仕様)
  - [ツールの動作概要](#ツールの動作概要)
  - [各データの取得について](#各データの取得について)
    - [プロトコルについて](#プロトコルについて)
    - [送信元, 宛先FW IP, VLAN IDについて](#送信元-宛先fw-ip-vlan-idについて)
    - [送信元, 宛先IPアドレスについて](#送信元-宛先ipアドレスについて)
    - [送信元, 宛先ポート](#送信元-宛先ポート)
    - [送信元NAT後IPについて](#送信元nat後ipについて)
    - [宛先NAT前IPについて](#宛先nat前ipについて)
    - [宛先NAT前ポートについて](#宛先nat前ポートについて)
    - [URL/ドメインについて](#urlドメインについて)
    - [アンチウイルスについて](#アンチウイルスについて)
    - [その他設定について](#その他設定について)
    - [期待値について](#期待値について)
    - [説明について](#説明について)
    - [その他のデータについて](#その他のデータについて)
  - [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)
  - [グループ要素の取得について](#グループ要素の取得について)
  - [送信元・宛先やサービスが複数指定されているポリシーの取得について](#送信元宛先やサービスが複数指定されているポリシーの取得について)
  - [出力するシナリオ数について](#出力するシナリオ数について)
  - [出力されるエラーメッセージについて](#出力されるエラーメッセージについて)
  - [その他出力されるメッセージについての補足](#その他出力されるメッセージについての補足)
  - [留意事項](#留意事項)

# ツールの仕様

以下にツールの仕様を記載します。

## ツールの動作概要

* ツールは以下の順序で動作します。
  1. コンフィグファイル(.conf)およびconfigディレクトリ内の`intf.csv`の読み込み
  2. シナリオの各データの取得
  3. シナリオの出力
* シナリオの各データは以下の設定値(名)を元に取得します。
  * プロトコル: 各ポリシーのサービスの設定値
  * 送信元FW IP: 各ポリシーの着信インターフェースの設定名 ★
  * 宛先FW IP: 各ポリシーの発信インターフェースの設定名 ★
  * 送信元VLAN ID: 各ポリシーの着信インターフェースの設定名 ★
  * 送信元IP: 各ポリシーの送信元の設定値
  * 送信元ポート: 各ポリシーのサービスの設定値
  * 送信元NAT IP: 各ポリシーのNATの設定値
  * 宛先VLAN ID: 各ポリシーの発信インターフェースの設定名 ★
  * 宛先IP: 各ポリシーの宛先の設定値
  * 宛先ポート: 各ポリシーのサービスの設定値
  * 宛先NAT IP: 各ポリシーの宛先の設定値
  * 宛先NAT前ポート: 各ポリシーの宛先&サービスの設定値
  * URL/domain: 各ポリシーのWebフィルタの設定値
  * 期待値: 各ポリシーのアクションの設定値
  * その他設定: VDOM、または各ポリシーのインスペクションモードの設定値 <br>
※ ★はconfigディレクトリ内の`intf.csv`から値を読み込みます。<br>
（指定したコンフィグファイル(.conf)からはインターフェース名以外は読み込みません）
* 各データの取得についての詳細は[各データの取得について](##各データの取得について)をご確認ください。
  * ポリシーの宛先がバーチャルIP/サーバーまたはVIPグループの場合は <br>[宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。
  * ポリシーの送信元・宛先やサービスがグループの場合は[グループ要素の取得について](##グループ要素の取得について)をご確認ください。
  * ポリシーの送信元・宛先やサービスが複数指定されている場合は <br>[アドレスやサービスが複数指定されている場合について](##アドレスやサービスが複数指定されている場合について)をご確認ください。
* 出力するシナリオ数についての詳細は[出力するシナリオ数について](##出力するシナリオ数について)をご確認ください。
* エラーによりシナリオが出力できなかった場合は、その理由が出力されます。
詳細は[出力されるエラーについて](##出力されるエラーについて)をご確認ください。

## 各データの取得について

以下にシナリオの各データの取得についての詳細を記載します。

### プロトコルについて

* シナリオに出力するプロトコルは下記の通りです。 <br>
※ NEEDLEWORKがバージョン12.0.0までに対応している、全てのプロトコルでシナリオの出力が可能です。
  * icmp
  * tcp
  * udp
  * dns(t)
  * http(s)
  * ftp(a)
  * imap
  * smtp
* ポリシーのサービスを元に、以下に倣ってプロトコルを取得します。
  * 宛先がバーチャルIP/サーバー,またはVIPグループではない場合
    * サービスに"ALL"が指定されている場合は、上記のプロトコルを全て出力します。
    * （サービスグループ内のサービス含め）下記に該当するサービスが指定された場合は、<br>右側に記載しているプロトコルをそれぞれ取得します。
      * "DNS": dnst, dns
      * "FTP": ftp, ftpa
      * "HTTP": http
      * "HTTPS": https
      * "IMAP": imap
      * "SMTP": smtp
    * 上記に該当しないサービスは、サービスに割り当てられているプロトコルを元に出力します。
      * コンフィグファイル(.conf)例（一部抜粋）:
    ```
    config firewall service custom
      edit "ALL_UDP"
          set category "General"
          set udp-portrange 1-65535
      next
      edit "PING"
          set category "Network Services"
          set protocol ICMP
          set icmptype 8
          unset icmpcode
      next
      edit "NTP"
          set category "Network Services"
          set tcp-portrange 123
          set udp-portrange 123
      next
    end
    ``` 
      * 出力例:
      ```
      "PING": icmp
      "NTP": tcp, udp
      "ALL_UDP": udp
      ``` 
  * 宛先がバーチャルIP/サーバーまたはVIPグループの場合
    * [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。

### 送信元, 宛先FW IP, VLAN IDについて

* ポリシーの送信、宛先インターフェースを元に、以下に倣って送信元, 宛先FW IP, VLAN IDを取得します。
  * configディレクトリ内の`intf.csv`に該当するインターフェースの情報が記載されているか確認します。
  * 該当するインターフェースの情報がconfigディレクトリ内の`intf.csv`に記載されている場合、その値を元に出力します。
  * 該当するインターフェースの情報がconfigディレクトリ内の`intf.csv`に記載されていない場合、エラーが出力されます。<br>
  ※ シナリオは生成されません。<br>
  以下にconfigディレクトリ内の`intf.csv`の記載例およびポリシー例に応じた出力例を記載します。<br>
  configディレクトリ内の`intf.csv`の記載例:
  ```
  inftname,address,subnetmask,vlanid
  lan,5.6.7.8,255.255.255.255,10
  lan2,172.16.0.20,255.255.255.0,
  wan1,10.0.0.10,255.255.255.0,
  ```
  ポリシー例（一部抜粋）: 
  ```
  set name "test2"
  set uuid c1adaa9c-5945-51ea-7e74-ff9fb0bd13e7
  set srcintf "lan"
  set dstintf "wan1"
  (略)
  ```
  出力例:
  ```
  送信元FW IP:　5.6.7.8
  送信元VLAN ID: 10
  宛先FW IP: 10.0.0.10
  宛先VLAN ID: 0
  ```

### 送信元, 宛先IPアドレスについて

ポリシーの送信元や宛先を元に、以下に倣って送信元, 宛先IPアドレスを取得します。

* 共通
  * サブネット
    * /32: そのIPアドレスを取得します。
    * /32以外: そのネットワーク内のネットワークアドレスとブロードキャストアドレスを除く、 <br>最初と最後のIPアドレスを取得します。 <br>
    例:
    * 10.10.10.0/24の場合、`10.10.10.1`, `10.10.10.254`を取得します。
  * IP範囲
    * IP範囲で指定した最初と最後のIPアドレスを取得します。 <br>
    例:
    * 10.10.20.1-10.10.20.20の場合、`10.10.20.1`, `10.10.20.20`を取得します。
  * all
    * スタティックルートやconfigディレクトリ内の`intf.csv`に情報がある場合、そのネットワーク内の全てのIPアドレスから、<br>ポリシーで使用しているIPアドレスを除いた始めのIPアドレスを取得します。<br>
    以下の場合、`10.0.0.2`を取得します。<br>
    * 例: 
   ルーティング: 10.0.0.0/24
   ポリシーで使用しているIP: 10.0.0.1/32, 10.0.0.16/28
    * スタティックルートやconfigディレクトリ内の`intf.csv`に情報がない場合、<br>インターフェース毎に固有のIPアドレスを取得します。
    * 例: lan2, wan2のスタティックルートやconfigディレクトリ内の`intf.csv`に情報がない場合の出力例
    ```
    lan2: 1.1.1.1
    wan2: 1.1.1.2
    ```
  * FQDN
    * 該当するFQDNを取得します。
    * **このシナリオでテストを実行する際は、FortiGateのDNSサーバーにNEEDLEWORKの管理IPを指定してください。** <br>
    ※ 指定しない場合は正常にテストが行えません。
* 宛先
   * バーチャルIP/サーバー, VIPグループ
   [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。

### 送信元, 宛先ポート

* ポリシーのサービスを元に、以下に倣って送信元, 宛先ポートを取得します。
* 宛先がバーチャルIP/サーバーやVIPグループではない場合
  * 送信元ポートが指定されていない場合は空文字をデータとして使用します。
  * 基本的にポリシーのサービスで指定しているポートを取得します。
  * 上記2つに該当しない場合について、以下に記載します。
    * ポリシーのサービスが"ALL"の場合、以下を取得します。
      * 送信元: 空文字をデータとして使用します。
      * 宛先:　以下を取得します。
        * icmp: 空文字
        * tcp: 80
        * udp: 53
        * その他プロトコル: 各サービスで指定しているポート
    * ポリシーのサービスがポートを範囲指定している場合、開始ポートのみを取得します。 <br>
    （範囲指定した全ポートを出力しません） <br>
      例：
      `set xxx-portrange 1000-2000`の場合、`1000`を出力
    * ポリシーのサービスがポートを複数指定している場合、始めに設定されているポートのみを取得します。 <br>
    （複数指定した全ポートを出力しません） <br>
      例1：`set xxx-portrange 100 200 300`の場合 <br>
      送信元: 空文字を出力 <br>
      宛先: `100`を出力 <br>
      例2: `set xxx-portrange 123:11122-11125 123:11155`の場合 <br>
      送信元: `11122`を出力 <br>
      宛先: `123`を出力 <br>
* 宛先がバーチャルIP/サーバーまたはVIPグループの場合
  * [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。

### 送信元NAT後IPについて

ポリシーのNATを元に、以下に倣って送信元NAT後IPを取得します。
* 宛先がバーチャルIP/サーバーやVIPグループではない場合
  * 送信元NAT有効の場合
    * 発信インターフェースのIP: configディレクトリ内の`intf.csv`を元に発信インターフェースのIPアドレスを取得します。 
    * ダイナミックIPプールを使う: 外部IPアドレス/範囲の始めのIPアドレスを取得します。 <br>
    以下コンフィグファイル(.conf)例の場合、`100.100.100.100`を取得します。 <br>
    コンフィグファイル(.conf)例（一部抜粋）:
    ```
    config firewall ippool
      edit "test_ippool"
          set startip 100.100.100.100
          set endip 100.100.100.105
      next
    end
    ```
  * 送信元NAT無効の場合 
    * 空文字をデータとして使用します。
* 宛先がバーチャルIP/サーバーまたはVIPグループの場合
  * [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。

### 宛先NAT前IPについて

ポリシーの宛先を元に、以下に倣って宛先NAT前IPを取得します。
* 宛先がバーチャルIP/サーバーまたはVIPグループではない場合
  * 空文字をデータとして使用します。
* 宛先がバーチャルIP/サーバーまたはVIPグループの場合
  * [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。

### 宛先NAT前ポートについて

ポリシーの宛先を元に、以下に倣って宛先NAT前ポートを取得します。
* 宛先がバーチャルIP/サーバーまたはVIPグループではない場合 
  * 空文字をデータとして使用します。
* 宛先がバーチャルIP/サーバーまたはVIPグループの場合
  * [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。

### URL/ドメインについて

ポリシーのウェブフィルタを元に、以下に倣ってURL/ドメインを取得します。
* 宛先がバーチャルIP/サーバーまたはVIPグループではない場合
  * ウェブフィルタが有効な場合 
    * 出力するシナリオのプロトコルがhttp(s)の場合
      * プロファイルで指定されたURLを取得します。
    * 出力するシナリオのプロトコルが上記以外の場合
      * 空文字をデータとして使用します。
  * ウェブフィルタが無効な場合 
    * 空文字をデータとして使用します。
* 宛先がバーチャルIP/サーバーまたはVIPグループの場合
  * [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。

### アンチウイルスについて

ポリシーのアンチウイルスを元に、以下に倣ってアンチウィルスのシナリオパラメーターを取得します。<br>
* 宛先がバーチャルIP/サーバーまたはVIPグループではない場合
  * 出力するシナリオのプロトコルがhttp(s), ftp(a), imap, smtpの場合 
  	* アンチウイルスが該当プロトコルで無効な場合
      * 空文字をデータとして使用します。
  	* アンチウイルスが該当プロトコルで有効な場合 
    	* `enable`を取得します。
  	* 以下にコンフィグファイル(.conf)例と出力例を記載します。<br>
  	コンフィグファイル(.conf)例（一部抜粋）:
      ```
      config antivirus profile
          edit "default"
              config http
                  set options scan avmonitor
              end
              config ftp
                  set options scan avmonitor
              end
              config imap
                  set options scan
                  set executables virus
              end
              config pop3
                  set options scan avmonitor
                  set executables virus
              end
              config smtp
                  set executables virus
              end
              set scan-mode legacy
          next
          edit "wifi-default"
              set comment "Default configuration for offloading WiFi traffic."
              config http
                  set options scan
              end
              config ftp
                  set options scan
              end
              config imap
                  set options scan
                  set executables virus
              end
              config pop3
                  set options scan
                  set executables virus
              end
              config smtp
                  set options scan
                  set executables virus
              end
              set scan-mode legacy
          next
      end
      ```
      出力例1: ポリシーのアンチウイルスで"default"が適用されている場合
      ```
      * 出力するシナリオのプロトコルがhttp(s)の場合: 
      * 空文字をデータとして使用します。
      * 出力するシナリオのプロトコルがftp(a)の場合: 
      * 空文字をデータとして使用します。
      * 出力するシナリオのプロトコルがimapの場合: `enable`を取得します。
      * 出力するシナリオのプロトコルがsmtpの場合: 
      * 空文字をデータとして使用します。
      * 出力するシナリオのプロトコルが上記以外の場合: `enable`を取得します。
      ```
      出力例2: ポリシーのアンチウイルスで"wifi-default"が適用されている場合
      ```
      * 出力するシナリオのプロトコルがhttp(s)の場合: `enable`を取得します。
      * 出力するシナリオのプロトコルがftp(a)の場合: `enable`を取得します。
      * 出力するシナリオのプロトコルがimapの場合: `enable`を取得します。
      * 出力するシナリオのプロトコルがsmtpの場合: `enable`を取得します。
      * 出力するシナリオのプロトコルが上記以外の場合: 
      * 空文字をデータとして使用します。
      ```
  * 出力するシナリオのプロトコルが上記以外の場合 
    * 空文字をデータとして使用します。
* 宛先がバーチャルIP/サーバーまたはVIPグループの場合
  * [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。

### その他設定について

以下の基準を元に、その他設定を取得します。
* 宛先がバーチャルIP/サーバーまたはVIPグループではない場合
  * インスペクションモードがプロキシモードの場合
    * 現在、下記の基準A〜Dの全てに該当している場合,`Proxy mode`を取得します。<br>
        **基準A:**
        * ポリシーまたはVDOMのインスペクションモードがプロキシベースである<br>
        **基準B:**
        * ポリシーのセキュリティプロファイルでアンチウイルスが有効である<br>
        **基準C:**
        * 出力するシナリオのプロトコルが`tcp, dnst, http(s), ftp(a), imap, smtp`である<br>
        **基準D:**
        * プロトコルオプションで有効なポートが宛先ポートとなるサービス(例: `FTP`)をポリシーで指定している<br>
        ※ プロトコルオプションの`default`で有効なポートは`"21", "25", "53", "80", "110", "119", "135", "143", "445"`です。
    * 上記以外の場合
      * 空文字をデータとして使用します。
  * インスペクションモードがフローモードの場合 
    * 空文字をデータとして使用します。
* 宛先がバーチャルIP/サーバーまたはVIPグループの場合
  * [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。

### 期待値について

ポリシーのアクションを元に、以下を上から優先的に確認し、該当するデータを取得します。
* 宛先がバーチャルIP/サーバーまたはVIPグループではない場合
  * URLフィルタまたはアンチウイルスでblockとなっている場合
    * `block`を取得します。
  * ポリシーのアクションがacceptとなっている場合
    * `pass`を取得します。
  * 上記以外の場合
    * `drop`を取得します。
* 宛先がバーチャルIP/サーバーまたはVIPグループの場合
  * [宛先がバーチャルIP/サーバーまたはVIPグループの場合について](#宛先がバーチャルipサーバーまたはvipグループの場合について)をご確認ください。
  

### 説明について

* 以下の2つを元に、データを取得します。
  * ポリシーの送信元または宛先にFQDNが使用されているか否か
  * 送信元NAT IPアドレスにダイナミックIPプールが使用されているか否か<br>

* 送信元NAT IPアドレスにダイナミックIPプールが使用されている＆FQDN指定有りのポリシーの場合 
  * `policy name = ポリシー名 src_nat_ip=a.b.c.d-A.B.C.D FQDN Policy`を取得します。
* 送信元NAT IPアドレスにダイナミックIPプールが使用されているポリシーの場合 
  * `policy name = ポリシー名 src_nat_ip=a.b.c.d-A.B.C.D`を取得します。
* FQDN指定有りのポリシーの場合 
  * `policy name = ポリシー名 FQDN Policy`を取得します。
* 上記該当しないポリシーの場合
  * `policy name = ポリシー名`を取得します。

### その他のデータについて

以下のデータは、ポリシーからデータの取得を行いません。
* 送信元インターフェイス
* 実機
* 宛先インターフェイス
* 最大実行回数
* タイムアウト

ヘッダーのみ出力し、ヘッダー以外のデータは空文字を取得および出力します。

## 宛先がバーチャルIP/サーバーまたはVIPグループの場合について

宛先がバーチャルIP/サーバーまたはVIPグループの場合、以下の取得データが`undefined`となる場合があります。
* プロトコル
* 宛先NAT前IP
* 宛先NAT前ポート
* 宛先IPアドレス
* 宛先ポート
* 送信元ポート
* URL/ドメイン
* アンチウイルス
* その他設定
* 期待値

以下にそれぞれのパターンのデータの取得について記載します。
* バーチャルIP/サーバーのオプションでサービスを指定している場合
  * 指定されているサービスとポリシーのサービスを元に各データを取得します。
  * プロトコルの取得例:
    | &nbsp;                 |           | VIPのサービス                                                 |            |           |           |
    | ---------------------- | --------- | ------------------------------------------------------------- | ---------- | --------- | --------- |
    |                        |           | "ALL"                                                         | "ALL_ICMP" | "NTP"     | "FTP"     |
    | **ポリシーのサービス** | "ALL"     | icmp, tcp, udp, dnst, dns, http, https, ftp, ftpa, imap, smtp | icmp       | tcp,udp   | ftp,ftpa  |
    | ^                      | "PING"    | icmp                                                          | icmp       | undefined | undefined |
    | ^                      | "NTP"     | tcp,udp                                                       | undefined  | tcp,udp   | undefined |
    | ^                      | "ALL_TCP" | tcp                                                           | undefined  | tcp       | tcp       |
  * その他のデータの取得:
  * 下記に該当している場合は、`undefined`を取得します。
    * 指定されているサービスまたはポリシーのサービスが"ALL_〇〇"または"PING"
      * もう一方のサービスとL4レベルで一致していない場合
    * 指定されているサービスまたはポリシーのサービスが"ALL"、"ALL_〇〇"、"PING"以外
      * もう一方のサービスが"ALL"、"ALL_〇〇"、"PING"のいずれかでかつL4レベルで一致していない場合
      * もう一方のサービスが"ALL"、"ALL_〇〇"、"PING"以外でかつL7レベルで一致していない場合
  * 上記に該当しない場合は右に記載のデータを取得します。
    * 宛先NAT前IP: 外部IPアドレス/範囲で指定されたIPアドレス
    * 宛先NAT前ポート: 外部サービスポートで指定されたポート <br>
    ※ オプションのフィルタが無効かつポートフォワードが有効な場合に限ります。 
    * 宛先IPアドレス: マップされたIPアドレス/範囲に設定したIPアドレス
    * 宛先ポート: ポートへマップで指定されたポート
    * 送信元ポート: 送信元ポートのデータ取得規則に則る
    * URL/ドメイン: URL/ドメインのデータ取得規則に則る
    * アンチウイルス: アンチウイルスのデータ取得規則に則る
    * その他設定: その他設定のデータ取得規則に則る
    * 期待値: 期待値のデータ取得規則に則る
* バーチャルIPのオプションでプロトコルを指定している場合
  * ポリシーのサービスと指定したプロトコルがL4レベルで一致した場合に出力します。 <br>
  以下例のように、ポリシーのサービスが"NTP"で指定したプロトコルがICMPの場合は、`undefined`となります。 <br>
  プロトコルの取得例:
    | &nbsp;                 |           | 指定したプロトコル |           |           |
    | ---------------------- | --------- | ------------------ | --------- | --------- |
    |                        |           | ICMP               | TCP       | UDP       |
    | **ポリシーのサービス** | "ALL"     | icmp               | tcp       | udp       |
    | ^                      | "PING"    | icmp               | undefined | undefined |
    | ^                      | "NTP"     | undefined          | tcp       | udp       |
    | ^                      | "ALL_TCP" | undefined          | tcp       | undefined |
  ※ undefinedの箇所はFortiGateでDropするため、現在はシナリオを**出力していません。**
  * その他のデータ
  * 下記に該当している場合は、`undefined`を取得します。
    * ポリシーのサービスが"ALL_〇〇"または"PING"
      * 指定したプロトコルとL4レベルで一致していない場合
    * ポリシーのサービスが"ALL"、"ALL_〇〇"、"PING"以外
      * ポリシーのサービスで指定しているポートとポートへマップで指定されたポートが一致していない場合
  * 上記に該当していない場合、右に記載のデータを取得します。
    * 宛先NAT前IP: 外部IPアドレス/範囲で指定されたIPアドレス
    * 宛先NAT前ポート: 外部サービスポートで指定されたポート <br>
    ※ オプションのフィルタが無効かつポートフォワードが有効な場合に限ります。 
    * 宛先IPアドレス: マップされたIPアドレス/範囲に設定したIPアドレス
    * 宛先ポート: ポートへマップで指定されたポート
    * 送信元ポート: 送信元ポートのデータ取得規則に則る
    * URL/ドメイン: URL/ドメインのデータ取得規則に則る
    * アンチウイルス: アンチウイルスのデータ取得規則に則る
    * その他設定: その他設定のデータ取得規則に則る
    * 期待値: 期待値のデータ取得規則に則る
* 上記に該当しない場合
  * 各データの出力パターンに則ります。

## グループ要素の取得について

送信元・宛先やサービスがグループの場合、グループに所属している全ての要素を出力します。 <br>
以下は例になります。 <br>
例1: 送信元に"testgrp"を指定した場合
* コンフィグファイル(.conf)例（一部抜粋）:
```
config firewall address
    edit "1.1.1.1"
        set uuid a315f698-5945-51ea-b22b-82ef6b9d23da
        set subnet 1.1.1.1 255.255.255.255
    next
    edit "gmail.com"
        set uuid 9d3059d0-74dc-51ec-54f9-f59c21fee783
        set type fqdn
        set fqdn "gmail.com"
    next
end
config firewall addrgrp
    edit "testgrp"
        set uuid 9d3087e8-74dc-51ec-15ac-6b17eb973370
        set member "gmail.com" "1.1.1.1"
    next
end
config firewall policy
    edit 110
        set name "test_address_grp"
        set uuid c1adaa9c-5945-51ea-7e74-ff9fb0bd123f1
        set srcintf "lan"
        set dstintf "wan1"
        set srcaddr "testgrp"
        set dstaddr "all"
        set action accept
        set schedule "always"
        set service "ALL"
        set logtraffic all
        set fsso disable
    next
end
```
* 出力例:
```
送信元IP: gmail.com, 1.1.1.1
```
例2: サービスに"Email Access"を指定した場合
* コンフィグファイル(.conf)例（一部抜粋）:
```
config firewall service custom
    edit "DNS"
        set category "Network Services"
        set tcp-portrange 53
        set udp-portrange 53
    next
    edit "IMAP"
        set category "Email"
        set tcp-portrange 143
    next
    edit "IMAPS"
        set category "Email"
        set tcp-portrange 993
    next
    edit "POP3"
        set category "Email"
        set tcp-portrange 110
    next
    edit "POP3S"
        set category "Email"
        set tcp-portrange 995
    next
    edit "SMTP"
        set category "Email"
        set tcp-portrange 25
    next
    edit "SMTPS"
        set category "Email"
        set tcp-portrange 465
    next
end
config firewall service group
    edit "Email Access"
        set member "DNS" "IMAP" "IMAPS" "POP3" "POP3S" "SMTP" "SMTPS"
    next
end
config firewall policy
    edit 111
        set name "test_service_grp"
        set uuid c1adaa9c-5945-51ea-7e74-ff9fb0bd145f2
        set srcintf "lan"
        set dstintf "wan1"
        set srcaddr "all"
        set dstaddr "all"
        set action accept
        set schedule "always"
        set service "Email Access"
        set logtraffic all
        set fsso disable
    next
end
```
* 出力例:
```
プロトコル: dnst,dns,imap,tcp,tcp,tcp,smtp,tcp
宛先ポート: 53,53,143,993,110,995,25,465
```

## 送信元・宛先やサービスが複数指定されているポリシーの取得について

送信元・宛先やサービスが複数指定されているポリシーは、全ての要素を出力します。 <br>
以下は例になります。 <br>
例1: 送信元に"1.1.1.1"と"gmail.com"を指定した場合 <br>
* コンフィグファイル(.conf)例（一部抜粋）:
```
config firewall address
    edit "1.1.1.1"
        set uuid a315f698-5945-51ea-b22b-82ef6b9d23da
        set subnet 1.1.1.1 255.255.255.255
    next
    edit "gmail.com"
        set uuid 9d3059d0-74dc-51ec-54f9-f59c21fee783
        set type fqdn
        set fqdn "gmail.com"
    next
end
config firewall policy
    edit 112
        set name "test_address_grp"
        set uuid c1adaa9c-5945-51ea-7e74-ff9fb0bd123f1
        set srcintf "lan"
        set dstintf "wan1"
        set srcaddr "1.1.1.1" "gmail.com"
        set dstaddr "all"
        set action accept
        set schedule "always"
        set service "ALL"
        set logtraffic all
        set fsso disable
    next
end
```
* 出力例:
```
送信元IP: gmail.com, 1.1.1.1
```

例2: サービスに"DNS"と"IMAP"を指定した場合
* コンフィグファイル(.conf)例（一部抜粋）:
```
config firewall service custom
    edit "DNS"
        set category "Network Services"
        set tcp-portrange 53
        set udp-portrange 53
    next
    edit "IMAP"
        set category "Email"
        set tcp-portrange 143
    next
end
config firewall policy
    edit 113
        set name "test_service_grp"
        set uuid c1adaa9c-5945-51ea-7e74-ff9fb0bd145f2
        set srcintf "lan"
        set dstintf "wan1"
        set srcaddr "all"
        set dstaddr "all"
        set action accept
        set schedule "always"
        set service "DNS" "IMAP"
        set logtraffic all
        set fsso disable
    next
end
```
* 出力例:
```
プロトコル: dnst,dns,imap
宛先ポート: 53,53,143,993
```

## 出力するシナリオ数について

* 出力するシナリオ数は、以下の式より算出できます。
* ウェブフィルタのデータが存在する場合
  * 出力するシナリオ数 = [送信元のデータ数 x 宛先のデータ数 x (サービスのデータ数 - 出力するプロトコルが`http`,`https`以外のサービスのデータ数) x ウェブフィルタのデータ数] - `undefined`がデータとなるシナリオ数
* ウェブフィルタのデータが存在しない場合
  * 出力するシナリオ数 = 送信元のデータ数 x 宛先のデータ数 x サービスのデータ数　- `undefined`がデータとなるシナリオ数
* 以下に例を記載します。
  * 例1:  
    データの取得例:
    ```
    プロトコル: dns,dnst
    送信元IPアドレス: 192.168.1.100
    宛先IPアドレス: 8.8.4.4
    ```
    出力するシナリオ数: 1 x 1 x 2 = 2
    シナリオの出力例:
    | protocol | src-fw  | src-vlan(option) | src-ip        | src-port(option) | src-nat-ip(option) | s-if(option) | is-receiver-physical(option) | dst-fw  | dst-vlan(option) | dst-nat-ip(option) | dst-nat-port(option) | dst-ip | dst-port | d-if(option) | url/domain(option) | anti-virus(option) | timeout(option) | try(option) | other-settings(option) | expect | description        |
    | -------- | ------- | ---------------- | ------------- | ---------------- | ------------------ | ------------ | ---------------------------- | ------- | ---------------- | ------------------ | -------------------- | ------ | -------- | ------------ | ------------------ | ------------------ | --------------- | ----------- | ---------------------- | ------ | ------------------ |
    | dnst     | 1.1.1.1 | &nbsp;           | 192.168.1.100 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.4  | 53           | &nbsp;             | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test |
    | dns      | 1.1.1.1 | &nbsp;           | 192.168.1.100 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.4  | 53           | &nbsp;             | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test |
  * 例2:  
    データの取得例:
    ```
    プロトコル: icmp
    送信元IPアドレス: 192.168.1.101, 192.168.1.102
    宛先IPアドレス: 8.8.4.5, 8.8.4.6
    ```
    出力するシナリオ数: 2 x 2 x 1 = 4
    シナリオの出力例:
    | protocol | src-fw  | src-vlan(option) | src-ip        | src-port(option) | src-nat-ip(option) | s-if(option) | is-receiver-physical(option) | dst-fw  | dst-vlan(option) | dst-nat-ip(option) | dst-nat-port(option) | dst-ip | dst-port | d-if(option) | url/domain(option) | anti-virus(option) | timeout(option) | try(option) | other-settings(option) | expect | description         |
    | -------- | ------- | ---------------- | ------------- | ---------------- | ------------------ | ------------ | ---------------------------- | ------- | ---------------- | ------------------ | -------------------- | ------ | -------- | ------------ | ------------------ | ------------------ | --------------- | ----------- | ---------------------- | ------ | ------------------- |
    | icmp     | 1.1.1.1 | &nbsp;           | 192.168.1.101 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.5  |              | &nbsp;             | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test2 |
    | icmp     | 1.1.1.1 | &nbsp;           | 192.168.1.101 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.6  |              | &nbsp;             | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test2 |
    | icmp     | 1.1.1.1 | &nbsp;           | 192.168.1.102 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.5  |              | &nbsp;             | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test2 |
    | icmp     | 1.1.1.1 | &nbsp;           | 192.168.1.102 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.6  |              | &nbsp;             | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test2 |
  * 例3:  
    データの取得例:
    ```
    プロトコル: http
    送信元IPアドレス: 192.168.1.103
    宛先IPアドレス: 8.8.4.6
    ウェブフィルタ: hoge.com, sub.hoge.com
    ```
    出力するシナリオ数: 1 x 1 x (1 - 0) x 2 = 2
    シナリオの出力例:
    | protocol | src-fw  | src-vlan(option) | src-ip        | src-port(option) | src-nat-ip(option) | s-if(option) | is-receiver-physical(option) | dst-fw  | dst-vlan(option) | dst-nat-ip(option) | dst-nat-port(option) | dst-ip | dst-port | d-if(option) | url/domain(option) | anti-virus(option) | timeout(option) | try(option) | other-settings(option) | expect | description         |
    | -------- | ------- | ---------------- | ------------- | ---------------- | ------------------ | ------------ | ---------------------------- | ------- | ---------------- | ------------------ | -------------------- | ------ | -------- | ------------ | ------------------ | ------------------ | --------------- | ----------- | ---------------------- | ------ | ------------------- |
    | http     | 1.1.1.1 | &nbsp;           | 192.168.1.103 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.6  | 80           | hoge.com           | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test3 |
    | http     | 1.1.1.1 | &nbsp;           | 192.168.1.103 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.6  | 80           | sub.hoge.com       | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test3 |
  * 例4:  
    データの取得例:
    ```
    プロトコル: ftp,ftpa,https
    送信元IPアドレス: 192.168.1.104, 192.168.1.105
    宛先IPアドレス: 8.8.4.7
    ウェブフィルタ: fuga.com, sub.fuga.com
    ```
    出力するシナリオ数: 2 x 1 x (3 - 2) x 2 = 4
    シナリオの出力例:
    | protocol | src-fw  | src-vlan(option) | src-ip        | src-port(option) | src-nat-ip(option) | s-if(option) | is-receiver-physical(option) | dst-fw  | dst-vlan(option) | dst-nat-ip(option) | dst-nat-port(option) | dst-ip | dst-port | d-if(option) | url/domain(option) | anti-virus(option) | timeout(option) | try(option) | other-settings(option) | expect | description         |
    | -------- | ------- | ---------------- | ------------- | ---------------- | ------------------ | ------------ | ---------------------------- | ------- | ---------------- | ------------------ | -------------------- | ------ | -------- | ------------ | ------------------ | ------------------ | --------------- | ----------- | ---------------------- | ------ | ------------------- |
    | https    | 1.1.1.1 | &nbsp;           | 192.168.1.104 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.7  | 443          | fuga.com           | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test4 |
    | https    | 1.1.1.1 | &nbsp;           | 192.168.1.104 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.7  | 443          | sub.fuga.com       | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test4 |
    | https    | 1.1.1.1 | &nbsp;           | 192.168.1.105 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.7  | 443          | fuga.com           | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test4 |
    | https    | 1.1.1.1 | &nbsp;           | 192.168.1.105 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | &nbsp;             | &nbsp;               | &nbsp; | 8.8.4.7  | 443          | sub.fuga.com       | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test4 |

## 出力されるエラーメッセージについて

出力されるエラーメッセージは以下の通りです。
* 「実行パスの取得に失敗しました」
  * ツール実行時に取得している実行パスの取得に失敗した際に表示されます。 <br>
  ツールを管理者権限で実行しているかをご確認ください。
* 「FortiGateのコンフィグファイル(.conf)のパスを指定してください」
  * FortiGateのコンフィグファイル(.conf)のパスが指定されていない場合に表示されます。<br>
  FortiGateのコンフィグファイル(.conf)のパスを指定してください。
* 「正しいフォーマットでconfigディレクトリ内の`intf.csv`を記載してください」
  * 正しいフォーマットでconfigディレクトリ内の`intf.csv`を記載していない場合に表示されます。 <br>
  下記の[事前準備](#事前準備)を参考にconfigディレクトリ内の`intf.csv`を記載ください。 <br>
      https://github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate/blob/master/README.md
* 「configディレクトリ内の`intf.csv`の%+v行目に記載している~は不正なxxです」
  * configディレクトリ内の`intf.csv`の%+v行目に記載している~がxxの値として不正な値の場合に表示されます。<br>
  * 下記を参考に正しい値を記載してください。<br>
  ※ ()内は`intf.csv`のヘッダー名です。<br>
    * IPアドレス(ipaddress): `A.B.C.D`の形式かつ正しいIPアドレス
    * VLAN ID(vlanid): `0~4094`
    * サブネットマスク(subnetmask): `A.B.C.D`の形式かつ正しいサブネットマスク
* 「△△の情報はconfigディレクトリ内の`intf.csv`に記載されていませんでした」
  * インターフェース名△△の情報がconfigディレクトリ内の`intf.csv`に記載されていない場合に表示されます。 <br>
  下記の[事前準備](#事前準備)を参考に必要情報をconfigディレクトリ内の`intf.csv`に記載ください。 <br>
      https://github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate/blob/master/README.md
* 「〇〇はサポートされていません」
  * 〇〇というサービス名がツールでICMP,TCP,UDPのいずれにも分類できなかった場合に表示されます。<br>ポリシーで設定されている〇〇がコンフィグファイル(.conf)でサービスとして定義されているかをご確認ください。
* 「CSVの生成に失敗しました」
* 「予期せぬエラーが発生しました」
* 「panic: runtime error: index out of range ~」
  * ツールが想定していない動作をしている可能性があります。
  * 上記3種類のエラーが出力された場合はお手数ですが、下記情報を添付したissueを作成お願いいたします。<br>
  ※ 情報の取り扱いにご注意ください。
    * ツールで使用したコンフィグファイル(.conf)
    * ツール実行時に出力されたログ

## その他出力されるメッセージについての補足

* 「ポリシー名=〇〇の送信元NAT IPアドレスはダイナミックIPプールなため、生成したシナリオとテスト結果が異なる可能性があります」
* 「本シナリオとテスト結果で送信元NAT IPアドレスが異なる場合は値を変更後、再度テストを実施してください」
  * ダイナミックIPプールから送信元NAT IPアドレスの値を取得した際に表示されます。
  * 以下の動作例のように、生成したシナリオとテスト結果で送信元NAT IPアドレスが異なる場合があります。<br>
  ```
  例: ポリシーで1.1.1.11-1.1.1.20が外部IP範囲となるダイナミックIPプールを使用している場合
  本ツール → 1.1.1.11をシナリオとして出力
  FortiGate → 1.1.1.11-1.1.1.20の中のいずれかのIP(例: 1.1.1.15)に変換する
  NEEDLEWORK → 送信元NATがシナリオに記載されたIPかどうかをテストする

  → NEEDLEWORKでは1.1.1.11の想定だが、実際は1.1.1.15であるため、想定と異なるため結果はDropになる
  →→ 生成したシナリオとテスト結果での送信元NATの値が異なる
  ```
  * 上記動作の場合、送信元NAT IPアドレスの値を変更後、再度テストを実施してください。
* 「ポリシー名=〇〇は送信元、または宛先IPアドレスにFQDNが指定されています」
* 「本シナリオでテストを実行する場合は、NEEDLEWORKのマニュアルの「送信元・宛先アドレスにFQDNを指定したシナリオでテストを実施する」をご確認ください」
  * ポリシー名が〇〇の送信元、または宛先IPアドレスにFQDNが指定されている場合に表示されます。
  * 本メッセージが表示された場合、下記のNEEDLEWORKの操作マニュアルの「送信元・宛先アドレスにFQDNを指定したシナリオでテストを実施する」をご確認ください。<br>
  https://support.needlework.jp/manual
* 「~はundefinedが存在するため出力をスキップします」
  * シナリオに`undefined`が存在する場合に表示されます。 <br>
  `undefined`が出力されるパターンについては、[宛先がバーチャルIP/サーバーまたはVIPグループの場合について](##宛先がバーチャルIP/サーバーまたはVIPグループの場合について)をご確認ください。 <br>
以下にメッセージが表示される場合のコンフィグファイル(.conf)例、データ例、シナリオ例を記載します。 <br>
コンフィグファイル(.conf)例（一部抜粋）: 
```
config firewall address
    edit "192.168.1.106"
        set uuid a315f698-5945-51ea-b22b-82ef6b9d23da
        set subnet 192.168.1.106 255.255.255.255
    next
end
config firewall vip
  edit "option_portforward"
    set uuid 63a0f124-702e-51eb-e021-eb94a4b6c122
    set extip 1.1.1.5
    set extintf "any"
    set portforward enable
    set mappedip "200.200.200.240"
    set protocol udp
    set extport 53
    set mappedport 53000
  next
end
config firewall policy
    edit 114
        set name "test_vip"
        set uuid c1adaa9c-5945-51ea-7e74-ff9fb0bd145f2
        set srcintf "lan"
        set dstintf "wan1"
        set srcaddr "all"
        set dstaddr "option_portforward"
        set action accept
        set schedule "always"
        set service "ALL"
        set logtraffic all
        set fsso disable
    next
end
```
データ例: 
```
プロトコル: undefined, undefined, udp, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
宛先NAT前IP: undefined, undefined, 1.1.1.5, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
宛先NAT前ポート: undefined, undefined, 53, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
宛先IPアドレス: undefined, undefined, 200.200.200.240, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
宛先ポート: undefined, undefined, 53000, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
送信元ポート: undefined, undefined, 空文字, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
URL/ドメイン: undefined, undefined, 空文字, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
アンチウイルス: undefined, undefined, 空文字, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
その他設定: undefined, undefined, 空文字, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
期待値: undefined, undefined, pass, undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined
```
シナリオ例: 
| protocol | src-fw  | src-vlan(option) | src-ip        | src-port(option) | src-nat-ip(option) | s-if(option) | is-receiver-physical(option) | dst-fw  | dst-vlan(option) | dst-nat-ip(option) | dst-nat-port(option) | dst-ip          | dst-port | d-if(option) | url/domain(option) | anti-virus(option) | timeout(option) | try(option) | other-settings(option) | expect | description            |
| -------- | ------- | ---------------- | ------------- | ---------------- | ------------------ | ------------ | ---------------------------- | ------- | ---------------- | ------------------ | -------------------- | --------------- | -------- | ------------ | ------------------ | ------------------ | --------------- | ----------- | ---------------------- | ------ | ---------------------- |
| udp      | 1.1.1.1 | &nbsp;           | 192.168.1.106 | &nbsp;           | &nbsp;             | &nbsp;       | &nbsp;                       | 2.2.2.2 | 22               | 1.1.1.5            | 53                   | 200.200.200.240 | 53000    | &nbsp;       | &nbsp;             | &nbsp;             | &nbsp;          | &nbsp;      | &nbsp;                 | pass   | policy name = test_vip |

## 留意事項

* 本リポジトリ内のconfigディレクトリ配下のFortiGateのコンフィグ(.conf)に含まれないポリシーの場合、 <br>
  正常にcsvが出力されない可能性がありますので、予めご了承ください。
