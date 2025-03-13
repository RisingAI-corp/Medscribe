```markdown
# Medscribe Repository – Project Setup & Development

## Table of Contents
- [Project Setup Instructions](#project-setup-instructions)
  - [Fork, Clone, & Upstream Setup](#fork-clone--upstream-setup)
  - [Create a Feature Branch](#create-a-feature-branch)
  - [Install Dependencies](#install-dependencies)
- [Development Server](#development-server)
- [Code Quality & Formatting](#code-quality--formatting)
- [Git Workflow](#git-workflow)
- [Troubleshooting](#troubleshooting)

---

## Project Setup Instructions

### Fork, Clone & Upstream Setup

1. **Fork the Repository**  
   Click the **Fork** button on GitHub to create your own copy.

2. **Clone Your Fork**  
   ```bash
   git clone https://github.com/<YourUsername>/Medscribe.git
   ```

3. **Set Upstream**  
   Configure the original repository as upstream:
   ```bash
   git remote add upstream https://github.com/RisingAI-corp/Medscribe.git
   ```
   Verify the remote configuration:
   ```bash
   git remote -v
   ```

### Create a Feature Branch

Always create a new branch for each feature or bug fix:
```bash
git checkout -b feature/your-feature-name
```

### Install Dependencies

Run the following command in the project root to install backend dependencies:
```bash
npm install   # or yarn install, as preferred
```

---

## Development Server

You have two options to start the development environment:

### Option 1: Run Both Servers Concurrently

From the project root, run:
```bash
npx concurrently "go run cmd/api/api.go" "cd MedscribeUI && npm run dev"
```
This command uses `concurrently` to start both the backend and frontend servers at the same time.

### Option 2: Run Servers in Separate Terminals

1. **Start the Backend Server:**  
   Open one terminal in the project root and run:
   ```bash
   go run cmd/api/api.go
   ```

2. **Start the Frontend Server:**  
   Open another terminal, navigate to the UI directory, and run:
   ```bash
   cd MedscribeUI && npm run dev
   ```

---

## Code Quality & Formatting

For backend linting, run:
```bash
golangci-lint run
```

---

## Git Workflow

After making changes:

1. **Sync with Upstream:**  
   ```bash
   git pull upstream main --rebase
   ```

2. **Push Your Feature Branch:**  
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Open a Pull Request:**  
   Create a PR to merge your feature branch into the upstream `main` branch.

---

## Troubleshooting

- **Node.js Version Error:**  
  If you encounter errors like:
  ```
  @eslint/js@9.10.0: The engine "node" is incompatible with this module. Expected version "^18.18.0 || ^20.9.0 || >=21.1.0". Got "18.13.0".
  ```
  upgrade your Node.js version to at least **18.18.0**.

  
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

## building the UI
- with vite
```bash
yarn build --mode production ##if you want to use prod base url

yarn build ##if you want to use localhost/local base  url

```
then just server the dist folder and index.html

## building and running docker image
```bash
docker build -t medscribe:latest . ## build
docker run -p 8080:8080 medscribe:latest ## run
```

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

---

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