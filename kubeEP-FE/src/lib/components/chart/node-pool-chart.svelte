<script>
	import FusionCharts from 'fusioncharts';
	import Timeseries from 'fusioncharts/fusioncharts.timeseries';
	import SvelteFusioncharts, { fcRoot } from 'svelte-fusioncharts';
	import { GetEventNodePoolStatistics } from '$lib/api/statistics.ts';
	import { onMount } from 'svelte';
	import moment from 'moment';

	export let id;
	export let name;
	export let maxNode = 0;

	let data = [];
	let error = null;
	let loaded = false;
	let maximumNodes = 0;

	onMount(async () => {
		fcRoot(FusionCharts, Timeseries);
		try {
			const response = await GetEventNodePoolStatistics(id);
			data = response.map((o) => [moment(o.created_at).format('YYYY-MM-DD hh:mm:ss A'), o.count]);
			maximumNodes = data.reduce(
				(prev, current) => (prev < current[1] ? current[1] : prev),
				maximumNodes
			);
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
			name: 'Node Count',
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
					text: `Node Pool ${name}`
				},
				subcaption: {
					text: `Updated Max Node ${maxNode}`
				},
				yAxis: [
					{
						plot: {
							value: `${name} Node Count`,
							type: 'line'
						},
						title: 'Node Count'
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
			<h3>Maximum Spawned Nodes : {maximumNodes}</h3>
		</div>
	</div>
{/if}

{#if loaded && error}
	<p>{error}</p>
{/if}
