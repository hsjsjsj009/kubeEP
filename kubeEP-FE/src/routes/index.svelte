<script context="module" lang="ts">
    import {GetRegisteredClusters} from "$lib/api/clusters";

    export async function load() {
        try {
            const data = await GetRegisteredClusters();
            return {
                props: {
                    isErr: false,
                    clustersData: data
                }
            };
        } catch (e) {
            return {
                props: {
                    isErr: true,
                    errData: e
                }
            };
        }
    }
</script>

<script lang="ts">
    import type {Cluster} from "$lib/api/type";

    export let isErr: boolean;
    export let clustersData: Cluster[] | null = null
    export let errData: never = null;
</script>

<div class="container flex flex-col flex-wrap content-center h-screen">
    <div class="flex flex-row justify-center content-center mt-5">
        <h1 class="font-bold text-center mx-[3rem]">Kube EP</h1>
    </div>
    <div class="flex flex-row justify-center content-center mt-2">
        <a class="rounded-3xl border-black border-2 p-2" href="/register/cluster">Register Clusters</a>
    </div>
    <table class="table-fixed mt-5 border-separate border border-slate-400">
        <thead>
            <tr>
                <th class="border border-slate-300 p-2"><h3>Cluster</h3></th>
                <th class="border border-slate-300 p-2"><h3>Datacenter</h3></th>
            </tr>
        </thead>
        <tbody>
        {#if isErr}
            <tr>
                <td colspan="2" class="border border-slate-300 p-2">Error : {errData}</td>
            </tr>
        {:else}
            {#each clustersData as cluster}
                <tr>
                    <td class="border border-slate-300 p-2"><h4>{cluster.name}</h4></td>
                    <td class="border border-slate-300 p-2 text-center"><h4>{cluster.datacenter_name} - {cluster.datacenter}</h4></td>
                </tr>
            {/each}
        {/if}
        </tbody>
    </table>
</div>
