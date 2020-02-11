package ssh

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the SSH check
// it implements the "check" interface
type Definition struct {
	ID       string        // unique identifier for this check
	Name     string        // a human-readable title for the check
	IP       string        // (required) IP of the host to run the ICMP check against
	Port     string        // (optional, default=22) The port to attempt an ssh connection on
	Username string        // (required) The user to login with over ssh
	Password string        // (required) The password for the user that you wish to login with
	Cmd      string        // (required) The command to execute once ssh connection established
	Timeout  time.Duration // (optional, default=5sec) How long to attempt a connection before terminating
}

func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	// Config SSH client
	config := &ssh.ClientConfig{
		User: d.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(d.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Create the ssh client
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", d.IP, d.Port), config)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// Create a session from the connection
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	defer session.Close()

	// Run a command
	output, err := session.CombinedOutput("/usr/bin/woami")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("OUtput: %s\n", output)
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, def []byte) error {

	// Set ID and Name
	d.ID = id
	d.Name = name

	// Unpack JSON definition
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Check for optional Port value
	if d.Port == "" {
		d.Port = "22"
	}

	// Check for optional timeout value
	if d.Timeout == 0*time.Second {
		d.Timeout = 5 * time.Second
	}

	// Check for missing fields
	missingFields := make([]string, 0)
	if d.IP == "" {
		missingFields = append(missingFields, "IP")
	}

	if d.Username == "" {
		missingFields = append(missingFields, "Username")
	}

	if d.Password == "" {
		missingFields = append(missingFields, "Password")
	}

	if d.Cmd == "" {
		missingFields = append(missingFields, "Cmd")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "ssh",
			Field: missingFields[0],
		}
	}
	return nil
}
