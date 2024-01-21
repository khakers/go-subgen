<script lang="ts">
  import { base } from "$app/paths";
  import {
    SkipToContent,
    HeaderNav,
    HeaderNavItem,
    HeaderNavMenu,
    Header,
    Content,
  } from "carbon-components-svelte";
  import "carbon-components-svelte/css/g90.css";
</script>

<Header
  platformName="Go-Subgen"
  href="https://github.com/khakers/Go-Subgen/tree/v1-dev"
>
  <svelte:fragment slot="skip-to-content">
    <SkipToContent />
  </svelte:fragment>
  <HeaderNav>
    <HeaderNavItem href="{base}/" text="Home" />
    <HeaderNavItem href="{base}/about" text="About" />
    <!-- <HeaderNavItem href="/" text="Link 3" />
    <HeaderNavMenu text="Menu">
      <HeaderNavItem href="/" text="Link 1" />
      <HeaderNavItem href="/" text="Link 2" />
      <HeaderNavItem href="/" text="Link 3" />
    </HeaderNavMenu> -->
  </HeaderNav>
</Header>
<Content>
  <slot />
</Content>
