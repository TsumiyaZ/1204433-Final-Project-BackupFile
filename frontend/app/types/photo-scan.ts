export type ScanMode = "source" | "destination";

export interface ScanOption {
  mode: ScanMode;
  label: string;
  description: string;
}
