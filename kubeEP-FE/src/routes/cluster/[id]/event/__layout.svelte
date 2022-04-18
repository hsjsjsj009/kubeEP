<script context="module" lang="ts">
    import {validate} from "uuid";

    export async function load({ params }) {
        const eventID = params.eventID;
        if (!validate(eventID)) {
            return {
                props: {
                    isErr: true,
                    errData: 'event id invalid'
                }
            };
        }
        return {
            props: {}
        }
    }
</script>

<script>
    export let isErr = false;
    export let errData = null;
</script>

{#if isErr}
    <div class="container flex flex-col flex-wrap content-center min-h-screen items-center">
        <p class="text-center">Error: {errData}</p>
    </div>
{:else }
    <slot/>
{/if}
