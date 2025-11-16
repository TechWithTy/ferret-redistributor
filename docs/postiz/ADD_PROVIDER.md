# Adding a New Social Provider

This guide walks through every touch point required to integrate a new social
provider end-to-end. The examples use **X** for names, but you can substitute
any provider identifier.

---

## Backend Steps

### 1. Settings DTO
1. Create the provider’s settings DTO at  
   `libraries/nestjs-libraries/src/dtos/posts/providers-settings/x-provider-settings.dto.ts`.  
   (Skip this if your provider has no additional settings.)
2. Register it in  
   `libraries/nestjs-libraries/src/dtos/posts/providers-settings/all.providers.settings.ts`.

### 2. Add to CreatePost discriminator
Open `libraries/nestjs-libraries/src/dtos/posts/create.post.dto.ts`, locate the
`@Discriminator` config, and append your mapping:

```ts
{ value: XProviderSettingsDto, name: 'x' },
```

### 3. Implement the provider class
1. Create `libraries/nestjs-libraries/src/integrations/social/x.provider.ts`.
2. Implement the required methods:
   ```ts
   import {
     AuthTokenDetails,
     PostDetails,
     PostResponse,
     SocialProvider,
   } from '@gitroom/nestjs-libraries/integrations/social/social.integrations.interface';

   export class XProvider implements SocialProvider {
     identifier = 'x';
     name = 'X';

     async refreshToken(refreshToken: string): Promise<AuthTokenDetails> { /* ... */ }
     async generateAuthUrl() { /* ... */ }
     async authenticate(params: { code: string; codeVerifier: string }) { /* ... */ }
     async post(
       id: string,
       accessToken: string,
       postDetails: PostDetails<XProviderSettingsDto>[],
     ): Promise<PostResponse[]> { /* ... */ }
   }
   ```
3. Add any custom helper methods (e.g., `organizations()`) if the frontend needs
   extra data such as available pages or orgs.

### 4. Register in the integration manager
Update `libraries/nestjs-libraries/src/integrations/integration.manager.ts` and
push the new class into either:
- `socialIntegrationList` (OAuth2 providers), or
- `articleIntegrationList` (token-based providers).

---

## Frontend Steps

### 1. Provider component
1. Create `apps/frontend/src/components/launches/providers/x/index.tsx` (or a
   similarly named folder).
2. Implement the preview/settings components and wrap them:
   ```tsx
   import { FC } from 'react';
   import { withProvider } from '@gitroom/frontend/components/launches/providers/high.order.provider';
   import { useSettings } from '@gitroom/frontend/components/launches/helpers/use.values';
   import { useIntegration } from '@gitroom/frontend/components/launches/helpers/use.integration';
   import { useCustomProviderFunction } from '@gitroom/frontend/components/launches/helpers/use.custom.provider.function';

   const XPreview: FC = () => {
     const { value } = useIntegration();
     const settings = useSettings();
     return (/* … */);
   };

   const XSettings: FC = () => {
     const form = useSettings();
     const customFunc = useCustomProviderFunction();
     // customFunc.get('organizations', { ... })
     return (/* … */);
   };

   export default withProvider(XSettings, XPreview, XProviderSettingsDto);
   ```

### 2. Register the provider
Update `apps/frontend/src/components/launches/providers/show.all.providers.tsx`
and append the new entry:

```tsx
{ identifier: 'x', component: XProviderComponent }
```

### 3. Assets
Add the provider icon/logo to the shared assets directory used by the provider
grid.

---

## Checklist

| Area | Action |
| --- | --- |
| DTOs | Provider settings DTO created and registered. |
| CreatePost | Discriminator entry added. |
| Provider class | `generateAuthUrl`, `authenticate`, `refreshToken`, and `post` implemented. |
| Integration manager | Provider added to the appropriate list. |
| Frontend | Settings + preview component created and registered. |
| Assets | Provider logo uploaded. |
| Custom hooks | Optional helper functions exposed and consumed via `useCustomProviderFunction`. |

Once every item above is complete, run the usual test suites (backend + frontend
lint/tests) and push the changes.





