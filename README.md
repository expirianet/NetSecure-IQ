# NetSecure-IQ

![NetSecure-IQ Dashboard Overview](image_dashboard_overview.png)

Un'innovativa soluzione SaaS (Software as a Service), multi-tenant e multi-site, progettata per rivoluzionare il monitoraggio della sicurezza fisica e cibernetica. In un'era in cui la convergenza di questi due domini è essenziale, NetSecure-IQ offre una supervisione costante e proattiva delle installazioni, garantendo visibilità completa e controllo granulare sugli asset di sicurezza.

## Indice

* [Introduzione](#introduzione)
* [Caratteristiche Principali](#caratteristiche-principali)
* [Architettura](#architettura)
* [Conformità e Sicurezza](#conformità-e-sicurezza)
* [Tecnologie Utilizzate](#tecnologie-utilizzate)
* [Come Iniziare (Per Sviluppatori)](#come-iniziare-per-sviluppatori)
    * [Prerequisiti](#prerequisiti)
    * [Setup dell'Ambiente di Sviluppo](#setup-dellambiente-di-sviluppo)
    * [Configurazione Database](#configurazione-database)
    * [Configurazione Servizi Backend Go](#configurazione-servizi-backend-go)
    * [Configurazione Frontend Vue.js](#configurazione-frontend-vuejs)
    * [Deploy su Cluster Kubernetes (Produzione)](#deploy-su-cluster-kubernetes-produzione)
    * [Esempio di docker-compose.yaml (per sviluppo standalone)](#esempio-di-docker-composeyaml-per-sviluppo-standalone)
    * [Simula un MikroTik (Per PoC/Test)](#simula-un-mikrotik-per-poctest)
* [Esecuzione dei Test](#esecuzione-dei-test)
* [Contribuire](#contribuire)
* [Licenza](#licenza)
* [Contatti](#contatti)

## Introduzione

NetSecure-IQ è una piattaforma dedicata al monitoraggio centralizzato e in tempo reale di sistemi di videosorveglianza (telecamere IP, NVR), allarmi e controlli di accesso. Si distingue per la sua integrazione nativa e intelligente tramite router MikroTik on-site, che consente un monitoraggio approfondito della rete.

## Caratteristiche Principali

* **Monitoraggio in Tempo Reale e Stato dei Dispositivi:** Visualizzazione immediata dello stato On/Off di telecamere, NVR, sensori e pannelli di allarme. Identificazione di indirizzi IP e MAC, monitoraggio del traffico di rete (RX/TX) e scansione periodica delle porte aperte e dei protocolli utilizzati sui dispositivi MikroTik.
* **Dashboard Interattive e Personalizzabili:** Creazione e configurazione di dashboard intuitive per visualizzare metriche pertinenti come grafici storici del traffico di rete, andamenti dell'uptime dei dispositivi, riepiloghi degli eventi di allarme e visualizzazione delle porte aperte nel tempo.
* **Notifiche Intelligenti e Avvisi Personalizzabili:** Avvisi immediati via email, SMS o push in caso di malfunzionamenti dei dispositivi, interruzioni di servizio, rilevamento di nuove porte aperte o traffico anomalo. Possibilità di configurare soglie personalizzate per l'attivazione degli avvisi.
* **Gestione Multi-tenant e Multi-site:** Ogni tenant (cliente) ha la propria area dati e configurazioni completamente isolate. Un tenant può gestire un numero illimitato di siti con una visione centralizzata di tutte le installazioni.
* **Reportistica Dettagliata:** Generazione di report storici per analisi retrospettive e audit, inclusi report sull'uptime dei dispositivi, cronologia degli allarmi e degli eventi, e analisi del traffico di rete.
* **Esperienza Utente Intuitiva (Frontend Vue.js):** Progettato per offrire un'esperienza utente pulita e ricca di informazioni, con una panoramica immediata dello stato generale del sistema, vista dettagliata per ogni sito/dispositivo, e un'interfaccia per la configurazione delle regole di notifica e la visualizzazione della cronologia degli avvisi.

## Architettura

La nostra architettura è robusta e scalabile, costruita su tecnologie moderne per garantire prestazioni e affidabilità.

![NetSecure-IQ Logical Diagram](NetSecure%20IQ%20-%20EN01-image-000.png)
*(Questo è un segnapato per il diagramma logico dell'architettura che è presente nel documento originale.)*

### Componenti Principali:

* **Backend:** Sviluppato in **Go** per garantire performance e scalabilità dei microservizi. Include servizi per Bootstrap/Provisioning, Data Ingestor, Data Query/API, Notifications/Alerting, e Configuration/Administration.
* **Frontend:** Realizzato con **Vue.js**, offre un'interfaccia utente intuitiva e reattiva per la gestione e la visualizzazione dei dati.
* **Orchestrazione:** Utilizza **Kubernetes** per la gestione e la scalabilità dei container.
* **Database:**
    * **InfluxDB:** Database time-series ottimizzato per dati di monitoraggio (uptime, traffico RX/TX, porte aperte).
    * **PostgreSQL:** Database relazionale per la gestione di dati strutturati come tenant, siti, MikroTik (MAC, chiavi WireGuard), utenti e regole di notifica.
* **Comunicazione Sicura:** Avviene tramite **WireGuard VPN tunnels** per la crittografia in transito, garantendo una connessione sicura tra i dispositivi MikroTik e la piattaforma.
* **Dispositivi Edge:** Router **MikroTik** (es. hAp Lite RB941-2nD) per la raccolta dati on-site e l'invio di metriche periodiche.

## Conformità e Sicurezza

Dalla fase di progettazione, abbiamo posto un'enfasi primaria sulla conformità normativa. NetSecure-IQ è sviluppato in aderenza ai principi del **GDPR (General Data Protection Regulation)** per la protezione dei dati personali e alla direttiva **NIS2 (Network and Information Systems Directive 2)** per la cybersecurity delle infrastrutture critiche.

Principi di sicurezza e privacy fondamentali includono: minimizzazione dei dati, crittografia (in transito via WireGuard VPN tunnels e a riposo), controlli granulari degli accessi, pianificazione della risposta agli incidenti (Incident Response) e audit trail dettagliati.

## Tecnologie Utilizzate

* **Backend:** Go
* **Frontend:** Vue.js
* **Orchestrazione:**
    * **Sviluppo:** Docker, Minikube (opzionale)
    * **Produzione:** Kubernetes (cluster a 3 nodi)
* **Database Time-Series:** InfluxDB
* **Database Relazionale:** PostgreSQL
* **VPN:** WireGuard
* **Dispositivi Edge:** MikroTik RouterOS
* **Containerizzazione:** Docker
* **Monitoraggio Cluster (Produzione):** Prometheus, Grafana (probabile)

## Come Iniziare (Per Sviluppatori)

Queste istruzioni ti guideranno nel setup di un ambiente di sviluppo locale per NetSecure-IQ.

### Prerequisiti

Assicurati di avere installato i seguenti strumenti:

* [Git](https://git-scm.com/downloads)
* [Go](https://golang.org/doc/install) (versione X.X o superiore)
* [Node.js](https://nodejs.org/en/download/) (versione X.X o superiore) e [Yarn](https://classic.yarnpkg.com/en/docs/install/) (consigliato per Vue.js)
* [Docker Desktop](https://www.docker.com/products/docker-desktop)
* [Minikube](https://minikube.sigs.k8s.io/docs/start/) (opzionale, per simulare un cluster Kubernetes locale)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) (strumento da riga di comando per interagire con Kubernetes)
* Un editor di codice (es. [VS Code](https://code.visualstudio.com/))

### Setup dell'Ambiente di Sviluppo

1.  **Clona il repository:**
    ```bash
    git clone [https://github.com/expirianet/NetSecure-IQ.git](https://github.com/expirianet/NetSecure-IQ.git)
    cd NetSecure-IQ
    ```

2.  **Configurazione per Sviluppo Standalone (con Docker e Minikube - Opzionale):**
    * Se *non* usi Minikube, puoi eseguire i servizi backend e frontend direttamente con Docker Compose. Crea un file `docker-compose.yaml` (vedi esempio sotto) e usa `docker-compose up`.
    * Se *usi* Minikube:
        * Avvia Minikube: `minikube start`
        * Configura `kubectl` per usare Minikube: `kubectl config use-context minikube`
        * Applica i deployment e i servizi Kubernetes (inclusi InfluxDB e PostgreSQL):
            ```bash
            kubectl apply -f config/kubernetes/database/postgresql-deployment.yaml
            kubectl apply -f config/kubernetes/database/influxdb-deployment.yaml
            # E gli eventuali Persistent Volume Claims (PVC)
            kubectl apply -f config/kubernetes/services/
            kubectl apply -f config/kubernetes/deployments/
            ```

3.  **Configurazione Database (InfluxDB e PostgreSQL):**

    * **InfluxDB:**
        * Per lo sviluppo, puoi usare l'immagine Docker ufficiale.
        * In produzione, considera Persistent Volumes e Persistent Volume Claims (PVC) per i dati.
        * Configura le variabili d'ambiente (es. `INFLUXDB_ADMIN_USER`, `INFLUXDB_ADMIN_PASSWORD`) nel deployment Kubernetes.
        * Esempio di deployment (in `config/kubernetes/database/influxdb-deployment.yaml`):
            ```yaml
            apiVersion: apps/v1
            kind: Deployment
            metadata:
              name: influxdb
            spec:
              # ...
              template:
                spec:
                  containers:
                    - name: influxdb
                      image: influxdb:latest
                      ports:
                        - containerPort: 8086
                      env:
                        - name: INFLUXDB_ADMIN_USER
                          value: "admin"
                        - name: INFLUXDB_ADMIN_PASSWORD
                          value: "your_password"
                      volumeMounts:
                        - name: influxdb-data
                          mountPath: /var/lib/influxdb
                  volumes:
                    - name: influxdb-data
                      # ... (Configura il Persistent Volume o EmptyDir)
            ```

    * **PostgreSQL:**
        * Per lo sviluppo, puoi usare l'immagine Docker ufficiale.
        * In produzione, *usa* Persistent Volumes e Persistent Volume Claims (PVC) per i dati.
        * Configura le variabili d'ambiente (es. `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`) nel deployment Kubernetes.
        * Esempio di deployment (in `config/kubernetes/database/postgresql-deployment.yaml`):
            ```yaml
            apiVersion: apps/v1
            kind: Deployment
            metadata:
              name: postgresql
            spec:
              # ...
              template:
                spec:
                  containers:
                    - name: postgresql
                      image: postgres:13
                      ports:
                        - containerPort: 5432
                      env:
                        - name: POSTGRES_USER
                          value: "netsecure_iq"
                        - name: POSTGRES_PASSWORD
                          value: "your_secure_password"
                        - name: POSTGRES_DB
                          value: "netsecure_iq_db"
                      volumeMounts:
                        - name: postgresql-data
                          mountPath: /var/lib/postgresql/data
                  volumes:
                    - name: postgresql-data
                      # ... (Configura il Persistent Volume o EmptyDir)
            ```

4.  **Configurazione Servizi Backend Go:**
    * Naviga nella directory del backend: `cd src/backend`
    * Installa le dipendenze Go: `go mod tidy`
    * (Opzionale) Configura le variabili d'ambiente necessarie (vedi `config/backend_env/`). Queste *devono* corrispondere alle configurazioni del database (es. indirizzi, username, password).
    * In sviluppo standalone, puoi eseguire i servizi direttamente con `go run main.go` (o il comando specifico per il servizio). In produzione, usa le immagini Docker e i deployment Kubernetes.

5.  **Configurazione Frontend Vue.js:**
    * Naviga nella directory del frontend: `cd src/frontend`
    * Installa le dipendenze Node.js: `yarn install` (o `npm install`)
    * (Opzionale) Configura le variabili d'ambiente necessarie (vedi `config/frontend_env/`). Queste *devono* corrispondere all'indirizzo del backend.
    * In sviluppo, avvia il server di sviluppo Vue.js: `yarn serve` (o `npm run serve`). In produzione, usa l'immagine Docker.

6.  **Deploy su Cluster Kubernetes (Produzione):**
    * Assicurati che `kubectl` sia configurato per puntare al tuo cluster Kubernetes di produzione (a 3 nodi).
    * Applica i deployment e i servizi Kubernetes (inclusi InfluxDB e PostgreSQL, se non già esistenti):
        ```bash
        kubectl apply -f config/kubernetes/database/postgresql-deployment.yaml
        kubectl apply -f config/kubernetes/database/influxdb-deployment.yaml
        # E gli eventuali Persistent Volume Claims (PVC)
        kubectl apply -f config/kubernetes/services/
        kubectl apply -f config/kubernetes/deployments/
        ```
    * Configura l'Ingress Controller per esporre i servizi.
    * Configura il WireGuard VPN Server sul cluster (vedi `config/kubernetes/` per i dettagli).

7.  **Esempio di `docker-compose.yaml` (per sviluppo standalone):**
    ```yaml
    version: "3.9"
    services:
      postgresql:
        image: postgres:13
        ports:
          - "5432:5432"
        volumes:
          - postgresql_data:/var/lib/postgresql/data
        environment:
          POSTGRES_USER: netsecure_iq
          POSTGRES_PASSWORD: your_secure_password
          POSTGRES_DB: netsecure_iq_db

      influxdb:
        image: influxdb:latest
        ports:
          - "8086:8086"
        volumes:
          - influxdb_data:/var/lib/influxdb
        environment:
          INFLUXDB_ADMIN_USER: admin
          INFLUXDB_ADMIN_PASSWORD: your_password

      backend:
        build: ./src/backend
        ports:
          - "8000:8000" # Esempio, adatta le porte
        environment:
          # ... (Variabili d'ambiente per il backend, inclusi i dettagli del database)
        depends_on:
          - postgresql
          - influxdb

      frontend:
        build: ./src/frontend
        ports:
          - "8080:8080" # Esempio, adatta le porte
        environment:
          # ... (Variabili d'ambiente per il frontend, incluso l'indirizzo del backend)
        depends_on:
          - backend

    volumes:
      postgresql_data:
      influxdb_data:
    ```

8.  **Simula un MikroTik (Per PoC/Test):**
    * Per la fase PoC o per i test, potresti dover eseguire manualmente o simulare gli script MikroTik per inviare dati al tuo backend.

## Esecuzione dei Test

* **Test Unitari:**
    * Backend: `go test ./...` da `src/backend`
    * Frontend: `yarn test:unit` (o `npm run test:unit`) da `src/frontend`
* **Test di Integrazione e E2E:**
    * Segui le istruzioni specifiche nella directory `tests/integration` e `tests/e2e`.

## Contribuire

Accogliamo con favore i contributi dei membri autorizzati del team. Tutti i contributi sono soggetti agli accordi di proprietà intellettuale aziendali.

## Licenza

Questo progetto è una proprietà aziendale di Experianet. Tutti i diritti sono riservati. Il codice sorgente non è open source e non può essere riprodotto, modificato o distribuito senza esplicita autorizzazione. L'utilizzo di NetSecure-IQ è regolamentato da specifici accordi di servizio e termini d'uso forniti dalla nostra azienda.

Per informazioni commerciali o partnership, si prega di contattarci.

## Contatti

Per domande o supporto, puoi contattare il team di sviluppo tramite:

* GitHub Issues: [https://github.com/expirianet/NetSecure-IQ/issues](https://github.com/expirianet/NetSecure-IQ/issues)
* Email: [info@expirianet.it](mailto:info@expirianet.it)

---

### File: `README.md` (English)

# NetSecure-IQ

![NetSecure-IQ Dashboard Overview](image_dashboard_overview.png)

An innovative SaaS (Software as a Service), multi-tenant, and multi-site solution designed to revolutionize physical and cybersecurity monitoring. In an era where the convergence of these two domains is essential, NetSecure-IQ offers constant and proactive supervision of installations, ensuring comprehensive visibility and granular control over security assets.

## Table of Contents

* [Introduction](#introduction)
* [Key Features](#key-features)
* [Architecture](#architecture)
* [Compliance and Security](#compliance-and-security)
* [Technologies Used](#technologies-used)
* [Getting Started (For Developers)](#getting-started-for-developers)
    * [Prerequisites](#prerequisites)
    * [Setting Up the Development Environment](#setting-up-the-development-environment)
    * [Database Configuration](#database-configuration)
    * [Configuring Go Backend Services](#configuring-go-backend-services)
    * [Configuring Vue.js Frontend](#configuring-vuejs-frontend)
    * [Deploy to Kubernetes Cluster (Production)](#deploy-to-kubernetes-cluster-production)
    * [Example docker-compose.yaml (for standalone development)](#example-docker-composeyaml-for-standalone-development)
    * [Simulate a MikroTik (For PoC/Testing)](#simulate-a-mikrotik-for-poctesting)
* [Running Tests](#running-tests)
* [Contributing](#contributing)
* [License](#license)
* [Contact](#contact)

## Introduction

NetSecure-IQ is a platform dedicated to centralized, real-time monitoring of video surveillance systems (IP cameras, NVRs), alarms, and access controls. A distinctive strength of NetSecure-IQ is its native and intelligent integration via on-site MikroTik routers, which allows for in-depth network monitoring.

## Key Features

* **Real-time Monitoring and Device Status:** Immediate On/Off status visualization for cameras, NVRs, sensors, and alarm panels. Identification of IP and MAC addresses, network traffic monitoring (RX/TX), and periodic scanning of open ports and protocols used on MikroTik devices.
* **Interactive and Customizable Dashboards:** Creation and configuration of intuitive dashboards to visualize relevant metrics such as historical network traffic graphs, device uptime trends, alarm event summaries, and open port visualization over time.
* **Intelligent Notifications and Customizable Alerts:** Immediate alerts via email, SMS, or push in case of device malfunctions, service interruptions, detection of new open ports, or anomalous traffic. Ability to configure custom thresholds for alert activation.
* **Multi-tenant and Multi-site Management:** Each tenant (client) has their data area and configurations completely isolated. A tenant can manage an unlimited number of sites with a centralized view of all installations.
* **Detailed Reporting:** Generation of historical reports for retrospective analysis and audits, including device uptime reports, alarm and event history, and network traffic analysis.
* **Intuitive User Experience (Vue.js Frontend):** The NetSecure IQ frontend, developed with Vue.js, aims to offer an intuitive, clean, and information-rich user experience. This includes an immediate overview of the overall system status with key indicators, a detailed view for each individual site or device, and an interface for configuring notification rules and viewing the history of generated alerts.

## Architecture

Our architecture is robust and scalable, built on modern technologies to ensure performance and reliability.

![NetSecure-IQ Logical Diagram](NetSecure%20IQ%20-%20EN01-image-000.png)
*(This is a placeholder for the logical architecture diagram from the original document.)*

### Core Components:

* **Backend:** Built on **Go** for high performance and scalability of microservices. This includes services for Bootstrap/Provisioning, Data Ingestor, Data Query/API, Notifications/Alerting, and Configuration/Administration.
* **Frontend:** Developed with **Vue.js**, offering an intuitive and reactive user interface for data management and visualization.
* **Orchestration:** Leverages **Kubernetes** for container management and scalability.
* **Databases:**
    * **InfluxDB:** A time-series database optimized for metric data like uptime, RX/TX traffic, and open ports.
    * **PostgreSQL:** A relational database for managing structured data such as tenants, sites, MikroTik (MAC, WG Keys), users, and notification rules.
* **Secure Communication:** Achieved through **WireGuard VPN tunnels** for in-transit encryption, ensuring a secure communication channel between MikroTik devices and the platform.
* **Edge Devices:** **MikroTik** routers (e.g., hAp Lite RB941-2nD) for on-site data collection and periodic metric submission.

## Compliance and Security

From the design phase, we have placed a primary focus on regulatory compliance. NetSecure-IQ is developed with adherence to **GDPR (General Data Protection Regulation)** principles for personal data protection and the **NIS2 (Network and Information Systems Directive 2)** directive for critical infrastructure cybersecurity. This offers peace of mind and reliability to our clients.

Fundamental security and privacy pillars include: data minimization, encryption (in transit via WireGuard VPN tunnels and at rest), granular access controls, Incident Response planning, and detailed audit trails.

## Technologies Used

* **Backend:** Go
* **Frontend:** Vue.js
* **Orchestration:**
    * **Development:** Docker, Minikube (optional)
    * **Production:** Kubernetes (3-node cluster)
* **Time-Series Database:** InfluxDB
* **Relational Database:** PostgreSQL
* **VPN:** WireGuard
* **Edge Devices:** MikroTik RouterOS
* **Containerization:** Docker
* **Cluster Monitoring (Production):** Prometheus, Grafana (likely)

## Getting Started (For Developers)

These instructions will guide you through setting up a local development environment for NetSecure-IQ.

### Prerequisites

Ensure you have the following tools installed:

* [Git](https://git-scm.com/downloads)
* [Go](https://golang.org/doc/install) (version X.X or higher)
* [Node.js](https://nodejs.org/en/download/) (version X.X or higher) and [Yarn](https://classic.yarnpkg.com/en/docs/install/) (recommended for Vue.js)
* [Docker Desktop](https://www.docker.com/products/docker-desktop)
* [Minikube](https://minikube.sigs.k8s.io/docs/start/) (optional, to simulate a local Kubernetes cluster)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) (command-line tool for interacting with Kubernetes)
* A code editor (e.g., [VS Code](https://code.visualstudio.com/))

### Setting Up the Development Environment

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/expirianet/NetSecure-IQ.git
    cd NetSecure-IQ
    ```

2.  **Standalone Development Setup (with Docker and Minikube - Optional):**
    * If you are *not* using Minikube, you can run the backend and frontend services directly with Docker Compose. Create a `docker-compose.yaml` file (see example below) and use `docker-compose up`.
    * If you *are* using Minikube:
        * Start Minikube: `minikube start`
        * Configure `kubectl` to use Minikube: `kubectl config use-context minikube`
        * Apply the Kubernetes deployments and services (including InfluxDB and PostgreSQL):
            ```bash
            kubectl apply -f config/kubernetes/database/postgresql-deployment.yaml
            kubectl apply -f config/kubernetes/database/influxdb-deployment.yaml
            # And any Persistent Volume Claims (PVCs)
            kubectl apply -f config/kubernetes/services/
            kubectl apply -f config/kubernetes/deployments/
            ```

3.  **Database Configuration (InfluxDB and PostgreSQL):**

    * **InfluxDB:**
        * For development, you can use the official Docker image.
        * In production, *use* Persistent Volumes and Persistent Volume Claims (PVCs) for the data.
        * Configure environment variables (e.g., `INFLUXDB_ADMIN_USER`, `INFLUXDB_ADMIN_PASSWORD`) in the Kubernetes deployment.
        * Example deployment (in `config/kubernetes/database/influxdb-deployment.yaml`):
            ```yaml
            apiVersion: apps/v1
            kind: Deployment
            metadata:
              name: influxdb
            spec:
              # ...
              template:
                spec:
                  containers:
                    - name: influxdb
                      image: influxdb:latest
                      ports:
                        - containerPort: 8086
                      env:
                        - name: INFLUXDB_ADMIN_USER
                          value: "admin"
                        - name: INFLUXDB_ADMIN_PASSWORD
                          value: "your_password"
                      volumeMounts:
                        - name: influxdb-data
                          mountPath: /var/lib/influxdb
                  volumes:
                    - name: influxdb-data
                      # ... (Configure the Persistent Volume or EmptyDir)
            ```

    * **PostgreSQL:**
        * For development, you can use the official Docker image.
        * In production, *always use* Persistent Volumes and Persistent Volume Claims (PVCs) for the data.
        * Configure environment variables (e.g., `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`) in the Kubernetes deployment.
        * Example deployment (in `config/kubernetes/database/postgresql-deployment.yaml`):
            ```yaml
            apiVersion: apps/v1
            kind: Deployment
            metadata:
              name: postgresql
            spec:
              # ...
              template:
                spec:
                  containers:
                    - name: postgresql
                      image: postgres:13
                      ports:
                        - containerPort: 5432
                      env:
                        - name: POSTGRES_USER
                          value: "netsecure_iq"
                        - name: POSTGRES_PASSWORD
                          value: "your_secure_password"
                        - name: POSTGRES_DB
                          value: "netsecure_iq_db"
                      volumeMounts:
                        - name: postgresql-data
                          mountPath: /var/lib/postgresql/data
                  volumes:
                    - name: postgresql-data
                      # ... (Configure the Persistent Volume or EmptyDir)
            ```

4.  **Configuring Go Backend Services:**
    * Navigate to the backend directory: `cd src/backend`
    * Install Go dependencies: `go mod tidy`
    * (Optional) Configure necessary environment variables (see `config/backend_env/`). These *must* match your database configurations (e.g., addresses, usernames, passwords).
    * In standalone development, you can run services directly with `go run main.go` (or the specific command for the service). In production, use the Docker images and Kubernetes deployments.

5.  **Configuring Vue.js Frontend:**
    * Navigate to the frontend directory: `cd src/frontend`
    * Install Node.js dependencies: `yarn install` (or `npm install`)
    * (Optional) Configure necessary environment variables (see `config/frontend_env/`). These *must* match your backend address.
    * In development, start the Vue.js development server: `yarn serve` (or `npm run serve`). In production, use the Docker image.

6.  **Deploy to Kubernetes Cluster (Production):**
    * Ensure `kubectl` is configured to point to your production Kubernetes cluster (3-node).
    * Apply the Kubernetes deployments and services (including InfluxDB and PostgreSQL, if not already existing):
        ```bash
        kubectl apply -f config/kubernetes/database/postgresql-deployment.yaml
        kubectl apply -f config/kubernetes/database/influxdb-deployment.yaml
        # And any Persistent Volume Claims (PVCs)
        kubectl apply -f config/kubernetes/services/
        kubectl apply -f config/kubernetes/deployments/
        ```
    * Configure the Ingress Controller to expose the services.
    * Configure the WireGuard VPN Server on the cluster (refer to `config/kubernetes/` for details).

7.  **Example `docker-compose.yaml` (for standalone development):**
    ```yaml
    version: "3.9"
    services:
      postgresql:
        image: postgres:13
        ports:
          - "5432:5432"
        volumes:
          - postgresql_data:/var/lib/postgresql/data
        environment:
          POSTGRES_USER: netsecure_iq
          POSTGRES_PASSWORD: your_secure_password
          POSTGRES_DB: netsecure_iq_db

      influxdb:
        image: influxdb:latest
        ports:
          - "8086:8086"
        volumes:
          - influxdb_data:/var/lib/influxdb
        environment:
          INFLUXDB_ADMIN_USER: admin
          INFLUXDB_ADMIN_PASSWORD: your_password

      backend:
        build: ./src/backend
        ports:
          - "8000:8000" # Example, adjust ports
        environment:
          # ... (Environment variables for the backend, including database details)
        depends_on:
          - postgresql
          - influxdb

      frontend:
        build: ./src/frontend
        ports:
          - "8080:8080" # Example, adjust ports
        environment:
          # ... (Environment variables for the frontend, including the backend address)
        depends_on:
          - backend

    volumes:
      postgresql_data:
      influxdb_data:
    ```

8.  **Simulate a MikroTik (For PoC/Testing):**
    * For the PoC phase or testing, you might need to manually run or simulate MikroTik scripts to send data to your backend.

## Running Tests

* **Unit Tests:**
    * Backend: `go test ./...` from `src/backend`
    * Frontend: `yarn test:unit` (or `npm run test:unit`) from `src/frontend`
* **Integration and E2E Tests:**
    * Follow specific instructions in the `tests/integration` and `tests/e2e` directories.

## Contributing

We welcome contributions from authorized team members. All contributions are subject to company intellectual property agreements.

## License

This project is a proprietary asset of Experianet. All rights reserved. The source code is not open source and may not be reproduced, modified, or distributed without explicit authorization. The use of NetSecure-IQ is governed by specific service agreements and terms of use provided by our company.

For commercial inquiries or partnerships, please contact us.

## Contact

For questions or support, you can reach the development team via:

* GitHub Issues: [https://github.com/expirianet/NetSecure-IQ/issues](https://github.com/expirianet/NetSecure-IQ/issues)
* Email: [info@expirianet.it](mailto:info@expirianet.it)

---
