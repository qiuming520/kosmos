package globalnodecontroller

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kosmos.io/kosmos/pkg/apis/kosmos/v1alpha1"
	"github.com/kosmos.io/kosmos/pkg/generated/clientset/versioned"
)

const (
	GlobalNodeStatusControllerName = "global-node-status-controller"
	DefaultStatusUpdateInterval    = 15 * time.Second
	ClientHeartbeatThreshold       = 10 * time.Second
)

type GlobalNodeStatusController struct {
	root           client.Client
	statusInterval time.Duration

	kosmosClient versioned.Interface
}

func NewGlobalNodeStatusController(
	root client.Client,
	kosmosClient versioned.Interface,
) *GlobalNodeStatusController {
	return &GlobalNodeStatusController{
		root:           root,
		statusInterval: DefaultStatusUpdateInterval,
		kosmosClient:   kosmosClient,
	}
}
func (c *GlobalNodeStatusController) Start(ctx context.Context) error {
	go wait.UntilWithContext(ctx, c.syncGlobalNodeStatus, c.statusInterval)

	<-ctx.Done()
	return nil
}
func (c *GlobalNodeStatusController) syncGlobalNodeStatus(ctx context.Context) {
	globalNodes := make([]*v1alpha1.GlobalNode, 0)
	//c.globalNodeLock.Lock()
	//defer c.globalNodeLock.Unlock()

	nodeList, err := c.kosmosClient.KosmosV1alpha1().GlobalNodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Errorf("Failed to fetch GlobalNodes: %v", err)
		return
	}
	for _, node := range nodeList.Items {
		nodeCopy := node.DeepCopy()
		globalNodes = append(globalNodes, nodeCopy)
	}

	err = c.updateGlobalNodeStatus(ctx, globalNodes)
	if err != nil {
		klog.Errorf("Failed to sync global node status: %v", err)
	}
}

func (c *GlobalNodeStatusController) updateGlobalNodeStatus(
	ctx context.Context,
	globalNodes []*v1alpha1.GlobalNode,
) error {
	for _, globalNode := range globalNodes {
		err := c.updateStatusForGlobalNode(ctx, globalNode)
		if err != nil {
			klog.Errorf("Failed to update status for global node %s: %v", globalNode.Name, err)
			return err
		}
	}
	return nil
}

func (c *GlobalNodeStatusController) updateStatusForGlobalNode(
	ctx context.Context,
	globalNode *v1alpha1.GlobalNode,
) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		currentNode, err := c.kosmosClient.KosmosV1alpha1().GlobalNodes().Get(ctx, globalNode.Name, metav1.GetOptions{})
		if err != nil {
			klog.Errorf("Failed to fetch the latest GlobalNode %s: %v", globalNode.Name, err)
			return err
		}

		if len(currentNode.Status.Conditions) == 0 {
			klog.Warningf("GlobalNode %s has no conditions, skipping status update", currentNode.Name)
			return nil
		}

		condition := currentNode.Status.Conditions[0]
		lastHeartbeatTime := condition.LastHeartbeatTime
		timeDiff := time.Since(lastHeartbeatTime.Time)

		statusType := "Ready"
		if timeDiff > ClientHeartbeatThreshold {
			statusType = "NotReady"
		}

		if string(condition.Type) != statusType {
			condition.Type = v1.NodeConditionType(statusType)
			condition.LastTransitionTime = metav1.NewTime(time.Now())

			currentNode.Status.Conditions[0] = condition

			_, err = c.kosmosClient.KosmosV1alpha1().GlobalNodes().UpdateStatus(ctx, currentNode, metav1.UpdateOptions{})
			if err != nil {
				if errors.IsConflict(err) {
					klog.Warningf("Conflict detected while updating status for GlobalNode %s, retrying...", globalNode.Name)
				} else {
					klog.Errorf("Failed to update status for GlobalNode %s: %v", globalNode.Name, err)
				}
				return err
			}

			klog.Infof("Successfully updated status for GlobalNode %s to %s", globalNode.Name, statusType)
		} else {
			klog.Infof("No status update required for GlobalNode %s, current status: %s", globalNode.Name, condition.Type)
		}

		return nil
	})
}
