<script context="module">
    export async function load({ params }) {
        return {
            props : {
                clusterID: params.id
            }
        }
    }
</script>


<script lang="ts">
    import {
        UpdateEventDataRequest,
        EventDetailedResponse,
        EventModifiedHPAConfigData,
        ModifiedHPAConfig
    } from "$lib/api/type";
    import {getContext, onMount} from "svelte";
    import {EditEvent, EventData} from "$lib/constant/context";
    import dayjs from "dayjs";
    import DatetimeInput from "$lib/components/input/datetime.svelte"
    import {Writable} from "svelte/store";
    import {GetClusterHPAs} from "$lib/api/clusters";
    import duration from "dayjs/plugin/duration.js"
    import {UpdateEvent} from "$lib/api/event";

    dayjs.extend(duration)

    export let clusterID;

    const eventData: EventDetailedResponse = getContext<EventDetailedResponse>(EventData);
    const today = dayjs()

    let loadingHPA = true;
    let errorHPA = null;
    let hpaList = [];

    let eventName = eventData.name
    let startDate = dayjs(eventData.start_time)
    let endDate = dayjs(eventData.end_time)
    let executeConfigAt = dayjs(eventData.execute_config_at)
    let calculateNodePool = eventData.calculate_node_pool

    let editEvent = getContext<Writable<boolean>>(EditEvent)
    $editEvent = true

    const modifiedHPA: Map<string,ModifiedHPAConfig> = {}
    for (const hpa of eventData.modified_hpa_configs) {
        modifiedHPA[`${hpa.name}-${hpa.namespace}`] = hpa
    }

    onMount(async () => {
        try {
            hpaList = await GetClusterHPAs(clusterID)
            hpaList = hpaList.map(o => {
                const key = `${o.name}-${o.namespace}`
                const existingHPA = modifiedHPA[key]
                if (existingHPA) {
                    return {
                        ...o,
                        min_replicas: existingHPA.min_replicas,
                        max_replicas: existingHPA.max_replicas,
                        isSelected: true,
                        default_min_replicas: o.min_replicas,
                        default_max_replicas: o.max_replicas
                    }
                }
                return {
                    ...o,
                    isSelected: false,
                    default_min_replicas: o.min_replicas,
                    default_max_replicas: o.max_replicas
                }
            })

        } catch (e) {
            errorHPA = e
        } finally {
            loadingHPA = false
        }
    })

    let loading = false
    let success = false
    let errUpdate = null

    const updateEvent = async (selectedHPA: EventModifiedHPAConfigData[]): Promise<void> => {
        loading = true
        success = false
        const req: UpdateEventDataRequest = {
            event_id: eventData.id,
            end_time: endDate.toJSON(),
            modified_hpa_configs: selectedHPA,
            name: eventName,
            start_time: startDate.toJSON(),
            calculate_node_pool: calculateNodePool,
            execute_config_at: executeConfigAt.toJSON()
        }
        try {
            await UpdateEvent(req)
        } catch (e) {
            errUpdate = e
            return
        } finally {
            loading = false
        }
        success = true
        setTimeout(() => {
            location.replace(`/cluster/${clusterID}/event/${eventData.id}`)
        }, 2000)
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

<div class="mt-2">
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
        <DatetimeInput id="event-end" minDate={endDate} bind:date={endDate}/>
    </div>
    <div class="mb-2">
        <label for="exec-at">Exec At : </label>
        <DatetimeInput id="exec-at" minDate={today} maxDate={startDate} bind:date={executeConfigAt}/>
    </div>
    <div class="mb-2">
        <label for="calculate-node-pool">Calculate Node Pool : </label>
        <input type="checkbox" id="calculate-node-pool" bind:checked={calculateNodePool}>
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
                <button on:click={() => updateEvent(selectedHPA)} class="mb-2" disabled={loading || success}>Update</button>
                {#if !loading}
                    {#if errUpdate}
                        <h5>Error : {errUpdate}</h5>
                    {/if}
                    {#if success}
                        <h5>Update success, you will be redirected to cluster main page in 2 seconds</h5>
                    {/if}
                {/if}
            {:else}
                <h4>Error Getting HPA</h4>
            {/if}
        {/if}
    </div>
</div>
