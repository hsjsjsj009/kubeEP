import API from '$lib/api/util';
import type {EventDetailedResponse, EventSimpleResponse} from '$lib/api/type';

export const GetEventDetailByID = async (id: string): Promise<EventDetailedResponse> => {
	return await API.get<EventDetailedResponse>(`/event/${id}`).then((res) => {
		return res.data.data;
	});
};

export const GetEventListByClusterID = async (clusterID: string): Promise<EventSimpleResponse[]> => {
	return await API.get<EventSimpleResponse[]>(`/event/list?cluster_id=${clusterID}`).then(res => res.data.data)
}
