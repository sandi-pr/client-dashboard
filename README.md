# CR3 Client Dashboard

## Install `microservice-create`
```bash
    $ sudo wget -O "/usr/local/bin/microservice-create" "https://gitlab.hexacode.org/cr3/client-dashboard/-/raw/main/microservice-create?ref_type=heads"
    $ sudo chmod +x "/usr/local/bin/microservice-create"
```

## Usage `microservice-create`
```bash
    $ microservice-create <folder> <namespace> <module_name> <module>
```
**Description:**
- `microservice-create`: for create microservice project
- `folder`: folder project
- `namespace`: namespace is microservice name and package name of frontend
- `module_name`: module_name is name of project
- `module`: module is golang module name

## Example `microservice-create`
```bash
    $ microservice-create my-microservice my-microservice "My Microservice" gitlab.hexacode.org/cr3/my-microservice
```

## Uninstall `microservice-create`
```bash
    $ sudo rm -rf "/usr/local/bin/microservice-create"
```

## Build docker image
```bash
    $ docker build -t client-dashboard .
```

## Run
```bash
    $ docker run -it --rm -p 8080:8080 --name client-dashboard client-dashboard
```

## Links
- Home Page

  **[GET]** `http://localhost:8080/client-dashboard/`

- Test Page

  **[GET]** `http://localhost:8080/client-dashboard/test`

- API Multiply

  **[POST]** `http://localhost:8080/client-dashboard/multiply`

- API Get Multiply Config

  **[GET]** `http://localhost:8080/client-dashboard/config/multiply`

- API Get Micro Service Config

  **[GET]** `http://localhost:8080/client-dashboard/configs`