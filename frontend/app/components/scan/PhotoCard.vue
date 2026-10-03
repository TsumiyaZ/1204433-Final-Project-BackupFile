<template>
  <article
    class="card overflow-hidden border bg-base-100 transition duration-200"
    :class="
      selectable && selected
        ? 'border-primary ring-2 ring-primary/20'
        : 'border-base-300 hover:border-primary/40'
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

      <span
        v-else-if="!previewFinished"
        class="loading loading-spinner loading-md text-primary"
      />

      <span v-else class="text-lg font-bold uppercase text-base-content/60">
        {{ photo.extension.replace(".", "") }}
      </span>
    </div>

    <div class="card-body gap-3 border-t border-base-300/70 p-4">
      <div class="flex items-start gap-3">
        <input
          v-if="selectable"
          type="checkbox"
          class="checkbox checkbox-primary mt-1"
          :checked="selected"
          :aria-label="`เลือก ${photo.filename}`"
          @change="handleSelection"
        />

        <span
          v-else
          class="badge badge-outline mt-0.5 border-primary/30 text-primary"
        >
          ปลายทาง
        </span>

        <div class="min-w-0">
          <h3 class="truncate font-semibold" :title="photo.filename">
            {{ photo.filename }}
          </h3>

          <p
            class="mt-1 break-all text-sm text-base-content/60"
            :title="photo.relativePath"
          >
            {{ photo.relativePath }}
          </p>

          <p class="mt-2 text-sm text-base-content/70">
            {{ formatFileSize(photo.size) }}
          </p>
        </div>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { dto } from "../../../wailsjs/go/models";

defineProps<{
  photo: dto.ScannedPhoto;
  selectable: boolean;
  selected: boolean;
  previewUrl: string;
  previewFinished: boolean;
}>();

const emit = defineEmits<{
  "update:selected": [selected: boolean];
}>();

function handleSelection(event: Event) {
  emit("update:selected", (event.target as HTMLInputElement).checked);
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
</script>
