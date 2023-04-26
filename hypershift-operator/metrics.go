package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	configv1 "github.com/openshift/api/config/v1"
	hyperv1 "github.com/openshift/hypershift/api/v1beta1"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	crmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

type hypershiftMetrics struct {
	clusterProxy                       *prometheus.GaugeVec
	clusterIdentityProviders           *prometheus.GaugeVec
	clusterLimitedSupportEnabled       *prometheus.GaugeVec
	clusterSilenceAlerts               *prometheus.GaugeVec
	hostedClusters                     *prometheus.GaugeVec
	hostedClustersWithFailureCondition *prometheus.GaugeVec
	hostedClustersNodePools            *prometheus.GaugeVec
	nodePools                          *prometheus.GaugeVec
	nodePoolsWithFailureCondition      *prometheus.GaugeVec
	nodePoolSize                       *prometheus.GaugeVec

	client crclient.Client

	log logr.Logger
}

func newMetrics(client crclient.Client, log logr.Logger) *hypershiftMetrics {
	return &hypershiftMetrics{
		clusterProxy: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Help: "Indicates cluster proxy state for each cluster",
			Name: "hypershift_cluster_proxy",
		}, []string{"namespace", "name", "proxy_http", "proxy_https", "proxy_trusted_ca"}),
		clusterIdentityProviders: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Help: "Indicates the number any identity provider in the fleet",
			Name: "hypershift_cluster_identity_providers",
		}, []string{"identity_provider"}),
		clusterLimitedSupportEnabled: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Help: "Indicates if limited support is enabled for each cluster",
			Name: "hypershift_cluster_limited_support_enabled",
		}, []string{"namespace", "name", "_id"}),
		clusterSilenceAlerts: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Help: "Indicates if alerts are silenced for each cluster",
			Name: "hypershift_cluster_silence_alerts",
		}, []string{"namespace", "name", "_id"}),
		hostedClusters: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "hypershift_hostedclusters",
			Help: "Number of HostedClusters by platform",
		}, []string{"platform"}),
		hostedClustersWithFailureCondition: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "hypershift_hostedclusters_failure_conditions",
			Help: "Total number of HostedClusters by platform with conditions in undesired state",
		}, []string{"platform", "condition"}),
		hostedClustersNodePools: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "hypershift_hostedcluster_nodepools",
			Help: "Number of NodePools associated with a given HostedCluster",
		}, []string{"cluster_name", "platform"}),
		nodePools: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "hypershift_nodepools",
			Help: "Number of NodePools by platform",
		}, []string{"platform"}),
		nodePoolsWithFailureCondition: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "hypershift_nodepools_failure_conditions",
			Help: "Total number of NodePools by platform with conditions in undesired state",
		}, []string{"platform", "condition"}),
		nodePoolSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "hypershift_nodepools_size",
			Help: "Number of replicas associated with a given NodePool",
		}, []string{"name", "platform"}),
		client: client,
		log:    log,
	}
}

func (m *hypershiftMetrics) Start(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := m.collect(ctx); err != nil {
				m.log.Error(err, "failed to collect metrics")
			}
		}
	}
}

func (m *hypershiftMetrics) collect(ctx context.Context) error {
	var clusters hyperv1.HostedClusterList
	if err := m.client.List(ctx, &clusters); err != nil {
		return fmt.Errorf("failed to list hostedclusters: %w", err)
	}
	m.observeHostedClusters(&clusters)
	var nodePools hyperv1.NodePoolList
	if err := m.client.List(ctx, &nodePools); err != nil {
		return fmt.Errorf("failed to list nodepools: %w", err)
	}
	if err := m.observeNodePools(ctx, &nodePools); err != nil {
		return err
	}
	return nil
}

