# Atom Fit

## Overview
Atom Fit is a project written in Go aimed at providing a robust API server for fitness applications. This repository contains the source code for the server, configuration files, and documentation.

## Features
- API server for fitness applications
- Database integration with PostgreSQL
- Environment-based configuration
- JWT authentication

## Installation
To install the project, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/stanislavCasciuc/atom-fit.git
    ```
2. Navigate to the project directory:
    ```sh
    cd atom-fit
    ```
3. Install dependencies:
    ```sh
    go mod download
    ```
4. Create a `.env` file in the root of the project with variables from main file:

5. Install migration tool:
    ```sh
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```
6. Run the migration:
    ```sh
    make migrate-up
    ```

## Usage
To use the project, run the following command:
Run the project:
 ```sh
 make run
## Environment Variables
The main file (`cmd/main/main.go`) relies on various environment variables to configure the application.
``````
## Contributing
Contributions are welcome! Please refer to the [Contributing Guidelines](CONTRIBUTING.md) for more information.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
``````
```
