import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getAvatarUrl(path?: string): string | undefined {
  if (!path) return undefined
  if (path.startsWith('http')) return path

  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'
  try {
    const url = new URL(baseUrl)
    // remove /api suffix if present for image root
    // Assuming upload root is at server root, not under /api
    return `${url.origin}${path}`
  } catch (e) {
    return `http://localhost:8080${path}`
  }
}