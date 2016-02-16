/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package auth

import (
	"golang.org/x/crypto/ssh"
	"time"

	"github.com/gravitational/teleport/lib/backend"
	"github.com/gravitational/teleport/lib/services"
)

// AccessPoint is a interface needed by nodes to control the access
// to the node, and provide heartbeats
type AccessPoint interface {
	// GetServers returns a list of registered servers
	GetServers() ([]services.Server, error)

	// UpsertServer registers server presence, permanently if ttl is 0 or
	// for the specified duration with second resolution if it's >= 1 second
	UpsertServer(s services.Server, ttl time.Duration) error

	// GetUserCAPub returns the user certificate authority public key
	GetUserCertificateAuthority() (*services.CertificateAuthority, error)

	// GetUserCAPub returns the host certificate authority public key
	GetHostCertificateAuthority() (*services.CertificateAuthority, error)

	// GetWebSessionsKeys returns a list of generated public keys

	// associated with user web session
	GetWebSessionsKeys(user string) ([]services.AuthorizedKey, error)

	// GetRemoteCerts returns a list of trusted remote certificates
	GetRemoteCertificates(certType, domainName string) ([]services.CertificateAuthority, error)

	// GetTrustedCerts returns a list of trusted certificates
	GetTrustedCertificates(certType string) ([]services.CertificateAuthority, error)

	// GetCertificateID Returns an ID of the certificate with the provided PublicKey value
	GetCertificateID(certType string, key ssh.PublicKey) (ID string, found bool, e error)

	// GetAllUserMappings Returns hashes of all the User Mapping entries. Used to copy all the data
	GetAllUserMappings() (hashes []string, e error)

	UserMappingExists(certificateID, teleportUser, osUser string) (bool, error)
}

type BackendAccessPoint struct {
	*services.CAService
	*services.PresenceService
	*services.ProvisioningService
	*services.WebService
}

func NewBackendAccessPoint(bk backend.Backend) *BackendAccessPoint {
	ap := BackendAccessPoint{}
	ap.CAService = services.NewCAService(bk)
	ap.PresenceService = services.NewPresenceService(bk)
	ap.ProvisioningService = services.NewProvisioningService(bk)
	ap.WebService = services.NewWebService(bk)

	return &ap
}
