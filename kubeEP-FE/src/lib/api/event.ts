import API from '$lib/api/util';
import type {
	EventCreationResponse,
	EventDataRequest,
	EventDetailedResponse,
	EventSimpleResponse,
	UpdateEventDataRequest
} from '$lib/api/type';

export const GetEventDetailByID = async (id: string): Promise<EventDetailedResponse> => {
	return await API.get<EventDetailedResponse>(`/event/${id}`).then((res) => {
		return res.data.data;
	});
};

export const GetEventListByClusterID = async (clusterID: string): Promise<EventSimpleResponse[]> => {
	return await API.get<EventSimpleResponse[]>(`/event/list?cluster_id=${clusterID}`).then(res => res.data.data)
}

export const CreateEvent = async (data: EventDataRequest): Promise<EventCreationResponse> => {
	return await API.post<EventDataRequest, EventCreationResponse>(`/event/register`, data).then(res => res.data.data)
}

export const UpdateEvent = async (data: UpdateEventDataRequest): Promise<EventCreationResponse> => {
	return await API.put<UpdateEventDataRequest, EventCreationResponse>(`/event/update`, data).then(res => res.data.data)
}

export const DeleteEvent = async (eventID: string): Promise<string> => {
	return await API.delete<void,string>(`/event/${eventID}`).then(res => res.data.data)
}
