<script lang="ts">
	export let msg: string | undefined;
</script>

{#if msg != null}
	<div class="flex items-center gap-2 rounded-md bg-white px-4 py-1 text-red-600">
		<span class="material-symbols-outlined">info</span>
		{msg}
	</div>
{/if}
