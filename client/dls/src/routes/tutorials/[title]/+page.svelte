<script context="module" lang="ts">
</script>

<script>
	import { currentTut, tutorialNav } from '@src/lib/stores/tutorials';
	import './app.css';
	import { makePath } from '@src/ts/helpers';

	$: next = tutorialNav.get()?.find((value) => value.link_title === $currentTut?.next);
	$: previous = tutorialNav.get()?.find((value) => value.link_title === $currentTut?.previous);
</script>

<!-- <svelte:body class="m-4" /> -->

<nav class="mx-4 my-8 flex justify-between text-white">
	<a
		class="bg-blue-500 p-4 text-xl text-white hover:underline"
		data-sveltekit-preload-code="off"
		data-sveltekit-preload-data="off"
		data-sveltekit-reload
		href={makePath('/')}>Home</a
	>
	<a
		class="bg-blue-500 p-4 text-xl text-white hover:underline"
		data-sveltekit-preload-code="off"
		data-sveltekit-preload-data="off"
		data-sveltekit-reload
		href={makePath('/profile')}>Start Creating</a
	>
	<div class="flex gap-4">
		<a
			class="bg-blue-500 p-4 text-xl text-white hover:underline"
			data-sveltekit-reload
			data-sveltekit-preload-code="off"
			data-sveltekit-preload-data="off"
			href={makePath('/signin')}>Sign in</a
		>
		<a
			class="bg-blue-500 p-4 text-xl text-white hover:underline"
			data-sveltekit-reload
			data-sveltekit-preload-code="off"
			data-sveltekit-preload-data="off"
			href={makePath('/signup')}>Sign up</a
		>
	</div>
</nav>

<div class="m-4">
	{@html $currentTut?.content}
</div>

<nav class="m-4 flex justify-between">
	<div>
		{#if previous != null}
			<a
				class="my-4 inline-block bg-blue-500 p-4 hover:underline"
				data-sveltekit-preload-data="off"
				href={makePath(`/tutorials/${previous.link_title}`)}>Previous Tutorial</a
			>
		{/if}
	</div>

	<div>
		{#if $currentTut?.link_title === 'index'}
			<a
				class="my-4 inline-block bg-blue-500 p-4 hover:underline"
				data-sveltekit-preload-data="off"
				href={makePath(`/tutorials/${next?.link_title}`)}>Next Tutorial</a
			>
		{:else}
			<a
				class="my-4 inline-block bg-blue-500 p-4 hover:underline"
				data-sveltekit-preload-data="off"
				href={makePath(`/quizes/${$currentTut?.id}`)}>Test You Knowledge</a
			>
		{/if}
	</div>
</nav>
