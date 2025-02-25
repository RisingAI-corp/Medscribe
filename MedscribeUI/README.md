# Medscribe Project Setup

## Overview

This project uses **Vite** for front-end development and **Node.js** for back-end services. We have configured **ESLint** and **Prettier** for linting and formatting, ensuring that the codebase remains clean and consistent across all contributors.

### Key Points:

- **ESLint** is used for linting and ensuring code quality.
- **Prettier** is used for automatic code formatting.
- **TypeScript** is used for type safety.
- **Vite** is used as a bundler for front-end development.
- We use **feature branching** for individual work and rebase onto `main` to maintain a clean commit history.

## Setup Instructions

### 1. Fork the Repository

1. Fork the repository on GitHub by clicking the **Fork** button at the top of the repo page.
2. Clone the forked repository:

   ```bash
   git clone https://github.com/YourUsername/Medscribe.git
   cd Medscribe
   ```

### 2. Set Upstream to Main Repository

After cloning your fork, set the original repository as the upstream:

```bash
git remote add upstream https://github.com/RisingAI-corp/Medscribe.git
```

Verify your remotes:

```bash
git remote -v
```

You should see:

```bash
origin    https://github.com/YourUsername/Medscribe.git (fetch)
origin    https://github.com/YourUsername/Medscribe.git (push)
upstream  https://github.com/RisingAI-corp/Medscribe.git (fetch)
upstream  https://github.com/RisingAI-corp/Medscribe.git (push)
```

### 3. Create a New Branch for Each Feature

Always create a new branch for every feature or bug fix you work on:

```bash
git checkout -b feature/new-feature
```

This keeps your work isolated and avoids conflicts.

### 4. Install Dependencies

Before you start working, ensure all dependencies are installed by running:

```bash
npm install
```

### 5. Run Development Server

To start the Vite development server, run:

```bash
npm run dev
```

### 6. Prettier and ESLint Setup

We have configured **Prettier** and **ESLint** to work together without conflicts.

#### Formatting on Save with Prettier:

1. Ensure you have the Prettier extension installed in your **VSCode**.
2. Add the following settings to your VSCode settings:

```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode"
}
```

#### Prettier Configuration:

Prettier is configured using `.prettierrc.json`:

```json
{
  "trailingComma": "all",
  "tabWidth": 2,
  "semi": true,
  "singleQuote": true,
  "printWidth": 120,
  "bracketSpacing": true
}
```

To format files from the command line:

```bash
npx prettier --write .
```

#### ESLint Configuration:

The **ESLint** configuration lives in the `eslint.config.mjs` file:

```js
import js from '@eslint/js';
import globals from 'globals';
import reactHooks from 'eslint-plugin-react-hooks';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  { ignores: ['dist', 'node_modules'] },
  {
    extends: [
      js.configs.recommended,
      ...tseslint.configs.strictTypeChecked,
      ...tseslint.configs.stylisticTypeChecked,
    ],
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
    plugins: {
      'react-hooks': reactHooks,
    },
    rules: {
      ...reactHooks.configs.recommended.rules,
    },
  },
);
```

To run ESLint manually and automatically fix errors:

```bash
npx eslint --fix .
```

This will apply ESLint fixes to your specified directory.

### 7. Sync Your Fork with the Upstream

Before pushing your changes, ensure your fork is up-to-date with the upstream `main` branch by using rebase (not merge):

```bash
git pull upstream main --rebase
```

This will ensure your branch has the latest changes from the `main` branch and prevent merge conflicts.

### 8. Push Your Changes

After making changes on your feature branch, push them to your fork:

```bash
git push origin feature/new-feature
```

### 9. Open a Pull Request

Once your changes are pushed, go to your GitHub fork and open a Pull Request (PR) to merge your feature branch into the upstream repository's `main` branch.

Make sure to provide a clear description of your changes.

---

## Troubleshooting

### Common Errors:

#### Node.js Version Incompatibility

If you encounter an error like:

```
@eslint/js@9.10.0: The engine "node" is incompatible with this module. Expected version "^18.18.0 || ^20.9.0 || >=21.1.0". Got "18.13.0".
```

You can resolve this by upgrading your Node.js version to a compatible version (>=18.18.0).

#### Prettier Conflicts with ESLint

If Prettier and ESLint are conflicting, ensure both are properly configured not to override each other. See the configuration above to make sure they are working harmoniously.

#### Running Prettier on Save

Ensure that your VSCode settings have `formatOnSave` enabled and that the Prettier extension is properly installed.

#### Storybook

command to run storybook
yarn storybook

### Helpful Links:

- [ESLint Documentation](https://eslint.org/docs/latest/)
- [Prettier Documentation](https://prettier.io/docs/en/index.html)
- [Vite Documentation](https://vitejs.dev/guide/)

---

## Commands Summary

- **Start development server**: `npm run dev`
- **Lint files**: `npx eslint .`
- **Fix linting issues**: `npx eslint --fix yourdirectory/`
- **Format files**: `npx prettier --write .`
- **Pull with rebase**: `git pull upstream main --rebase`
- **Push your changes**: `git push origin feature/new-feature`
- **Runs storybook**: `yarn storybook`
