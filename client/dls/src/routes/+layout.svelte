<script>
	import { navigating } from '$app/stores';
	import LoadingSpinner from '@src/lib/components/LoadingSpinner.svelte';
	import { loading } from '@src/lib/stores/loading';

	$: console.log('Loading changed to: ', $loading);
	$: loading.set(!!$navigating);
</script>

{#if $loading}
	<div
		class="fixed bottom-0 right-0 z-50 grid h-screen w-screen items-center justify-center bg-black/50"
	>
		<LoadingSpinner />
	</div>
{/if}
<slot />
