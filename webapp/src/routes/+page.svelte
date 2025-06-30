<script>
	import FileUploader from "$lib/FileUploader.svelte";

</script>


<FileUploader 
  uploadUrl="/api/uploadFile"
  acceptedTypes=".html,.zip"
  maxFileSize={40 * 1024 * 1024}
  multiple={true}
  onfilesselected={(e) => console.log('Files selected:', e.files)}
  onuploadcomplete={(e) => console.log('Upload complete!')}
  onuploaderror={(e) => console.log('Error:', e.error)}
/>
