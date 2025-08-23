<script lang="ts">
  interface Props {
    uploadUrl?: string;
    acceptedTypes?: string;
    maxFileSize?: number;
    multiple?: boolean;
    disabled?: boolean;
    additionalHeaders?: Record<string, string>;
    onfilesselected?: (event: { files: File[] }) => void;
    onfileuploaded?: (event: { file: File; result: any }) => void;
    onuploadcomplete?: (event: { fileCount: number }) => void;
    onuploaderror?: (event: { error: string }) => void;
  }
  
  let {
    uploadUrl = '/api/upload',
    acceptedTypes = '',
    maxFileSize = 20 * 1024 * 1024, // 20MB default
    multiple = false,
    disabled = false,
    additionalHeaders = {},
    onfilesselected,
    onfileuploaded,
    onuploadcomplete,
    onuploaderror
  }: Props = $props();
  
  let fileInput: HTMLInputElement;
  let selectedFiles: File[] = $state([]);
  let uploading: boolean = $state(false);
  let uploadProgress: number = $state(0);
  let error: string = $state('');
  let link: string = $state('');
  let uploadedFiles: string[] = $state([]);
  let skippedFiles: string[] = $state([]);
  
  function handleFileSelect(event: Event): void {
    const target = event.target as HTMLInputElement;
    if (!target.files) return;
    
    const files = Array.from(target.files);
    selectedFiles = files;
    error = '';
    link = '';
    
    // Validate file sizes
    const oversizedFiles = files.filter(file => file.size > maxFileSize);
    if (oversizedFiles.length > 0) {
      error = `File(s) too large. Maximum size: ${formatFileSize(maxFileSize)}`;
      return;
    }
    
    onfilesselected?.({ files });
  }
  
  function handleDragOver(event: DragEvent): void {
    event.preventDefault();
    if (event.dataTransfer) {
      event.dataTransfer.dropEffect = 'copy';
    }
  }
  
  function handleDrop(event: DragEvent): void {
    event.preventDefault();
    if (!event.dataTransfer) return;
    
    const files = Array.from(event.dataTransfer.files);
    
    // Update file input
    const dt = new DataTransfer();
    files.forEach(file => dt.items.add(file));
    fileInput.files = dt.files;
    
    selectedFiles = files;
    error = '';
    link = '';
    
    onfilesselected?.({ files });
  }
  
  async function uploadFiles(): Promise<void> {
    if (selectedFiles.length === 0) {
      error = 'Please select files to upload';
      return;
    }
    
    uploading = true;
    uploadProgress = 0;
    error = '';
    link = '';
    
    try {
      for (let i = 0; i < selectedFiles.length; i++) {
        const file = selectedFiles[i];
        const formData = new FormData();
        formData.append('file', file);
        
        const headers = new Headers();
        Object.entries(additionalHeaders).forEach(([key, value]) => {
          headers.set(key, value);
        });

        const response = await fetch(uploadUrl, {
          method: 'POST',
          headers,
          body: formData
        });
        
        if (!response.ok) {
          throw new Error(`Upload failed: ${response.statusText}`);
        }
        
        const result = await response.json();
        uploadProgress = ((i + 1) / selectedFiles.length) * 100;
        link = result.link
        
        // Store uploaded and skipped files
        uploadedFiles = result.uploadedFiles || [];
        skippedFiles = result.skippedFiles || [];

        onfileuploaded?.({ file, result });
      }
      
      const uploadedCount = selectedFiles.length;
      selectedFiles = [];
      fileInput.value = '';
      
      onuploadcomplete?.({ fileCount: uploadedCount });
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
      error = errorMessage;
      onuploaderror?.({ error: errorMessage });
    } finally {
      uploading = false;
      uploadProgress = 0;
    }
  }
  
  function removeFile(index: number): void {
    selectedFiles = selectedFiles.filter((_, i) => i !== index);
    
    // Update file input
    const dt = new DataTransfer();
    selectedFiles.forEach(file => dt.items.add(file));
    fileInput.files = dt.files;
    
    if (selectedFiles.length === 0) {
      link = '';
      error = '';
    }
  }
  
  function formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }
  
  function clearAll(): void {
    selectedFiles = [];
    fileInput.value = '';
    error = '';
    link = '';
    uploadProgress = 0;
  }
  
  function handleKeyDown(event: KeyboardEvent): void {
    if (event.key === 'Enter') {
      fileInput.click();
    }
  }
</script>

