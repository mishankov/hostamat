<script>
	import FileUploader from "$lib/FileUploader.svelte";

</script>


<FileUploader 
  uploadUrl="/files/upload"
  acceptedTypes=".html"
  multiple={true}
  onfilesselected={(e) => console.log('Files selected:', e.files)}
  onuploadcomplete={(e) => console.log('Upload complete!')}
  onuploaderror={(e) => console.log('Error:', e.error)}
/>
