<script context="module" lang="ts">
    import {GetEventListByClusterID} from "$lib/api/event";
    import {GetClusterSimpleData} from "$lib/api/clusters";

    /** @type {import('./events.svelte').Load} */
    export async function load({ params }) {
        const clusterID = params.id
        try {
            const reqEventList = GetEventListByClusterID(clusterID)
            const reqClusterData = GetClusterSimpleData(clusterID)
            const data = await Promise.all([reqEventList, reqClusterData]).then(values => {
                return values
            })
            return {
                props: {
                    isErr: false,
                    eventList: data[0],
                    clusterData: data[1]
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
    import type {Cluster, EventSimpleResponse} from "$lib/api/type";

    export let isErr = false;
    export let errData = null;
    export let eventList: EventSimpleResponse[] = null;
    export let clusterData: Cluster = null
</script>

<div class="container flex flex-col flex-wrap content-center min-h-screen items-center">
    {#if isErr}
        <p class="text-center">Error: {errData}</p>
    {:else}
        <h2 class="text-center mt-2">Cluster</h2>
        <h3 class="text-center font-bold mb-2">{clusterData.name}</h3>
        <div>
            <h3 class="text-center mb-2">List Event</h3>
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
                            <td class="border border-slate-300 p-2"><h4>{event.name}</h4></td>
                            <td class="border border-slate-300 p-2"><h4>{new Date(event.start_time).toLocaleString()}</h4></td>
                            <td class="border border-slate-300 p-2"><h4>{new Date(event.end_time).toLocaleString()}</h4></td>
                            <td class="border border-slate-300 p-2"><h4>{event.status}</h4></td>
                        </tr>
                    {/each}
                {/if}
                </tbody>
            </table>
        </div>
    {/if}
</div>