func setupMetrics(mgr manager.Manager) error {
	metrics := newMetrics(mgr.GetClient(), mgr.GetLogger().WithName("metrics"))
	if err := crmetrics.Registry.Register(metrics.clusterProxy); err != nil {
		return fmt.Errorf("failed to to register clusterProxy metric: %w", err)
	}
	if err := crmetrics.Registry.Register(metrics.clusterIdentityProviders); err != nil {
		return fmt.Errorf("failed to to register clusterIdentityProviders metric: %w", err)
	}
	if err := crmetrics.Registry.Register(metrics.clusterLimitedSupportEnabled); err != nil {
		return fmt.Errorf("failed to to register clusterLimitedSupportEnabled metric: %w", err)
	}
	if err := crmetrics.Registry.Register(metrics.clusterSilenceAlerts); err != nil {
		return fmt.Errorf("failed to to register clusterSilenceAlerts metric: %w", err)
	}
	if err := crmetrics.Registry.Register(metrics.hostedClusters); err != nil {
		return fmt.Errorf("failed to to register hostedClusters metric: %w", err)
	}
	if err := crmetrics.Registry.Register(metrics.hostedClustersWithFailureCondition); err != nil {
		return fmt.Errorf("failed to to register hostedClustersWithCondition metric: %w", err)
	}
	if err := crmetrics.Registry.Register(metrics.hostedClustersNodePools); err != nil {
		return fmt.Errorf("failed to to register hostedClustersNodePools metric: %w", err)
	}
	if err := crmetrics.Registry.Register(metrics.nodePools); err != nil {
		return fmt.Errorf("failed to to register nodePools metric: %w", err)
	}
	if err := crmetrics.Registry.Register(metrics.nodePoolSize); err != nil {
		return fmt.Errorf("failed to to register nodePoolSize metric: %w", err)
	}
	if err := mgr.Add(metrics); err != nil {
		return fmt.Errorf("failed to add metrics runnable to manager: %w", err)
	}

	return nil
}

var expectedHCConditionStates = map[hyperv1.ConditionType]bool{
	hyperv1.HostedClusterAvailable:          true,
	hyperv1.HostedClusterProgressing:        false,
	hyperv1.IgnitionEndpointAvailable:       true,
	hyperv1.UnmanagedEtcdAvailable:          true,
	hyperv1.ValidHostedClusterConfiguration: true,
	hyperv1.SupportedHostedCluster:          true,
	hyperv1.ClusterVersionSucceeding:        true,
	hyperv1.ClusterVersionUpgradeable:       true,
	hyperv1.ReconciliationActive:            true,
	hyperv1.ValidOIDCConfiguration:          true,
}

func (m *hypershiftMetrics) observeHostedClusters(hostedClusters *hyperv1.HostedClusterList) {
	hcCount := newLabelCounter()
	hcByConditions := newLabelCounter()

	identityProvidersCounter := map[configv1.IdentityProviderType]float64{
		configv1.IdentityProviderTypeBasicAuth:     0,
		configv1.IdentityProviderTypeGitHub:        0,
		configv1.IdentityProviderTypeGitLab:        0,
		configv1.IdentityProviderTypeGoogle:        0,
		configv1.IdentityProviderTypeHTPasswd:      0,
		configv1.IdentityProviderTypeKeystone:      0,
		configv1.IdentityProviderTypeLDAP:          0,
		configv1.IdentityProviderTypeOpenID:        0,
		configv1.IdentityProviderTypeRequestHeader: 0,
	}

	for _, hc := range hostedClusters.Items {
		// Collect proxy metric.
		var proxyHTTP, proxyHTTPS, proxyTrustedCA string
		if hc.Spec.Configuration != nil && hc.Spec.Configuration.Proxy != nil {
			if hc.Spec.Configuration.Proxy.HTTPProxy != "" {
				proxyHTTP = "1"
			}
			if hc.Spec.Configuration.Proxy.HTTPSProxy != "" {
				proxyHTTPS = "1"
			}
			if hc.Spec.Configuration.Proxy.TrustedCA.Name != "" {
				proxyTrustedCA = "1"
			}
			m.clusterProxy.WithLabelValues(hc.Namespace, hc.Name, proxyHTTP, proxyHTTPS, proxyTrustedCA).Set(1)
		} else {
			m.clusterProxy.WithLabelValues(hc.Namespace, hc.Name, proxyHTTP, proxyHTTPS, proxyTrustedCA).Set(0)
		}

		// Group identityProviders by type.
		if hc.Spec.Configuration != nil && hc.Spec.Configuration.OAuth != nil {
			for _, identityProvider := range hc.Spec.Configuration.OAuth.IdentityProviders {
				identityProvidersCounter[identityProvider.Type] = identityProvidersCounter[identityProvider.Type] + 1
			}
		}

		// Collect limited support metric.
		if _, ok := hc.Labels[hyperv1.LimitedSupportLabel]; ok {
			m.clusterLimitedSupportEnabled.WithLabelValues(hc.Namespace, hc.Name, hc.Spec.ClusterID).Set(1)
		} else {
			m.clusterLimitedSupportEnabled.WithLabelValues(hc.Namespace, hc.Name, hc.Spec.ClusterID).Set(0)
		}

		// Collect silence alerts metric
		if _, ok := hc.Labels[hyperv1.SilenceClusterAlertsLabel]; ok {
			m.clusterSilenceAlerts.WithLabelValues(hc.Namespace, hc.Name, hc.Spec.ClusterID).Set(1)
		} else {
			m.clusterSilenceAlerts.WithLabelValues(hc.Namespace, hc.Name, hc.Spec.ClusterID).Set(0)
		}

		platform := string(hc.Spec.Platform.Type)
		hcCount.Add(platform)
		for _, cond := range hc.Status.Conditions {
			expectedState, known := expectedHCConditionStates[hyperv1.ConditionType(cond.Type)]
			if !known {
				continue
			}
			if expectedState {
				if cond.Status == metav1.ConditionFalse {
					hcByConditions.Add(platform, "not_"+cond.Type)
				}
			} else {
				if cond.Status == metav1.ConditionTrue {
					hcByConditions.Add(platform, cond.Type)
				}
			}
		}
	}

	// Collect identityProvider metric.
	for identityProvider, count := range identityProvidersCounter {
		m.clusterIdentityProviders.WithLabelValues(string(identityProvider)).Set(count)
	}

	for key, count := range hcCount.Counts() {
		labels := counterKeyToLabels(key)
		m.hostedClusters.WithLabelValues(labels...).Set(float64(count))
	}

	for key, count := range hcByConditions.Counts() {
		labels := counterKeyToLabels(key)
		m.hostedClustersWithFailureCondition.WithLabelValues(labels...).Set(float64(count))
	}
}

