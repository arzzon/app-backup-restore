## Application Backup and Restore Tool

### Overview
The Kubernetes Backup and Restore Tool is a utility for managing backups of Kubernetes resources, such as Pods, Deployments, ConfigMaps, and PersistentVolumeClaims (PVCs). It allows users to define, create, and restore backups of application resources within a Kubernetes cluster.

### Features
- Define application namespaces and resources to be backed up.
- Create backups of specified resources.
- Restore backups to the same or different namespaces.

### Usage
1. Define Application

   To define an application for backup, use the /application/ endpoint with a PUT request, providing the namespace and name of the application.

Example:

    curl -X PUT -d '{"namespace": "test-mariadb", "name": "my-mariadb-app"}' http://localhost:8080/application/

2. Make a Backup
   
   To create a backup of an application, use the /backup/ endpoint with a PUT request, specifying the application ID.

Example:

    curl -X PUT -d '{"app": "<app_id>"}' http://localhost:8080/backup/

3. Restore a Backup
   
   To restore a backup to a namespace, use the /restore/ endpoint with a PUT request, providing the namespace and backup ID.

Example:
 
    curl -X PUT -d '{"namespace": "test-mariadb", "backupId": "<backup-id>"}' http://localhost:8080/restore/

### Installation

#### Build:
    make build

#### Run the application:
    ./backup-restore-tool