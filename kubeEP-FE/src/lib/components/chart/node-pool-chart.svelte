<script>
    import FusionCharts from 'fusioncharts';
    import Timeseries from 'fusioncharts/fusioncharts.timeseries';
    import SvelteFusioncharts, {fcRoot} from "svelte-fusioncharts";
    import {GetEventNodePoolStatistics} from "$lib/api/statistics.ts";
    import {onMount} from "svelte";
    import moment from "moment";

    fcRoot(FusionCharts, Timeseries);

    export let id;
    export let name;

    let data = [];
    let error = null;
    let loaded = false;

    onMount(async () => {
        try {
            data = await GetEventNodePoolStatistics(id)
            data = data.map(o => [moment(o.created_at).format("YYYY-MM-DD hh:mm:ss A"), o.count])
            loaded = true
        } catch (e) {
            error = e;
        }
    })

    const schema = [
        {
            "name": "Time",
            "type": "date",
            "format": "%Y-%m-%d %I:%M:%S %p"
        },
        {
            "name": "Node Count",
            "type": "number"
        }
    ]

    const getChartConfig = (data,schema) => {
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
    <SvelteFusioncharts {...getChartConfig(data, schema)} />
{/if}

{#if loaded && error}
    <p>{error}</p>
{/if}
