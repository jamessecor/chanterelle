/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_ADDRESS: string;
  // add more env variables here as needed
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
