# Project Setup

This repository is set up using **Vite** for fast development, **ESLint** for code linting, and **Prettier** for consistent code formatting.

## Getting Started

### Installation

To install the necessary dependencies, run the following command:

```bash
npm install
```

### Development

To start the Vite development server, use:

```bash
npm run dev
```

This will start the application and serve it at `http://localhost:3000`.

## Linting and Formatting

We use **ESLint** to ensure code quality and **Prettier** for consistent formatting. ESLint is integrated with Prettier to avoid conflicts between linting and formatting rules.

### Running ESLint

To manually check for linting errors, run:

```bash
npx eslint .
```

### Running Prettier

To format your codebase using Prettier, run:

```bash
npx prettier --write .
```

This will format all files based on the rules defined in the `.prettierrc.json` file.

### Format on Save (VSCode)

To enable **Prettier** to format your files automatically when you save in **VSCode**:

1. Install the **Prettier - Code Formatter** extension from the VSCode Marketplace.
2. Open VSCode settings (`Cmd + ,` or `Ctrl + ,`).
3. Search for `editor.formatOnSave` and enable it.
4. Set Prettier as the default formatter by searching for `editor.defaultFormatter` and choosing **Prettier**.

This ensures that your code will be formatted automatically on save.

## Configuration

### ESLint

ESLint is configured using the **flat config** format and integrates with TypeScript and React. The ESLint configuration file (`.eslintrc.mjs`) ensures that Prettier and ESLint work together without conflicts.

### Prettier

Prettier's rules are defined in the `.prettierrc.json` file. Here's an example of the current configuration:

```json
{
  "trailingComma": "all",
  "tabWidth": 2,
  "semi": true,
  "singleQuote": true,
  "printWidth": 100,
  "bracketSpacing": true
}
```

## Scripts

Here are the key scripts you can run:

- **Start the Development Server**: 
  ```bash
  npm run dev
  ```
- **Run ESLint**: 
  ```bash
  npx eslint .
  ```
- **Run Prettier to Format Code**: 
  ```bash
  npx prettier --write .
  ```

## Troubleshooting

- **Node Version Issues**: Ensure that your Node.js version is compatible with the project’s dependencies (Node.js >= 18.18.0). Use [nvm](https://github.com/nvm-sh/nvm) to switch between versions if needed.
  
- **ESLint Conflicts with Prettier**: Ensure that Prettier rules are applied by checking that `eslint-config-prettier` is included in your ESLint configuration.

---

### Useful Links

- [Vite Documentation](https://vitejs.dev/guide/)
- [ESLint Documentation](https://eslint.org/docs/latest/)
- [Prettier Documentation](https://prettier.io/docs/en/index.html)
