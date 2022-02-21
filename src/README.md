
## Project Structure

- `src` contains source code.
    - `manager` contains code of job manager
        - `main`
            - `manager.go` entry point of job manager
            - `manager.json` configuration of job manager
            - `webroot` front end of job manager
    - `worker` contains code of job worker
        - `main`
            - `worker.go` entry point of job worker
            - `worker.json` configuration of job worker
    - `common` contains common code shared by all modules
    - `test` contains api examples in .http format, can be run in IntelliJ IDEs
    - `lib` contains go examples interact with different tools
    - `util` contains database drivers and cron expression parse tool
- `software` contains linux programs
- `server-deployment-conf` contains server deployment configuration on Linux