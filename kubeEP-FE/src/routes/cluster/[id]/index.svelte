<script context="module" lang="ts">
    export async function load({ params }) {
        const clusterID = params.id
        return {
            redirect: `/cluster/${clusterID}/hpa`,
            status: 302
        }
    }
</script>
