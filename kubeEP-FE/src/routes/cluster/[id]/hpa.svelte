<script context="module" lang="ts">
    /** @type {import('./hpa.svelte').Load} */
    // eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
    export async function load({params}) {
        const clusterID = params.id
        return {
            props: {
                clusterID
            }
        }
    }
</script>

<script lang="ts">
    import type {SimpleHPA} from "$lib/api/type";
    import {onMount} from "svelte";
    import {GetClusterHPAs} from "$lib/api/clusters";

    export let clusterID = "";
    let isErr = false
    let hpaList: SimpleHPA[]
    let errData: never = null;
    let loading = true

    onMount(async () => {
        loading = true
        try {
            hpaList = await GetClusterHPAs(clusterID)
        }catch (e) {
            isErr = true
            errData = e
        }finally {
            loading = false
        }

    })
</script>

{#if loading}
    <p class="text-center">Loading....</p>
{:else}
    {#if isErr}
        <p class="text-center">Error: {errData}</p>
    {:else}
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
            {#each hpaList as hpa}
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
    {/if}
{/if}