var expectedNPConditionStates = map[string]bool{
	hyperv1.NodePoolValidReleaseImageConditionType:  true,
	hyperv1.NodePoolValidPlatformImageType:          true,
	hyperv1.NodePoolValidMachineConfigConditionType: true,
	hyperv1.NodePoolReadyConditionType:              true,
	hyperv1.NodePoolUpdatingVersionConditionType:    false,
	hyperv1.NodePoolUpdatingConfigConditionType:     false,
}

func (m *hypershiftMetrics) observeNodePools(ctx context.Context, nodePools *hyperv1.NodePoolList) error {
	npByCluster := newLabelCounter()
	npCount := newLabelCounter()
	npByCondition := newLabelCounter()
	for _, np := range nodePools.Items {
		hc := &hyperv1.HostedCluster{}
		hc.Namespace = np.Namespace
		hc.Name = np.Spec.ClusterName
		hcPlatform := ""
		if err := m.client.Get(ctx, crclient.ObjectKeyFromObject(hc), hc); err == nil {
			hcPlatform = string(hc.Spec.Platform.Type)
			npByCluster.Add(crclient.ObjectKeyFromObject(hc).String(), hcPlatform)
		} else {
			m.log.Error(err, "cannot get hosted cluster for nodepool", "nodepool", crclient.ObjectKeyFromObject(&np).String())
		}
		platform := string(np.Spec.Platform.Type)
		npCount.Add(platform)

		for _, cond := range np.Status.Conditions {
			expectedState, known := expectedNPConditionStates[cond.Type]
			if !known {
				continue
			}
			if expectedState {
				if cond.Status == corev1.ConditionFalse {
					npByCondition.Add(platform, "not_"+cond.Type)
				}
			} else {
				if cond.Status == corev1.ConditionTrue {
					npByCondition.Add(platform, cond.Type)
				}
			}
		}
		m.nodePoolSize.WithLabelValues(crclient.ObjectKeyFromObject(&np).String(), platform).Set(float64(np.Status.Replicas))
	}
	for key, count := range npByCluster.Counts() {
		labels := counterKeyToLabels(key)
		m.hostedClustersNodePools.WithLabelValues(labels...).Set(float64(count))
	}
	for key, count := range npCount.Counts() {
		labels := counterKeyToLabels(key)
		m.nodePools.WithLabelValues(labels...).Set(float64(count))
	}
	for key, count := range npByCondition.Counts() {
		labels := counterKeyToLabels(key)
		m.nodePoolsWithFailureCondition.WithLabelValues(labels...).Set(float64(count))
	}
	return nil
}

type labelCounter struct {
	counts map[string]int
}

func newLabelCounter() *labelCounter {
	return &labelCounter{
		counts: map[string]int{},
	}
}

func (c *labelCounter) Add(values ...string) {
	key := strings.Join(values, "|")
	c.counts[key]++
}

func (c *labelCounter) Counts() map[string]int {
	return c.counts
}

func counterKeyToLabels(key string) []string {
	return strings.Split(key, "|")
}
