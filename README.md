# G-Hopper

Welcome to G-Hopper! This is a fullstack web application built with Go and React. 

## Table of Contents
- [G-Hopper](#g-hopper)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Setup development environment](#setup-development-environment)
  - [Setting Up Samples Database (Important)](#setting-up-samples-database-important)
    - [prerequisites](#prerequisites-1)
    - [Steps](#steps)
  - [Available Commands](#available-commands)
    - [Development Commands](#development-commands)
    - [Production Commands](#production-commands)
    - [Utility Commands](#utility-commands)
  - [Project Structure](#project-structure)
  - [Environment Variables](#environment-variables)
    - [Development (.env.development)](#development-envdevelopment)
    - [Production](#production)
  - [License](#license)

## Prerequisites

Before you begin, ensure you have the following installed:
- Docker ([Mac](https://docs.docker.com/desktop/setup/install/mac-install/), [Windows](https://docs.docker.com/desktop/setup/install/windows-install/))
- Docker Compose (Included with Docker Desktop)
- Make 
  - **Mac:** `brew install make`
  - **Windows:** Use [Chocolatey](https://chocolatey.org/install) to install `make` with `choco install make` (You could also use WSL to run the project)
  - **Linux:** Linux distributions usually come with `make` pre-installed. If not, you can install it with `sudo apt install make`
- kubectl (Optional for production deployment only)

## Setup development environment

1. Clone the repository:
   ```bash
   git clone https://github.com/Emeruem-Kennedy1/ghopper.git
   cd ghopper
   ```

2. Create your environment file:
   ```bash
   cp .env.example .env.development
   ```

3. Configure your environment variables in `.env.development`
    SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET are required. You can get them from the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications) after creating a new application. 
    
    The SAMPLES_DB variables are for establishing a connection to a database that holds the songs and samples. You should setup a local database for development. See the next section for instructions on setting up the samples database.

4. On the spotify developer dashboard, add `http://localhost:9797/auth/spotify/callback` as a redirect URI for your application. If you are setting a different port for the backend, you should change the port in the redirect URI.

5. Start the development environment:
    ```bash
    make dev
    ```
    This command will:
    - Build and start all necessary containers
    - Set up the development environment with hot-reloading
    - Start the frontend on http://localhost:51920 (or whatever the environment variable `FRONTEND_PORT` is set to)
    - Start the backend on http://localhost:9797 (or whatever the environment variable `BACKEND_PORT` is set to)

    Other useful development commands:
    ```bash
    # Stop the development environment
    make dev-down

    # View logs
    make dev-logs

    # Access container shells
    make backend-shell
    make frontend-shell
    ```

## Setting Up Samples Database (Important)

### prerequisites
- MariaDB or any MySQL database (Note: If you don't use MariaDB you might need to do you might need to do more setup)
- Node.js
- Prisma CLI (`npm install -g prisma`)

### Steps

1. Create a folder (not in the project directory) for the samples database:
   ```bash
   mkdir samples-db
   cd samples-db
   ```
2. Start MySQL and create a new database:
   ```bash
   mysql -u root -p
   ```
   ```sql
   CREATE DATABASE ghopper;
   ```
3. Create a user and grant privileges:
   ```sql
    CREATE USER 'user'@'%' IDENTIFIED BY 'password';
    GRANT ALL PRIVILEGES ON ghopper.* TO 'user'@'%';
    FLUSH PRIVILEGES; 
    ```
4. Create a prisma folder and add a schema.prisma file:
    ```bash
    mkdir prisma
    cd prisma
    touch schema.prisma
    ```
5. Add the following to the schema.prisma file:
    ```python
    generator client {
    provider = "prisma-client-js"
    }

    datasource db {
    provider = "mysql"
    url      = "mysql://user:password@localhost:3306/ghopper"
    }

    model Artist {
    id             Int           @id @default(autoincrement())
    name           String        @unique
    createdAt      DateTime      @default(now())
    updatedAt      DateTime      @updatedAt
    songs          SongArtist[]
    }

    model Song {
    id            Int          @id @default(autoincrement())
    title         String
    releaseYear   Int?
    createdAt     DateTime     @default(now())
    updatedAt     DateTime     @updatedAt
    artists       SongArtist[]
    genres        Genre[]      @relation("SongToGenre")
    samplesUsed   Sample[]     @relation("SampleUsedInSong")
    sampledInSongs Sample[]    @relation("SongSampledIn")
    }

    model SongArtist {
    id        Int      @id @default(autoincrement())
    song      Song     @relation(fields: [songId], references: [id])
    songId    Int
    artist    Artist   @relation(fields: [artistId], references: [id])
    artistId  Int
    isMainArtist Boolean
    createdAt DateTime @default(now())
    updatedAt DateTime @updatedAt

    @@unique([songId, artistId, isMainArtist])
    @@index([songId])
    @@index([artistId])
    }

    model Genre {
    id        Int      @id @default(autoincrement())
    name      String   @unique
    createdAt DateTime @default(now())
    updatedAt DateTime @updatedAt
    songs     Song[]   @relation("SongToGenre")
    }

    model Sample {
    id              Int      @id @default(autoincrement())
    originalSongId  Int      @map("original_song_id")
    sampledInSongId Int      @map("sampled_in_song_id")
    createdAt       DateTime @default(now())
    updatedAt       DateTime @updatedAt
    originalSong    Song     @relation("SongSampledIn", fields: [originalSongId], references: [id])
    sampledInSong   Song     @relation("SampleUsedInSong", fields: [sampledInSongId], references: [id])

    @@unique([originalSongId, sampledInSongId])
    @@index([originalSongId])
    @@index([sampledInSongId])
    }
    ```

6. Run the following commands to generate the Prisma client and apply the schema to the database:
    ```bash
    prisma generate
    prisma db push
    ```

7. Add some data to the database. 
8. Setup the environment variables in the `.env.development` file:
    ```bash
    SAMPLES_DB_USER=user
    SAMPLES_DB_PASSWORD=password
    SAMPLES_DB_HOST=localhost
    SAMPLES_DB_PORT=3306
    SAMPLES_DB_NAME=ghopper
    ```


## Production (optional)

To deploy to production:

1. Build and push Docker images:
   ```bash
   make deploy
   ```

2. Apply Kubernetes configurations:
   ```bash
   make apply-all
   ```

**Note:** I will add a more detailed guide on everything that needs to be setup in a video tutorial.

## Available Commands

Here's a complete list of available make commands:

### Development Commands
- `make dev` - Start the development environment
- `make dev-down` - Stop the development environment
- `make dev-logs` - Show development logs
- `make backend-shell` - Access backend container shell
- `make frontend-shell` - Access frontend container shell

### Production Commands
- `make deploy` - Build and push Docker images
- `make apply-secrets` - Apply Kubernetes secrets
- `make apply-app` - Apply Kubernetes applications
- `make apply-all` - Apply all Kubernetes configurations

### Utility Commands
- `make clean` - Clean up Docker resources
- `make prune` - Remove unused Docker resources
- `make help` - Display help message

## Project Structure

```
.
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── Dockerfile.prod
│   └── ...
├── frontend/
│   ├── src/
│   ├── Dockerfile.prod
│   └── ...
├── k3s/
│   ├── apps/
│   └── config/
├── docker-compose.dev.yml
├── Makefile
└── README.md
```

## Environment Variables

See `.env.example` for a list of environment variables required for the project.

### Development (.env.development)
copy `.env.example` to .env.development and fill in the values
```bash
cp .env.example .env.development
```


### Production
Production secrets are managed through Kubernetes secrets. See `k3s/config/secrets.yaml`. Copy `k3s/config/secrets-template.yaml` to `k3s/config/secrets.yaml` and fill in the values.
```bash
cp k3s/config/secrets-template.yaml k3s/config/secrets.yaml
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
