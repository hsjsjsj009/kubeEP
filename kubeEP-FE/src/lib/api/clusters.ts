import type {Cluster, ClusterDetailResponse, SimpleHPA} from "$lib/api/type";
import API from "$lib/api/util";

export const GetRegisteredClusters = async (): Promise<Cluster[]> => {
    return await API.get<Cluster[]>(`/cluster/list`).then((res) => {
        return res.data.data;
    });
};

export const GetClusterHPAs = async (clusterID: string): Promise<SimpleHPA[]> => {
    return await API.get<SimpleHPA[]>(`/cluster/${clusterID}/hpa`).then(res => res.data.data)
}

export const GetClusterSimpleData = async (clusterID: string): Promise<Cluster> => {
    return await API.get<Cluster>(`/cluster/${clusterID}`).then(res => res.data.data)
}
