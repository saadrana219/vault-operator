/*
Copyright 2018 The vault-operator Authors

Commercial software license.
*/

// This file was automatically generated by informer-gen

package internalinterfaces

import (
	versioned "github.com/coreos-inc/vault-operator/pkg/generated/clientset/versioned"
	runtime "k8s.io/apimachinery/pkg/runtime"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

type NewInformerFunc func(versioned.Interface, time.Duration) cache.SharedIndexInformer

// SharedInformerFactory a small interface to allow for adding an informer without an import cycle
type SharedInformerFactory interface {
	Start(stopCh <-chan struct{})
	InformerFor(obj runtime.Object, newFunc NewInformerFunc) cache.SharedIndexInformer
}
