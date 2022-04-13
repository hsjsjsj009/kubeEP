<script context="module" lang="ts">
	import { GetEventDetailByID } from '$lib/api/event';
	import { validate } from 'uuid';

	/** @type {import('./[id]').Load} */
	// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
	export async function load({ params }) {
		const eventID = params.id;
		if (!validate(eventID)) {
			return {
				props: {
					isErr: true,
					errData: 'event id invalid'
				}
			};
		}

		try {
			const data = await GetEventDetailByID(eventID);
			return {
				props: {
					isErr: false,
					eventData: data
				}
			};
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
	import { browser } from '$app/env';
	import type { EventDetailedResponse } from '$lib/api/type';

	export let isErr = false;
	export let eventData: EventDetailedResponse;
	export let errData = null;

	const startTime = new Date(eventData.start_time);
	const endTime = new Date(eventData.end_time);
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
		{#if eventData.status === 'SUCCESS' && browser}
			<div class="mt-2">
				<h2 class="font-bold">Monitoring</h2>
				<div class="flex mt-1 mb-3">
					<div class="flex-1 mb-3 mx-3 overflow-y-auto max-h-[75vh]">
						{#await import('$lib/components/chart/node-pool-chart.svelte')}
							<h3>Loading Component...</h3>
						{:then c}
							{#each eventData.updated_node_pools as updatedNodePool}
								<svelte:component
									this={c.default}
									id={updatedNodePool.id}
									name={updatedNodePool.node_pool_name}
									maxNode={updatedNodePool.max_node}
								/>
							{/each}
						{:catch e}
							<h3>Error Loading Component : {e}</h3>
						{/await}
					</div>
					<div class="flex-1 mb-3 mx-3 overflow-y-auto max-h-[75vh]">
						{#await import('$lib/components/chart/hpa-chart.svelte')}
							<h3>Loading Component...</h3>
						{:then c}
							{#each eventData.modified_hpa_configs as modifiedHPAConfig}
								<svelte:component
									this={c.default}
									name={modifiedHPAConfig.name}
									namespace={modifiedHPAConfig.namespace}
									minPods={modifiedHPAConfig.min_replicas}
									maxPods={modifiedHPAConfig.max_replicas}
									id={modifiedHPAConfig.id}
								/>
							{/each}
						{:catch e}
							<h3>Error Loading Component : {e}</h3>
						{/await}
					</div>
				</div>
			</div>
		{/if}
	</div>
{/if}
