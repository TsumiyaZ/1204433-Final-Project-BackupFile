<template>
  <section
    class="rounded-[1.5rem] border border-base-300 bg-base-100 p-5"
    aria-live="polite"
  >
    <template v-if="result">
      <h2 class="text-lg font-bold">ย้ายรูปเสร็จแล้ว</h2>

      <dl class="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-4">
        <div
          v-for="stat in summary"
          :key="stat.label"
          class="rounded-xl bg-base-200 p-3"
        >
          <dt class="text-sm text-base-content/70">{{ stat.label }}</dt>
          <dd class="text-2xl font-bold" :class="stat.className">
            {{ stat.value }}
          </dd>
        </div>
      </dl>

      <p class="mt-3 text-sm text-base-content/70">
        ใช้เวลา {{ (result.durationMs / 1000).toFixed(1) }} วินาที
      </p>
    </template>

    <template v-else>
      <div class="flex items-center justify-between gap-4">
        <h2 class="text-lg font-bold">กำลังย้ายรูป</h2>
        <span class="text-sm font-semibold text-primary">
          {{ progress?.current ?? 0 }} / {{ progress?.total ?? total }}
        </span>
      </div>

      <progress
        class="progress progress-primary mt-4 h-3 w-full"
        :value="progress?.percent ?? 0"
        max="100"
        :aria-label="`ความคืบหน้า ${progress?.percent ?? 0}%`"
      />

      <p
        v-if="progress"
        class="mt-3 truncate text-sm text-base-content/70"
        :title="progress.item.filename"
      >
        ล่าสุด: {{ progress.item.filename }}
      </p>
    </template>
  </section>
</template>

<script setup lang="ts">
import { computed } from "vue";

import type { dto } from "../../../wailsjs/go/models";
import type { MoveProgress } from "../../types/photo-move";

const props = defineProps<{
  total: number;
  progress: MoveProgress | null;
  result: dto.MovePhotosResult | null;
}>();

const summary = computed(() => {
  if (!props.result) return [];

  return [
    { label: "ทั้งหมด", value: props.result.requested, className: "" },
    { label: "ย้ายสำเร็จ", value: props.result.moved, className: "text-success" },
    { label: "ข้าม", value: props.result.skipped, className: "text-warning" },
    { label: "ล้มเหลว", value: props.result.failed, className: "text-error" },
  ];
});
</script>
