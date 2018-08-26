# restic-api
This exposes the common functions of Restic with REST APIs. Originally created to be used with Schelly.

# Usage



# REST Endpoints

  - ```POST /backups```
    - When invoked, triggers a new backup
    - Request body: json 
      
      ```
        {
           source: [path to file or dir],
           name: [name of the collection]
        }
      ```
    - Response body: json 
     
      ```
        {
           id:{alphanumeric-backup-id},
           status:{backup-status}
           message:{backy2-message}
        }
      ```
      - status must be one of:
          - 'pending' - backup is not finished yet
          - 'canceled' - backup was canceled
          - 'success' - backup has completed successfuly
          - 'error' - there was an error on backup
      
    - Status code 201 if created successfuly

  - ```GET {webhook-url}/{backup-id}```
    - Invoked when Schelly wants to query a specific backup instance
    - Request header: ```{webhook-headers}```
    - Response body: json
    
       ```
         {
           id:{id},
           status:{backup-status},
           message:{backend message}
         }
       ```
    - Status code: 200 if found, 404 if not found

  - ```DELETE {webhook-url}/{backup-id}```
    - Invoked when Schelly wants to trigger a new backup
    - Request body: json ```{webhook-delete-body}```
    - Request header: ```{webhook-headers}```
    - Response body: json 
     
      ```
        {
           id:{alphanumeric-backup-id},
           status:{backup-status}
           message:{backend-message}
        }
      ```
      
    - Status code 200 if deleted successfuly
