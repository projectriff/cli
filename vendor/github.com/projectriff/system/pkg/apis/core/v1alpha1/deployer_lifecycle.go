/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	knapis "github.com/knative/pkg/apis"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	DeployerConditionReady                                = knapis.ConditionReady
	DeployerConditionDeploymentReady knapis.ConditionType = "DeploymentReady"
	DeployerConditionServiceReady    knapis.ConditionType = "ServiceReady"
)

var deployerCondSet = knapis.NewLivingConditionSet(
	DeployerConditionDeploymentReady,
	DeployerConditionServiceReady,
)

func (ds *DeployerStatus) GetObservedGeneration() int64 {
	return ds.ObservedGeneration
}

func (ds *DeployerStatus) IsReady() bool {
	return deployerCondSet.Manage(ds).IsHappy()
}

func (*DeployerStatus) GetReadyConditionType() knapis.ConditionType {
	return DeployerConditionReady
}

func (ds *DeployerStatus) GetCondition(t knapis.ConditionType) *knapis.Condition {
	return deployerCondSet.Manage(ds).GetCondition(t)
}

func (ds *DeployerStatus) InitializeConditions() {
	deployerCondSet.Manage(ds).InitializeConditions()
}

func (ds *DeployerStatus) MarkDeploymentNotOwned(name string) {
	deployerCondSet.Manage(ds).MarkFalse(DeployerConditionDeploymentReady, "NotOwned",
		"There is an existing Deployment %q that we do not own.", name)
}

func (ds *DeployerStatus) PropagateDeploymentStatus(cds *appsv1.DeploymentStatus) {
	var ac *appsv1.DeploymentCondition
	for _, c := range cds.Conditions {
		if c.Type == appsv1.DeploymentAvailable {
			ac = &c
			break
		}
	}
	if ac == nil {
		return
	}
	switch {
	case ac.Status == corev1.ConditionUnknown:
		deployerCondSet.Manage(ds).MarkUnknown(DeployerConditionDeploymentReady, ac.Reason, ac.Message)
	case ac.Status == corev1.ConditionTrue:
		deployerCondSet.Manage(ds).MarkTrue(DeployerConditionDeploymentReady)
	case ac.Status == corev1.ConditionFalse:
		deployerCondSet.Manage(ds).MarkFalse(DeployerConditionDeploymentReady, ac.Reason, ac.Message)
	}
}

func (ds *DeployerStatus) MarkServiceNotOwned(name string) {
	deployerCondSet.Manage(ds).MarkFalse(DeployerConditionServiceReady, "NotOwned",
		"There is an existing Service %q that we do not own.", name)
}

func (ds *DeployerStatus) PropagateServiceStatus(ss *corev1.ServiceStatus) {
	// services don't have meaningful status
	deployerCondSet.Manage(ds).MarkTrue(DeployerConditionServiceReady)
}
