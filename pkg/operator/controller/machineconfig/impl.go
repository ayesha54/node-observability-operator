/*
Copyright 2022.

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

package machineconfigcontroller

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type defaultImpl struct{}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . impl
type impl interface {
	InClusterConfig() (*rest.Config, error)
	NewForConfig(c *rest.Config) (*kubernetes.Clientset, error)
	ClientGet(context.Context, client.Client, client.ObjectKey, client.Object) error
	ClientList(context.Context, client.Client, client.ObjectList, ...client.ListOption) error
	ClientStatusUpdate(context.Context, client.Client, client.Object, ...client.UpdateOption) error
	ClientUpdate(context.Context, client.Client, client.Object, ...client.UpdateOption) error
	ClientCreate(context.Context, client.Client, client.Object, ...client.CreateOption) error
	ClientDelete(context.Context, client.Client, client.Object, ...client.DeleteOption) error
	ClientPatch(context.Context, client.Client, client.Object, client.Patch, ...client.PatchOption) error
}

func (d *defaultImpl) NewForConfig(
	c *rest.Config,
) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(c)
}
func (d *defaultImpl) InClusterConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}
func (*defaultImpl) NewClient(mgr ctrl.Manager) (client.Client, error) {
	return client.New(mgr.GetConfig(), client.Options{})
}
func (*defaultImpl) ManagerGetClient(m manager.Manager) client.Client {
	return m.GetClient()
}
func (*defaultImpl) ManagerGetEventRecorderFor(
	m manager.Manager, name string,
) record.EventRecorder {
	return m.GetEventRecorderFor(name)
}
func (*defaultImpl) NewControllerManagedBy(
	m manager.Manager,
	name string,
	p1 func(obj runtime.Object) bool,
	p2 func(obj runtime.Object) bool,
	r reconcile.Reconciler,
) error {
	return ctrl.NewControllerManagedBy(m).
		Named(name).
		WithEventFilter(predicate.And(
			resource.NewPredicates(p1),
			resource.NewPredicates(p2),
		)).
		For(&corev1.Pod{}).
		Complete(r)
}
func (*defaultImpl) ClientGet(
	ctx context.Context, c client.Client, key client.ObjectKey, obj client.Object,
) error {
	return c.Get(ctx, key, obj)
}
func (*defaultImpl) ClientStatusUpdate(
	ctx context.Context, c client.Client, obj client.Object, opts ...client.UpdateOption,
) error {
	return c.Status().Update(ctx, obj, opts...)
}
func (*defaultImpl) ClientList(
	ctx context.Context, c client.Client, list client.ObjectList, opts ...client.ListOption,
) error {
	return c.List(ctx, list, opts...)
}
func (*defaultImpl) ClientUpdate(
	ctx context.Context, c client.Client, obj client.Object, opts ...client.UpdateOption,
) error {
	return c.Update(ctx, obj, opts...)
}
func (*defaultImpl) ClientCreate(
	ctx context.Context, c client.Client, obj client.Object, opts ...client.CreateOption,
) error {
	return c.Create(ctx, obj, opts...)
}
func (*defaultImpl) ClientDelete(
	ctx context.Context, c client.Client, obj client.Object, opts ...client.DeleteOption,
) error {
	return c.Delete(ctx, obj, opts...)
}
func (*defaultImpl) ClientPatch(
	ctx context.Context, c client.Client, obj client.Object, patch client.Patch, opts ...client.PatchOption,
) error {
	return c.Patch(ctx, obj, patch, opts...)
}
