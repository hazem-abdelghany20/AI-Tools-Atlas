<template>
  <div class="tool-profile-skeleton space-y-6">
    <!-- Hero Skeleton -->
    <section class="bg-dark-surface rounded-lg p-6 md:p-8">
      <div class="flex flex-col md:flex-row-reverse gap-6 items-start">
        <!-- Logo Skeleton -->
        <div class="w-24 h-24 md:w-32 md:h-32 rounded-xl skeleton"></div>

        <!-- Content Skeleton -->
        <div class="flex-1 w-full">
          <!-- Category & Badges -->
          <div class="flex gap-2 mb-3">
            <div class="w-20 h-6 rounded-full skeleton"></div>
            <div class="w-24 h-6 rounded-full skeleton"></div>
          </div>

          <!-- Name -->
          <div class="w-3/4 h-10 rounded skeleton mb-3"></div>

          <!-- Tagline -->
          <div class="w-1/2 h-6 rounded skeleton mb-4"></div>

          <!-- Rating -->
          <div class="flex items-center gap-4 mb-6">
            <div class="w-32 h-6 rounded skeleton"></div>
            <div class="w-20 h-6 rounded-full skeleton"></div>
          </div>

          <!-- Buttons -->
          <div class="flex gap-3">
            <div class="w-32 h-12 rounded-lg skeleton"></div>
            <div class="w-24 h-12 rounded-lg skeleton"></div>
          </div>
        </div>
      </div>
    </section>

    <!-- Overview Skeleton -->
    <section class="bg-dark-surface rounded-lg p-6 md:p-8 space-y-6">
      <!-- Best For -->
      <div>
        <div class="w-24 h-6 rounded skeleton mb-2"></div>
        <div class="w-full h-16 rounded skeleton"></div>
      </div>

      <!-- Description -->
      <div>
        <div class="w-28 h-6 rounded skeleton mb-3"></div>
        <div class="space-y-2">
          <div class="w-full h-4 rounded skeleton"></div>
          <div class="w-full h-4 rounded skeleton"></div>
          <div class="w-3/4 h-4 rounded skeleton"></div>
        </div>
      </div>

      <!-- Use Cases -->
      <div>
        <div class="w-40 h-6 rounded skeleton mb-3"></div>
        <div class="flex gap-2">
          <div class="w-24 h-10 rounded-lg skeleton"></div>
          <div class="w-28 h-10 rounded-lg skeleton"></div>
          <div class="w-20 h-10 rounded-lg skeleton"></div>
        </div>
      </div>

      <!-- Platforms -->
      <div>
        <div class="w-32 h-6 rounded skeleton mb-3"></div>
        <div class="flex gap-3">
          <div class="w-24 h-10 rounded-lg skeleton"></div>
          <div class="w-24 h-10 rounded-lg skeleton"></div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.skeleton {
  background: linear-gradient(90deg, #0A0B10 25%, #1F2933 50%, #0A0B10 75%);
  background-size: 200% 100%;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
