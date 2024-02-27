# AppArmor


![image](https://github.com/ChrisDevOpsOrg/kubernetes-cks/blob/main/images/syscall.png)



# Profile Modes
- Unconfined
  - Process can escape
- Complain
  - Process can escape but it will be logged
- Enforce
  - Process cannot escape


---
# Main Command
- Show all profiles
  ```sh
  aa-status
  ```

- generate a new profile(smart wrapper around aa-logprof)
  ```sh
  aa-genprof
  ```

- put profile in complain mode
  ```sh
  aa-complain
  ```


- put profile in enforce mode
  ```sh
  aa-enforce
  ```

- Update the profile if app produced some more usage logs(syslog)
  ```sh
  aa-logprof
  ```

