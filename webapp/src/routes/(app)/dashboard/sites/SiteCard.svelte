<script lang="ts">
	import { goto } from '$app/navigation';
	import { Alert, Button, Card, IconButton, Link, Modal, Toolbar } from '$lib/components';
	import { DotsVertical, Warning } from '$lib/icons';

	import FileUploader from './components/FileUploader.svelte';
	import UpdateSite from './UpdateSite.svelte';
	import { type Site, type SiteWarning } from '$lib/types';
	import { deleteSite } from '$lib/api-client';

	interface Props {
		site: Site;
		loading: boolean;
		errorMessage: string;
		onsiteupdate: () => Promise<void>;
	}

	let {
		site,
		loading = $bindable(false),
		errorMessage = $bindable(''),
		onsiteupdate = async () => {}
	}: Props = $props();

	let uploadModal = $state<Modal>();
	let deleteModal = $state<Modal>();
	let warningsModal = $state<Modal>();
	let fileUploader = $state<FileUploader>();
	let refreshAfterUploadModalClose = $state(false);
	let menuOpen = $state(false);
	let menuRef = $state<HTMLDivElement>();
	let openUpdateSite = $state<(() => void) | null>(null);

	const uploadTitleId = $derived(`upload-site-title-${site.id}`);
	const uploadDescId = $derived(`upload-site-description-${site.id}`);
	const deleteTitleId = $derived(`delete-site-title-${site.id}`);
	const deleteDescId = $derived(`delete-site-description-${site.id}`);
	const deleteButtonId = $derived(`delete-site-confirm-${site.id}`);
	const warningsTitleId = $derived(`site-warnings-title-${site.id}`);
	const warningsDescId = $derived(`site-warnings-description-${site.id}`);

	const resetState = () => {
		loading = false;
	};

	const warningCopy = (warning: SiteWarning) => {
		switch (warning.code) {
			case 'missing_home_page':
				return {
					title: 'Главная страница не найдена',
					message:
						'Загрузите файл index.html в корень сайта или ZIP-архива, чтобы главный адрес сайта открывался.'
				};
		}
	};

	const handleCardClick = () => {
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(`/dashboard/sites/${site.id}`);
	};

	const toggleMenu = (e: Event) => {
		e.stopPropagation();
		menuOpen = !menuOpen;
	};

	const closeMenu = () => {
		menuOpen = false;
	};

	const handleMenuAction = (e: Event, action: () => void) => {
		e.stopPropagation();
		closeMenu();
		action();
	};

	$effect(() => {
		if (!menuOpen || !menuRef) return;

		const handleClickOutside = (event: MouseEvent) => {
			if (menuRef && !menuRef.contains(event.target as Node)) {
				closeMenu();
			}
		};

		document.addEventListener('click', handleClickOutside, { capture: true });
		return () => document.removeEventListener('click', handleClickOutside, { capture: true });
	});
</script>

