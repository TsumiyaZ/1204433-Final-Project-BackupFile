<template>
  <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
    <MovePhotoCard
      v-for="photo in photos"
      :key="photo.path"
      :photo="photo"
      :preview-url="previewUrls[photo.path] ?? ''"
      :checked="checkedPaths.includes(photo.path)"
      :disabled="disabled"
      :status="statusOf(photo.path)"
      :result="results[photo.path]"
      @update:checked="updateChecked(photo.path, $event)"
    />
  </div>
</template>

<script setup lang="ts">
import type { dto } from "../../../wailsjs/go/models";
import type { MoveStatus } from "../../types/photo-move";

import MovePhotoCard from "./MovePhotoCard.vue";

const props = defineProps<{
  photos: dto.ScannedPhoto[];
  previewUrls: Record<string, string>;
  checkedPaths: string[];
  disabled: boolean;
  results: Record<string, dto.MoveItemResult>;
  statusOf: (path: string) => MoveStatus;
}>();

const emit = defineEmits<{
  "update:checkedPaths": [paths: string[]];
}>();

function updateChecked(path: string, checked: boolean) {
  const nextPaths = new Set(props.checkedPaths);

  if (checked) {
    nextPaths.add(path);
  } else {
    nextPaths.delete(path);
  }

  emit("update:checkedPaths", Array.from(nextPaths));
}
</script>
