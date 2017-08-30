package operator

import (
	"context"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/coreos-inc/vault-operator/pkg/spec"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func (v *Vaults) run(ctx context.Context) {
	source := cache.NewListWatchFromClient(
		v.vaultsCRCli.RESTClient(),
		spec.VaultResourcePlural,
		v.namespace,
		fields.Everything())

	v.queue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "vault-operator")
	v.indexer, v.informer = cache.NewIndexerInformer(source, &spec.Vault{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc:    v.onAddVault,
		UpdateFunc: v.onUpdateVault,
		DeleteFunc: v.onDeleteVault,
	}, cache.Indexers{})

	defer v.queue.ShutDown()

	logrus.Info("starting Vaults controller")
	go v.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), v.informer.HasSynced) {
		logrus.Error("Timed out waiting for caches to sync")
		return
	}

	const numWorkers = 1
	for i := 0; i < numWorkers; i++ {
		go wait.Until(v.runWorker, time.Second, ctx.Done())
	}

	<-ctx.Done()
	logrus.Info("stopping Vaults controller")
}

func (v *Vaults) onAddVault(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		panic(err)
	}
	v.queue.Add(key)
}

func (v *Vaults) onUpdateVault(oldObj, newObj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err != nil {
		panic(err)
	}
	v.queue.Add(key)
}

func (v *Vaults) onDeleteVault(obj interface{}) {
	vr, ok := obj.(*spec.Vault)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			panic(fmt.Sprintf("unknown object from Vault delete event: %#v", obj))
		}
		vr, ok = tombstone.Obj.(*spec.Vault)
		if !ok {
			panic(fmt.Sprintf("Tombstone contained object that is not a Vault: %#v", obj))
		}
	}

	cancel := v.ctxCancels[vr.Name]
	cancel()
	delete(v.ctxCancels, vr.Name)

	// IndexerInformer uses a delta queue, therefore for deletes we have to use this
	// key function.
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		panic(err)
	}
	v.toDelete[key] = vr
	v.queue.Add(key)
}
