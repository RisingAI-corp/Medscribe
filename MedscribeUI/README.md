### UI README (`MedscribeUI/README.md`)

````markdown
# MedscribeUI – Frontend Setup & Philosophy

## Table of Contents

- [Introduction](#introduction)
- [Project Philosophy](#project-philosophy)
- [Setup Instructions](#setup-instructions)
- [Running the UI](#running-the-ui)
- [Code Quality & Linting](#code-quality--linting)
- [Component & Story Organization](#component--story-organization)
- [State Management Approach](#state-management-approach)
- [API Handling with React Query](#api-handling-with-react-query)
- [Commands Summary](#commands-summary)

---

## Introduction

Welcome to MedscribeUI! This guide outlines the setup for our frontend project and explains the design philosophy behind our choices.

---

## Project Philosophy

- **Vite for Speed:**  
  We chose [Vite](https://vitejs.dev/) because of its incredible speed and efficient development experience. Research more on your own if interested.
- **Component-Centric Story Organization:**  
  Instead of placing all Storybook stories in one massive directory, we include stories alongside the components that use them. This makes finding and maintaining stories much more intuitive.

- **State-Management Approach**  
  Our state management is based on atoms:

  - **Central Atom:**  
    A single, central atom holds the main state.
  - **Derived Atoms for Layouts:**  
    Each layout component has a derived atom that extracts only the necessary data from the central atom. This derived atom is responsible for providing data to its layout and its presentational (dumb) children.  
    _Benefits:_ Faster search, better organization, and more intuitive.

---

## Setup Instructions

1. **Install Dependencies**  
   Ensure you have Node.js installed, then run:
   ```bash
   npm install   # or yarn install
   ```
````

2. **Configure Prettier in VSCode**  
   To enable formatting on save, add the following to your VSCode settings:
   ```json
   {
     "editor.formatOnSave": true,
     "editor.defaultFormatter": "esbenp.prettier-vscode"
   }
   ```

---

## Running the UI

Start the frontend development server with:

```bash
npm run dev
```

---

## Code Quality & Linting

- **ESLint in VSCode:**  
  Make sure you have the ESLint extension installed. To run ESLint from the command line:
  ```bash
  npx eslint --fix .
  ```
- **VSCode Linting:**  
  The integrated linting in VSCode will help you catch errors as you type.

---

## Component & Story Organization

We organize stories directly within the folders of their respective components. This approach:

- Keeps the project structure clean.
- Makes searching for specific stories faster and more intuitive.

## API Handling with React Query

We use [React Query](https://react-query-v3.tanstack.com/) because it simplifies API request handling, caching, and state synchronization. It makes fetching and updating data effortless. Take some time to explore its benefits on your own!

---

## Commands Summary

- **Install Dependencies:**
  ```bash
  npm install   # or yarn install
  ```
- **Start the UI Development Server:**
  ```bash
  npm run dev
  ```
- **Run ESLint (with auto-fix):**
  ```bash
  npx eslint --fix .
  ```
- **Run Storybook (if needed):**  
  If you need to view UI components in isolation:
  ```bash
  yarn storybook
  ```

```

---
```
