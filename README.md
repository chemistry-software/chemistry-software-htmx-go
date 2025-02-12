# chemistry-software-htmx-go
A simple web app using Go, HTMX and Tailwind CSS

## Installation

This guide assumes you have [Go](https://go.dev/dl/) and [Node.js](https://nodejs.org/) installed already.

### Windows

1.  **Clone the repository:**
    Open a command prompt and run:
    ```cmd
    git clone https://github.com/chemistry-software/chemistry-software-htmx-go.git
    cd chemistry-software-htmx-go
    ```

2.  **Install npm dependencies:**
    In the project directory, run:
    ```cmd
    npm install
    ```
    This will install Tailwind CSS and other necessary development tools.

3.  **Run the development server:**
    In the project directory, run:
    ```cmd
    npm run dev
    ```
    This command starts both the Go backend server and the Tailwind CSS build process in development mode.

### Linux

1.  **Clone the repository:**
    Open a terminal and run:
    ```bash
    git clone https://github.com/chemistry-software/chemistry-software-htmx-go.git
    cd chemistry-software-htmx-go
    ```

2.  **Install npm dependencies:**
    In the project directory, run:
    ```bash
    npm install
    ```
    This will install Tailwind CSS and other necessary development tools.

3.  **Run the development server:**
    In the project directory, run:
    ```bash
    npm run dev
    ```
    This command starts both the Go backend server and the Tailwind CSS build process in development mode.

## Usage

After running the development server with `npm run dev`, navigate to `http://localhost:8080` in your web browser.

**Developing and Making Changes**

*   **Go Backend:**  Modify the Go files (`.go`) in the project. The server should automatically restart when you save changes.
*   **HTML Templates:** Edit the HTML templates in the `templates/` directory (`.html`). Refresh your browser to see changes.
*   **CSS Styling:**  Edit `static/css/input.css`. Tailwind CSS will automatically rebuild `static/css/tailwind.css` and your browser should hot-reload the CSS changes.
*   **JavaScript:** Modify JavaScript files in `static/js/`. Refresh your browser to see changes.

**Important Notes:**

*   **Node.js and npm are required** to build the CSS with Tailwind CSS. They are only needed during development, not for running the final Go application.
*   The `npm run dev` command starts both the Go server and the Tailwind CSS builder concurrently, making development easier.
*   Make sure you have Go and Node.js installed and configured correctly before starting.

# Copyright

Copyright 2025 [chemistry.software](https://chemistry.software)

Licensed under the [MIT License](https://opensource.org/licenses/MIT).

