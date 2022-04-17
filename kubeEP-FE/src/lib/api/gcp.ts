import type {
    DatacenterResponse, GCPCluster,
    GCPClusterResponse,
    GCPRegisterClustersRequest,
    RegisterGCPDatacenterRequest
} from "$lib/api/type";
import API from "$lib/api/util";

export const RegisterGCPDatacenter = async (data: RegisterGCPDatacenterRequest): Promise<DatacenterResponse> => {
    return await API.post<RegisterGCPDatacenterRequest,DatacenterResponse>(`/gcp/register/datacenter`, data).then((res) => {
        return res.data.data;
    });
};

export const GetGCPClusters = async (datacenterID: string): Promise<GCPClusterResponse> => {
    return await API.get<GCPClusterResponse>(`/gcp/clusters?datacenter_id=${datacenterID}`).then(res => res.data.data)
}

export const RegisterGCPClusters = async (data: GCPRegisterClustersRequest): Promise<GCPCluster[]> => {
    return await API.post<GCPRegisterClustersRequest, GCPCluster[]>(`/gcp/register/clusters`, data).then(res => res.data.data)
}
