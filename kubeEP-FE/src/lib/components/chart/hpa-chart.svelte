<script>
    import FusionCharts from 'fusioncharts';
    import Timeseries from 'fusioncharts/fusioncharts.timeseries';
    import SvelteFusioncharts, {fcRoot} from "svelte-fusioncharts";
    import {GetEventHPAStatistics} from "$lib/api/statistics.ts";
    import {onMount} from "svelte";
    import moment from "moment";

    fcRoot(FusionCharts, Timeseries);

    export let id;
    export let name;
    export let namespace;

    let data = [];
    let error = null;
    let loaded = false;

    onMount(async () => {
        try {
            data = await GetEventHPAStatistics(id)
            data = [
                ...data.map(o => [moment(o.created_at).format("YYYY-MM-DD hh:mm:ss A"),"Replicas", o.replicas]),
                ...data.map(o => [moment(o.created_at).format("YYYY-MM-DD hh:mm:ss A"),"Ready Replicas", o.ready_replicas]),
                ...data.map(o => [moment(o.created_at).format("YYYY-MM-DD hh:mm:ss A"),"Available Replicas", o.available_replicas]),
                ...data.map(o => [moment(o.created_at).format("YYYY-MM-DD hh:mm:ss A"),"Unavailable Replicas", o.unavailable_replicas])
            ]
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
            "name": "Type",
            "type": "string"
        },
        {
            "name": "Pods",
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
                    text: `HPA ${name} - Namespace ${namespace} `
                },
                series: "Type",
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
    <SvelteFusioncharts {...getChartConfig(data, schema)} />
{/if}

{#if loaded && error}
    <p>{error}</p>
{/if}
