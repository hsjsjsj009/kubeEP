<script lang="ts">

    import {Datacenter} from "$lib/types/datacenter.ts";
    import {browser} from "$app/env"
    import { JSONEditor } from 'svelte-jsoneditor'
    import type {GCPRegisterClustersRequest, RegisterGCPDatacenterRequest} from "$lib/api/type";
    import {GetGCPClusters, RegisterGCPClusters, RegisterGCPDatacenter} from "$lib/api/gcp"
    import {goto} from "$app/navigation";

    let selectedDatacenter = Datacenter.GCP
    let datacenterName

    let saContent = {
        text: undefined,
        json: {
            "type": "service_account",
            "project_id": "project_id",
            "private_key_id": "private_key_id",
            "private_key": "private_key",
            "client_email": "client_email",
            "client_id": "client_id",
            "auth_uri": "https://accounts.google.com/o/oauth2/auth",
            "token_uri": "https://oauth2.googleapis.com/token",
            "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
            "client_x509_cert_url": "client_x509_cert_url"
        }
    }

    let mode = "code"

    let loading = false
    let registerClusterDone = false
    let registerClusterLoading = false
    let registerClusterError = null
    let gcpClusters = []
    let fetchGCPClusters = false
    let fetchGCPClusterError = null
    let datacenterID;

    // eslint-disable-next-line @typescript-eslint/ban-types
    const gcpFetchCluster = async (name: string, saCreds: Object): Promise<void> => {
        fetchGCPClusters = true
        loading = true
        try {
            const registerGCPCluster: RegisterGCPDatacenterRequest = {
                is_temporary: true,
                name,
                sa_key_credentials: saCreds
            }
            const registerGCPResponse = await RegisterGCPDatacenter(registerGCPCluster)
            datacenterID = registerGCPResponse.datacenter_id
            const gcpClustersResponse = await GetGCPClusters(datacenterID)
            gcpClusters = gcpClustersResponse.clusters.map(o => {
                return {
                    ...o,
                    isPicked: false
                }
            })
        } catch (e) {
            fetchGCPClusterError = e
        } finally {
            loading = false
        }
    }

    $: pickedClusters = gcpClusters.filter(o => o.isPicked).map(o => o.name)

    const gcpRegisterClusters = async (clusters: string[], datacenterID: string): Promise<void> => {
        registerClusterLoading = true
        try {
            const data: GCPRegisterClustersRequest = {
                clusters_name: clusters,
                datacenter_id: datacenterID,
                is_datacenter_temporary: true
            }
            await RegisterGCPClusters(data)
            registerClusterDone = true
            setTimeout(() => {
                goto("/")
            }, 2000)
        } catch (e) {
            registerClusterError = e
        } finally {
            registerClusterLoading = false
        }
    }
</script>

<svelte:head>
    <title>Register Clusters</title>
</svelte:head>

<div class="container flex flex-col flex-wrap content-center min-h-screen justify-center">
    <h2 class="font-bold text-center">Register Cluster</h2>
    <div class="mt-3">
        <label for="datacenter">Datacenter :</label>
        <select bind:value={selectedDatacenter} id="datacenter" class="text-lg rounded">
            {#each Object.values(Datacenter) as value}
                <option value={value}>{value}</option>
            {/each}
        </select>
    </div>
    <div class="mt-1">
        <label for="datacenter-name">Datacenter Name :</label>
        <input type="text" class="border-black border rounded px-1" bind:value={datacenterName} id="datacenter-name">
    </div>
    {#if selectedDatacenter === Datacenter.GCP}
        {#if browser}
            <p>SA Credentials :</p>
            <div class="mt-1 mb-1">
                <JSONEditor bind:content={saContent} bind:mode={mode} />
            </div>
            <div class="mt-1">
                <button on:click={() => gcpFetchCluster(datacenterName, mode === "code" ? JSON.parse(saContent.text) : saContent.json)} disabled={loading} class="rounded-3xl border-black border-2 p-2 disabled:border-gray-500 disabled:text-gray-500">
                    {#if loading}
                        Loading...
                    {:else}
                        Fetch Cluster
                    {/if}
                </button>
            </div>
            {#if fetchGCPClusters}
                <div class="mt-2 flex flex-col">
                    <h3 class="font-bold text-center">List GCP Cluster</h3>
                    <table class="table-fixed mt-5 border-separate border border-slate-400">
                        <thead>
                            <tr>
                                <th class="border border-slate-300 p-2 w-3"><h4>Cluster</h4></th>
                                <th class="border border-slate-300 p-2 w-1"><h4>Add Cluster?</h4></th>
                            </tr>
                        </thead>
                        <tbody>
                            {#if loading}
                                <tr>
                                    <td colspan="2" class="border border-slate-300 p-2">Loading...</td>
                                </tr>
                            {:else}
                                {#if fetchGCPClusterError}
                                    <tr>
                                        <td colspan="2" class="border border-slate-300 p-2">Error : {fetchGCPClusterError}</td>
                                    </tr>
                                {:else}
                                    {#each gcpClusters as cluster}
                                        <tr>
                                            <td class="border border-slate-300 p-2 text-center"><h4>{cluster.name}</h4></td>
                                            <td class="border border-slate-300 p-2 text-center">
                                                <input type="checkbox" bind:checked={cluster.isPicked}>
                                            </td>
                                        </tr>
                                    {/each}
                                {/if}
                            {/if}
                        </tbody>
                    </table>
                    {#if !loading && !fetchGCPClusterError}
                        <button on:click={() => gcpRegisterClusters(pickedClusters, datacenterID)} disabled={registerClusterLoading} class="rounded-3xl border-black border-2 p-2 mt-3 max-w-fit disabled:border-gray-500 disabled:text-gray-500">
                            {#if registerClusterLoading}
                                Loading...
                            {:else}
                                Register Clusters
                            {/if}
                        </button>
                        {#if registerClusterDone}
                            <p>Register cluster done, will be redirected to homepage in 2 seconds</p>
                        {/if}
                    {/if}
                </div>
            {/if}
        {:else}
            <h5>Loading....</h5>
        {/if}
    {/if}
</div>
