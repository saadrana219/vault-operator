/*
Copyright 2018 The vault-operator Authors

Commercial software license.
*/

// This file was automatically generated by informer-gen

package v1alpha1

import (
	internalinterfaces "github.com/coreos-inc/vault-operator/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// VaultServices returns a VaultServiceInformer.
	VaultServices() VaultServiceInformer
}

type version struct {
	internalinterfaces.SharedInformerFactory
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory) Interface {
	return &version{f}
}

// VaultServices returns a VaultServiceInformer.
func (v *version) VaultServices() VaultServiceInformer {
	return &vaultServiceInformer{factory: v.SharedInformerFactory}
}
