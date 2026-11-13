# Source Structure

- `api/`: HTTP client, request/response types, and feature API modules.
- `assets/`: static images, fonts, and icons bundled by Vite.
- `components/ui/`: reusable primitive UI components.
- `components/common/`: shared application components.
- `components/layout/`: navigation and layout building blocks.
- `composables/`: reusable Composition API state and behavior.
- `config/`: application and environment configuration.
- `constants/`: shared constant values and enumerations.
- `data/`: domain types and category helpers; no seeded application data.
- `directives/`: reusable Vue directives.
- `layouts/`: page-level layout compositions.
- `pages/`: route views.
- `plugins/`: Vue app plugins and registration.
- `services/`: domain workflows that coordinate APIs and stores.
- `stores/`: Pinia state stores.
- `types/`: shared TypeScript domain and app types.
- `utils/`: framework-independent helpers.

Keep feature-specific code close to its domain. Promote it to a shared directory when more than one feature needs it.
