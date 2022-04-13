<script>
	import FusionCharts from 'fusioncharts';
	import Timeseries from 'fusioncharts/fusioncharts.timeseries';
	import SvelteFusioncharts, { fcRoot } from 'svelte-fusioncharts';
	import { GetEventHPAStatistics } from '$lib/api/statistics.ts';
	import { onMount } from 'svelte';
	import moment from 'moment';

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
				replicas: response.reduce(
					(prev, current) => (prev < current.replicas ? current.replicas : prev),
					0
				),
				readyReplicas: response.reduce(
					(prev, current) => (prev < current.replicas ? current.ready_replicas : prev),
					0
				),
				availableReplicas: response.reduce(
					(prev, current) => (prev < current.replicas ? current.available_replicas : prev),
					0
				),
				unavailableReplicas: response.reduce(
					(prev, current) => (prev < current.replicas ? current.unavailable_replicas : prev),
					0
				)
			};
			data = [
				...response.map((o) => [
					moment(o.created_at).format('YYYY-MM-DD hh:mm:ss A'),
					'Replicas',
					o.replicas
				]),
				...response.map((o) => [
					moment(o.created_at).format('YYYY-MM-DD hh:mm:ss A'),
					'Ready Replicas',
					o.ready_replicas
				]),
				...response.map((o) => [
					moment(o.created_at).format('YYYY-MM-DD hh:mm:ss A'),
					'Available Replicas',
					o.available_replicas
				]),
				...response.map((o) => [
					moment(o.created_at).format('YYYY-MM-DD hh:mm:ss A'),
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
					text: `Min Pods ${minPods} - Max Pods ${maxPods}`
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
	<h3>Loading...</h3>
{/if}

{#if data.length > 0 && !error && loaded}
	<div class="mb-2">
		<SvelteFusioncharts {...getChartConfig(data, schema)} />
		<div class="text-left">
			<h3>Maximum Replicas : {maximumPodsStatistic.replicas}</h3>
			<h3>Maximum Ready Replicas : {maximumPodsStatistic.readyReplicas}</h3>
			<h3>Maximum Available Replicas : {maximumPodsStatistic.availableReplicas}</h3>
			<h3>Maximum Unavailable Replicas : {maximumPodsStatistic.unavailableReplicas}</h3>
		</div>
	</div>
{/if}

{#if loaded && error}
	<p>{error}</p>
{/if}
