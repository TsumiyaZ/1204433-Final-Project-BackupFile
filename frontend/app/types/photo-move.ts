import type { dto } from "../../wailsjs/go/models";

export type MoveStatus = "pending" | "moved" | "skipped" | "failed";

export interface MoveProgress {
  current: number;
  total: number;
  percent: number;
  item: dto.MoveItemResult;
}

export interface BackupError {
  message: string;
}
