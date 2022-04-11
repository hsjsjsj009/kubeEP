import type { HPAStatus, NodePoolStatus } from '$lib/api/type';
import API from '$lib/api/util';

export const GetEventNodePoolStatistics = async (
	updatedNodePoolID: string
): Promise<NodePoolStatus[]> => {
	return await API.get<NodePoolStatus[]>(`/event/status/node-pool/${updatedNodePoolID}`).then(
		(res) => res.data.data
	);
};

export const GetEventHPAStatistics = async (scheduledHPAConfigID: string): Promise<HPAStatus[]> => {
	return await API.get<HPAStatus[]>(`/event/status/hpa/${scheduledHPAConfigID}`).then(
		(res) => res.data.data
	);
};
