<script lang="ts">
	import '@src/app.css';
	import { _Problem } from './+page';
	import { afterUpdate, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { apiURL } from '@src/ts/api/helpers';
	import { loading } from '@src/lib/stores/loading';
	import ErrorMsg from '@comps/ErrorMsg.svelte';
	import { makePath, getAuthHeader, setProjectID } from '@src/ts/helpers';

	loading.set(true);
	export let data;

	let projects = data.projects;
	let newProjectName = '';
	let error: string | undefined = undefined;
	onMount(async () => {
		if (data.problem === _Problem.NotLoggedIn) {
			await goto(makePath('/signin'));
			return;
		}
		loading.set(false);
	});
</script>

<div class="mx-[10rem] flex flex-col">
	<nav class=" my-8 flex justify-between text-white">
		<a
			class="bg-blue-500 p-4 text-xl text-white hover:underline"
			data-sveltekit-preload-code="off"
			data-sveltekit-preload-data="off"
			data-sveltekit-reload
			href={makePath('/')}>Home</a
		>
		<a
			class="bg-blue-500 p-4 text-xl text-white hover:underline"
			data-sveltekit-reload
			data-sveltekit-preload-code="off"
			data-sveltekit-preload-data="off"
			href={makePath('/tutorials/index')}>Tutorials</a
		>
	</nav>
	<div class="my-20 flex justify-between">
		<h1 class="text-4xl">Select Project</h1>
		<form
			class="flex flex-col gap-5"
			on:submit={async () => {
				loading.set(true);
				const res = await fetch(apiURL('dbproject-create/create'), {
					body: JSON.stringify({ name: newProjectName }),
					method: 'POST',
					credentials: 'include',
					headers: getAuthHeader()
				});
				if (res.status === 409) {
					// conflict
					error = 'A project with this name already exists';
				} else if (res.status === 401) {
					// unauthorized
					await goto(makePath('/signin'));
					return;
				} else if (res.ok) {
					const json = await res.json();
					const id = json['id'];
					const name = newProjectName;
					projects.unshift({ id, name });
					projects = projects;
				}
				loading.set(false);
			}}
		>
			<div>
				<input class="rounded-md p-4 text-neutral-950" type="text" bind:value={newProjectName} />
				<input class="p-4 hover:bg-neutral-600" type="submit" value="Create Project" />
			</div>
			{#if error != null}
				<ErrorMsg msg={error} />
			{/if}
		</form>
	</div>
	<div class="flex flex-wrap gap-10">
		{#if projects.length > 0}
			{#each projects as project (project.id)}
				<button
					class="rounded-md border border-neutral-600 p-4 hover:bg-neutral-600"
					on:click={async () => {
						loading.set(true);
						// await fetch(apiURL('dbproject/select'), {
						// 	method: 'POST',
						// 	body: JSON.stringify({
						// 		projectID: project.id
						// 	}),
						// 	credentials: 'include',
						// 	headers: getAuthHeader()
						// });
						// document.cookie = `projectID=${project.id.toString()};path=/;SameSite=None`;
						// console.log('Set projectID cookie');
						// console.log('Cookies: ', document.cookie);
						setProjectID(project.id);
						await goto(makePath('/dls'));
					}}>{project.name}</button
				>
			{/each}
		{:else}
			<span class="text-neutral-500">No project created</span>
		{/if}
	</div>
</div>