<div class="file-upload-container">
  <div 
    class="drop-zone" 
    class:drag-over={false}
    ondragover={handleDragOver}
    ondrop={handleDrop}
    role="button"
    tabindex="0"
    onclick={() => fileInput.click()}
    onkeydown={handleKeyDown}
  >
    <div class="drop-zone-content">
      <svg class="upload-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
        <polyline points="17,8 12,3 7,8"/>
        <line x1="12" y1="3" x2="12" y2="15"/>
      </svg>
      <p class="drop-text">
        <strong>Click to select files</strong> or drag and drop
      </p>
      <p class="drop-subtext">
        {#if acceptedTypes}
          Accepted: {acceptedTypes}
        {:else}
          All file types accepted
        {/if}
        • Max size: {formatFileSize(maxFileSize)}
      </p>
    </div>
  </div>
  
  <input
    bind:this={fileInput}
    type="file"
    {multiple}
    accept={acceptedTypes}
    onchange={handleFileSelect}
    style="display: none;"
    {disabled}
  />
  
  {#if selectedFiles.length > 0}
    <div class="selected-files">
      <div class="files-header">
        <h3>Selected Files ({selectedFiles.length})</h3>
        <button type="button" class="clear-btn" onclick={clearAll}>
          Clear All
        </button>
      </div>
      
      <div class="files-list">
        {#each selectedFiles as file, index}
          <div class="file-item">
            <div class="file-info">
              <span class="file-name">{file.name}</span>
              <span class="file-size">{formatFileSize(file.size)}</span>
            </div>
            <button 
              type="button" 
              class="remove-btn"
              onclick={() => removeFile(index)}
              aria-label="Remove file"
            >
              ×
            </button>
          </div>
        {/each}
      </div>
    </div>
  {/if}
  
  {#if uploading}
    <div class="progress-container">
      <div class="progress-bar">
        <div class="progress-fill" style="width: {uploadProgress}%"></div>
      </div>
      <span class="progress-text">Uploading... {Math.round(uploadProgress)}%</span>
    </div>
  {/if}
  
  {#if error}
    <div class="message error">
      {error}
    </div>
  {/if}
  
  {#if link}
    <div class="message success">
      <a href="https://{link}" target="_blank">{link}</a>
      {#if uploadedFiles.length > 0}
        <div class="uploaded-files">
          <h4>Uploaded Files:</h4>
          <ul>
            {#each uploadedFiles as file}
              <li>{file}</li>
            {/each}
          </ul>
        </div>
      {/if}
      {#if skippedFiles.length > 0}
        <div class="skipped-files">
          <h4>Skipped Files:</h4>
          <ul>
            {#each skippedFiles as file}
              <li>{file}</li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
  {/if}
  
  <div class="actions">
    <button 
      type="button" 
      class="upload-btn"
      onclick={uploadFiles}
      disabled={disabled || uploading || selectedFiles.length === 0}
    >
      {#if uploading}
        Uploading...
      {:else}
        Upload {selectedFiles.length > 0 ? `(${selectedFiles.length})` : ''}
      {/if}
    </button>
  </div>
</div>

<style>
  .file-upload-container {
    max-width: 500px;
    margin: 0 auto;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
  
  .drop-zone {
    border: 2px dashed #cbd5e1;
    border-radius: 8px;
    padding: 2rem;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s ease;
    background-color: #f8fafc;
  }
  
  .drop-zone:hover {
    border-color: #3b82f6;
    background-color: #eff6ff;
  }
  
  .drop-zone:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 2px;
  }
  
  .drop-zone-content {
    pointer-events: none;
  }
  
  .upload-icon {
    width: 48px;
    height: 48px;
    margin: 0 auto 1rem;
    color: #64748b;
  }
  
  .drop-text {
    margin: 0 0 0.5rem;
    color: #1e293b;
    font-size: 1rem;
  }
  
  .drop-subtext {
    margin: 0;
    color: #64748b;
    font-size: 0.875rem;
  }
  
  .selected-files {
    margin-top: 1.5rem;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    overflow: hidden;
  }
  
  .files-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background-color: #f1f5f9;
    border-bottom: 1px solid #e2e8f0;
  }
  
  .files-header h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
    color: #1e293b;
  }
  
  .clear-btn {
    padding: 0.25rem 0.75rem;
    font-size: 0.875rem;
    color: #64748b;
    background: none;
    border: 1px solid #cbd5e1;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .clear-btn:hover {
    color: #ef4444;
    border-color: #ef4444;
  }
  
  .files-list {
    max-height: 200px;
    overflow-y: auto;
  }
  
  .file-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    border-bottom: 1px solid #f1f5f9;
  }
  
  .file-item:last-child {
    border-bottom: none;
  }
  
  .file-info {
    flex: 1;
    min-width: 0;
  }
  
  .file-name {
    display: block;
    font-weight: 500;
    color: #1e293b;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .file-size {
    font-size: 0.875rem;
    color: #64748b;
  }
  
  .remove-btn {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: none;
    border: none;
    border-radius: 50%;
    cursor: pointer;
    color: #64748b;
    font-size: 18px;
    line-height: 1;
    transition: all 0.2s ease;
  }
  
  .remove-btn:hover {
    background-color: #fee2e2;
    color: #ef4444;
  }
  
  .progress-container {
    margin-top: 1rem;
  }
  
  .progress-bar {
    width: 100%;
    height: 8px;
    background-color: #e2e8f0;
    border-radius: 4px;
    overflow: hidden;
  }
  
  .progress-fill {
    height: 100%;
    background-color: #3b82f6;
    transition: width 0.3s ease;
  }
  
  .progress-text {
    display: block;
    margin-top: 0.5rem;
    font-size: 0.875rem;
    color: #64748b;
    text-align: center;
  }
  
  .message {
    margin-top: 1rem;
    padding: 0.75rem 1rem;
    border-radius: 6px;
    font-size: 0.875rem;
  }
  
  .message.error {
    background-color: #fef2f2;
    color: #dc2626;
    border: 1px solid #fecaca;
  }
  
  .message.success {
    background-color: #f0fdf4;
    color: #16a34a;
    border: 1px solid #bbf7d0;
  }
  
  .uploaded-files,
  .skipped-files {
    margin-top: 1rem;
  }
  
  .uploaded-files h4,
  .skipped-files h4 {
    margin: 0.5rem 0;
    font-size: 0.875rem;
    color: #1e293b;
  }
  
  .uploaded-files ul,
  .skipped-files ul {
    margin: 0;
    padding-left: 1.25rem;
    font-size: 0.875rem;
    color: #64748b;
  }
  
  .skipped-files {
    color: #d97706;
  }
  
  .actions {
    margin-top: 1.5rem;
    text-align: center;
  }
  
  .upload-btn {
    padding: 0.75rem 2rem;
    background-color: #3b82f6;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .upload-btn:hover:not(:disabled) {
    background-color: #2563eb;
  }
  
  .upload-btn:disabled {
    background-color: #cbd5e1;
    cursor: not-allowed;
  }
</style>
