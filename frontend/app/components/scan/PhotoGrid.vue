<template>
  <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
    <PhotoCard
      v-for="photo in photos"
      :key="photo.path"
      :photo="photo"
      :selectable="scanMode === 'source'"
      :selected="selectedPaths.includes(photo.path)"
      :preview-url="previewUrls[photo.path] ?? ''"
      :preview-finished="previewFinished[photo.path] ?? false"
      @update:selected="updatePhotoSelection(photo.path, $event)"
    />
  </div>
</template>

<script setup lang="ts">
import type { dto } from "../../../wailsjs/go/models";
import type { ScanMode } from "../../types/photo-scan";

import PhotoCard from "./PhotoCard.vue";

const props = defineProps<{
  photos: dto.ScannedPhoto[];
  scanMode: ScanMode;
  selectedPaths: string[];
  previewUrls: Record<string, string>;
  previewFinished: Record<string, boolean>;
}>();

const emit = defineEmits<{
  "update:selectedPaths": [paths: string[]];
}>();

function updatePhotoSelection(path: string, selected: boolean) {
  const nextPaths = new Set(props.selectedPaths);

  if (selected) {
    nextPaths.add(path);
  } else {
    nextPaths.delete(path);
  }

  emit("update:selectedPaths", Array.from(nextPaths));
}
</script>
