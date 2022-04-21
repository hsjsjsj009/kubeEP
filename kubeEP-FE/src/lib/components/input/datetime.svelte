<script>
    import dayjs from 'dayjs'

    export let format = 'YYYY-MM-DDTHH:mm'
    export let date = dayjs()
    export let minDate = null;
    export let maxDate = null;

    let internal

    const input = (x) => (internal = dayjs(x).format(format))
    const output = (x) => (date = dayjs(x, format))

    $: input(date)
    $: output(internal)
</script>

<input type="datetime-local" max={maxDate ? maxDate.format(format) : null} min={minDate ? minDate.format(format): null} bind:value={internal} {...$$props}>
