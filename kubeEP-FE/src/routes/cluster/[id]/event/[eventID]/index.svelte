<script lang="ts">
	import { browser } from '$app/env';
	import type { EventDetailedResponse } from '$lib/api/type';
	import {EventStatus} from "$lib/api/type";
	import {getContext} from "svelte";
	import {EditEvent, EventData} from "$lib/constant/context";
	import type {Writable} from "svelte/store";

	const eventData: EventDetailedResponse = getContext<EventDetailedResponse>(EventData);

	let editEvent = getContext<Writable<boolean>>(EditEvent)
	$editEvent = false


</script>

{#if eventData}
	{#if eventData.status !== EventStatus.Success}
		<div class="flex flex-row justify-center content-center mt-2">
			<div class="mx-2 flex flex-col items-center content-center">
				<h3 class="text-center">List Modified HPA</h3>
				<table class="table-fixed border-separate border border-slate-400 text-center">
					<thead>
					<tr>
						<th class="border border-slate-300 p-2"><h3>Name</h3></th>
						<th class="border border-slate-300 p-2"><h3>Namespace</h3></th>
						<th class="border border-slate-300 p-2"><h3>Min Replicas</h3></th>
						<th class="border border-slate-300 p-2"><h3>Max Replicas</h3></th>
					</tr>
					</thead>
					<tbody>
					{#each eventData.modified_hpa_configs as hpaConfig}
						<tr>
							<td class="border border-slate-300 p-2"><h4>{hpaConfig.name}</h4></td>
							<td class="border border-slate-300 p-2"><h4>{hpaConfig.namespace}</h4></td>
							<td class="border border-slate-300 p-2"><h4>{hpaConfig.min_replicas}</h4></td>
							<td class="border border-slate-300 p-2"><h4>{hpaConfig.max_replicas}</h4></td>
						</tr>
					{/each}
					</tbody>
				</table>
			</div>
			<div class="mx-2 items-center content-center flex flex-col">
				<h3 class="text-center">List Updated Node Pool</h3>
				<table class="table-fixed border-separate border border-slate-400 text-center">
					<thead>
					<tr>
						<th class="border border-slate-300 p-2"><h3>Name</h3></th>
						<th class="border border-slate-300 p-2"><h3>Max Node</h3></th>
					</tr>
					</thead>
					<tbody>
					{#if eventData.updated_node_pools.length === 0}
						<tr>
							<td class="border border-slate-300 p-2" colspan="2"><h4 class="text-center">Empty</h4></td>
						</tr>
					{:else }
						{#each eventData.updated_node_pools as updatedNodePool}
							<tr>
								<td class="border border-slate-300 p-2"><h4>{updatedNodePool.node_pool_name}</h4></td>
								<td class="border border-slate-300 p-2"><h4>{updatedNodePool.max_node}</h4></td>
							</tr>
						{/each}
					{/if}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
	{#if eventData.status === EventStatus.Success && browser}
		<div class="mt-2">
			<h3 class="font-bold">Monitoring</h3>
			<div class="flex mt-1 mb-3">
				<div class="flex-1 mb-3 mx-3 overflow-y-auto max-h-[75vh]">
					{#await import('$lib/components/chart/node-pool-chart.svelte')}
						<h4>Loading Component...</h4>
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
						<h4>Error Loading Component : {e}</h4>
					{/await}
				</div>
				<div class="flex-1 mb-3 mx-3 overflow-y-auto max-h-[75vh]">
					{#await import('$lib/components/chart/hpa-chart.svelte')}
						<h4>Loading Component...</h4>
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
						<h4>Error Loading Component : {e}</h4>
					{/await}
				</div>
			</div>
		</div>
	{/if}
{/if}
