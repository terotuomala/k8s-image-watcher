package controller

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	apps_v1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/terotuomala/k8s-image-watcher/pkg/client"
	"github.com/terotuomala/k8s-image-watcher/pkg/config"
	slack "github.com/terotuomala/k8s-image-watcher/pkg/handler"
)

const (
	deployment  = "deployment"
	daemonset   = "daemonset"
	statefulset = "statefulset"
)

type Controller struct {
	kubeClient   kubernetes.Interface
	informer     cache.SharedIndexInformer
	resourceType string
	slackClient  *slack.SlackNotifier
}

func Create(conf *config.Config, slackClient *slack.SlackNotifier) {
	kubeClient := client.GetClient()

	if conf.Resource.Deployment {
		deploymentInformer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					if conf.Namespace == "" {
						return kubeClient.AppsV1().Deployments("").List(context.Background(), options)
					}
					return kubeClient.AppsV1().Deployments(conf.Namespace).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					if conf.Namespace == "" {
						return kubeClient.AppsV1().Deployments("").Watch(context.Background(), options)
					}
					return kubeClient.AppsV1().Deployments(conf.Namespace).Watch(context.Background(), options)
				},
			},
			&apps_v1.Deployment{},
			0,
			cache.Indexers{},
		)

		c := addNewEventHandler(kubeClient, deploymentInformer, deployment, slackClient)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Start(stopCh)
	}

	if conf.Resource.DaemonSet {
		daemonsetInformer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1().DaemonSets(conf.Namespace).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1().DaemonSets(conf.Namespace).Watch(context.Background(), options)
				},
			},
			&apps_v1.DaemonSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := addNewEventHandler(kubeClient, daemonsetInformer, daemonset, slackClient)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Start(stopCh)
	}

	if conf.Resource.StatefulSet {
		statefulsetInformer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1().StatefulSets(conf.Namespace).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1().StatefulSets(conf.Namespace).Watch(context.Background(), options)
				},
			},
			&apps_v1.StatefulSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := addNewEventHandler(kubeClient, statefulsetInformer, statefulset, slackClient)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Start(stopCh)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func addNewEventHandler(client kubernetes.Interface, informer cache.SharedIndexInformer, resourceType string, slackClient *slack.SlackNotifier) *Controller {
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			switch resourceType {
			case deployment:
				resourceObject := obj.(*apps_v1.Deployment)
				log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("Found %s: %s/%s", resourceType, resourceObject.GetNamespace(), resourceObject.Name)

			case daemonset:
				resourceObject := obj.(*apps_v1.DaemonSet)

				log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("Found %s: %s/%s", resourceType, resourceObject.GetNamespace(), resourceObject.Name)
			case statefulset:
				resourceObject := obj.(*apps_v1.StatefulSet)
				log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("Found %s: %s/%s", resourceType, resourceObject.GetNamespace(), resourceObject.Name)

			default:
				log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("Did not match any")
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			switch resourceType {
			case deployment:
				newResourceObject := newObj.(*apps_v1.Deployment)
				oldResourceObject := oldObj.(*apps_v1.Deployment)

				oldImage := oldResourceObject.Spec.Template.Spec.Containers[0].Image
				newImage := newResourceObject.Spec.Template.Spec.Containers[0].Image

				if hasImageChanged(oldImage, newImage) {
					log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("%s: %s/%s image updated from %s to %s", resourceType, oldResourceObject.GetNamespace(), oldResourceObject.Name, oldImage, newImage)
					if slackClient != nil {
						message := fmt.Sprintf("%s: %s/%s image updated from %s to %s", resourceType, oldResourceObject.GetNamespace(), oldResourceObject.Name, oldImage, newImage)
						if err := slackClient.SendMessage(message); err != nil {
							log.Errorf("Failed to send notification to Slack: %v", err)
						}
					}
				}

			case daemonset:
				newResourceObject := newObj.(*apps_v1.DaemonSet)
				oldResourceObject := oldObj.(*apps_v1.DaemonSet)

				oldImage := oldResourceObject.Spec.Template.Spec.Containers[0].Image
				newImage := newResourceObject.Spec.Template.Spec.Containers[0].Image

				if hasImageChanged(oldImage, newImage) {
					log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("%s: %s/%s image updated from %s to %s", resourceType, oldResourceObject.GetNamespace(), oldResourceObject.Name, oldImage, newImage)
				}

			case statefulset:
				newResourceObject := newObj.(*apps_v1.StatefulSet)
				oldResourceObject := oldObj.(*apps_v1.StatefulSet)

				oldImage := oldResourceObject.Spec.Template.Spec.Containers[0].Image
				newImage := newResourceObject.Spec.Template.Spec.Containers[0].Image

				if hasImageChanged(oldImage, newImage) {
					log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("%s: %s/%s image updated from %s to %s", resourceType, oldResourceObject.GetNamespace(), oldResourceObject.Name, oldImage, newImage)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			switch resourceType {
			case deployment:
				resourceObject := obj.(*apps_v1.Deployment)
				log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("%s %s/%s deleted", resourceType, resourceObject.GetNamespace(), resourceObject.Name)

			case daemonset:
				resourceObject := obj.(*apps_v1.DaemonSet)
				log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("%s %s/%s deleted", resourceType, resourceObject.GetNamespace(), resourceObject.Name)

			case statefulset:
				resourceObject := obj.(*apps_v1.StatefulSet)
				log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("%s %s/%s deleted", resourceType, resourceObject.GetNamespace(), resourceObject.Name)
			}
		},
	})

	return &Controller{
		kubeClient:   client,
		informer:     informer,
		resourceType: resourceType,
		slackClient:  slackClient,
	}
}

func (c *Controller) Start(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()

	log.WithFields(log.Fields{"pkg": "contoller.go"}).Infof("Starting to watch %ss", c.resourceType)
	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
}

func hasImageChanged(oldImage, newImage string) bool {
	return oldImage != newImage
}
