package v1

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ConfigMapLifecycle interface {
	Create(obj *v1.ConfigMap) (*v1.ConfigMap, error)
	Remove(obj *v1.ConfigMap) (*v1.ConfigMap, error)
	Updated(obj *v1.ConfigMap) (*v1.ConfigMap, error)
}

type configMapLifecycleAdapter struct {
	lifecycle ConfigMapLifecycle
}

func (w *configMapLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1.ConfigMap))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *configMapLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1.ConfigMap))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *configMapLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1.ConfigMap))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewConfigMapLifecycleAdapter(name string, clusterScoped bool, client ConfigMapInterface, l ConfigMapLifecycle) ConfigMapHandlerFunc {
	adapter := &configMapLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.ConfigMap) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
