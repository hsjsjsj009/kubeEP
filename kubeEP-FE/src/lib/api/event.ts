import API from '$lib/api/util';
import type { EventDetailedResponse } from '$lib/api/type';

export const GetEventDetailByID = async (id: string): Promise<EventDetailedResponse> => {
	return await API.get<EventDetailedResponse>(`/event/${id}`).then((res) => {
		return res.data.data;
	});
};
