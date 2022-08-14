# Threat Modeling Report for K8s based - Test BooksStore Application



**Threats are majorly based on K8s Security, JWT Tokens & Application Security Fundamentals.** 
**(Assumptions - Books & User records are stored in Postgres DB)**




|Sr No                                                                                                                                                                   |Threat                                                                                        |Severity                 |Breach Category                  |Mitigation             |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------|-------------------------|---------------------------------|-----------------------|
|1                                                                                                                                                                       |A malicious user can extract sensitive information stored in ConfigMap.  A malicious user can extract additional information from DBs or can extract the user/admin information by misusing sensitive information stored in ConfigMap.|High                     |Information Disclosure           |Ensure that we must have a cloud vault integration (like hashicorp) in place and all the sensitive information should be stored in the cloud vault.  To decrease the severity level to medium, we should deploy the sensitive information as a secret rather than deploying it as a ConfigMap on K8s cluster.   But, please note that deploying the sensitive information as a secret will only ‘base64 encode’ the sensitive information/secrets on the k8s cluster, which could still be exploited.|
|2                                                                                                                                                                       |A malicious user can extract the sensitive information as the application configuration-files are running in ‘Default’ namespace.  A malicious user can assign PersistentVolumeClaim to any other deployment (running in default namespace).|Medium                   |Information Disclosure           |Ensure that ConfigMaps, Secrets, Services, Deployments and PersistentVolumeClaims are deployed in a specific namespace rather than running them in the default namespace. This ensures that only specified namespaced deployments can access the resources.   Ensure appropriate labels are assigned to each resource for restricting access to the resources, based on specific labels only.  Ensure that we have RBAC policies with appropriate Role and RoleBinding configured for the admins/users/apps accessing specified namespace.|
|3                                                                                                                                                                       |A malicious user can infiltrate into deployment files, if A malicious user has access to Default namespace and therefore has access to Default ServiceAccount.  A malicious user can delete and re-execute the deployment files with malfunctioned changes having adverse effects.|Medium                   |Tampering                        |Ensure that deployment files are configured based on their specific service accounts.   Service accounts should be mapped to the specified namespace.  Service accounts should be bound with RBAC permissions by applying appropriate Role and RoleBinding|
|4                                                                                                                                                                       |If A malicious user has shell access into the container (deployment files), then A malicious user can make rogue/unauthorized API calls to K8s cluster (Kube API server) and can extract sensitive cluster specific information.|High                     |Information Disclosure           |Disable shell access for the containers.  Remove binary executable files from the container to ensure that files cannot be re-ran, once the containers are up and running.  Ensure that while deploying service accounts for each deployment (within a specific namespace), the automount service account token is set to false. automountServiceAccountToken: false  This ensures that Kube-API calls are blocked from the container to the K8s cluster.|
|5                                                                                                                                                                       |A malicious user can perform MITM or can extract the sensitive information if TLS is not enabled for the Books website.|High                     |Information Disclosure, Tampering|Implement Secure Ingress resource by specifying a Secret that contains a TLS private key and certificate. The Ingress resource only supports a single TLS port, 443, and assumes TLS termination at the ingress point (traffic to the Service and its Pods is in plaintext).|
|6                                                                                                                                                                       |A malicious user can make system calls / kernel calls from containers, if containers have root privileges.|High                     |Information Disclosure, Tampering Spoofing, Elevation of Privilege|Ensure that containers are running with non-root privileges.  Ensure that below is set for deployment files, to restrict root calls from running containers:  runAsNonRoot: true privileged: false allowPrivilegeEscalation: false readOnlyRootFilesystem: true  capabilities:  drop:  - all As we are setting the containers to runAsNonRoot: true, we must even set the security contexts for deployments as non root. For eg:   securityContext:  runAsUser: 1001  runAsGroup: 2000  fsGroup: 3000 Use Seccomp profiles- Seccomp stands for secure computing mode and has been a feature of the Linux kernel since version 2.6.12. It can be used to sandbox the privileges of a process, restricting the calls (system calls) it is able to make from userspace into the kernel.  Refer: https://v1-23.docs.kubernetes.io/docs/tutorials/security/seccomp/  Use AppArmor profiles - AppArmor is a Linux kernel security module that supplements the standard Linux user and group based permissions to confine programs to a limited set of resources. AppArmor can be configured for any application to reduce its potential attack surface and provide greater in-depth defense  Refer: https://v1-23.docs.kubernetes.io/docs/tutorials/security/apparmor/|
|7                                                                                                                                                                       |A malicious user can perform a DOS attack by sending large no. of API requests to containers.  This would lead the containers to perform ‘noisy neighbors attack’ and can exhaust all the CPUs and Memory of other running pods/containers.|High                     |Denial of Service                |Ensure that you must assign resource quota policies at the namespace level, in order to limit overconsumption of the CPU and memory resources a pod is allowed to consume.|
|8                                                                                                                                                                       |A malicious user can take advantage of misconfigured security-configurations within the yaml deployment files, running on the k8s cluster.    Security misconfigurations lead to vulnerabilities within the system. Please run the kubesec tool to identify all the security misconfigurations.|High                     |Information Disclosure, Tampering, Spoofing, Denial of Service, Repudiation, Elevation of Privilege|Ensure that yaml files are scanned with tools like Kubesec (https://github.com/controlplaneio/kubesec) to secure the yaml file configurations.   Security misconfigurations report can be extracted by running the command:   “docker run -i kubesec/kubesec:v2 scan /dev/stdin < k8s/specified.yaml ”  “docker run -i kubesec/kubesec:v2 scan /dev/stdin < k8s/postgres.yaml ” Based on the findings, fix all the security misconfigurations present in the yaml files.|
|9                                                                                                                                                                       |A malicious user can exploit the high vulnerability found in postgres:12 docker image. Vulnerability was found with the help of trivy scanning tool. |High                     |Information Disclosure, Tampering|Ensure that docker images are scanned for vulnerable components using an image scanning tool like (https://github.com/aquasecurity/trivy).   Vulnerability report can be extracted by running the command:    “docker run aquasec/trivy image postgres:12”  Based on the findings, do not use the vulnerable docker images. Use alternative images which do not have any vulnerabilities.|
|10                                                                                                                                                                      |A malicious user can extract the service-account tokens or manually created secrets from the ETCD database if the secrets are not stored in encrypted format.  A malicious user has to make a simple API call to the kube API server to extract the secret details. Eg:  “ETCDCTL_API=3 etcdctl get /registry/secrets/default/secret1 [...] &#124; hexdump -C”|High                     |Information Disclosure           |Ensure that all the secrets (including service accounts) are stored in encrypted format in the ETCD database.  Refer https://v1-23.docs.kubernetes.io/docs/tasks/administer-cluster/encrypt-data/ to enable encryption of all secrets.  Do not forget to recreate all the secrets again, once the encryption is enabled. Otherwise, none of your existing secrets would be recognized by the k8s cluster.   Refer below command to recreate secrets:  Refer kubectl get secrets --all-namespaces -o json &#124; kubectl replace -f -|
|11                                                                                                                                                                      |A malicious user can perform MITM or can extract the sensitive information if TLS is not enabled for the Kubernetes Dashboard Access.|High                     |Information Disclosure, Tampering|Limit exposure of Kubernetes Dashboard. Disable public access via the internet.   Ensure the Dashboard Service Account is not open and accessible to users. Configure the login page and enable RBAC. Configure TLS for accessing Dashboard.|
|12                                                                                                                                                                      |A malicious user can exploit k8s-cluster security misconfigurations and can perform a variety of attacks based on misconfigured settings.   Security misconfigurations lead to vulnerabilities within the system. Please run the kube-bench tool to identify all the security misconfigurations.|High                     |Information Disclosure, Tampering, Spoofing, Denial of Service, Repudiation, Elevation of Privilege|Ensure that all the Kubernetes nodes (Master node & worker nodes. In our case, we only have Master node) are scanned separately with the Kube-bench tool to identify security misconfigurations.   Make sure to fix all the security misconfigurations based on the findings/report.  Refer: https://github.com/aquasecurity/kube-bench  Harder the cluster nodes to avoid security attacks   Refer: Secure the Master node configuration files  Secure the Worker node configuration files  Also, harden the other components of K8s cluster. We can refer to the same redhat reference for it.|
|13                                                                                                                                                                      |A malicious user could perform an unauthorized action and can deny its occurrence after completion of that unauthorized action.|Medium                   |Repudiation                      |In the absence of proper auditing in place for all the actions on k8s cluster, it might not be possible to prove the occurrence of an event, and also trace the events leading to an untoward/malicious event.  Ensure that we implement Falco as a service on the K8s cluster and configure necessary alerts for auditing of activities on K8s cluster.  Some examples of events that should trigger an alert would include: A shell is run inside a container  A container mounts a sensitive path from the host such as /proc  A sensitive file is unexpectedly read in a running container such as /etc/shadow  An outbound network connection is established   Refer Falco Project: https://falco.org/docs/|
|14                                                                                                                                                                      |A malicious user can pull/run rogue/malicious docker images which are not listed in the organization's registry.   A malicious user can execute the latest or beta tagged images and invite zero day attacks on k8s cluster.|High                     |Information Disclosure, Tampering, Spoofing, Denial of Service, Repudiation, Elevation of Privilege|Ensure that we have mechanisms in place to restrict developers to only run whitelisted images in Organization’s private repository.  We can implement Open Policy Agent (OPA) to restrict developers not to use: Images with latest/beta/insecure tag etc  Pull images from company’s local registry only   OPA is highly scalable. We can also configure OPA to enforce other restrictions on k8s cluster.  Refer: https://github.com/open-policy-agent/opa|
|15                                                                                                                                                                      |A malicious user can create and execute the deployment files with root privileges.   A malicious user can mount root file system/directories into the containers|High                     |Information Disclosure, Elevation of Privilege|Implement Pod Security Policies (PSP) at K8s cluster, to restrict admins/developers from running the applications with root privileges. We can implement PSP to resolve below security aspects:  Do not run application processes as root.  Do not allow privilege escalation.  Use a read-only root filesystem.  Do not use the host network or process space.  Drop unused and unnecessary Linux capabilities.  Use SELinux options for more fine-grained process controls.  Give each application its own Kubernetes Service Account.  Do not mount the service account credentials in a container if it does not need to access the Kubernetes API.  Refer: https://v1-23.docs.kubernetes.io/docs/concepts/security/pod-security-policy/  https://cloud.redhat.com/blog/12-kubernetes-configuration-best-practices#2-use-pod-security-policies-to-prevent-risky-containerspods-from-being-used|
|16                                                                                                                                                                      |A malicious user can shell into the application container & can perform root permissive/ escalated tasks to get hold of the host kernel (bare metal kernel access) on which the K8s cluster is running.|Medium                   |Information Disclosure, Elevation of Privilege, Spoofing|Ensure that applications are running as sandboxed pods to increase the level of security within K8s cluster.   Sandboxed pods use a different runtime class environment (refer https://v1-23.docs.kubernetes.io/docs/concepts/containers/runtime-class/) rather than the one getting used by K8s cluster.   This separates/sandbox the kernel level access for pods/containers running on K8s cluster.   Refer https://gvisor.dev/docs/ https://katacontainers.io/docs/|
|17                                                                                                                                                                      |A malicious user can spin up a container & can execute privileged tasks (as root).  A malicious user can abuse security misconfigurations within a Docker file & can perform a variety of attacks based on misconfigured settings.|High                     |Information Disclosure, Elevation of Privilege,|Ensure that Dockerfiles do not run with a root user.  Ensure that Dockerfiles are scanned with Docker Benchmark script. Ensure that all the security misconfigurations/findings are addressed, found in docker benchmark scan report.  Refer: https://github.com/docker/docker-bench-security|
|18                                                                                                                                                                      |A malicious user can extract the hard-coded JWT Signing Key information, used for signing the JWT tokens.  A malicious user can tamper/manipulate the JWT Token and then resign it with the extracted JWT signing key.   A malicious user can make unauthorized changes to the Books Store app (hardcoded IDs). After making the changes A malicious user can re-sign the token as he/she is already aware of JWT Signing Key.  A malicious user can extract/manipulate the hard-coded integer value of Booksl ID or other hardcoded IDs.|High                     |Information Disclosure, Tampering, Repudiation|Ensure that JWT Signing Key is not hard-coded in the code.   Instead, integrate a cloud vault solution like hashicorp and extract the signing key on the fly by making an API call to the hashicorp vault whenever needed (for signing JWT tokens). The symmetric key used here for signing the JWT tokens should be randomly generated, non-guessable and should be a minimum of 128 bit long.   Though not recommended (incase of legacy applications), to decrease the severity of vulnerability to medium, we can define a variable here and we can deploy the actual JWT signing key as a ‘secret’ in the kubernetes cluster. Ensure that the signing key is an opaque string (minimum 128 bit long) that is randomly generated and is non-guessable to the users/admin.  Ensure that instead of hardcoding the Book IDs (and other detials) in code, make use of variables and pass the value of variables when the containers are up & running as a ConfigMap. We can make use of ConfigMaps here. Ensure that Book IDs are randomly generated and are non-guessable to the users/admin.   Also, ensure that identification of Admin should be based on non-guessable opaque values, rather than keeping it as a simple bool value & hard-coding it in code . Make sure that the admin token should never be enabled by default. Infact, we should create an admin role and assign it to specific/required users.|
|19                                                                                                                                                                      |A malicious user can get hold of the symmetric key used for signing the JWT tokens. A malicious user can recreate/manipulate the JWT tokens once he/she is aware of JWT signing key|High                     |Information Disclosure, Tampering, Repudiation|Ensure that we give priority for using asymmetric algorithms like ES256 (minimum), RS384 (minimum) etc.  However, for business reasons, if we need to use symmetric algorithms, ensure that we use strong algorithms like HS384. Ensure that the symmetric key is shared via a vault service and should be rotated on regular intervals of time.  Ensure that the symmetric key (for HS algorithm) or private key (for ES/RS algorithms) is not hard-coded in the code (as explained in Vul 18). For signing the JWT tokens, integrate cloud vault service (like Hashicorp) and extract the signing key from the vault on the fly to sign the tokens.|
|20                                                                                                                                                                      |A user session will never end as JWT token is set to never expire.   A malicious user can extract the user information and can perform unauthorized actions using extracted information.|High                     |Spoofing                         |Do not use JWT for maintaining user sessions. Or set the JWT expiration time to be very less, maybe 30 seconds or less, to minimize the threat.   Implement a token block list that will be used to mimic the "logout" feature that exists with traditional session management systems.|
|21                                                                                                                                                                      |A malicious user might try to change the field “alg” in the token header for “none”.   A vulnerable application, after checking the JWT header and detecting “alg”: “none”, will accept this token without any verification as if it were legitimate, and as a result the tampered token will be accepted by the server.|High                     |Spoofing, Tampering              |Keep a white list of authorized algorithms on the application side and dismiss all tokens having a signature algorithm that is different from the one authorized on the server;  It is recommended to work with one algorithm only, e.g., HS256 or RS256.|
|22                                                                                                                                                                      |A malicious user can take advantage of weak configurations of the JWT, which might render the entire implementation weak.|Medium                   |Spoofing, Tampering, Lateral Movement|Always check against a whitelist the iss claim. When using the JWT you should be sure that it has been issued by someone you expected to issue it.  Always check the aud claim in the token and confront it with a whitelist. The server should expect that the token has been issued for an audience, which the server is part of. It should reject any requests that contain tokens intended for different audiences. This helps to mitigate attack vectors where one resource server would obtain a genuine Access Token intended for it, and then use it to gain access to resources on a different resource server, which would not normally be available to the original server.  Next, the consumer has to check the reserved "exp" and "nbf" claims to ensure that the JWT is valid.|
|23                                                                                                                                                                      |A malicious user can create its own JWT token by passing malicious values, as the token signatures are not getting verified. |High                     |Spoofing, Tampering, Non-repudiation|Ensure to use func (*Parser) Parse function, which would Parse, validate, and return a token. keyFunc will receive the parsed token and should return the key for validating. If everything is kosher, err will be nil.  Refer https://pkg.go.dev/github.com/dgrijalva/jwt-go#Parser.Parse|
|24                                                                                                                                                                      |A malicious user can access the Books website with a self created never expiring JWT Token. The authentication logic or session logic of the app is very weak and is based on manually passing a JWT token as an authentication reference in Authorization Header (Authorization Bearer).  A malicious user can also create and pass an admin token along with the API request in the Authorization Header.  JWT token is handled manually and can be misused for creating multiple unverified connections.|High                     |Spoofing                         |Do not use JWT for maintaining user sessions.   Implement the logic of username/password based authentication and we can maintain a securely designed User DB in Postgres DB.  Ensure that User Passwords are not stored in plain-text in Postgres DB   OR stored as a weak hashed value (signed using a weak hashing algorithm like MD5 or SHA1. Instead use SHA256 minimum)  OR using the same IV or salt for all the password hashes, while storing it in Postgres DB.  Refer below for implementing Authentication Logic:  https://mattermost.com/blog/how-to-build-an-authentication-microservice-in-golang-from-scratch/   Ensure that every user input is validated. Implement strong input validation using an allowlist. The validation must be done at the server side before using the input for any further operations.  Implement Roles in the application and assign the roles to specific users, instead of passing the Admin JWT tokens. Users (User IDs or Unique IDs or User Email) could be validated based on assigned permissions and can be granted with admin privileges/activities (after successful user permissions/role verification).|
|25                                                                                                                                                                      |A malicious user can perform MITM or can extract the sensitive information in the absence of TLS enabled connection for all the exposed APIs.|High                     |Information Disclosure,Tampering |Implement TLS 1.2 or above for all the exposed APIs to the end users.    A very simple reference for setting up a Go TLS server - https://joeldare.com/how-to-create-an-https-tls-server-in-go But this is not sufficient as it is using a self signed certificate. We need to purchase a certificate from trusted third party vendors like GlobalSign, Entrust etc for a registered domain.|
|26                                                                                                                                                                      |A malicious user can execute system calls from the app container as the code makes use of sig calls SIGINT, SIGTERM to bring down the app or to perform advanced attack techniques.|Medium                   |Tampering                        |Block the use of system calls from code.|
|27                                                                                                                                                                      |A malicious user can perform MITM or can extract the sensitive information if the connections between clients, proxies and databases are not encrypted in motion.|High                     |Information Disclosure,Tampering |The connections between clients, proxies and databases should be encrypted with Transport Layer Security (TLS 1.2 or above) to protect data in motion.|
|28                                                                                                                                                                      |Postgres DB credentials are constant and are not rotated. If A malicious user get hold of credentials, he can exploit the DB and can get hold of user sensitive information  There is no input validation in place for DB strings "host=%s port=% user=%s password=%s dbname=%s sslmode=disable  Therefore a user could pass any random values to perform a variety of attacks like DOS, SQL Injection, Access information leaked from error logs etc|High                     |Information Disclosure           |Ensure that DB credentials are rotated on regular intervals of time and shouldnt be kept constant. In the production environment, LDAP based authentication should be enabled.  Ensure that every user input is validated. Implement strong input validation using an allowlist. The validation must be done at the server side before using the input for any further operations.|
|29                                                                                                                                                                      |A malicious user might be able to perform injection attacks against PostgresDB by sending crafted input.  Injection flaws occur when untrusted data is sent to an interpreter as part of a command or query. The attacker’s malicious data can trick the interpreter into executing unintended commands or accessing data without proper authorization.  Therefore a user could pass any random values to perform a variety of attacks like DOS, SQL Injection, Access information leaked from error logs etc|High                     |Information Disclosure,Tampering |Implement parameterized queries based on rule-based configuration, to intercept and block queries based on multiple parameters (e.g., user or syntax) to, for example, prevent malicious attacks like injection from deleting all rows in a table or accessing restricted tables/columns.  Instead of using string formatting or concatenation to assemble the query, you use a placeholder for the parameters  Implement input validation of any data that is stored in PostgresDB.  Implement strong input validation using an allowlist.  Refer https://go.dev/doc/database/sql-injection  https://www.stackhawk.com/blog/golang-sql-injection-guide-examples-and-prevention/|
|30                                                                                                                                                                      |A malicious user can extract the sensitive user information from the Postgres DB if the encryption at rest is not enabled.|High                     |Information Disclosure           |Implement encryption at rest in Postgres DB. Use strong ciphers for encryption eg. AES256-GCM. Tables could be encrypted with Advanced Encryption Standard (AES) algorithms to protect data at rest.|
|31                                                                                                                                                                      |A malicious user might be able to steal data from PostgresDB and then misuse it, if data masking is not enabled.|Medium                   |Information Disclosure           |Implement data masking techniques which can be configured to hide sensitive data (e.g., PII/SPI)|
|32                                                                                                                                                                      |In the absence of proper logging of transactions, it might not be possible to prove the occurrence of an event, and also trace the events leading to an untoward event.  A malicious user might deny performing an action in PostgresDB|Medium                   |Repudiation                      |Enable auditing to track all database events – connections, queries (DML/DDL/DCL) and tables accessed – logging the time, username and host, database and operation, and more. In addition to local files, remote files are supported via syslog and rsyslog – often to aggregate database events from multiple servers and/or restrict access.|
|33                                                                                                                                                                      |A malicious user might send unlimited requests to the APIs to overwhelm the APIs, and potentially crash the APIs|High                     |Denial of Service                |Implement rate-limiting and throttling of APIs.   Restrict multiple login attempts by defining the number of allowed login attempts in a given period of time on Books Website (to avoid brute force attacks).   Implement a mechanism of captcha identification/puzzle solving to avoid bot attacks.  Implement High-Availability, and Redundancy in the system.  Refer https://gobyexample.com/rate-limiting|
|34                                                                                                                                                                      |A malicious user might be able to gain additional privileges and perform unauthorized actions w.r.t app & its functioning.  If there are any insecure direct object references, an attacker might be able to provide crafted input to access objects directly, and bypass access control.|High                     |Elevation of Privilege           |Implement Object level authorisation checks (proper parameter validation) in every function that accesses a data source using an untrusted input.  Do not rely on IDs that the client sends. Use IDs stored in the session object instead.  Check authorization for each client request to access the database.  Use random IDs that cannot be guessed (UUIDs).  Refer: https://www.appknox.com/blog/understanding-insecure-direct-object-references-idor|
|35                                                                                                                                                                      |An abnormal action performed by the adversary might lead to an exception that can leak internal working details of the system The adversary can then leverage such helpful error messages to gain an understanding of the working of the system, and then craft further attacks against the system.|Medium                   |Information Disclosure           |Implement exception handling. All exceptions must be handled gracefully. During any exceptions occurring, the system must provide the user with a generic error message, and never internal details of the error, or detailed error messages.|
|36                                                                                                                                                                      |A malicious user might be able to steal sensitive information if Logs contain PII or sensitive data  |High                     |Information Disclosure           |Do not log any sensitive data or PII that can be used to uniquely identify any user or customer.  Ensure that logging should not include any sensitive information.|
|37                                                                                                                                                                      |A malicious user might be able to perform injection attacks against the Books app APIs by sending crafted input to the system (POST Requests) & APIs.|High                     |Tampering                        |Implement strong input validation using an allowlist. The validation must be done at the server side before using the input for any further operations.  Refer https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html|
|38                                                                                                                                                                      |A malicious user might be able to steal sensitive data or launch crafted attacks due to missing security configuration of HTTP/API responses|Medium                   |Tampering                        |Include the following Secure headers in all HTTP responses: Permissions Policy HTTP Strict Transport Security X-Frame-Options X-Content-Type-Options Content-Security-Policy X-Permitted-Cross-Domain-Policies Referrer-Policy Clear-Site-Data Cross-Origin-Embedder-Policy Cross-Origin-Opener-Policy Cross-Origin-Resource-Policy|
|39                                                                                                                                                                      |In the absence of proper logging of transactions, it might not be possible to prove the occurrence of an event, and also trace the events leading to an untoward event.  A malicious user might deny performing an action after completion of the transaction (via API request/Response)|Medium                   |Repudiation                      |Implement proper auditing and logging. All transactions, whether originating from the user, or originating from the server should be logged. Logging must be performed at all servers /processes involved, and not just on some specific servers|

