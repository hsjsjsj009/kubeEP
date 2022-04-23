<script context="module" lang="ts">
    export async function load({params}) {
        const clusterID = params.id
        return {
            props : {
                isErr: false,
                clusterID
            }
        }
    }
</script>
<script lang="ts">
    import {GetClusterHPAs} from "$lib/api/clusters";
    import {onMount} from "svelte";
    import dayjs from "dayjs";
    import duration from "dayjs/plugin/duration"
    import DatetimeInput from "$lib/components/input/datetime.svelte"
    import {CreateEvent} from "$lib/api/event";
    import type {EventDataRequest, EventModifiedHPAConfigData} from "$lib/api/type";
    import {goto} from "$app/navigation";

    export let clusterID;

    let loadingHPA = true;
    let errorHPA = null;
    let hpaList = [];
    let eventName = "";
    dayjs.extend(duration)
    onMount(async () => {
        try {
            hpaList = await GetClusterHPAs(clusterID)
            hpaList = hpaList.map(o => (
                {
                    ...o,
                    isSelected: false,
                    default_min_replicas: o.min_replicas,
                    default_max_replicas: o.max_replicas
                }
            ))
        } catch (e) {
            errorHPA = e
        } finally {
            loadingHPA = false
        }
    })
    let today = dayjs()
    let startDate = dayjs()
    let endDate = startDate.add(dayjs.duration({days:1}))
    let loading = false
    let success = false
    let errCreate = null

    const createEvent = async (selectedHPA: EventModifiedHPAConfigData[]): Promise<void> => {
        loading = true
        success = false
        const req: EventDataRequest = {
            cluster_id: clusterID,
            end_time: endDate.toJSON(),
            modified_hpa_configs: selectedHPA,
            name: eventName,
            start_time: startDate.toJSON()
        }
        try {
            await CreateEvent(req)
        } catch (e) {
            errCreate = e
            return
        } finally {
            loading = false
        }
        success = true
        setTimeout(() => {
            goto(`/cluster/${clusterID}/event/list`)
        }, 2000)
    }

    $: {
        if (startDate.isAfter(endDate) || startDate.isSame(endDate)) {
            endDate = startDate.add(dayjs.duration({days:1}))
        }
    }

    $: selectedHPA = hpaList.filter(
        o => o.isSelected
    ).map(
        o => (
            {
                name: o.name,
                namespace: o.namespace,
                min_replicas:o.min_replicas,
                max_replicas:o.max_replicas
            }
        )
    )
</script>

<div class="mb-2">
    <label for="event-name">Name : </label>
    <input type="text" id="event-name" bind:value={eventName}>
</div>
<div class="mb-2">
    <label for="event-start">Start Time : </label>
    <DatetimeInput minDate={today} id="event-start" bind:date={startDate}/>
</div>
<div class="mb-2">
    <label for="event-end">End Time : </label>
    <DatetimeInput id="event-end" minDate={startDate.add(dayjs.duration({minutes:1}))} bind:date={endDate}/>
</div>
<div class="mb-2 text-center">
    {#if loadingHPA}
        <h3>Loading HPA List</h3>
    {/if}
    {#if !loadingHPA }
        {#if !errorHPA}
            <h3>List HPA</h3>
            <table class="mb-2 table-fixed border-separate border border-slate-400 text-center">
                <thead>
                <tr>
                    <th class="border border-slate-300 p-2"><h3>Name</h3></th>
                    <th class="border border-slate-300 p-2"><h3>Namespace</h3></th>
                    <th class="border border-slate-300 p-2"><h3>Min<br>Replicas</h3></th>
                    <th class="border border-slate-300 p-2"><h3>Max<br>Replicas</h3></th>
                    <th class="border border-slate-300 p-2"><h3>Current<br>Replicas</h3></th>
                    <th class="border border-slate-300 p-2"><h3>Add?</h3></th>
                </tr>
                </thead>
                <tbody>
                {#each hpaList as hpa}
                    <tr>
                        <td class="border border-slate-300 p-2"><h4>{hpa.name}</h4></td>
                        <td class="border border-slate-300 p-2"><h4>{hpa.namespace}</h4></td>
                        <td class="border border-slate-300 p-2">
                            {#if hpa.isSelected}
                                <input class="w-[4rem] text-center" type="number" min={hpa.default_min_replicas} bind:value={hpa.min_replicas}>
                            {:else}
                                <h4 class="text-center">{hpa.default_min_replicas}</h4>
                            {/if}
                        </td>
                        <td class="border border-slate-300 p-2">
                            {#if hpa.isSelected}
                                <input class="w-[4rem] text-center" type="number" min={hpa.default_max_replicas} bind:value={hpa.max_replicas}>
                            {:else}
                                <h4 class="text-center">{hpa.default_max_replicas}</h4>
                            {/if}
                        </td>
                        <td class="border border-slate-300 p-2"><h4>{hpa.current_replicas}</h4></td>
                        <td class="border border-slate-300 p-2"><input type="checkbox" bind:checked={hpa.isSelected}></td>
                    </tr>
                {/each}
                </tbody>
            </table>
            <button on:click={() => createEvent(selectedHPA)} class="mb-2" disabled={loading || success}>Create</button>
            {#if !loading}
                {#if errCreate}
                    <h5>Error : {errCreate}</h5>
                {/if}
                {#if success}
                    <h5>Creation success, you will be redirected to cluster main page in 2 seconds</h5>
                {/if}
            {/if}
        {:else}
            <h4>Error Getting HPA</h4>
        {/if}
    {/if}
</div>
