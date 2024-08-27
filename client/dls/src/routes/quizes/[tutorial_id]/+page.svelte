<script context="module" lang="ts">
	export enum AnswerStatus {
		Correct,
		Incorrect,
		UnAttempted
	}
</script>

<script lang="ts">
	import { onMount, setContext } from 'svelte';
	import './app.css';

	import Mcq from '@src/lib/components/Mcq/Mcq.svelte';
	import { makePath } from '@src/ts/helpers';

	export let data;
	let answerStatuses = new Array<AnswerStatus>(data.mcqs.length).fill(AnswerStatus.UnAttempted);
	setContext('quiz-answers', answerStatuses);
	let nCorrect = 0;
	let nIncorrect = 0;
	let nTotal = data.mcqs.length;

	let submitted = false;

	onMount(() => {
		const radioInputs = document.querySelectorAll<HTMLInputElement>('input[type=radio]');
		for (const input of radioInputs) {
			input.checked = false;
			console.log(input);
		}
		window.onbeforeunload = function () {
			window.scrollTo(0, 0);
		};
	});
</script>

<div class="m-4">
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
	<h1>{data.tutorial.display_title} - Quiz</h1>
	<div>
		{#each data.mcqs as mcq, idx}
			<Mcq {mcq} {idx} />
		{/each}
	</div>
	<div>
		<button
			class="my-4 inline-block bg-blue-500 p-4 hover:underline"
			disabled={submitted}
			on:click={submitted
				? () => {}
				: () => {
						if (answerStatuses.indexOf(AnswerStatus.UnAttempted) != -1) {
							alert('Please attempt all questions.');
							return;
						}
						nCorrect = answerStatuses.reduce((acc, cur) => {
							if (cur === AnswerStatus.Correct) {
								return acc + 1;
							}
							return acc;
						}, 0);
						nIncorrect = answerStatuses.reduce((acc, cur) => {
							if (cur === AnswerStatus.Incorrect) {
								return acc + 1;
							}
							return acc;
						}, 0);
						submitted = true;
					}}>Submit</button
		>
	</div>

	{#if submitted}
		<div>
			<h2>Your Results</h2>
			<div>
				<div>
					<div>Total</div>
					<div>{nTotal}</div>
				</div>
				<hr />
				<div style="color: green;">
					<div>Correct</div>
					<div>{nCorrect}</div>
				</div>
				<hr />
				<div style="color: red;">
					<div>Incorrect</div>
					<div>{nIncorrect}</div>
				</div>
			</div>
			<div
				style="display:flex; justify-content: space-between; margin-inline: 1rem; font-size: 1rem;"
			>
				<button
					class="my-4 inline-block bg-blue-500 p-4 hover:underline"
					on:click={() => {
						window.location.reload();
					}}>Retake Quiz</button
				>
				{#if data.tutorial.next != null}
					<a
						class="my-4 inline-block bg-blue-500 p-4 hover:underline"
						href={makePath(`/tutorials/${data.tutorial.next}`)}>Next Tutorial</a
					>
				{:else}
					<a
						data-sveltekit-reload
						class="my-4 inline-block bg-blue-500 p-4 hover:underline"
						href={makePath(`/tutorials/index`)}>Goto Home</a
					>
				{/if}
			</div>
		</div>
	{/if}
</div>
