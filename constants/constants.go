package constants

const ConfigDir = "/etc/tdservice/"
const ConfigFile = "config.yml"

const TLSCertFile = "cert.pem"
const TLSKeyFile = "key.pem"

const PIDFile = "tdservice.pid"

// State represents whether or not a daemon is running or not
type State bool

const (
	// Stopped is the default nil value, indicating not running
	Stopped State = false
	// Running means the daemon is active
	Running State = true
)
