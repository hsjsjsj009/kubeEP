<script context="module" lang="ts">
    import {GetEventDetailByID} from "$lib/api/event";
    import { validate } from "uuid";
    import NodePoolChart from "$lib/components/chart/node-pool-chart.svelte";
    import HPAChart from "$lib/components/chart/hpa-chart.svelte";

    /** @type {import('./[id]').Load} */
    // eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
    export async function load({ params }) {
        const eventID = params.id
        if (!validate(eventID)) {
            return {
                props: {
                    isErr: true,
                    errData: "event id invalid"
                }
            }
        }

        try {
            const data = await GetEventDetailByID(eventID)
            return {
                props: {
                    isErr: false,
                    eventData: data
                }
            };
        }
        catch (e) {
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
    import type {EventDetailedResponse} from "$lib/api/type";

    export let isErr = false;
    export let eventData: EventDetailedResponse;
    export let errData = null;

    const startTime = new Date(eventData.start_time)
    const endTime = new Date(eventData.end_time)
</script>

<svelte:head>
    {#if isErr}
        <title>Event</title>
    {:else}
        <title>Event {eventData.name}</title>
    {/if}
</svelte:head>

{#if isErr}
<div class="container mx-auto text-center">
    {errData}
</div>
{:else}
<div class="container mx-auto text-center">
    <h1>Event</h1>
    <h1>Name : {eventData.name}</h1>
    <h1>Start Time : {startTime.toLocaleString()}</h1>
    <h1>End Time : {endTime.toLocaleString()}</h1>
    <h1>Status : {eventData.status}</h1>
    {#if eventData.status === "SUCCESS"}
        <div class="mt-2">
            <h2 class="font-bold">Monitoring</h2>
            <div class="flex mt-1">
                <div class="flex-1">
                    {#each eventData.updated_node_pools as updatedNodePool}
                        <NodePoolChart id={updatedNodePool.id} name={updatedNodePool.node_pool_name} />
                    {/each}
                </div>
                <div class="flex-1">
                    {#each eventData.modified_hpa_configs as modifiedHPAConfig}
                        <HPAChart name={modifiedHPAConfig.name} namespace={modifiedHPAConfig.namespace} id={modifiedHPAConfig.id} />
                    {/each}
                </div>
            </div>
        </div>
    {/if}
</div>

{/if}

