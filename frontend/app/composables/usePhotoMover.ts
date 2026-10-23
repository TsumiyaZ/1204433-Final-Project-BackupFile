import { computed, onMounted, onUnmounted, ref } from "vue";

import { MovePhotos } from "../../wailsjs/go/controller/BackupController";
import { dto } from "../../wailsjs/go/models";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import type { BackupError, MoveProgress, MoveStatus } from "../types/photo-move";
import { useSelectedPhotos } from "./useSelectedPhotos";

export function usePhotoMover() {
  const { selectedPhotos, selectedPreviews, clearSelection } =
    useSelectedPhotos();

  const checkedPaths = ref<string[]>(
    selectedPhotos.value.map((photo) => photo.path),
  );
  const statusByPath = ref<Record<string, dto.MoveItemResult>>({});
  const progress = ref<MoveProgress | null>(null);
  const result = ref<dto.MovePhotosResult | null>(null);
  const isMoving = ref(false);
  const errorMessage = ref("");

  const canMove = computed(
    () => !isMoving.value && !result.value && checkedPaths.value.length > 0,
  );
  const allChecked = computed(
    () =>
      selectedPhotos.value.length > 0 &&
      checkedPaths.value.length === selectedPhotos.value.length,
  );

  function statusOf(path: string): MoveStatus {
    const status = statusByPath.value[path]?.status;
    return status === "moved" || status === "skipped" || status === "failed"
      ? status
      : "pending";
  }

  function toggleAll() {
    checkedPaths.value = allChecked.value
      ? []
      : selectedPhotos.value.map((photo) => photo.path);
  }

  function getErrorMessage(error: unknown): string {
    return error instanceof Error ? error.message : String(error);
  }

  async function move() {
    if (!canMove.value) return;

    isMoving.value = true;
    errorMessage.value = "";
    progress.value = null;
    statusByPath.value = {};

    try {
      result.value = await MovePhotos(
        new dto.MovePhotosRequest({ paths: [...checkedPaths.value] }),
      );

      for (const item of result.value.items) {
        statusByPath.value[item.sourcePath] = item;
      }

      // Photos are gone from Source now, so the Scan selection is stale.
      clearSelection();
    } catch (error) {
      errorMessage.value = getErrorMessage(error);
    } finally {
      isMoving.value = false;
    }
  }

  const stopListeners: Array<() => void> = [];

  onMounted(() => {
    stopListeners.push(
      EventsOn("backup:progress", (event: MoveProgress) => {
        progress.value = event;
        statusByPath.value[event.item.sourcePath] = event.item;
      }),
      EventsOn("backup:error", (event: BackupError) => {
        errorMessage.value = event.message;
      }),
    );
  });

  onUnmounted(() => {
    for (const stop of stopListeners) stop();
  });

  return {
    selectedPhotos,
    selectedPreviews,
    checkedPaths,
    statusByPath,
    progress,
    result,
    isMoving,
    errorMessage,
    canMove,
    allChecked,
    statusOf,
    toggleAll,
    move,
  };
}
