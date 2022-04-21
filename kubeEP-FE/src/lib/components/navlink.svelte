<script lang="ts">
    import {page} from '$app/stores'

    export let href: string;
</script>

<a {href} class="rounded-3xl border-2 p-2 mx-2" class:border-black={$page.url.pathname === href} class:text-gray-500={$page.url.pathname !== href}>
    <slot/>
</a>
