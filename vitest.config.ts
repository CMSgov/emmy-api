import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    include: ['test/**/*.test.ts'],
    exclude: ['**/node_modules/**', '**/dist/**'],

    // Match Jest's `clearMocks: true` + slightly stricter isolation
    clearMocks: true,
    mockReset: true,
    restoreMocks: true,

    environment: 'node',
  },
});
