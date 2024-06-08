<template>
  <Toolbar>
    <template #start>
      <div class="flex gap-1 max-h-3rem align-content-center">
        <router-link to="/">
          <Button>
            <i class="fa fa-home" />
          </Button>
        </router-link>
        <router-link to="/users">
          <Button>
            <i class="fa fa-user" />
          </Button>
        </router-link>
          <Button>
            <i class="fa fa-sitemap" />
          </Button>
        </div>
      </template>

    <template #end>
      <IconField iconPosition="left">
        <InputIcon>
          <i class="fa fa-search" />
        </InputIcon>
        <InputText placeholder="Search" />
      </IconField>
    </template>

  </Toolbar>
</template>
<script setup lang="ts">

</script>
