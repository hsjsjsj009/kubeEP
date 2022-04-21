<script context="module" lang="ts">
    import {validate} from "uuid";
    import {GetClusterSimpleData} from "$lib/api/clusters";
    import type {Cluster} from "$lib/api/type";

    export async function load({ params }) {
        const clusterID = params.id
        if (!validate(clusterID)) {
            return {
                props: {
                    isErr: true,
                    errData: 'cluster id invalid'
                }
            }
        }
        try {
            const data = await GetClusterSimpleData(clusterID)
            data.id = clusterID

            return {
                props: {
                    clusterData: data,
                }
            }
        }catch (e) {
            return {
                props: {
                    isErr: true,
                    errData: e
                }
            }
        }
    }
</script>

<script lang="ts">
    import NavLink from '$lib/components/navlink.svelte'
    import {setContext} from "svelte";
    import {ClusterData} from "$lib/constant/context";

    export let isErr = false;
    export let errData = null;
    export let clusterData: Cluster;

    if (clusterData) {
        setContext<Cluster>(ClusterData, clusterData)
    }
</script>

<svelte:head>
    {#if clusterData}
        <title>Cluster {clusterData.name}</title>
    {/if}
</svelte:head>

{#if isErr}
    <div class="container flex flex-col flex-wrap content-center min-h-screen items-center">
        <p class="text-center">Error: {errData}</p>
    </div>
{:else }
    <div class="container flex flex-col flex-wrap content-center min-h-screen items-center">
        <h2 class="text-center mt-2">Cluster</h2>
        <h3 class="text-center font-bold mb-2">{clusterData.name}</h3>
        <div class="flex flex-col items-center">
            <div class="flex flex-row justify-center mb-3">
                <NavLink href="/cluster/{clusterData.id}/hpa"><h3>List HPA</h3></NavLink>
                <NavLink href="/cluster/{clusterData.id}/event/list"><h3>List Events</h3></NavLink>
                <NavLink href="/cluster/{clusterData.id}/event/create"><h3>Create Event</h3></NavLink>
            </div>
            <slot />
        </div>
    </div>
{/if}