<div class="site-card-wrapper">
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<div class="card-click-target" onclick={handleCardClick}>
		<Card floating class="site-card">
			<div class="card-inner">
				<div class="site-info">
					<div class="site-name-row">
						<h3 class="site-name" title={site.name}>{site.name}</h3>
					</div>
					<span class="site-domain" title={site.domain}>{site.domain}</span>
				</div>

				<div class="toolbar-wrapper" bind:this={menuRef}>
					<Toolbar>
						{#if site.warnings?.length > 0}
						<IconButton
							label="Предупреждения сайта"
							title="Предупреждения сайта"
							tone="warning"
							grouped
							onclick={(e: Event) => {
								e.stopPropagation();
								warningsModal?.open();
							}}
						>
							<Warning size={25} />
						</IconButton>
					{/if}
					<IconButton label="Действия" title="Действия" onclick={toggleMenu} grouped>
						<DotsVertical size={25} />
					</IconButton>
					</Toolbar>

					{#if menuOpen}
						<div class="context-menu" role="menu">
							<a
								href="https://{site.domain}"
								target="_blank"
								rel="noopener noreferrer"
								class="menu-item"
								role="menuitem"
								onclick={(e) => e.stopPropagation()}
							>
								<span>Открыть сайт</span>
							</a>
							<button
								class="menu-item"
								type="button"
								role="menuitem"
								onclick={(e) => handleMenuAction(e, () => uploadModal?.open())}
							>
								<span>Загрузить файлы</span>
							</button>
							<button
								class="menu-item"
								type="button"
								role="menuitem"
								onclick={(e) => handleMenuAction(e, () => openUpdateSite?.())}
							>
								<span>Изменить название</span>
							</button>
							<button
								class="menu-item"
								type="button"
								role="menuitem"
								onclick={(e) => {
									e.stopPropagation();
									closeMenu();
									navigator.clipboard.writeText(`https://${site.domain}`);
								}}
							>
								<span>Копировать URL</span>
							</button>
							<div class="menu-divider" role="separator"></div>
							<button
								class="menu-item danger"
								type="button"
								role="menuitem"
								onclick={(e) => handleMenuAction(e, () => deleteModal?.open())}
							>
								<span>Удалить</span>
							</button>
						</div>
					{/if}
				</div>
			</div>
		</Card>
	</div>

	<UpdateSite
		siteName={site.name}
		siteId={site.id}
		onsuccess={async () => {
			resetState();
			await onsiteupdate();
		}}
		ontrigger={(open) => {
			openUpdateSite = open;
		}}
	/>

	<Modal
		bind:this={uploadModal}
		tone="neutral"
		title={`Загрузить файлы для "${site.name}"`}
		panelWidth="33.75rem"
		titleId={uploadTitleId}
		descriptionId={uploadDescId}
		initialFocusSelector=".drop-zone"
		closeCallback={() => {
			fileUploader?.resetUploadState();
			if (refreshAfterUploadModalClose) {
				refreshAfterUploadModalClose = false;
				void onsiteupdate();
			}
		}}
	>
		<div class="upload-modal">
			<Alert variant="info" title="Информация о загрузке сайтов">
				<p>При загрузке файлов все старые файлы будут удалены</p>
				<p id={uploadDescId}>
					Разрешённые типы файлов <Link href="/about" target="_blank" alwaysUnderline
						>на странице</Link
					>
				</p>
			</Alert>
			<FileUploader
				bind:this={fileUploader}
				onclose={() => uploadModal?.close()}
				onuploadcomplete={() => {
					refreshAfterUploadModalClose = true;
				}}
				multiple={true}
				uploadUrl="/api/uploadFile"
				siteId={site.id}
				additionalHeaders={{
					'H-Site-Id': site.id
				}}
			/>
		</div>
	</Modal>

	<Modal
		bind:this={deleteModal}
		tone="danger"
		title="Подтверждение удаления"
		titleId={deleteTitleId}
		descriptionId={deleteDescId}
		initialFocusSelector={'#' + deleteButtonId}
	>
		<div class="delete-modal">
			<p id={deleteDescId} class="delete-modal-text">
				Вы уверены, что хотите удалить сайт "{site.name}"?
			</p>
			<div class="delete-actions">
				<Button
					id={deleteButtonId}
					tone="danger"
					emphasis="solid"
					onclick={async () => {
						const { data, error } = await deleteSite(site.id);

						if (error) {
							errorMessage = error;
							deleteModal?.close();
							return;
						}

						if (data) {
							resetState();
							deleteModal?.close();
							await onsiteupdate();
						}
					}}>Удалить</Button
				>
				<Button tone="danger" onclick={() => deleteModal?.close()}>Отмена</Button>
			</div>
		</div>
	</Modal>

	<Modal
		bind:this={warningsModal}
		tone="neutral"
		title="Предупреждения сайта"
		panelWidth="32rem"
		titleId={warningsTitleId}
		descriptionId={warningsDescId}
	>
		<div class="warnings-modal">
			<p id={warningsDescId} class="warnings-modal-text">
				Проверьте эти предупреждения перед тем, как делиться сайтом.
			</p>

			{#each site.warnings ?? [] as warning (warning.code)}
				<Alert variant="warning" title={warningCopy(warning).title}>
					<p>{warningCopy(warning).message}</p>
				</Alert>
			{/each}
		</div>
	</Modal>
</div>

<style>
	.site-card-wrapper {
		position: relative;
	}

	.card-click-target {
		cursor: pointer;
	}

	:global(.site-card) {
		--card-surface: var(--color-background-object);
	}

	.card-inner {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-md);
		padding: var(--space-lg);
		min-width: 0;
	}

	.site-info {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		min-width: 0;
		flex: 1;
	}

	.site-name-row {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		min-width: 0;
	}

	.site-name {
		font-size: 1.05rem;
		font-weight: 700;
		line-height: 1.3;
		color: var(--color-text-regular);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		min-width: 0;
	}

	.site-domain {
		font-size: 0.82rem;
		font-weight: 500;
		color: var(--color-text-secondary);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		min-width: 0;
	}

	.toolbar-wrapper {
		position: relative;
		flex-shrink: 0;
	}

	.context-menu {
		position: absolute;
		top: calc(100% + 0.25rem);
		right: 0;
		z-index: 10;
		min-width: 12rem;
		background: var(--color-background-object);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		box-shadow: var(--shadow-lg);
		padding: var(--space-xs);
		display: flex;
		flex-direction: column;
		gap: 0.1rem;
		animation: menu-in 160ms var(--ease-out-quart) both;
	}

	@keyframes menu-in {
		0% {
			opacity: 0;
			transform: translateY(-0.25rem) scale(0.98);
		}
		100% {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	.menu-item {
		all: unset;
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		padding: 0.55rem 0.7rem;
		border-radius: var(--radius-sm);
		font-size: 0.88rem;
		font-weight: 500;
		color: var(--color-text-regular);
		cursor: pointer;
		white-space: nowrap;
		transition:
			background-color 120ms ease-out,
			color 120ms ease-out;
	}

	.menu-item:hover {
		background: var(--color-neutral-surface);
	}

	.menu-item.danger {
		color: var(--color-danger-text);
	}

	.menu-item.danger:hover {
		background: var(--color-danger-surface);
	}

	.menu-divider {
		height: 1px;
		background: var(--color-border);
		margin: 0.15rem 0.4rem;
	}

	.upload-modal {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		min-width: 0;
		width: 100%;
	}

	.delete-modal {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}

	.warnings-modal {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
	}

	.warnings-modal-text {
		color: var(--color-text-secondary);
	}

	.delete-modal-text {
		color: var(--color-danger-text);
	}

	.delete-actions {
		display: flex;
		gap: var(--space-sm);
	}
</style>
