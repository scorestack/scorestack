module github.com/scorestack/scorestack/dynamicbeat

go 1.15

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.2.0+incompatible
	github.com/Shopify/sarama => github.com/elastic/sarama v1.19.1-0.20200629123429-0e7b69039eec
	github.com/cucumber/godog => github.com/cucumber/godog v0.8.1
	github.com/docker/docker => github.com/docker/engine v0.0.0-20191113042239-ea84732a7725
	github.com/docker/go-plugins-helpers => github.com/elastic/go-plugins-helpers v0.0.0-20200207104224-bdf17607b79f
	github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20
	github.com/dop251/goja_nodejs => github.com/dop251/goja_nodejs v0.0.0-20171011081505-adff31b136e6
	github.com/fsnotify/fsevents => github.com/elastic/fsevents v0.0.0-20181029231046-e1d381a4d270
	github.com/fsnotify/fsnotify => github.com/adriansr/fsnotify v0.0.0-20180417234312-c9bbe1f46f1d
	github.com/google/gopacket => github.com/adriansr/gopacket v1.1.18-0.20200327165309-dd62abfa8a41
	github.com/insomniacslk/dhcp => github.com/elastic/dhcp v0.0.0-20200227161230-57ec251c7eb3 // indirect
	github.com/scorestack/scorestack => github.com/scorestack/scorestack v0.5.1
	github.com/tonistiigi/fifo => github.com/containerd/fifo v0.0.0-20190816180239-bda0ff6ed73c
	golang.org/x/tools => golang.org/x/tools v0.0.0-20200602230032-c00d67ef29d0 // release 1.14
)

require (
	github.com/Bowery/prompt v0.0.0-20190916142128-fa8279994f75 // indirect
	github.com/akavel/rsrc v0.9.0 // indirect
	github.com/dchest/safefile v0.0.0-20151022103144-855e8d98f185 // indirect
	github.com/dlclark/regexp2 v1.2.1 // indirect
	github.com/dop251/goja v0.0.0-20200929101608-beb0a9a01fbc // indirect
	github.com/dop251/goja_nodejs v0.0.0-20200811150831-9bc458b4bbeb // indirect
	github.com/elastic/beats/v7 v7.0.0-alpha2.0.20201001160411-075696e6a142
	github.com/elastic/go-elasticsearch/v7 v7.9.0
	github.com/elastic/go-sysinfo v1.4.0 // indirect
	github.com/emersion/go-imap v1.0.5
	github.com/fatih/color v1.9.0 // indirect
	github.com/go-ldap/ldap/v3 v3.2.3
	github.com/go-ping/ping v0.0.0-20200918120429-e8ae07c3cec8
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hirochachacha/go-smb2 v1.0.3
	github.com/jlaffaye/ftp v0.0.0-20200812143550-39e3779af0db
	github.com/josephspurrier/goversioninfo v1.2.0 // indirect
	github.com/kardianos/govendor v1.0.9 // indirect
	github.com/magefile/mage v1.10.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/miekg/dns v1.1.31
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/mitchellh/gox v1.0.1
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/oneNutW0nder/winrm v0.0.0-20200403191630-928a10cb3c1e
	github.com/pierrre/gotestcover v0.0.0-20160517101806-924dca7d15f0
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/reviewdog/reviewdog v0.10.2
	github.com/tsg/go-daemon v0.0.0-20200207173439-e704b93fd89b
	go.elastic.co/ecszap v0.2.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201001193750-eb9a90e9f9cb
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/net v0.0.0-20200930145003-4acb6c075d10 // indirect
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20201001104356-43ebab892c4c
	gosrc.io/xmpp v0.5.1
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
	howett.net/plist v0.0.0-20200419221736-3b63eb3a43b5 // indirect
)
