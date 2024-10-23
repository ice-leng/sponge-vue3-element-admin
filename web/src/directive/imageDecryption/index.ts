import { Directive, DirectiveBinding } from "vue";

export const srcDesc: Directive = {
  mounted(el: HTMLImageElement, binding: DirectiveBinding) {
    const fetchOptions: RequestInit = {};

    async function fetchAndSetImage(el: HTMLImageElement, url: string) {
      try {
        const res = await fetch(url, fetchOptions);
        if (!res.ok) {
          throw new Error(`Failed to fetch image: ${res.statusText}`);
        }
        const blob = await res.blob();
        const file = blob.slice(3, blob.size); // Slice the blob if necessary
        el.src = URL.createObjectURL(file);
      } catch (error) {
        console.error("Error fetching the image:", error);
      }
    }

    fetchAndSetImage(el, binding.value);
  },
  updated(el, binding) {
    const fetchOptions: RequestInit = {};
    async function fetchAndSetImage(el: HTMLImageElement, url: string) {
      try {
        const res = await fetch(url, fetchOptions);
        if (!res.ok) {
          throw new Error(`Failed to fetch image: ${res.statusText}`);
        }
        const blob = await res.blob();
        const file = blob.slice(3, blob.size); // Slice the blob if necessary
        el.src = URL.createObjectURL(file);
      } catch (error) {
        console.error("Error fetching the image:", error);
      }
    }

    fetchAndSetImage(el, binding.value);
  },
};
