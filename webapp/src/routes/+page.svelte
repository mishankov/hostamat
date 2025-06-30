<script>
	import FileUploader from "$lib/FileUploader.svelte";

</script>


<FileUploader 
  uploadUrl="/api/uploadFile"
  acceptedTypes=".html,.zip"
  multiple={true}
  onfilesselected={(e) => console.log('Files selected:', e.files)}
  onuploadcomplete={(e) => console.log('Upload complete!')}
  onuploaderror={(e) => console.log('Error:', e.error)}
/>
