<script lang="ts">
  import { Menubar } from "bits-ui";
</script>

<Menubar.Root>
  <Menubar.Menu>
    <Menubar.Trigger>API</Menubar.Trigger>
    <Menubar.Portal>
      <Menubar.Content>
        <Menubar.Item>API Documentation</Menubar.Item>
        <Menubar.Item>API Status</Menubar.Item>
      </Menubar.Content>
    </Menubar.Portal>
  </Menubar.Menu>

  <Menubar.Menu>
    <Menubar.Trigger>Help</Menubar.Trigger>
    <Menubar.Portal>
      <Menubar.Content>
        <Menubar.Item>Documentation</Menubar.Item>
        <Menubar.Item>Support</Menubar.Item>
        <Menubar.Item>Feedback</Menubar.Item>
      </Menubar.Content>
    </Menubar.Portal>
  </Menubar.Menu>

  <Menubar.Menu>
    <Menubar.Trigger>Pricing</Menubar.Trigger>
    <Menubar.Portal>
      <Menubar.Content>
        <Menubar.Item>Plans</Menubar.Item>
        <Menubar.Item>Features</Menubar.Item>
      </Menubar.Content>
    </Menubar.Portal>
  </Menubar.Menu>

  <Menubar.Menu>
    <Menubar.Trigger>Profile</Menubar.Trigger>
    <Menubar.Portal>
      <Menubar.Content>
        <Menubar.Item>Account Settings</Menubar.Item>
        <Menubar.Item>Logout</Menubar.Item>
      </Menubar.Content>
    </Menubar.Portal>
  </Menubar.Menu>
</Menubar.Root>
