package v3

import (
	"context"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/controller"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	WorkloadGroupVersionKind = schema.GroupVersionKind{
		Version: "v3",
		Group:   "project.cattle.io",
		Kind:    "Workload",
	}
	WorkloadResource = metav1.APIResource{
		Name:         "workloads",
		SingularName: "workload",
		Namespaced:   true,

		Kind: WorkloadGroupVersionKind.Kind,
	}
)

type WorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workload
}

type WorkloadHandlerFunc func(key string, obj *Workload) error

type WorkloadLister interface {
	List(namespace string, selector labels.Selector) (ret []*Workload, err error)
	Get(namespace, name string) (*Workload, error)
}

type WorkloadController interface {
	Informer() cache.SharedIndexInformer
	Lister() WorkloadLister
	AddHandler(handler WorkloadHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type WorkloadInterface interface {
	ObjectClient() *clientbase.ObjectClient
	Create(*Workload) (*Workload, error)
	Get(name string, opts metav1.GetOptions) (*Workload, error)
	Update(*Workload) (*Workload, error)
	Delete(name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*WorkloadList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() WorkloadController
}

type workloadLister struct {
	controller *workloadController
}

func (l *workloadLister) List(namespace string, selector labels.Selector) (ret []*Workload, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*Workload))
	})
	return
}

func (l *workloadLister) Get(namespace, name string) (*Workload, error) {
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    WorkloadGroupVersionKind.Group,
			Resource: "workload",
		}, name)
	}
	return obj.(*Workload), nil
}

type workloadController struct {
	controller.GenericController
}

func (c *workloadController) Lister() WorkloadLister {
	return &workloadLister{
		controller: c,
	}
}

func (c *workloadController) AddHandler(handler WorkloadHandlerFunc) {
	c.GenericController.AddHandler(func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*Workload))
	})
}

type workloadFactory struct {
}

func (c workloadFactory) Object() runtime.Object {
	return &Workload{}
}

func (c workloadFactory) List() runtime.Object {
	return &WorkloadList{}
}

func (s *workloadClient) Controller() WorkloadController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.workloadControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(WorkloadGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &workloadController{
		GenericController: genericController,
	}

	s.client.workloadControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type workloadClient struct {
	client       *Client
	ns           string
	objectClient *clientbase.ObjectClient
	controller   WorkloadController
}

func (s *workloadClient) ObjectClient() *clientbase.ObjectClient {
	return s.objectClient
}

func (s *workloadClient) Create(o *Workload) (*Workload, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*Workload), err
}

func (s *workloadClient) Get(name string, opts metav1.GetOptions) (*Workload, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*Workload), err
}

func (s *workloadClient) Update(o *Workload) (*Workload, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*Workload), err
}

func (s *workloadClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *workloadClient) List(opts metav1.ListOptions) (*WorkloadList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*WorkloadList), err
}

func (s *workloadClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

func (s *workloadClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}
