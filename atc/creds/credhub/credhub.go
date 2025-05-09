package credhub

import (
	"path"
	"time"

	"github.com/concourse/concourse/atc/creds"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/lager/v3"
)

type CredHubAtc struct {
	CredHub *LazyCredhub
	logger  lager.Logger
	prefix  string
}

// NewSecretLookupPaths defines how variables will be searched in the underlying secret manager
func (c CredHubAtc) NewSecretLookupPaths(teamName string, pipelineName string, allowRootPath bool) []creds.SecretLookupPath {
	lookupPaths := []creds.SecretLookupPath{}
	if len(pipelineName) > 0 {
		lookupPaths = append(lookupPaths, creds.NewSecretLookupWithPrefix(path.Join(c.prefix, teamName, pipelineName)+"/"))
	}
	lookupPaths = append(lookupPaths, creds.NewSecretLookupWithPrefix(path.Join(c.prefix, teamName)+"/"))
	if allowRootPath {
		lookupPaths = append(lookupPaths, creds.NewSecretLookupWithPrefix(c.prefix+"/"))
	}
	return lookupPaths
}

// Get retrieves the value and expiration of an individual secret
func (c CredHubAtc) Get(secretPath string) (any, *time.Time, bool, error) {
	var cred credentials.Credential
	var found bool
	var err error

	cred, found, err = c.findCred(secretPath)
	if err != nil {
		c.logger.Error("unable to retrieve credhub secret", err)
		return nil, nil, false, err
	}

	if !found {
		return nil, nil, false, nil
	}

	return cred.Value, nil, true, nil
}

func (c CredHubAtc) findCred(path string) (credentials.Credential, bool, error) {
	var cred credentials.Credential
	var err error

	ch, err := c.CredHub.CredHub()
	if err != nil {
		return cred, false, err
	}

	results, err := ch.FindByPartialName(path)
	if err != nil {
		return cred, false, err
	}

	// same as https://github.com/cloudfoundry/credhub-cli/blob/main/commands/find.go#L22
	if len(results.Credentials) == 0 {
		return cred, false, nil
	}

	cred, err = ch.GetLatestVersion(path)
	if err != nil {
		return cred, false, err
	}

	return cred, true, nil
}
