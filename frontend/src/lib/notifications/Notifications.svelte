<script>
	// @ts-nocheck

	import { toaster } from './toaster.js';

	export let themes = {
		danger: '#bb2124',
		success: '#22bb33',
		warning: '#f0ad4e',
		info: '#5bc0de',
		default: '#aaaaaa'
	};

	export let timeout = 3000;
	export let sessionKey = 'byk-toasts';

	/**
	 * @type {any[]}
	 */
	let toasts = [];

	/**
	 * @param {HTMLLIElement} _node
	 */
	function animateOut(_node, { delay = 0, duration = 1000 }) {
		return {
			delay,
			duration,
			css: (/** @type {number} */ t) => `opacity: ${(t - 0.7) * 1}; transform-origin: top right;`
		};
	}
	/**
	 * @param {CustomEvent} event
	 */
	function createToast({ detail }) {
		const { message, type, options = {} } = detail;
		const background = themes[type] || themes.default;
		const persist = options.persist;
		const computedTimeout = options.persist ? 0 : options.timeout || timeout;
		const id = Math.random()
			.toString(36)
			.replace(/[^a-z]+/g, '');

		try {
			sessionStorage.setItem(
				sessionKey,
				JSON.stringify([
					...JSON.parse(sessionStorage.getItem(sessionKey) || '[]'),
					{ ...detail, id }
				])
			);
		} catch (e) {
			console.log('Cannot use session storage ' + e.message);
		}

		toasts = [
			{
				id,
				message,
				background,
				persist,
				timeout: computedTimeout,
				width: '100%'
			},
			...toasts
		];
	}

	/**
	 * @param {{ persist: any; id: any; }} toast
	 */
	function maybePurge(toast) {
		if (!toast.persist) purge(toast.id);
	}

	/**
	 * @param {any} id
	 */
	function purge(id) {
		const filter = (/** @type {{ id: any; }} */ t) => t.id !== id;
		toasts = toasts.filter(filter);
		try {
			sessionStorage.setItem(
				sessionKey,
				JSON.stringify(JSON.parse(sessionStorage.getItem(sessionKey) || '[]').filter(filter))
			);
		} catch (e) {
			console.log('Cannot use session storage ' + e.message);
		}
	}
</script>

<ul class="toasts" use:toaster={sessionKey} on:notify={createToast}>
	{#each toasts as toast (toast.id)}
		<li class="toast" style="background: {toast.background};" out:animateOut|global>
			{#if toast.persist}
				<button class="close" on:click={() => purge(toast.id)}> ✕ </button>
			{/if}
			<div class="content">
				{toast.message}
			</div>
			<div
				class="progress"
				style="animation-duration: {toast.timeout}ms;"
				on:animationend={() => maybePurge(toast)}
			></div>
		</li>
	{/each}
</ul>

<style>
	.toasts {
		list-style: none;
		position: fixed;
		top: 0;
		right: 0;
		padding: 0;
		margin: 0;
		z-index: 9999;
	}

	.toasts > .toast {
		display: flex;
		align-items: center;
		position: relative;
		margin: 1vh 1vw;
		min-width: 98vw;
		position: relative;
		animation: animate-in 600ms forwards;
		color: #fff;
		min-height: 8vh;
	}

	.toasts > .toast > button {
		position: absolute;
		font-size: 18px;
		right: 0;
		margin: 6px;
		color: #fff;
		outline: none;
		border: 0;
		background-color: transparent;
	}

	.toasts > .toast > .content {
		padding: 1vw;
		display: flex;
		font-weight: 500;
		margin-right: 20px;
	}

	.toasts > .toast > .progress {
		position: absolute;
		bottom: 0;
		background-color: rgb(0, 0, 0, 0.3);
		height: 6px;
		width: 100%;
		animation-name: shrink;
		animation-timing-function: linear;
		animation-fill-mode: forwards;
	}

	.toasts > .toast:before,
	.toasts > .toast:after {
		content: '';
		position: absolute;
		z-index: -1;
		top: 50%;
		bottom: 0;
		left: 1vw;
		right: 1vw;
		border-radius: 100px / 10px;
	}

	.toasts > .toast:after {
		right: 1vw;
		left: auto;
		transform: skew(8deg) rotate(3deg);
	}

	@keyframes animate-in {
		0%,
		60%,
		75%,
		90%,
		to {
			-webkit-animation-timing-function: cubic-bezier(0.215, 0.61, 0.355, 1);
			animation-timing-function: cubic-bezier(0.215, 0.61, 0.355, 1);
		}

		0% {
			opacity: 0;
			transform: translate3d(3000px, 0, 0);
		}

		60% {
			opacity: 1;
			transform: translate3d(-25px, 0, 0);
		}

		75% {
			transform: translate3d(10px, 0, 0);
		}

		90% {
			transform: translate3d(-5px, 0, 0);
		}

		to {
			transform: none;
		}
	}

	@keyframes shrink {
		0% {
			width: 98vw;
		}
		100% {
			width: 0;
		}
	}

	@media (min-width: 480px) {
		@keyframes animate-in {
			0%,
			60%,
			75%,
			90%,
			to {
				-webkit-animation-timing-function: cubic-bezier(0.215, 0.61, 0.355, 1);
				animation-timing-function: cubic-bezier(0.215, 0.61, 0.355, 1);
			}

			0% {
				opacity: 0;
				transform: translate3d(3000px, 0, 0);
			}

			60% {
				opacity: 1;
				transform: translate3d(-25px, 0, 0);
			}

			75% {
				transform: translate3d(10px, 0, 0);
			}

			90% {
				transform: translate3d(-5px, 0, 0);
			}

			to {
				transform: none;
			}
		}

		@keyframes shrink {
			0% {
				width: 40vw;
			}
			100% {
				width: 0;
			}
		}
	}

	@media screen and (min-width: 600px) {
		.toasts > .toast {
			min-width: 40vw;
			min-height: auto;
		}

		.toasts > .toast > .content {
			justify-content: flex-start;
		}
	}
</style>
