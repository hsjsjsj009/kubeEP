<script context="module" lang="ts">
    import {GetClusterDetailData} from "$lib/api/clusters";

    /** @type {import('./index.svelte').Load} */
    // eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
    export async function load({ params }) {
        const clusterID = params.id
        try {
            const data = await GetClusterDetailData(clusterID)
            return {
                props: {
                    isErr: false,
                    clusterData: data
                }
            }
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
    import type {ClusterDetailResponse} from "$lib/api/type";

    export let isErr = false
    export let clusterData: ClusterDetailResponse
    export let errData: never = null;

</script>

<div class="container flex flex-col flex-wrap content-center min-h-screen items-center">
    {#if isErr}
        <p class="text-center">Error: {errData}</p>
    {:else}
        <h2 class="text-center mt-2">Cluster</h2>
        <h3 class="text-center font-bold mb-2">{clusterData.cluster.name}</h3>
        <div>
            <h3 class="text-center mb-2">List HPA</h3>
            <table class="table-fixed border-separate border border-slate-400 text-center">
                <thead>
                <tr>
                    <th class="border border-slate-300 p-2"><h3>Name</h3></th>
                    <th class="border border-slate-300 p-2"><h3>Namespace</h3></th>
                    <th class="border border-slate-300 p-2"><h3>Min<br>Replicas</h3></th>
                    <th class="border border-slate-300 p-2"><h3>Max<br>Replicas</h3></th>
                    <th class="border border-slate-300 p-2"><h3>Current<br>Replicas</h3></th>
                </tr>
                </thead>
                <tbody>
                {#each clusterData.hpa_list as hpa}
                    <tr>
                        <td class="border border-slate-300 p-2"><h4>{hpa.name}</h4></td>
                        <td class="border border-slate-300 p-2"><h4>{hpa.namespace}</h4></td>
                        <td class="border border-slate-300 p-2"><h4>{hpa.min_replicas}</h4></td>
                        <td class="border border-slate-300 p-2"><h4>{hpa.max_replicas}</h4></td>
                        <td class="border border-slate-300 p-2"><h4>{hpa.current_replicas}</h4></td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>
