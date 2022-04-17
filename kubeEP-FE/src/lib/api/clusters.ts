import type {Cluster, ClusterDetailResponse} from "$lib/api/type";
import API from "$lib/api/util";

export const GetRegisteredClusters = async (): Promise<Cluster[]> => {
    return await API.get<Cluster[]>(`/cluster/list`).then((res) => {
        return res.data.data;
    });
};

export const GetClusterDetailData = async (clusterID: string): Promise<ClusterDetailResponse> => {
    return await API.get<ClusterDetailResponse>(`/cluster/${clusterID}/detail`).then(res => res.data.data)
}

export const GetClusterSimpleData = async (clusterID: string): Promise<Cluster> => {
    return await API.get<Cluster>(`/cluster/${clusterID}`).then(res => res.data.data)
}
