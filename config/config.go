// Config client for Spring Cloud Configuration server.
package config

import (
	"errors"
	"fmt"
	"github.com/ContainX/go-utils/encoding"
	"github.com/ContainX/go-utils/envsubst"
	"github.com/ContainX/go-utils/httpclient"
	"os"
	"strings"
)

const (
	// EnvConfigProfile holds the name of the environment variable (CONFIG_PROFILE)
	// which is used during runtime lookups
	EnvConfigProfile = "CONFIG_PROFILE"
	// The configuration server URI environment variable (CONFIG_SERVER_URI)
	// ex. http://host:8888
	EnvConfigServerURI = "CONFIG_SERVER_URI"
	// UriDefault is the default URI to the configuration server
	UriDefault = "http://localhost:8888"
	// ProfileDefault is the default profile
	ProfileDefault = "default"
	// LabelDefault is the initial SCM branch
	LabelDefault = "master"
	// Format is {uri}/{label}/{name}-{profile}.type
	configPathFmt = "%s/%s/%s-%s.%s"
	extJSON       = "json"
	extPROP       = "properties"
	extYAML       = "yml"
)

var (
	NameNotDeclaredErr = errors.New("Name must be declared")
	FileNotDeclaredErr = errors.New("Filename must have a value")
)

type ConfigClient interface {
	// Fetch queries the remote configuration service and populates the
	// target value
	Fetch(target interface{}) error

	// FetchWithSubstitution fetches a remote config, substitutes environment variables
	// and writes it to the target
	FetchWithSubstitution(target interface{}) error

	// Fetch queries the remote configuration service and populates
	// a map of kv strings.   This call flattens hierarchical values
	// into flattened form.  Example:  datasource.mysql.user
	FetchAsMap() (map[string]string, error)

	// Fetch queries the remote configuration service and returns
	// the result as a JSON string
	FetchAsJSON() (string, error)

	// Fetch queries the remote configuration service and returns
	// the result as a YAML string
	FetchAsYAML() (string, error)

	// Fetch queries the remote configuration service and returns
	// the result as a Properties string
	FetchAsProperties() (string, error)

	// Bootstrap returns a reference to the current bootstrap settings
	Bootstrap() *Bootstrap
}

type client struct {
	bootstrap *Bootstrap
}

// Bootstrap is the properties needed to fetch a remote configuration from
// spring cloud configuration server.
type Bootstrap struct {

	// The URI of the remote server (default http://localhost:8888).
	URI string `json:"uri"`

	// Context is used as the base URI an optional /refresh endpoint which can by
	// exposed to update an anonynmous function when configuration changes
	Context string `json:"context"`

	// Profile represents the default to use when fetching remote configuration (comma-separated).
	// Default is "default".
	//
	// Note: During runtime the config client looks for the presence of an environment
	// variable called CONFIG_PROFILE.  If this is defined it overwrites this value.
	Profile string `json:"profile"`

	// Name of application used to fetch remote properties.
	Name string `json:"name"`

	// Label name to use to pull remote configuration properties. The default is set
	// on the server (generally "master" for a git based server).
	Label string `json:"label"`

	// The username to use (HTTP Basic) when contacting the remote server.
	Username string `json:"username,omitempty"`

	// The password to use (HTTP Basic) when contacting the remote server.
	Password string `json:"password,omitempty"`
}

// New creates a new ConfigClient based on b Bootstrap
// Error will be thrown if Name is not set
func New(b Bootstrap) (ConfigClient, error) {
	if b.Name == "" {
		return nil, NameNotDeclaredErr
	}

	b.URI = defaultVal(b.URI, UriDefault)
	b.Profile = defaultVal(b.Profile, ProfileDefault)
	b.Label = defaultVal(b.Label, LabelDefault)

	client := &client{
		bootstrap: &b,
	}
	return client, nil
}

func LoadFromFile(filename string) (ConfigClient, error) {
	if filename == "" {
		return nil, FileNotDeclaredErr
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	if encoder, err := encoding.NewEncoderFromFileExt(filename); err == nil {
		config := &Bootstrap{}
		if e := encoder.UnMarshal(f, config); e != nil {
			return nil, e
		}
		config.PopulateDefaultsIfEmpty()
		if config.Name == "" {
			return nil, NameNotDeclaredErr
		}
		client := &client{
			bootstrap: config,
		}
		return client, nil
	} else {
		return nil, err
	}
}

func (b *Bootstrap) PopulateDefaultsIfEmpty() {
	if b.URI == "" {
		b.URI = UriDefault
	}
	if b.Profile == "" {
		b.Profile = ProfileDefault
	}
	if b.Label == "" {
		b.Label = LabelDefault
	}
}

// defaultVal returns "d" if "s" aka source has an empty value
func defaultVal(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

// Profile returns the resolved value of the profile.  It will
// either use the CLUSTER_PROFILE env variable or fallback to
// the specified default value as a fallback
func (c *client) resolveProfile() string {
	if v := os.Getenv(EnvConfigProfile); v != "" {
		return v
	}
	return c.bootstrap.Profile
}

// URI returns the resolved value of the config server URI.  It will
// either use the CONFIG_SERVER_URI env variable or fallback to
// the specified default value as a fallback
func (c *client) resolveURI() string {
	if v := os.Getenv(EnvConfigServerURI); v != "" {
		return v
	}
	return c.bootstrap.URI
}

// Fetch queries the remote configuration service and populates the
// target value
func (c *client) Fetch(target interface{}) error {
	uri := c.buildRequestURI(extJSON)
	resp := httpclient.Get(uri, target)
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (c *client) FetchWithSubstitution(target interface{}) error {
	content, err := c.FetchAsYAML()
	if err != nil {
		return err
	}

	enc, _ := encoding.NewEncoder(encoding.YAML)

	if err = enc.UnMarshalStr(content, target); err != nil {
		return err
	}
	return nil
}

func (c *client) FetchAsMap() (map[string]string, error) {
	uri := c.buildRequestURI(extPROP)
	resp := httpclient.Get(uri, nil)
	if resp.Error != nil {
		return nil, resp.Error
	}

	m := map[string]string{}
	for _, line := range strings.Split(resp.Content, "\n") {
		kv := strings.Split(line, ":")
		if len(kv) == 2 {
			m[kv[0]] = strings.TrimSpace(kv[1])
		}
	}
	return m, nil
}

func (c *client) FetchAsProperties() (string, error) {
	return c.fetchAsString(extPROP)
}

func (c *client) FetchAsJSON() (string, error) {
	return c.fetchAsString(extJSON)
}

func (c *client) FetchAsYAML() (string, error) {
	return c.fetchAsString(extYAML)
}

func (c *client) fetchAsString(extension string) (string, error) {
	uri := c.buildRequestURI(extension)

	resp := httpclient.Get(uri, nil)
	content := resp.Content
	if resp.Error == nil {
		content = envsubst.Substitute(strings.NewReader(content), false, func(s string) string {
			return os.Getenv(s)
		})
	}
	return content, resp.Error
}

func (c *client) Bootstrap() *Bootstrap {
	return c.bootstrap
}

// Builds the final request URI for fetching a remote configuration.
// The returned URI is in the format of : {uri}/{label}/{name}-{profile}.json
func (c *client) buildRequestURI(t string) string {
	return fmt.Sprintf(configPathFmt, c.resolveURI(), c.bootstrap.Label, c.bootstrap.Name, c.resolveProfile(), t)
}
