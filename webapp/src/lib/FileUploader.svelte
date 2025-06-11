<!-- FileUploader.svelte -->
<script>
    export let uploadUrl = "/api/upload";  // Default endpoint
    export let acceptedTypes = "*/*";       // Default accepted file types
    export let maxSizeMB = 10;              // Default max file size (MB)

    let selectedFile = null;
    let isUploading = false;
    let progress = 0;
    let uploadError = null;
    let uploadSuccess = false;

    // Handle file selection
    function handleFileChange(e) {
        const file = e.target.files[0];
        if (!file) return;

        // Validate file size
        if (file.size > maxSizeMB * 1024 * 1024) {
            uploadError = `File too large (max ${maxSizeMB}MB)`;
            selectedFile = null;
            return;
        }

        selectedFile = file;
        uploadError = null;
    }

    // Perform file upload
    async function uploadFile() {
        if (!selectedFile) return;

        isUploading = true;
        progress = 0;
        uploadError = null;
        uploadSuccess = false;

        const formData = new FormData();
        formData.append('file', selectedFile);

        try {
            const xhr = new XMLHttpRequest();
            
            // Progress tracking
            xhr.upload.addEventListener('progress', e => {
                if (e.lengthComputable) {
                    progress = Math.round((e.loaded / e.total) * 100);
                }
            });

            // Handle completion
            xhr.addEventListener('load', () => {
                if (xhr.status >= 200 && xhr.status < 300) {
                    uploadSuccess = true;
                } else {
                    uploadError = `Upload failed: ${xhr.statusText}`;
                }
                isUploading = false;
            });

            // Handle errors
            xhr.addEventListener('error', () => {
                uploadError = "Network error";
                isUploading = false;
            });

            xhr.open("POST", uploadUrl);
            xhr.send(formData);
        } catch (err) {
            uploadError = err.message;
            isUploading = false;
        }
    }

    // Reset component state
    function reset() {
        selectedFile = null;
        isUploading = false;
        progress = 0;
        uploadError = null;
        uploadSuccess = false;
    }
</script>

<div class="upload-container">
    <label class="file-input">
        <input 
            type="file" 
            accept={acceptedTypes}
            on:change={handleFileChange}
            disabled={isUploading}
        />
        Choose File
    </label>

    {#if selectedFile}
        <div class="file-info">
            <span>{selectedFile.name} ({Math.round(selectedFile.size / 1024)} KB)</span>
            {#if !isUploading}
                <button class="clear-btn" on:click={reset} title="Remove file">×</button>
            {/if}
        </div>
    {/if}

    <button 
        class="upload-btn" 
        on:click={uploadFile}
        disabled={!selectedFile || isUploading}
    >
        {isUploading ? 'Uploading...' : 'Upload'}
    </button>

    {#if isUploading}
        <div class="progress-bar">
            <div class="progress-fill" style={`width: ${progress}%`} />
        </div>
        <div class="progress-text">{progress}%</div>
    {/if}

    {#if uploadSuccess}
        <div class="success-message">✓ Upload successful!</div>
    {/if}

    {#if uploadError}
        <div class="error-message">{uploadError}</div>
    {/if}
</div>

<style>
    .upload-container {
        display: flex;
        flex-direction: column;
        gap: 12px;
        max-width: 500px;
        font-family: sans-serif;
    }

    .file-input {
        display: inline-block;
        padding: 8px 15px;
        background: #f0f0f0;
        border-radius: 4px;
        cursor: pointer;
        border: 1px solid #ccc;
    }

    .file-input:hover {
        background: #e0e0e0;
    }

    .file-input input {
        display: none;
    }

    .file-info {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 8px;
        background: #f8f8f8;
        border-radius: 4px;
        border-left: 3px solid #4285f4;
    }

    .clear-btn {
        background: transparent;
        border: none;
        font-size: 1.2em;
        cursor: pointer;
        color: #888;
    }

    .clear-btn:hover {
        color: #ff3d00;
    }

    .upload-btn {
        padding: 10px 20px;
        background: #4285f4;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: background 0.3s;
    }

    .upload-btn:hover:not(:disabled) {
        background: #3367d6;
    }

    .upload-btn:disabled {
        background: #999;
        cursor: not-allowed;
    }

    .progress-bar {
        height: 8px;
        background: #e0e0e0;
        border-radius: 4px;
        overflow: hidden;
    }

    .progress-fill {
        height: 100%;
        background: #4285f4;
        transition: width 0.3s;
    }

    .progress-text {
        text-align: center;
        font-size: 0.85em;
        color: #666;
    }

    .success-message {
        color: #2e7d32;
        padding: 8px;
        background: #edf7ed;
        border-radius: 4px;
    }

    .error-message {
        color: #d32f2f;
        padding: 8px;
        background: #fdecea;
        border-radius: 4px;
    }
</style>
