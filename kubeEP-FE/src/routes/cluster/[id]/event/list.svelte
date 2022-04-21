<script context="module" lang="ts">
    import {GetEventListByClusterID} from "$lib/api/event";

    /** @type {import('./list.svelte').Load} */
    export async function load({ params }) {
        const clusterID = params.id
        try {
            const data = await GetEventListByClusterID(clusterID)
            return {
                props: {
                    isErr: false,
                    eventList: data,
                    clusterID
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
    import type {EventSimpleResponse} from "$lib/api/type";

    export let isErr = false;
    export let errData = null;
    export let eventList: EventSimpleResponse[] = null
    export let clusterID = ""
</script>

{#if isErr}
    <p class="text-center">Error: {errData}</p>
{:else}
    <table class="table-fixed border-separate border border-slate-400 text-center">
        <thead>
        <tr>
            <th class="border border-slate-300 p-2"><h3>Name</h3></th>
            <th class="border border-slate-300 p-2"><h3>Start Time</h3></th>
            <th class="border border-slate-300 p-2"><h3>End Time</h3></th>
            <th class="border border-slate-300 p-2"><h3>Status</h3></th>
        </tr>
        </thead>
        <tbody>
        {#if eventList.length === 0}
            <tr>
                <td class="border border-slate-300 p-2" colspan="4"><h4>Empty</h4></td>
            </tr>
        {:else}
            {#each eventList as event}
                <tr>
                    <td class="border border-slate-300 p-2"><a href="/cluster/{clusterID}/event/{event.id}"><h4>{event.name}</h4></a></td>
                    <td class="border border-slate-300 p-2"><h4>{new Date(event.start_time).toLocaleString()}</h4></td>
                    <td class="border border-slate-300 p-2"><h4>{new Date(event.end_time).toLocaleString()}</h4></td>
                    <td class="border border-slate-300 p-2"><h4>{event.status}</h4></td>
                </tr>
            {/each}
        {/if}
        </tbody>
    </table>
{/if}
