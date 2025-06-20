<script>
	import FileUploader from "$lib/FileUploader.svelte";

</script>

<FileUploader 
  uploadUrl="/files/upload"
  acceptedTypes=".html"
  maxSizeMB={5}
/>