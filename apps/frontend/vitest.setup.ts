import '@testing-library/jest-dom/vitest';
import { afterEach, vi } from 'vitest';
import { cleanup } from '@testing-library/react';

// Unmount React trees and clear mocks between tests for isolation.
afterEach(() => {
  cleanup();
  vi.clearAllMocks();
});

// jsdom does not implement matchMedia, which next-themes and some UI components
// query. Provide a no-op implementation.
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }),
});

// jsdom lacks ResizeObserver, used by some Radix primitives.
global.ResizeObserver = class {
  observe() {}
  unobserve() {}
  disconnect() {}
};
