<template>
  <div class="app-shell min-h-screen bg-dark-bg text-white" dir="rtl">
    <HeaderNav />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import HeaderNav from './HeaderNav.vue'
</script>

<style scoped>
.app-shell {
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
}
</style>
