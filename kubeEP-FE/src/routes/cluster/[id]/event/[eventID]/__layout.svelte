<script context="module" lang="ts">
    import {validate} from "uuid";

    export async function load({ params }) {
        const eventID = params.eventID;
        const clusterID = params.id
        if (!validate(eventID)) {
            return {
                props: {
                    isErr: true,
                    errData: 'event id invalid'
                }
            };
        }
        return {
            props: {
                clusterID,
                eventID
            }
        }
    }
</script>

<script lang="ts">
    import {EventDetailedResponse} from "$lib/api/type";
    import {onMount, setContext} from "svelte";
    import {EditEvent, EventData} from "$lib/constant/context";
    import NavLink from "$lib/components/navlink.svelte"
    import {Writable, writable} from "svelte/store";
    import {GetEventDetailByID, DeleteEvent} from "$lib/api/event";


    export let isErr = false;
    export let errData = null;
    export let eventID;
    export let clusterID;

    let loading = false

    let editEvent = writable<boolean>(true)

    let startTime
    let endTime

    let eventData: EventDetailedResponse;

    $: if (eventData) {
        setContext<EventDetailedResponse>(EventData, eventData)
        startTime = new Date(eventData.start_time);
        endTime = new Date(eventData.end_time);
    }

    onMount(async () => {
        if (!isErr) {
            loading = true
            try {
                const data = await GetEventDetailByID(eventID);
                if (data.cluster.id !== clusterID) {
                    isErr = true
                    errData = 'cluster id mismatch'
                    return
                }
                eventData = data
            } catch (e) {
                isErr = true
                errData = e
            } finally {
                loading = false
            }
        }
    })

    let loadingDelete = false
    let errDelete = null
    let infoMessage = null

    const deleteEvent = async () => {
        loadingDelete = true
        try {
            await DeleteEvent(eventID)
            infoMessage = "Success delete event, will be redirect to event list page in 2 seconds"
            setTimeout(() => {
                location.replace(`/cluster/${clusterID}/event/list`)
            },2000)
        } catch (e) {
            errDelete = e
        } finally {
            loadingDelete = false
        }
    }



    setContext<Writable<boolean>>(EditEvent, editEvent)
 </script>

<svelte:head>
    {#if eventData}
        <title>Event {eventData.name}</title>
    {/if}
</svelte:head>

{#if isErr}
    <div class="container flex flex-col flex-wrap content-center min-h-screen items-center">
        <p class="text-center">Error: {errData}</p>
    </div>
{:else if loading}
    <div class="container flex flex-col flex-wrap content-center min-h-screen items-center">
        <p class="text-center">Loading...</p>
    </div>
{:else if eventData && !loading }
    <div class="container mx-auto text-center">
        <h2>Event</h2>
        {#if !$editEvent}
        <h4>Name : {eventData.name}</h4>
        <h4>Start Time : {startTime.toLocaleString()}</h4>
        <h4>End Time : {endTime.toLocaleString()}</h4>
        <h4>Status : {eventData.status}</h4>
        <h4>Calculate Node Pool : {eventData.calculate_node_pool ? "Yes" : "No"}</h4>
        <button on:click={deleteEvent} class="mb-2 rounded-3xl border-2 p-2 mx-2"><h4>Delete</h4></button>
        {#if errDelete}
            <p class="mb-2">Error Delete : {errDelete}</p>
        {/if}
        {#if infoMessage}
            <p class="mb-2">{infoMessage}</p>
        {/if}
        {/if}
        <div class="flex flex-row justify-center content-center mt-2">
            <NavLink href="/cluster/{eventData.cluster.id}/event/{eventData.id}"><h3>Detail</h3></NavLink>
            <NavLink href="/cluster/{eventData.cluster.id}/event/{eventData.id}/edit"><h3>Edit</h3></NavLink>
        </div>
        <slot/>
    </div>
{/if}
