<script>
	import FusionCharts from 'fusioncharts';
	import Timeseries from 'fusioncharts/fusioncharts.timeseries';
	import SvelteFusioncharts, { fcRoot } from 'svelte-fusioncharts';
	import { GetEventHPAStatistics } from '$lib/api/statistics.ts';
	import { onMount } from 'svelte';
	import dayjs from 'dayjs';

	fcRoot(FusionCharts, Timeseries);

	export let id;
	export let name;
	export let namespace;
	export let minPods = 0;
	export let maxPods = 0;

	let data = [];
	let error = null;
	let loaded = false;
	let maximumPodsStatistic = {};

	onMount(async () => {
		try {
			const response = await GetEventHPAStatistics(id);
			maximumPodsStatistic = {
				replicas: Math.max(...response.map(o => o.replicas)),
				readyReplicas: Math.max(...response.map(o => o.ready_replicas)),
				availableReplicas: Math.max(...response.map(o => o.available_replicas)),
				unavailableReplicas: Math.max(...response.map(o => o.unavailable_replicas))
			};
			data = [
				...response.map((o) => [
					dayjs(o.created_at).format('YYYY-MM-DD hh:mm:ss A'),
					'Replicas',
					o.replicas
				]),
				...response.map((o) => [
					dayjs(o.created_at).format('YYYY-MM-DD hh:mm:ss A'),
					'Ready Replicas',
					o.ready_replicas
				]),
				...response.map((o) => [
					dayjs(o.created_at).format('YYYY-MM-DD hh:mm:ss A'),
					'Available Replicas',
					o.available_replicas
				]),
				...response.map((o) => [
					dayjs(o.created_at).format('YYYY-MM-DD hh:mm:ss A'),
					'Unavailable Replicas',
					o.unavailable_replicas
				])
			];
			loaded = true;
		} catch (e) {
			error = e;
		}
	});

	const schema = [
		{
			name: 'Time',
			type: 'date',
			format: '%Y-%m-%d %I:%M:%S %p'
		},
		{
			name: 'Type',
			type: 'string'
		},
		{
			name: 'Pods',
			type: 'number'
		}
	];

	const getChartConfig = (data, schema) => {
		const fusionDataStore = new FusionCharts.DataStore(),
			fusionTable = fusionDataStore.createDataTable(data, schema);

		return {
			type: 'timeseries',
			width: '100%',
			height: 450,
			renderAt: 'chart-container',
			dataSource: {
				data: fusionTable,
				caption: {
					text: `HPA ${name} - Namespace ${namespace} `
				},
				subcaption: {
					text: `Updated Min Pods ${minPods} - Updated Max Pods ${maxPods}`
				},
				series: 'Type',
				yAxis: [
					{
						plot: {
							value: `HPA ${name} - Namespace ${namespace} Pod Count`,
							type: 'line'
						},
						title: 'Pod Count'
					}
				]
			}
		};
	};
</script>

{#if !loaded}
	<h4>Loading...</h4>
{/if}

{#if data.length > 0 && !error && loaded}
	<div class="mb-2">
		<SvelteFusioncharts {...getChartConfig(data, schema)} />
		<div class="text-left">
			<h4>Maximum Spawned Replicas : {maximumPodsStatistic.replicas}</h4>
			<h4>Maximum Spawned Ready Replicas : {maximumPodsStatistic.readyReplicas}</h4>
			<h4>Maximum Spawned Available Replicas : {maximumPodsStatistic.availableReplicas}</h4>
			<h4>Maximum Spawned Unavailable Replicas : {maximumPodsStatistic.unavailableReplicas}</h4>
		</div>
	</div>
{/if}

{#if loaded && error}
	<p>{error}</p>
{/if}
