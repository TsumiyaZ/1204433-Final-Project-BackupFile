<template>
  <article
    class="card overflow-hidden border bg-base-100 transition duration-200"
    :class="
      checked && status === 'pending'
        ? 'border-primary ring-2 ring-primary/20'
        : 'border-base-300'
    "
  >
    <div
      class="flex aspect-video items-center justify-center overflow-hidden bg-base-200"
    >
      <img
        v-if="previewUrl"
        :src="previewUrl"
        :alt="photo.filename"
        class="h-full w-full object-cover"
        loading="lazy"
        decoding="async"
      />

      <span v-else class="text-lg font-bold uppercase text-base-content/60">
        {{ photo.extension.replace(".", "") }}
      </span>
    </div>

    <div class="card-body gap-3 border-t border-base-300/70 p-4">
      <div class="flex items-start gap-3">
        <input
          type="checkbox"
          class="checkbox checkbox-primary mt-1"
          :checked="checked"
          :disabled="disabled"
          :aria-label="`เลือก ${photo.filename}`"
          @change="handleChange"
        />

        <div class="min-w-0 flex-1">
          <h3 class="truncate font-semibold" :title="photo.filename">
            {{ photo.filename }}
          </h3>

          <p
            class="mt-1 break-all text-sm text-base-content/60"
            :title="photo.relativePath"
          >
            {{ photo.relativePath }}
          </p>

          <div class="mt-3 flex flex-wrap items-center gap-2">
            <span class="badge" :class="statusBadge.className">
              {{ statusBadge.label }}
            </span>

            <span
              v-if="result?.message && status !== 'moved'"
              class="text-sm text-base-content/70"
            >
              {{ result.message }}
            </span>
          </div>

          <p
            v-if="status === 'moved' && result?.filename !== photo.filename"
            class="mt-2 text-sm text-base-content/70"
          >
            บันทึกเป็น {{ result?.filename }}
          </p>
        </div>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from "vue";

import type { dto } from "../../../wailsjs/go/models";
import type { MoveStatus } from "../../types/photo-move";

const props = defineProps<{
  photo: dto.ScannedPhoto;
  previewUrl: string;
  checked: boolean;
  disabled: boolean;
  status: MoveStatus;
  result?: dto.MoveItemResult;
}>();

const emit = defineEmits<{
  "update:checked": [checked: boolean];
}>();

const statusBadge = computed(() => {
  switch (props.status) {
    case "moved":
      return { label: "ย้ายแล้ว", className: "badge-success" };
    case "skipped":
      return { label: "ข้าม", className: "badge-warning" };
    case "failed":
      return { label: "ล้มเหลว", className: "badge-error" };
    default:
      return {
        label: props.checked ? "รอย้าย" : "ไม่ย้าย",
        className: "badge-ghost",
      };
  }
});

function handleChange(event: Event) {
  emit("update:checked", (event.target as HTMLInputElement).checked);
}
</script>
