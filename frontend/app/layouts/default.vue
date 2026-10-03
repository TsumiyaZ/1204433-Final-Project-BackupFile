<template>
  <div data-theme="cream" class="min-h-dvh bg-base-200 text-base-content">
    <header
      class="sticky top-0 z-40 border-b border-base-300/80 bg-base-100/95 backdrop-blur"
    >
      <nav class="navbar mx-auto min-h-16 max-w-7xl px-4 sm:px-6">
        <div class="navbar-start">
          <NuxtLink
            to="/"
            class="flex min-h-12 items-center gap-3 rounded-xl px-2 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
            aria-label="ไปหน้าแรก"
          >
            <span
              class="grid size-10 place-items-center rounded-xl bg-primary text-primary-content"
              aria-hidden="true"
            >
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
                class="size-6"
              >
                <path d="M4 7.5A2.5 2.5 0 0 1 6.5 5h2l1.2 1.5h7.8A2.5 2.5 0 0 1 20 9v8.5a2.5 2.5 0 0 1-2.5 2.5h-11A2.5 2.5 0 0 1 4 17.5z" />
                <path d="m8 16 2.4-2.6a1 1 0 0 1 1.5 0l1.1 1.2 1.4-1.5a1 1 0 0 1 1.5 0L18 15.5" />
                <circle cx="15.5" cy="10" r="1.2" />
              </svg>
            </span>

            <span>
              <span class="block text-base font-bold leading-tight">Photo Backup</span>
              <span class="block text-xs text-base-content/55">จัดเก็บรูปอย่างเป็นระเบียบ</span>
            </span>
          </NuxtLink>
        </div>

        <div class="navbar-end gap-1 sm:gap-2">
          <NuxtLink
            v-for="item in navigation"
            :key="item.to"
            :to="item.to"
            class="btn min-h-11 border-0 px-3 sm:px-4"
            :class="isActive(item.to) ? 'btn-primary' : 'btn-ghost'"
          >
            {{ item.label }}
          </NuxtLink>
        </div>
      </nav>
    </header>

    <div>
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute();

const navigation = [
  { label: "หน้าแรก", to: "/" },
  { label: "ตั้งค่า", to: "/setup" },
  { label: "สแกนรูป", to: "/scan" },
];

function isActive(path: string) {
  return route.path === path;
}
</script>
