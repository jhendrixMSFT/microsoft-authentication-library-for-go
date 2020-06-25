// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package requests

import (
	"sync"

	"github.com/AzureAD/microsoft-authentication-library-for-go/src/internal/msalbase"
)

var instanceDiscoveryCache = map[string]*instanceDiscoveryMetadata{}
var instanceDiscoveryCacheInitOnce sync.Once

func initInstanceDiscoveryCache() {

}

type IAadInstanceDiscovery interface {
	GetMetadataEntry(authorityInfo *msalbase.AuthorityInfo) (*instanceDiscoveryMetadata, error)
}

type AadInstanceDiscovery struct {
	webRequestManager IWebRequestManager
}

func CreateAadInstanceDiscovery() *AadInstanceDiscovery {
	instanceDiscoveryCacheInitOnce.Do(initInstanceDiscoveryCache)
	return &AadInstanceDiscovery{}
}

func (d *AadInstanceDiscovery) doInstanceDiscoveryAndCache(authorityInfo *msalbase.AuthorityInfo) (*instanceDiscoveryMetadata, error) {
	discoveryResponse, err := d.webRequestManager.GetAadinstanceDiscoveryResponse(authorityInfo)
	if err != nil {
		return nil, err
	}

	for _, metadataEntry := range discoveryResponse.Metadata {
		metadataEntry.TenantDiscoveryEndpoint = discoveryResponse.TenantDiscoveryEndpoint
		for _, aliasedAuthority := range metadataEntry.Aliases {
			instanceDiscoveryCache[aliasedAuthority] = &metadataEntry
		}
	}

	instanceDiscoveryCache[authorityInfo.Host] = createInstanceDiscoveryMetadata(authorityInfo.Host, discoveryResponse.TenantDiscoveryEndpoint)
	return d.GetMetadataEntry(authorityInfo)
}

func (d *AadInstanceDiscovery) GetMetadataEntry(authorityInfo *msalbase.AuthorityInfo) (*instanceDiscoveryMetadata, error) {
	if metadata, ok := instanceDiscoveryCache[authorityInfo.Host]; ok {
		return metadata, nil
	}

	metadata, err := d.doInstanceDiscoveryAndCache(authorityInfo)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
