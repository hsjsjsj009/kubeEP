import type {Cluster} from "$lib/api/type";
import API from "$lib/api/util";

export const GetRegisteredClusters = async (): Promise<Cluster[]> => {
    return await API.get<Cluster[]>(`/clusters`).then((res) => {
        return res.data.data;
    });
};


