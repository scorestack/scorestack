package checktypes

import (
	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/dns"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/ftp"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/git"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/http"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/icmp"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/imap"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/ldap"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/mssql"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/mysql"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/noop"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/postgresql"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/smb"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/smtp"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/ssh"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/vnc"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/winrm"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes/xmpp"
	"go.uber.org/zap"
)

func GetCheckType(c check.Config) check.Check {
	var def check.Check
	switch c.Type {
	case "noop":
		def = &noop.Definition{}
	case "http":
		def = &http.Definition{}
	case "icmp":
		def = &icmp.Definition{}
	case "ssh":
		def = &ssh.Definition{}
	case "dns":
		def = &dns.Definition{}
	case "ftp":
		def = &ftp.Definition{}
	case "ldap":
		def = &ldap.Definition{}
	case "vnc":
		def = &vnc.Definition{}
	case "imap":
		def = &imap.Definition{}
	case "smtp":
		def = &smtp.Definition{}
	case "winrm":
		def = &winrm.Definition{}
	case "xmpp":
		def = &xmpp.Definition{}
	case "mysql":
		def = &mysql.Definition{}
	case "smb":
		def = &smb.Definition{}
	case "postgresql":
		def = &postgresql.Definition{}
	case "mssql":
		def = &mssql.Definition{}
	case "git":
		def = &git.Definition{}
	default:
		zap.S().Warnf("check id %s had an invalid type: %s", c.ID, c.Type)
		def = &noop.Definition{}
	}

	return def
}
