<script lang="ts">
  import FileUploader from "./FileUploader.svelte";
  
  export interface Site {
    id: string;
    subdomain: string;
    domain: string;
    createdAt: string;
  }

  export let site: Site;
</script>

<div class="site-uploader">
  <h3>Upload files for {site.domain}</h3>
  <FileUploader 
    uploadUrl="/api/uploadFile"
    additionalHeaders={{ 'H-Site-Id': site.id }}
    acceptedTypes=".html,.zip"
    maxFileSize={40 * 1024 * 1024}
    multiple={true}
    onfilesselected={(e) => console.log('Files selected:', e.files)}
    onuploadcomplete={(e) => console.log('Upload complete!')}
    onuploaderror={(e) => console.log('Error:', e.error)}
  />
</div>

<style>
  .site-uploader {
    border: 1px solid #eee;
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
  }
</style>
