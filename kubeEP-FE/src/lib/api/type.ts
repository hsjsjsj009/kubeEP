import type { Datacenter } from '$lib/types/datacenter';

export interface Cluster {
	id: string;
	name: string;
	datacenter: Datacenter
	datacenter_name: string
}

export interface ModifiedHPAConfig {
	id: string;
	name: string;
	namespace: string;
	min_replicas: number;
	max_replicas: number;
}

export enum EventStatus {
	Failed = 'FAILED',
	Success = 'SUCCESS',
	Executing = 'EXECUTING',
	Prescaled = 'PRESCALED',
	Watching = 'WATCHING',
	Pending = 'PENDING'
}

export interface UpdatedNodePool {
	id: string;
	node_pool_name: string;
	max_node: number
}

export interface EventSimpleResponse {
	id: string;
	name: string;
	start_time: string;
	end_time: string;
	status: EventStatus;
}

export interface EventDetailedResponse extends EventSimpleResponse {
	created_at: string;
	updated_at: string;
	cluster: Cluster;
	modified_hpa_configs: ModifiedHPAConfig[];
	updated_node_pools: UpdatedNodePool[];
	execute_config_at: string;
	calculate_node_pool: boolean;
}

export interface NodePoolStatus {
	created_at: string;
	count: number;
}

export interface HPAStatus {
	created_at: string;
	replicas: number;
	available_replicas: number;
	ready_replicas: number;
	unavailable_replicas: number;
}

export interface RegisterGCPDatacenterRequest {
	name: string;
	// eslint-disable-next-line @typescript-eslint/ban-types
	sa_key_credentials: Object;
	is_temporary: boolean
}

export interface DatacenterResponse {
	datacenter_id: string
	is_temporary: boolean
}

export interface GCPCluster extends Cluster {
	location: string
}

export interface GCPClusterResponse {
	clusters: GCPCluster[]
	is_temporary_datacenter: boolean
}

export interface GCPRegisterClustersRequest {
	clusters_name: string[]
	datacenter_id: string
	is_datacenter_temporary: boolean
}

export interface SimpleHPA {
	name: string
	namespace: string
	min_replicas: number | null
	max_replicas: number
	current_replicas: number

}

export interface ClusterDetailResponse {
	cluster: Cluster
	hpa_list: SimpleHPA[]
}

export interface EventModifiedHPAConfigData {
	name: string
	namespace: string
	min_replicas: number
	max_replicas: number
}

export interface EventDataRequest {
	name: string
	start_time: string
	end_time: string
	cluster_id: string
	modified_hpa_configs: EventModifiedHPAConfigData[]
	calculate_node_pool: boolean
	execute_config_at: string
}

export interface EventCreationResponse {
	event_id: string
}

export interface UpdateEventDataRequest {
	name: string
	start_time: string
	end_time: string
	event_id: string
	modified_hpa_configs: EventModifiedHPAConfigData[],
	calculate_node_pool: boolean
	execute_config_at: string
}

