<template>
  <div class="rounded-[1.5rem] border border-base-300 bg-base-100 p-5 sm:p-6">
    <div class="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
      <div class="min-w-0">
        <div class="flex flex-wrap items-center gap-3">
          <span
            class="badge min-h-7 border-primary/30 bg-primary/8 font-semibold text-primary"
          >
            {{ scanMode === "source" ? "ต้นทาง" : "ปลายทาง" }}
          </span>

          <h2 class="text-xl font-bold">
            {{ scanMode === "source" ? "Scan Source" : "Scan Destination" }}
          </h2>
        </div>

        <p
          class="mt-3 break-all text-sm leading-6 text-base-content/60"
          :title="currentPath"
        >
          {{ currentPath || `ยังไม่ได้กำหนด ${folderLabel}` }}
        </p>
      </div>

      <div class="flex shrink-0 flex-col gap-3 sm:flex-row">
        <NuxtLink
          to="/setup"
          class="btn btn-outline min-h-12 border-primary/35 text-primary"
        >
          ตั้งค่าโฟลเดอร์
        </NuxtLink>

        <button
          type="button"
          class="btn btn-primary min-h-12 min-w-44"
          :disabled="!canScan"
          @click="emit('scan')"
        >
          <span v-if="isScanning" class="loading loading-spinner loading-sm" />
          {{ isScanning ? "กำลังสแกน..." : `สแกน ${folderLabel}` }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

import type { ScanMode } from "../../types/photo-scan";

const props = defineProps<{
  scanMode: ScanMode;
  currentPath: string;
  canScan: boolean;
  isScanning: boolean;
}>();

const emit = defineEmits<{
  scan: [];
}>();

const folderLabel = computed(() =>
  props.scanMode === "source" ? "Source" : "Destination",
);
</script>
