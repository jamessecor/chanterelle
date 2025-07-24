/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_ADDRESS: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
